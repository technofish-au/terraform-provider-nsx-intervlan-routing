// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func NewClient(server string, username string, password string, insecure bool, opts ...ClientOption) (*Client, error) {

	s, e := url.Parse(server)
	if e != nil {
		panic(e)
	}
	var svr string
	if s.Scheme == "" {
		if s.Scheme != "https" && !insecure {
			svr = "https://" + server
		} else {
			svr = "http://" + server
		}
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
	if client.Client == nil {
		client.Client = &http.Client{}
	}

	err := GetDefaultHeaders(&client, username, password)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func GetDefaultHeaders(c *Client, username string, password string) error {
	XSRF_TOKEN := "X-XSRF-TOKEN"
	path := c.Server + "/api/session/create?j_username=" + username + "&j_password=" + password

	// Call session create
	req, err := http.NewRequest(http.MethodPost, path, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return fmt.Errorf("failed to create session: %s", err)
	}

	response, err := c.Client.Do(req)
	if err != nil || response == nil {
		return fmt.Errorf("failed to create session: %s", err)
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("failed to create session: status code %d", response.StatusCode)
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
		if strings.EqualFold(XSRF_TOKEN, k) {
			c.XsrfToken = v[0]
		}
	}

	response.Body.Close()

	return nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
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

	req.Header.Add("X-XSRF-TOKEN", c.XsrfToken)
	req.Header.Add("Cookie", c.Session)

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return nil
}

func (c *Client) DeleteSegmentPort(ctx context.Context, segment_id string, port_id string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteSegmentPortRequest(c.Server, segment_id, port_id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func NewDeleteSegmentPortRequest(server string, segment_id string, port_id string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/policy/api/v1/infra/segments/%s/ports/%s", segment_id, port_id)
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

func (c *Client) ListSegmentPorts(ctx context.Context, segment_id string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListSegmentPortsRequest(c.Server, segment_id)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}

	return c.Client.Do(req)
}

func NewListSegmentPortsRequest(server string, segment_id string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/policy/api/v1/infra/segments/" + segment_id + "/ports"
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

func (c *Client) GetSegmentPort(ctx context.Context, segment_id string, port_id string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetSegmentPortRequest(c.Server, segment_id, port_id)

	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func NewGetSegmentPortRequest(server string, segment_id string, port_id string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/policy/api/v1/infra/segments/" + segment_id + "/ports/" + port_id
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
