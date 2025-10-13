// Copyright (c) Technofish Consulting Pty Ltd
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertSegmentPortToTF(segment ApiSegmentPort) SegmentPort {
	var segmentPort SegmentPort
	var addressBindings []PortAddressBinding

	for _, address := range segment.AddressBindings {
		addressBindings = append(addressBindings, PortAddressBinding{
			IpAddress:  types.StringValue(address.IpAddress),
			MacAddress: types.StringValue(address.MacAddress),
			VlanId:     types.StringValue(address.VlanId),
		})
	}
	segmentPort.AddressBindings = addressBindings
	segmentPort.AdminState = types.StringValue(segment.AdminState)

	addresses := types.StringValue("")
	if segment.Attachment.AllocateAddresses != "" {
		addresses = types.StringValue(segment.Attachment.AllocateAddresses)
	}
	segmentPort.Attachment = PortAttachment{
		AllocateAddresses: addresses,
		AppId:             types.StringValue(segment.Attachment.AppId),
		ContextId:         types.StringValue(segment.Attachment.ContextId),
		Id:                types.StringValue(segment.Attachment.Id),
		TrafficTag:        types.Int32Value(segment.Attachment.TrafficTag),
		Type:              types.StringValue(segment.Attachment.Type),
	}

	if segment.Description != "" {
		segmentPort.Description = types.StringValue(segment.Description)
	}

	if segment.DisplayName != "" {
		segmentPort.DisplayName = types.StringValue(segment.DisplayName)
	}

	// Not an optional field
	segmentPort.Id = types.StringValue(segment.Id)

	if segment.ParentPath != "" {
		segmentPort.ParentPath = types.StringValue(segment.ParentPath)
	}

	if segment.Path != "" {
		segmentPort.Path = types.StringValue(segment.Path)
	}

	if segment.RelativePath != "" {
		segmentPort.RelativePath = types.StringValue(segment.RelativePath)
	}

	// Also not an optional field
	segmentPort.ResourceType = types.StringValue(segment.ResourceType)

	return segmentPort
}

func ConvertTFToSegmentPort(segment SegmentPort) ApiSegmentPort {
	var segmentPort ApiSegmentPort
	var addressBindings []ApiPortAddressBinding

	for _, address := range segment.AddressBindings {
		addressBindings = append(addressBindings, ApiPortAddressBinding{
			IpAddress:  address.IpAddress.ValueString(),
			MacAddress: address.MacAddress.ValueString(),
			VlanId:     address.VlanId.ValueString(),
		})
	}
	segmentPort.AddressBindings = addressBindings
	segmentPort.AdminState = segment.AdminState.ValueString()

	addresses := ""
	if !segment.Attachment.AllocateAddresses.IsUnknown() {
		addresses = segment.Attachment.AllocateAddresses.ValueString()
	}
	segmentPort.Attachment = ApiPortAttachment{
		AllocateAddresses: addresses,
		AppId:             segment.Attachment.AppId.ValueString(),
		ContextId:         segment.Attachment.ContextId.ValueString(),
		Id:                segment.Attachment.Id.ValueString(),
		TrafficTag:        segment.Attachment.TrafficTag.ValueInt32(),
		Type:              segment.Attachment.Type.ValueString(),
	}

	segmentPort.Description = segment.Description.ValueString()
	segmentPort.DisplayName = segment.DisplayName.ValueString()
	segmentPort.Id = segment.Id.ValueString()
	segmentPort.ParentPath = segment.ParentPath.ValueString()
	segmentPort.Path = segment.Path.ValueString()
	segmentPort.RelativePath = segment.RelativePath.ValueString()
	segmentPort.ResourceType = segment.ResourceType.ValueString()

	return segmentPort
}
