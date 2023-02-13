# KHole 

Sends alerts from crashed K8s pods to Slack. This is a simple tool that can be used to monitor your Kubernetes cluster.

```
apiVersion: core.cloudnativeskunkworks.io/v1alpha1
kind: KHoleConfiguration
metadata:
  name: kholeconfiguration
spec:
  output:
    slack:
      channelID: ""
      token: ""
```

## Example

When a pod crashes or doesn't work properly, KHole will send a message to the configured outputs.

<img src="images/01.png" width="800px;">

An example Slack message

<img src="images/02.png" width="800px;">

### Installation

```
git clone https://github.com/cloud-native-skunkworks/khole.git
make install 
make deploy
```

Apply the CR with your details

```
kubectl apply -f config/samples/core_v1alpha1_kholeconfiguration.yaml
```

### Development

This is a typical kubebuilder project. 
It has no webhook so you can just `go run` as long as you have a valid `KUBECONFIG` env. You'll need to install the custom resources.

If you want to build the docker image just override the `IMG` env variable in the make file e.g. `IMG=your-registry/khole:latest make docker-build`