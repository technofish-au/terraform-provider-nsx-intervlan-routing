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
	AddressBindings []ApiPortAddressBinding `json:"address_bindings"`
	AdminState      string                  `json:"admin_state"`
	Attachment      ApiPortAttachment       `json:"attachment"`
	Description     string                  `json:"description"`
	DisplayName     string                  `json:"display_name"`
	Id              string                  `json:"id"`
	ParentPath      string                  `json:"parent_path"`
	Path            string                  `json:"path"`
	RelativePath    string                  `json:"relative_path"`
	ResourceType    string                  `json:"resource_type"`
}

type ApiPortAddressBinding struct {
	IpAddress  string `json:"ip_address"`
	MacAddress string `json:"mac_address"`
	VlanId     string `json:"vlan_id"`
}

type ApiPortAttachment struct {
	AllocateAddresses string `json:"allocate_addresses"`
	AppId             string `json:"app_id"`
	ContextId         string `json:"context_id"`
	Id                string `json:"id"`
	TrafficTag        int32  `json:"traffic_tag"`
	Type              string `json:"type"`
}
