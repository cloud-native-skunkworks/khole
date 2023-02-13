package output

import (
	"fmt"

	corev1alpha1 "github.com/cloud-native-skunkworks/khole/api/v1alpha1"
	"github.com/slack-go/slack"
	corev1 "k8s.io/api/core/v1"
)

func sendSlackAlert(pod *corev1.Pod, token string, channelID string, message string) error {
	// Create a new client to slack by giving token
	client := slack.New(token, slack.OptionDebug(false))
	// Create the Slack attachment that we will send to the channel
	attachment := slack.Attachment{
		Pretext: "KHole helper message",
		Text:    "Warning, a pod is in a bad state",
		// Color Styles the Text, making it possible to have like Warnings etc.
		Color: "#36a64f",
		// Fields are Optional extra data!
		Fields: []slack.AttachmentField{
			{
				Title: "Name",
				Value: pod.Name,
			},
			{
				Title: "Namespace",
				Value: pod.Namespace,
			},
			{
				Title: "Status",
				Value: message,
			},
		},
	}
	// PostMessage will send the message away.
	// First parameter is just the channelID, makes no sense to accept it
	_, timestamp, err := client.PostMessage(
		channelID,
		// uncomment the item below to add a extra Header to the message, try it out :)
		//slack.MsgOptionText("New message from bot", false),
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Message sent at %s", timestamp)
	return nil
}

func SendAlert(pod *corev1.Pod,
	kholeconfiguration *corev1alpha1.KHoleConfiguration,
	message string) error {

	if kholeconfiguration.Spec.Output.Slack.Token != "" &&
		kholeconfiguration.Spec.Output.Slack.ChannelID != "" {
		err := sendSlackAlert(pod, kholeconfiguration.Spec.Output.Slack.Token,
			kholeconfiguration.Spec.Output.Slack.ChannelID,
			message)
		if err != nil {
			return err
		}
	}
	return nil
}
