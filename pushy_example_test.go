package pushy_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cyberhck/pushy"
	"gopkg.in/jarcoal/httpmock.v1"
)

func ExampleGetDefaultAPIEndpoint() {
	fmt.Print(pushy.GetDefaultAPIEndpoint())
	// Output:
	// https://api.pushy.me
}

func setupNotifyStuff() func() {
	httpmock.Activate()
	apiToken := "API_TOKEN"
	endpoint := fmt.Sprintf("https://api.pushy.me/push?api_key=%s", apiToken)
	expectedResponse := `{"success":true, "id":"some_id"}`
	httpmock.RegisterResponder("POST", endpoint, httpmock.NewStringResponder(200, expectedResponse))
	return httpmock.DeactivateAndReset
}
func setupNotifyDeletionStuff() func() {
	httpmock.Activate()
	apiToken := "API_TOKEN"
	endpoint := fmt.Sprintf("https://api.pushy.me/pushes/some_id?api_key=%s", apiToken)
	expectedResponse := `{"success":true, "id":"some_id"}`
	httpmock.RegisterResponder("DELETE", endpoint, httpmock.NewStringResponder(http.StatusOK, expectedResponse))
	return httpmock.DeactivateAndReset
}
func setupDeviceInfoStuff() func() {
	httpmock.Activate()
	expectedResponse := `
{
  "device": {
    "date": 1445207358,
    "platform": "android"
  },
  "subscriptions": [
    "news",
    "media"
  ],
  "presence": {
    "online": true,
    "last_active": {
      "date": 1464006925,
      "seconds_ago": 215
    }
  },
  "pending_notifications": [
    {
      "id": "5742fe0407c3674e226892f9",
      "date": 1464008196,
      "payload": {
        "message": "Hello World!"
      },
      "expiration": 1466600196
    }
  ]
}
`
	endpoint := fmt.Sprintf("https://api.pushy.me/devices/DEVICE_ID?api_key=%s", "API_TOKEN")
	httpmock.RegisterResponder("GET", endpoint, httpmock.NewStringResponder(http.StatusOK, expectedResponse))
	return httpmock.DeactivateAndReset
}
func setupDevicePresenceStuff() func() {
	httpmock.Activate()
	expectedResponse := `{
  "presence": [
    {
      "id": "a6f36efb913f1def30c6",
      "online": false,
      "last_active": 1429406442
    }
  ]
}`
	endpoint := "https://api.pushy.me/devices/presence?api_key=API_TOKEN"
	httpmock.RegisterResponder("POST", endpoint, httpmock.NewStringResponder(http.StatusOK, expectedResponse))
	return httpmock.DeactivateAndReset
}
func setupNotificationStatusStuff() func() {
	httpmock.Activate()
	expectedResponse := `
{
  "push": {
    "date": 1464003935,
    "payload": {
      "message": "Hello World!"
    },
    "expiration": 1466595935,
    "pending_devices": [
      "fe8f7b2c102e883e5b41d2"
    ]
  }
}
`
	endpoint := "https://api.pushy.me/pushes/PUSH_ID?api_key=API_TOKEN"
	httpmock.RegisterResponder("GET", endpoint, httpmock.NewStringResponder(http.StatusOK, expectedResponse))
	return httpmock.DeactivateAndReset
}
func setupSubscribeToTopicStuff() func() {
	httpmock.Activate()
	expectedResponse := `{"success": true}`
	endpoint := "https://api.pushy.me/devices/subscribe?api_key=API_TOKEN"
	httpmock.RegisterResponder("POST", endpoint, httpmock.NewStringResponder(http.StatusOK, expectedResponse))
	return httpmock.DeactivateAndReset
}
func setupUnSubscribeFromTopicStuff() func() {
	httpmock.Activate()
	expectedResponse := `{"success": true}`
	endpoint := "https://api.pushy.me/devices/unsubscribe?api_key=API_TOKEN"
	httpmock.RegisterResponder("POST", endpoint, httpmock.NewStringResponder(http.StatusOK, expectedResponse))
	return httpmock.DeactivateAndReset
}

func ExamplePushy_NotifyDevice() {
	cleaner := setupNotifyStuff()
	defer cleaner()
	sdk := pushy.Create("API_TOKEN", pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(100 * time.Millisecond))
	status, _, _ := sdk.NotifyDevice(pushy.SendNotificationRequest{})
	fmt.Println(status.Success)
	fmt.Println(status.ID)
	// Output:
	// true
	// some_id
}

func ExamplePushy_DeleteNotification() {
	cleaner := setupNotifyDeletionStuff()
	defer cleaner()
	sdk := pushy.Create("API_TOKEN", pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(100 * time.Millisecond))
	status, _, _ := sdk.DeleteNotification("some_id")
	fmt.Print(status.Success)
	// Output:
	// true
}

func ExamplePushy_DeviceInfo() {
	cleaner := setupDeviceInfoStuff()
	defer cleaner()
	sdk := pushy.Create("API_TOKEN", pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(10 * time.Millisecond))

	res, _, _ := sdk.DeviceInfo("DEVICE_ID")
	fmt.Println(res.Presence.Online)
	fmt.Println(res.Presence.LastActive.SecondsAgo)
	// Output:
	// true
	// 215
}

func ExamplePushy_DevicePresence() {
	cleaner := setupDevicePresenceStuff()
	defer cleaner()
	sdk := pushy.Create("API_TOKEN", pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(10 * time.Millisecond))
	presence, _, _ := sdk.DevicePresence("DEVICE_ID")
	fmt.Println(presence.Presence[0].ID)
	fmt.Println(presence.Presence[0].Online)
	fmt.Println(presence.Presence[0].LastActive)
	// Output:
	// a6f36efb913f1def30c6
	// false
	// 1429406442
}

func ExamplePushy_NotificationStatus() {
	cleaner := setupNotificationStatusStuff()
	defer cleaner()
	sdk := pushy.Create("API_TOKEN", pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(10 * time.Millisecond))
	status, _, _ := sdk.NotificationStatus("PUSH_ID")
	fmt.Println(status.Push.Expiration)
	fmt.Println(status.Push.Date)
	fmt.Println(status.Push.Payload) // todo isn't working for some reason.
	// Output:
	// 1466595935
	// 1464003935
	// map[message:Hello World!]
}

func ExamplePushy_SubscribeToTopic() {
	cleaner := setupSubscribeToTopicStuff()
	defer cleaner()
	sdk := pushy.Create("API_TOKEN", pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(10 * time.Millisecond))
	subscription, _, _ := sdk.SubscribeToTopic("DEVICE_ID", "TOPIC")
	fmt.Println(subscription.Success)
	// Output:
	// true
}

func ExamplePushy_UnsubscribeFromTopic() {
	cleaner := setupUnSubscribeFromTopicStuff()
	defer cleaner()
	sdk := pushy.Create("API_TOKEN", pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(10 * time.Millisecond))
	subscription, _, _ := sdk.UnsubscribeFromTopic("DEVICE_ID", "TOPIC")
	fmt.Println(subscription.Success)
	// Output:
	// true
}
