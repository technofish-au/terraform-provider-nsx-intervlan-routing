// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function.
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	Server         string
	XsrfToken      string
	Session        string
	Client         HttpRequestDoer
	RequestEditors []RequestEditorFn
}

type ClientOption func(*Client) error

type ClientInterface interface {
	DeleteSegmentPort(string) (*http.Response, error)
	ListSegmentPorts(string) (*ListSegmentPortsResponse, error)
	GetSegmentPort(string, string) (*SegmentPort, error)
	PatchSegmentPort(string, string) (*bool, error)
}

type ListSegmentPortsRequest struct {
	SegmentId string `json:"segment_id"`
}

type ListSegmentPortsResponse struct {
	Results       []SegmentPort `json:"results"`
	ResultCount   int           `json:"result_count"`
	SortBy        string        `json:"sort_by"`
	SortAscending bool          `json:"sort_ascending"`
}

type PatchSegmentPortRequest struct {
	SegmentId   string      `json:"segment_id"`
	PortId      string      `json:"port_id"`
	SegmentPort SegmentPort `json:"segment_port"`
}

type PortAddressBindingEntry struct {
	IpAddress  string `json:"ip_address"`
	MacAddress string `json:"mac_address"`
	VlanId     string `json:"vlan_id"`
}

type PortAttachment struct {
	AllocateAddresses string `json:"allocate_addresses"`
	AppId             string `json:"app_id"`
	ContextId         string `json:"context_id"`
	Id                string `json:"id"`
	TrafficTag        string `json:"traffic_tag"`
	Type              string `json:"type"`
}

type SegmentPort struct {
	AddressBindings []PortAddressBindingEntry `json:"address_bindings"`
	AdminState      string                    `json:"admin_state"`
	Attachment      PortAttachment            `json:"attachment"`
	Description     string                    `json:"description"`
	DisplayName     string                    `json:"display_name"`
	Id              string                    `json:"id"`
	ResourceType    string                    `json:"resource_type"`
}
