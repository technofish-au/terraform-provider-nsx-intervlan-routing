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
		var pab PortAddressBinding
		if address.IpAddress != "" {
			pab.IpAddress = types.StringValue(address.IpAddress)
		}
		if address.MacAddress != "" {
			pab.MacAddress = types.StringValue(address.MacAddress)
		}
		if address.VlanId != "" {
			pab.VlanId = types.StringValue(address.VlanId)
		}
		addressBindings = append(addressBindings, pab)
	}
	segmentPort.AddressBindings = addressBindings

	segmentPort.AdminState = types.StringValue(segment.AdminState)

	var attachment PortAttachment
	if segment.Attachment.AllocateAddresses != "" {
		attachment.AllocateAddresses = types.StringValue(segment.Attachment.AllocateAddresses)
	}
	if segment.Attachment.AppId != "" {
		attachment.AppId = types.StringValue(segment.Attachment.AppId)
	}
	if segment.Attachment.ContextId != "" {
		attachment.ContextId = types.StringValue(segment.Attachment.ContextId)
	}
	if segment.Attachment.Id != "" {
		attachment.Id = types.StringValue(segment.Attachment.Id)
	}
	if segment.Attachment.TrafficTag >= 0 {
		attachment.TrafficTag = types.Int32Value(segment.Attachment.TrafficTag)
	}
	if segment.Attachment.Type != "" {
		attachment.Type = types.StringValue(segment.Attachment.Type)
	}
	segmentPort.Attachment = attachment

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
