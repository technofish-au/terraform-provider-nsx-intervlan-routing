// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "terraform-provider-nsx-intervlan-routing/client"

type SegmentPort struct {
	AddressBindings client.PortAddressBindingEntry `tfsdk:"address_bindings"`
	AdminState      string                         `tfsdk:"admin_state"`
	Attachment      client.PortAttachment          `tfsdk:"attachment"`
	Description     string                         `tfsdk:"description"`
	DisplayName     string                         `tfsdk:"display_name"`
	Id              string                         `tfsdk:"id"`
	ResourceType    string                         `tfsdk:"resource_type"`
}
