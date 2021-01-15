package pushy

import (
	"io"
	"net/http"
)

// IPushyClient interface to implement to qualify as a Pushy Client
type IPushyClient interface {
	SetHTTPClient(client IHTTPClient)
	GetHTTPClient() IHTTPClient
	DeviceInfo(deviceID string) (*DeviceInfo, *Error, error)
	DevicePresence(deviceID ...string) (*DevicePresenceResponse, *Error, error)
	NotificationStatus(pushID string) (*NotificationStatus, *Error, error)
	DeleteNotification(pushID string) (*SimpleSuccess, *Error, error)
	SubscribeToTopic(deviceID string, topics ...string) (*SimpleSuccess, *Error, error)
	UnsubscribeFromTopic(token string, topics ...string) (*SimpleSuccess, *Error, error)
	NotifyDevice(request SendNotificationRequest) (*NotificationResponse, *Error, error)
}

// Pushy is a basic struct with two configs: APIToken and APIEndpoint
// implements IPushyClient interface
type Pushy struct {
	APIToken    string
	APIEndpoint string
	httpClient  IHTTPClient
}

// Error are simple error responses returned from pushy if request isn't valid
type Error struct {
	Error string `json:"error"`
}

// Device is basic representation of a device
type Device struct {
	Date     int    `json:"date"`
	Platform string `json:"platform"`
}

// DeviceInfo is a basic structure which has additional info
type DeviceInfo struct {
	Device        Device   `json:"device"`
	Subscriptions []string `json:"subscriptions"`
	Presence      struct {
		Online     bool `json:"online"`
		LastActive struct {
			Date       int `json:"date"`
			SecondsAgo int `json:"seconds_ago"`
		} `json:"last_active"`
	} `json:"presence"`
	PendingNotifications []Notification `json:"pending_notifications"`
}

// Notification is a basic representation of a notification
type Notification struct {
	ID      string      `json:"id"`
	Date    int         `json:"date"`
	Payload interface{} `json:"payload"`
}

// DevicePresenceResponse is representation of device presence response from pushy
type DevicePresenceResponse struct {
	Presence []Presence `json:"presence"`
}

// Presence is a basic representation of a device's presence
type Presence struct {
	ID         string `json:"id"`
	Online     bool   `json:"online"`
	LastActive int    `json:"last_active"`
}

// NotificationStatus is a basic status info of a Notification
type NotificationStatus struct {
	Push struct {
		Date           int         `json:"date"`
		Payload        interface{} `json:"payload"`
		Expiration     int         `json:"expiration"`
		PendingDevices []string    `json:"pending_devices"`
	} `json:"push"`
}

// SimpleSuccess is a response from pushy when our request is accepted
type SimpleSuccess struct {
	Success bool `json:"success"`
}

// DeviceSubscriptionRequest is representation of a request we send for subscribing a new device to topic
type DeviceSubscriptionRequest struct {
	Token  string   `json:"token"`
	Topics []string `json:"topics"`
}

// SendNotificationRequest is representation of data to be sent to pushy service to create new notification
type SendNotificationRequest struct {
	To                  []string               `json:"to"`
	Data                map[string]interface{} `json:"data"`
	TimeToLive          int                    `json:"time_to_live"`
	IOSMutableContent   bool                   `json:"mutable_content"`
	IOSContentAvailable bool                   `json:"content_available"`
	IOSNotification     *IOSNotification       `json:"notification"`
}

// IOSNotification is a basic data for notification for iOS devices
// it's internally called notification,
// is represented as IOSNotification to communicate that this only applies for iOS devices
type IOSNotification struct {
	Body         string   `json:"body"`
	Badge        int      `json:"badge"`
	Sound        string   `json:"sound"`
	Title        string   `json:"title"`
	Category     string   `json:"category"`
	LocKey       string   `json:"loc_key"`
	LocArgs      []string `json:"loc_args"`
	TitleLocKey  string   `json:"title_loc_key"`
	TitleLocArgs []string `json:"title_loc_args"`
}

// NotificationResponse is a simple response from server when a new notification is created
type NotificationResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

// DevicePresenceRequest is representation of request needed to get information of device(s)'s presence
type DevicePresenceRequest struct {
	Tokens []string `json:"tokens"`
}

// IHTTPClient is signature needed for an object to be acceptable as a http client.
type IHTTPClient interface {
	Get(string) (*http.Response, error)
	Post(string, string, io.Reader) (*http.Response, error)
	Do(*http.Request) (*http.Response, error)
}
