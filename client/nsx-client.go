// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

func setupLogging(debug bool) {
	var ll logrus.Level
	if debug {
		ll = logrus.DebugLevel
	} else {
		ll = logrus.InfoLevel
	}
	logrus.SetLevel(ll)
}

func NewClient(server string, username string, password string, insecure bool, debug bool, opts ...ClientOption) (*Client, error) {
	setupLogging(debug)
	logrus.Debug("Creating new NSX API Client")
	// Ensure we have a scheme set for the endpoint.
	s, e := url.Parse(server)
	if e != nil {
		logrus.Errorf("Error parsing server URL: %s, exiting", e)
		panic(e)
	}
	var svr string
	if s.Scheme == "" {
		logrus.Debug("Using default https scheme for server")
		svr = "https://" + server
	}

	// create a client with sane default values
	client := Client{
		Server: svr,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}

	// create httpClient, if not already present
	tr := &http.Transport{}
	if insecure {
		logrus.Debug("Insecure mode enabled. Skipping remote certificate verification")
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	if client.Client == nil {
		logrus.Debug("Client is not instantiated. Creating client.Client with the Transport configuration specified")
		client.Client = &http.Client{Transport: tr}
	}

	logrus.Debug("Client created. Calling GetDefaultHeaders function")
	err := GetDefaultHeaders(&client, username, password)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func GetDefaultHeaders(c *Client, username string, password string) error {
	logrus.Debug("Starting the GetDefaultHeaders function call")
	XsrfToken := "X-XSRF-TOKEN"

	path := c.Server + "/api/session/create"
	logrus.Debugf("Session Create URI is %s", path)

	data := url.Values{}
	data.Set("j_username", username)
	data.Set("j_password", password)

	body := bytes.NewBufferString(data.Encode())

	// Call session create
	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		logrus.Debugf("Failed to create session %s", err)
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := c.Client.Do(req)
	if err != nil || response == nil {
		logrus.Debugf("Failed to create session %s", err)
		return err
	}
	logrus.Debugf("Response header is %s", response.Header)

	if response.StatusCode != 200 {
		logrus.Debugf("Request responded with a non-200 status code. Code is: %d", response.StatusCode)
		return fmt.Errorf("status code %d", response.StatusCode)
	}

	// Go over the headers
	for k, v := range response.Header {
		if strings.EqualFold("Set-Cookie", k) {
			r, _ := regexp.Compile("JSESSIONID=.*?;")
			result := r.FindString(v[0])
			if result != "" {
				c.Session = result
			}
		}
		if strings.EqualFold(XsrfToken, k) {
			c.XsrfToken = v[0]
		}
	}

	err = response.Body.Close()
	if err != nil {
		return err
	}

	logrus.Debug("Successfully completed the GetDefaultHeaders function call")
	return nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
//func WithHTTPClient(doer HttpRequestDoer) ClientOption {
//	return func(c *Client) error {
//		c.Client = doer
//		return nil
//	}
//}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
//func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
//	return func(c *Client) error {
//		c.RequestEditors = append(c.RequestEditors, fn)
//		return nil
//	}
//}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	logrus.Debugf("RequestEditors are: %v", c.RequestEditors)
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}

	req.Header.Add("User-Agent", "Swagger-Codegen/1.0.0/go")

	return nil
}

func (c *Client) DeleteSegmentPort(ctx context.Context, segmentId string, portId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteSegmentPortRequest(c.Server, segmentId, portId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}

	req.Header.Add("X-XSRF-TOKEN", c.XsrfToken)
	req.Header.Add("Cookie", c.Session)

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.Client.Do(req)
}

func NewDeleteSegmentPortRequest(server string, segmentId string, portId string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/policy/api/v1/infra/segments/%s/ports/%s", segmentId, portId)
	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) ListSegmentPorts(ctx context.Context, segmentId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	logrus.Debug(fmt.Sprintf("ListSegmentPorts called with segment ID: %s", segmentId))
	req, err := NewListSegmentPortsRequest(&c.Server, segmentId)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}

	req.Header.Add("X-XSRF-TOKEN", c.XsrfToken)
	req.Header.Add("Set-Cookie", c.Session)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		logrus.Errorf("Failed to list segment ports %s", err)
		return nil, err
	}

	logrus.Debugf("ListSegmentPorts response: %v", resp)
	return resp, nil
}

func NewListSegmentPortsRequest(server *string, segmentId string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(*server)
	if err != nil {
		logrus.Errorf("Failed to parse the server %s", err)
		return nil, err
	}

	operationPath := "/policy/api/v1/infra/segments/" + segmentId + "/ports"
	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		logrus.Errorf("Failed to parse the full URL %s", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		logrus.Errorf("Failed to create the new http request %s", err)
		return nil, err
	}

	logrus.Debugf("Created the request as %v", req)
	return req, nil
}

func (c *Client) GetSegmentPort(ctx context.Context, segmentId string, portId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetSegmentPortRequest(c.Server, segmentId, portId)

	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}

	req.Header.Add("X-XSRF-TOKEN", c.XsrfToken)
	req.Header.Add("Cookie", c.Session)

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.Client.Do(req)
}

func NewGetSegmentPortRequest(server string, segmentId string, portId string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/policy/api/v1/infra/segments/" + segmentId + "/ports/" + portId
	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) PatchSegmentPort(ctx context.Context, body PatchSegmentPortRequest, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPatchSegmentPortRequest(c.Server, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}

	req.Header.Add("X-XSRF-TOKEN", c.XsrfToken)
	req.Header.Add("Cookie", c.Session)

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.Client.Do(req)
}

func NewPatchSegmentPortRequest(server string, body PatchSegmentPortRequest) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/policy/api/v1/infra/segments/" + body.SegmentId + "/ports/" + body.PortId
	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	buf, err := json.Marshal(body.SegmentPort)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)

	req, err := http.NewRequest(http.MethodPatch, queryURL.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	return req, nil
}
