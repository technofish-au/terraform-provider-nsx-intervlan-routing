// Copyright (c) Technofish Consulting Pty Ltd.
// SPDX-License-Identifier: MPL-2.0

package helpers

type ListSegmentPortsRequest struct {
	SegmentId string `json:"segment_id"`
}

type ListSegmentPortsResponse struct {
	Results       []ApiSegmentPort `json:"results"`
	ResultCount   int              `json:"result_count"`
	SortBy        string           `json:"sort_by"`
	SortAscending bool             `json:"sort_ascending"`
}

type PatchSegmentPortRequest struct {
	SegmentId      string         `json:"segment_id"`
	PortId         string         `json:"port_id"`
	ApiSegmentPort ApiSegmentPort `json:"segment_port"`
}

type ApiSegmentPort struct {
	AddressBindings []ApiPortAddressBinding `json:"address_bindings,omitempty"`
	AdminState      string                  `json:"admin_state,omitempty"`
	Attachment      ApiPortAttachment       `json:"attachment,omitempty"`
	Description     string                  `json:"description,omitempty"`
	DisplayName     string                  `json:"display_name,omitempty"`
	Id              string                  `json:"id,omitempty"`
	ParentPath      string                  `json:"parent_path,omitempty"`
	Path            string                  `json:"path,omitempty"`
	RelativePath    string                  `json:"relative_path,omitempty"`
	ResourceType    string                  `json:"resource_type,omitempty"`
}

type ApiPortAddressBinding struct {
	IpAddress  string `json:"ip_address,omitempty"`
	MacAddress string `json:"mac_address,omitempty"`
	VlanId     string `json:"vlan_id,omitempty"`
}

type ApiPortAttachment struct {
	AllocateAddresses string `json:"allocate_addresses,omitempty"`
	AppId             string `json:"app_id,omitempty"`
	ContextId         string `json:"context_id,omitempty"`
	Id                string `json:"id,omitempty"`
	TrafficTag        int32  `json:"traffic_tag,omitempty"`
	Type              string `json:"type,omitempty"`
}
