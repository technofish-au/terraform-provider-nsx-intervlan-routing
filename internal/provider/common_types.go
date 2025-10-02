// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SegmentPort struct {
	AddressBindings []PortAddressBinding `tfsdk:"address_bindings"`
	AdminState      types.String         `tfsdk:"admin_state"`
	Attachment      PortAttachment       `tfsdk:"attachment"`
	Description     types.String         `tfsdk:"description"`
	DisplayName     types.String         `tfsdk:"display_name"`
	Id              types.String         `tfsdk:"id"`
	ParentPath      types.String         `tfsdk:"parent_path"`
	Path            types.String         `tfsdk:"path"`
	RelativePath    types.String         `tfsdk:"relative_path"`
	ResourceType    types.String         `tfsdk:"resource_type"`
}

type PortAddressBinding struct {
	IpAddress  types.String `tfsdk:"ip_address"`
	MacAddress types.String `tfsdk:"mac_address"`
	VlanId     types.String `tfsdk:"vlan_id"`
}

type PortAttachment struct {
	AllocateAddresses types.String `tfsdk:"allocate_addresses"`
	AppId             types.String `tfsdk:"app_id"`
	ContextId         types.String `tfsdk:"context_id"`
	Id                types.String `tfsdk:"id"`
	TrafficTag        types.Int32  `tfsdk:"traffic_tag"`
	Type              types.String `tfsdk:"type"`
}
