package provider

import (
	"terraform-provider-nsx-intervlan-routing/client"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertSegmentPortToTfSdk(port client.SegmentPort) (SegmentPort, error) {
	var portAddressBindings []PortAddressBinding
	for _, address := range port.AddressBindings {
		portAddressBindings = append(portAddressBindings, PortAddressBinding{
			IpAddress:  types.StringValue(address.IpAddress),
			MacAddress: types.StringValue(address.MacAddress),
			VlanId:     types.StringValue(address.VlanId),
		})
	}
	attachment := PortAttachment{
		AllocateAddresses: types.StringValue(port.Attachment.AllocateAddresses),
		AppId:             types.StringValue(port.Attachment.AppId),
		ContextId:         types.StringValue(port.Attachment.ContextId),
		Id:                types.StringValue(port.Attachment.Id),
		TrafficTag:        types.Int32Value(port.Attachment.TrafficTag),
		Type:              types.StringValue(port.Attachment.Type),
	}
	segmentPort := SegmentPort{
		AddressBindings: portAddressBindings,
		AdminState:      types.StringValue(port.AdminState),
		Attachment:      attachment,
		Description:     types.StringValue(port.Description),
		DisplayName:     types.StringValue(port.DisplayName),
		Id:              types.StringValue(port.Id),
		ParentPath:      types.StringValue(port.ParentPath),
		Path:            types.StringValue(port.Path),
		RelativePath:    types.StringValue(port.RelativePath),
		ResourceType:    types.StringValue(port.ResourceType),
	}

	return segmentPort, nil
}

func ConvertSegmentPortToClient(port SegmentPort) (client.SegmentPort, error) {
	var portAddressBindings []client.PortAddressBindingEntry
	for _, address := range port.AddressBindings {
		portAddressBindings = append(portAddressBindings, client.PortAddressBindingEntry{
			IpAddress:  address.IpAddress.ValueString(),
			MacAddress: address.MacAddress.ValueString(),
			VlanId:     address.VlanId.ValueString(),
		})
	}
	attachment := client.PortAttachment{
		AllocateAddresses: port.Attachment.AllocateAddresses.ValueString(),
		AppId:             port.Attachment.AppId.ValueString(),
		ContextId:         port.Attachment.ContextId.ValueString(),
		Id:                port.Attachment.Id.ValueString(),
		TrafficTag:        port.Attachment.TrafficTag.ValueInt32(),
		Type:              port.Attachment.Type.ValueString(),
	}

	segmentPort := client.SegmentPort{
		AddressBindings: portAddressBindings,
		AdminState:      port.AdminState.ValueString(),
		Attachment:      attachment,
		Description:     port.Description.ValueString(),
		DisplayName:     port.DisplayName.ValueString(),
		Id:              port.Id.ValueString(),
		ParentPath:      port.ParentPath.ValueString(),
		Path:            port.Path.ValueString(),
		RelativePath:    port.RelativePath.ValueString(),
		ResourceType:    port.ResourceType.ValueString(),
	}

	return segmentPort, nil
}
