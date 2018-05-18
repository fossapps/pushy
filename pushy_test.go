package pushy_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cyberhck/pushy"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

type endpoint struct {
	method string
	url    string
}

// region SDK related settings
func TestCreate(t *testing.T) {
	table := []struct {
		token    string
		endpoint string
	}{
		{
			token:    "token",
			endpoint: "http://example.com",
		},
		{
			token:    "token2",
			endpoint: "http://api.example.com",
		},
	}
	for _, data := range table {
		sdk := pushy.Create(data.token, data.endpoint)
		if sdk.APIToken != data.token {
			t.Error("api token not set")
		}
		if sdk.APIEndpoint != data.endpoint {
			t.Error("endpoint is not set")
		}
	}
}

func TestGetDefaultApiEndpoint(t *testing.T) {
	assert.Equal(t, "https://api.pushy.me", pushy.GetDefaultAPIEndpoint())
}

func TestGetDefaultHTTPClient(t *testing.T) {
	client := pushy.GetDefaultHTTPClient(1 * time.Millisecond)
	if client == nil {
		t.Error("client is nil")
	}
}

func TestPushy_SetHTTPClient(t *testing.T) {
	sdk := pushy.Create("token", pushy.GetDefaultAPIEndpoint())
	client := pushy.GetDefaultHTTPClient(3 * time.Millisecond)
	sdk.SetHTTPClient(client)
	newClient := sdk.GetHTTPClient()
	if newClient != client {
		t.Error("SetHTTPClient should use exact same http client provided by user")
	}
}

func TestPushy_GetHTTPClient(t *testing.T) {
	sdk := pushy.Create("token", pushy.GetDefaultAPIEndpoint())
	client := pushy.GetDefaultHTTPClient(3 * time.Millisecond)
	sdk.SetHTTPClient(client)
	newClient := sdk.GetHTTPClient()
	if newClient != client {
		t.Error("GetHTTPClient should return underlying http client")
	}
}

func TestPushy_GetHTTPClient2(t *testing.T) {
	sdk := pushy.Create("token", pushy.GetDefaultAPIEndpoint())
	client := sdk.GetHTTPClient()
	if client != nil {
		t.Error("should not use any http client which isn't set by user")
	}
}

// endregion

func TestEverythingHandlesBadRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	apiToken := "API_TOKEN"
	deviceToken := "DEVICE"
	endpoints := getEndpointsDefinitions()
	body := `{"error":"not found / bad token"}`
	for _, endpoint := range endpoints {
		url := fmt.Sprintf("%s%s", "https://api.pushy.me", endpoint.url)
		httpmock.RegisterResponder(endpoint.method, url, httpmock.NewStringResponder(http.StatusBadRequest, body))
	}
	sdk := pushy.Create(apiToken, pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(10 * time.Millisecond))

	deviceInfo, pushyErr, err := sdk.DeviceInfo(deviceToken)
	Assert := assert.New(t)
	Assert.Contains(err.Error(), "400")
	Assert.NotNil(pushyErr)
	Assert.Nil(deviceInfo)

	devicePresence, pushyErr, err := sdk.DevicePresence("TOKEN")
	Assert.Contains(err.Error(), "400")
	Assert.NotNil(pushyErr)
	Assert.Nil(devicePresence)

	notificationStatus, pushyErr, err := sdk.NotificationStatus("TOKEN")
	Assert.Contains(err.Error(), "400")
	Assert.NotNil(pushyErr)
	Assert.Nil(notificationStatus)

	deleteNotification, pushyErr, err := sdk.DeleteNotification("TOKEN")
	Assert.Contains(err.Error(), "400")
	Assert.NotNil(pushyErr)
	Assert.Nil(deleteNotification)

	subscription, pushyErr, err := sdk.SubscribeToTopic("S", "topic")
	Assert.Contains(err.Error(), "400")
	Assert.NotNil(pushyErr)
	Assert.Nil(subscription)

	unSubscription, pushyErr, err := sdk.UnsubscribeFromTopic("S", "topic")
	Assert.Contains(err.Error(), "400")
	Assert.NotNil(pushyErr)
	Assert.Nil(unSubscription)

	notifyDevice, pushyErr, err := sdk.NotifyDevice(pushy.SendNotificationRequest{})
	Assert.Contains(err.Error(), "400")
	Assert.NotNil(pushyErr)
	Assert.Nil(notifyDevice)
}

func TestEverythingHandlesNetworkError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	apiToken := "API_TOKEN"
	deviceToken := "DEVICE"
	endpoints := getEndpointsDefinitions()
	for _, endpoint := range endpoints {
		url := fmt.Sprintf("%s%s", "https://api.pushy.me", endpoint.url)
		httpmock.RegisterResponder(endpoint.method, url, httpmock.NewErrorResponder(errors.New("ERR CONN RESET")))
	}
	sdk := pushy.Create(apiToken, pushy.GetDefaultAPIEndpoint())

	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(10 * time.Millisecond))

	deviceInfo, pushyErr, err := sdk.DeviceInfo(deviceToken)
	Assert := assert.New(t)
	Assert.Contains(err.Error(), "ERR CONN RESET")
	Assert.Nil(pushyErr)
	Assert.Nil(deviceInfo)

	devicePresence, pushyErr, err := sdk.DevicePresence("TOKEN")
	Assert.Contains(err.Error(), "ERR CONN RESET")
	Assert.Nil(pushyErr)
	Assert.Nil(devicePresence)

	notificationStatus, pushyErr, err := sdk.NotificationStatus("TOKEN")
	Assert.Contains(err.Error(), "ERR CONN RESET")
	Assert.Nil(pushyErr)
	Assert.Nil(notificationStatus)

	deleteNotification, pushyErr, err := sdk.DeleteNotification("TOKEN")
	Assert.Contains(err.Error(), "ERR CONN RESET")
	Assert.Nil(pushyErr)
	Assert.Nil(deleteNotification)

	invalidUrlParameter, pushyErr, err := sdk.DeleteNotification("TOKE%%%%%N")
	Assert.Contains(err.Error(), "invalid URL escape")
	Assert.Nil(pushyErr)
	Assert.Nil(invalidUrlParameter)

	subscription, pushyErr, err := sdk.SubscribeToTopic("S", "topic")
	Assert.Contains(err.Error(), "ERR CONN RESET")
	Assert.Nil(pushyErr)
	Assert.Nil(subscription)

	unSubscription, pushyErr, err := sdk.UnsubscribeFromTopic("S", "topic")
	Assert.Contains(err.Error(), "ERR CONN RESET")
	Assert.Nil(pushyErr)
	Assert.Nil(unSubscription)

	notifyDevice, pushyErr, err := sdk.NotifyDevice(pushy.SendNotificationRequest{})
	Assert.Contains(err.Error(), "ERR CONN RESET")
	Assert.Nil(pushyErr)
	Assert.Nil(notifyDevice)
}

// region API communication
func TestPushy_DeviceInfo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	apiToken := "API_TOKEN"
	deviceToken := "DEVICE"
	endpoint := fmt.Sprintf("https://api.pushy.me/devices/%s?api_key=%s", deviceToken, apiToken)
	expectedResponse := `
{
  "device": {
    "date": 1000,
    "platform": "android"
  },
  "subscriptions": [
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
	httpmock.RegisterResponder("GET", endpoint, httpmock.NewStringResponder(200, expectedResponse))
	Assert := assert.New(t)
	sdk := pushy.Create(apiToken, pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(100 * time.Millisecond))
	info, pushyError, err := sdk.DeviceInfo(deviceToken)
	Assert.Nil(pushyError)
	Assert.Nil(err)
	Assert.Equal("android", info.Device.Platform)
	Assert.Equal(1000, info.Device.Date)
	Assert.Equal("media", info.Subscriptions[0])
	Assert.Equal(true, info.Presence.Online)
	Assert.Equal(215, info.Presence.LastActive.SecondsAgo)
	Assert.Equal("5742fe0407c3674e226892f9", info.PendingNotifications[0].ID)
}

func TestPushy_DevicePresence(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	apiToken := "API_TOKEN"
	deviceToken := "DEVICE"
	endpoint := fmt.Sprintf("https://api.pushy.me/devices/presence?api_key=%s", apiToken)
	expectedResponse := `
{
  "presence": [
    {
      "id": "a6f36efb913f1def30c6",
      "online": false,
      "last_active": 1429406442
    }
  ]
}
`
	httpmock.RegisterResponder("POST", endpoint, httpmock.NewStringResponder(200, expectedResponse))
	Assert := assert.New(t)
	sdk := pushy.Create(apiToken, pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(100 * time.Millisecond))

	info, _, _ := sdk.DevicePresence(deviceToken)
	Assert.Equal(false, info.Presence[0].Online)
	Assert.Equal("a6f36efb913f1def30c6", info.Presence[0].ID)
	Assert.Equal(1429406442, info.Presence[0].LastActive)
}

func TestPushy_NotificationStatus(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	apiToken := "API_TOKEN"
	endpoint := fmt.Sprintf("https://api.pushy.me/pushes/PUSH_ID?api_key=%s", apiToken)
	expectedResponse := `
{
  "push": {
    "date": 100,
    "payload": {
      "message": "Hello World!"
    },
    "expiration": 105,
    "pending_devices": [
      "device_id"
    ]
  }
}
`
	httpmock.RegisterResponder("GET", endpoint, httpmock.NewStringResponder(200, expectedResponse))
	Assert := assert.New(t)
	sdk := pushy.Create(apiToken, pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(100 * time.Millisecond))

	status, _, _ := sdk.NotificationStatus("PUSH_ID")
	Assert.Equal(100, status.Push.Date)
	payloadMap, _ := status.Push.Payload.(map[string]interface{})
	Assert.Contains(payloadMap["message"], "Hello World!")
	Assert.Equal(105, status.Push.Expiration)
	Assert.Equal("device_id", status.Push.PendingDevices[0])
}

func TestPushy_DeleteNotification(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	apiToken := "API_TOKEN"
	endpoint := fmt.Sprintf("https://api.pushy.me/pushes/PUSH_ID?api_key=%s", apiToken)
	expectedResponse := `{"success":true}`
	httpmock.RegisterResponder("DELETE", endpoint, httpmock.NewStringResponder(200, expectedResponse))
	sdk := pushy.Create(apiToken, pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(100 * time.Millisecond))

	status, _, _ := sdk.DeleteNotification("PUSH_ID")
	assert.Equal(t, true, status.Success)
}

func TestPushy_SubscribeToTopic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	apiToken := "API_TOKEN"
	endpoint := fmt.Sprintf("https://api.pushy.me/devices/subscribe?api_key=%s", apiToken)
	expectedResponse := `{"success":true}`
	httpmock.RegisterResponder("POST", endpoint, httpmock.NewStringResponder(200, expectedResponse))
	Assert := assert.New(t)
	sdk := pushy.Create(apiToken, pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(100 * time.Millisecond))

	status, _, _ := sdk.SubscribeToTopic("TOKEN", "topic")
	Assert.Equal(true, status.Success)
}

func TestPushy_UnsubscribeFromTopic(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	apiToken := "API_TOKEN"
	endpoint := fmt.Sprintf("https://api.pushy.me/devices/unsubscribe?api_key=%s", apiToken)
	expectedResponse := `{"success":true}`
	httpmock.RegisterResponder("POST", endpoint, httpmock.NewStringResponder(200, expectedResponse))
	Assert := assert.New(t)
	sdk := pushy.Create(apiToken, pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(100 * time.Millisecond))

	status, _, _ := sdk.UnsubscribeFromTopic("TOKEN", "topic")
	Assert.Equal(true, status.Success)
}

func TestPushy_NotifyDevice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	apiToken := "API_TOKEN"
	endpoint := fmt.Sprintf("https://api.pushy.me/push?api_key=%s", apiToken)
	expectedResponse := `{"success":true, "id":"some_id"}`
	httpmock.RegisterResponder("POST", endpoint, httpmock.NewStringResponder(200, expectedResponse))
	Assert := assert.New(t)
	sdk := pushy.Create(apiToken, pushy.GetDefaultAPIEndpoint())
	sdk.SetHTTPClient(pushy.GetDefaultHTTPClient(100 * time.Millisecond))

	status, _, _ := sdk.NotifyDevice(pushy.SendNotificationRequest{})
	Assert.Equal(true, status.Success)
	Assert.Equal("some_id", status.ID)
}

// endregion

func getEndpointsDefinitions() []endpoint {
	return []endpoint{
		{
			method: "GET",
			url:    "/devices/DEVICE?api_key=API_TOKEN",
		},
		{
			method: "POST",
			url:    "/devices/presence?api_key=API_TOKEN",
		},
		{
			method: "GET",
			url:    "/pushes/TOKEN?api_key=API_TOKEN",
		},
		{
			method: "DELETE",
			url:    "/pushes/TOKEN?api_key=API_TOKEN",
		},
		{
			method: "POST",
			url:    "/devices/subscribe?api_key=API_TOKEN",
		},
		{
			method: "POST",
			url:    "/devices/unsubscribe?api_key=API_TOKEN",
		},
		{
			method: "POST",
			url:    "/push?api_key=API_TOKEN",
		},
	}
}
