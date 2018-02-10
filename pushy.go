// Package pushy can be used to communicate with pushy service easily,
// saving time which would require to figure out the data type, endpoint, http method and shape of request to make
package pushy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Create is a helper method to initialize a simple Pushy struct
//  &Pushy{ApiToken: "token", ApiEndpoint: "https://api.pushy.me"}
// can be used
func Create(APIToken string, APIEndpoint string) *Pushy {
	return &Pushy{
		APIToken:    APIToken,
		APIEndpoint: APIEndpoint,
	}
}

// GetDefaultAPIEndpoint returns the default api endpoint where they're hosted.
func GetDefaultAPIEndpoint() string {
	return "https://api.pushy.me"
}

// GetDefaultHTTPClient returns a httpClient with configured timeout
func GetDefaultHTTPClient(timeout time.Duration) IHTTPClient {
	client := http.Client{
		Timeout: timeout,
	}
	return IHTTPClient(&client)
}

// SetHTTPClient sets a http client, it's useful when you're using sandboxed env like appengine
// this is required to do, pushy won't automatically use the default http client.
func (p *Pushy) SetHTTPClient(client IHTTPClient) {
	p.httpClient = client
}

// GetHTTPClient returns the client which is being used with pushy
func (p *Pushy) GetHTTPClient() IHTTPClient {
	return p.httpClient
}

// DeviceInfo returns information about a particular device
func (p Pushy) DeviceInfo(deviceID string) (*DeviceInfo, *Error, error) {
	url := p.APIEndpoint + "/devices/" + deviceID + "?api_key=" + p.APIToken
	var errResponse *Error
	var info *DeviceInfo
	err := get(p.httpClient, url, &info, &errResponse)
	return info, errResponse, err
}

// DevicePresence returns data about presence of a data
func (p *Pushy) DevicePresence(deviceID ...string) (*DevicePresenceResponse, *Error, error) {
	url := p.APIEndpoint + "/devices/presence?api_key=" + p.APIToken
	var devicePresenceResponse *DevicePresenceResponse
	var pushyErr *Error
	err := post(p.httpClient, url, DevicePresenceRequest{Tokens: deviceID}, &devicePresenceResponse, &pushyErr)
	return devicePresenceResponse, pushyErr, err
}

// NotificationStatus returns status of a particular notification
func (p *Pushy) NotificationStatus(pushID string) (*NotificationStatus, *Error, error) {
	url := fmt.Sprintf(p.APIEndpoint+"/pushes/%s?api_key=%s", pushID, p.APIToken)
	var errResponse *Error
	var status *NotificationStatus
	err := get(p.httpClient, url, &status, &errResponse)
	return status, errResponse, err
}

// DeleteNotification deletes a created notification
func (p *Pushy) DeleteNotification(pushID string) (*SimpleSuccess, *Error, error) {
	url := fmt.Sprintf(p.APIEndpoint+"/pushes/%s?api_key=%s", pushID, p.APIToken)
	var success *SimpleSuccess
	var pushyErr *Error
	err := del(p.httpClient, url, &success, &pushyErr)
	return success, pushyErr, err
}

// SubscribeToTopic subscribes a particular device to topics (when you want to do from backend)
func (p *Pushy) SubscribeToTopic(deviceID string, topics ...string) (*SimpleSuccess, *Error, error) {
	url := fmt.Sprintf(p.APIEndpoint+"/devices/subscribe?api_key=%s", p.APIToken)
	request := DeviceSubscriptionRequest{
		Token:  deviceID,
		Topics: topics,
	}
	var success *SimpleSuccess
	var pushyErr *Error
	err := post(p.httpClient, url, request, &success, &pushyErr)
	return success, pushyErr, err
}

// UnsubscribeFromTopic un subscribes a particular device from topics (when you want to do from backend)
func (p *Pushy) UnsubscribeFromTopic(token string, topics ...string) (*SimpleSuccess, *Error, error) {
	url := fmt.Sprintf(p.APIEndpoint+"/devices/unsubscribe?api_key=%s", p.APIToken)
	request := DeviceSubscriptionRequest{
		Token:  token,
		Topics: topics,
	}
	var success *SimpleSuccess
	var pushyErr *Error
	err := post(p.httpClient, url, request, &success, &pushyErr)
	return success, pushyErr, err
}

// NotifyDevice sends notification data to devices
func (p *Pushy) NotifyDevice(request SendNotificationRequest) (*NotificationResponse, *Error, error) {
	url := fmt.Sprintf(p.APIEndpoint+"/push?api_key=%s", p.APIToken)
	var success *NotificationResponse
	var pushyErr *Error
	err := post(p.httpClient, url, request, &success, &pushyErr)
	return success, pushyErr, err
}

func get(client IHTTPClient, url string, positiveResponse interface{}, errResponse interface{}) error {
	response, err := client.Get(url)
	if err != nil {
		positiveResponse = nil
		errResponse = nil
		return err
	}
	defer response.Body.Close()
	body := response.Body
	if response.StatusCode >= 400 {
		positiveResponse = nil
		json.NewDecoder(body).Decode(&errResponse)
		return errors.New(strconv.Itoa(response.StatusCode) + " " + response.Status)
	}
	// decode positiveResponse
	errResponse = nil
	json.NewDecoder(body).Decode(positiveResponse)
	return nil
}

func post(client IHTTPClient, url string, body interface{}, posRes interface{}, errRes interface{}) error {
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(body)
	response, err := client.Post(url, "application/json", buffer)
	if err != nil {
		posRes = nil
		errRes = nil
		return err
	}
	defer response.Body.Close()
	b := response.Body
	if response.StatusCode >= 400 {
		posRes = nil
		json.NewDecoder(b).Decode(&errRes)
		return errors.New(strconv.Itoa(response.StatusCode) + " " + response.Status)
	}
	errRes = nil
	json.NewDecoder(b).Decode(posRes)
	return nil
}

func del(client IHTTPClient, url string, posRes interface{}, errRes interface{}) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err // I don't think there's anything which can result in this error.
	}
	response, err := client.Do(req)
	if err != nil {
		posRes = nil
		errRes = nil
		return err
	}
	defer response.Body.Close()
	b := response.Body
	if response.StatusCode >= 400 {
		posRes = nil
		json.NewDecoder(b).Decode(&errRes)
		return errors.New(strconv.Itoa(response.StatusCode) + " " + response.Status)
	}
	errRes = nil
	json.NewDecoder(b).Decode(posRes)
	return nil
}
