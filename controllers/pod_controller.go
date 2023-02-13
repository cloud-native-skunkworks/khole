/*
Copyright 2023 AlexsJones.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/cloud-native-skunkworks/khole/pkg/output"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PodReconciler reconciles a Pod object
type PodReconciler struct {
	client.Client
	Scheme                          *runtime.Scheme
	KHoleConfigurationReconcilerRef *KHoleConfigurationReconciler
}

//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=pods/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Pod object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *PodReconciler) checkAndSend(pod *corev1.Pod,
	ctx context.Context, l logr.Logger, message string) (ctrl.Result, error) {

	if r.KHoleConfigurationReconcilerRef != nil && r.KHoleConfigurationReconcilerRef.lastKnownConfiguration != nil {
		// Check if the pod has already been alerted
		if _, ok := pod.Annotations["khole.io/alerted"]; ok {
			return ctrl.Result{}, nil
		}
		// Fire alert
		if err := output.SendAlert(
			pod, r.KHoleConfigurationReconcilerRef.
				lastKnownConfiguration, message); err != nil {
			l.Error(err, "Unable to send alert")
			return ctrl.Result{}, err
		}
		// Annotate that the pod has been alerted
		pod.Annotations["khole.io/alerted"] = "true"
		if err := r.Update(ctx, pod); err != nil {
			l.Error(err, "Unable to annotate the pod")
			return ctrl.Result{}, err
		}
	} else {
		l.Info("Unable to send alert, no configuration found")
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}
func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// TODO(user): your logic here
	pod := &corev1.Pod{}
	err := r.Get(ctx, req.NamespacedName, pod)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// get container status
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.State.Waiting != nil {
			if containerStatus.State.Waiting.Reason == "PodInitializing" {
				continue
			}
			r.checkAndSend(pod, ctx, l, containerStatus.State.Waiting.Reason)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(r)
}
