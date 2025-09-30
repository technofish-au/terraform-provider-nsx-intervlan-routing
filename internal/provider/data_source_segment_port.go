// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"terraform-provider-nsx-intervlan-routing/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSourceWithConfigure = &SegmentPortDataSource{}
	_ datasource.DataSource              = &SegmentPortDataSource{}
)

func NewSegmentPortDataSource() datasource.DataSource {
	return &SegmentPortDataSource{}
}

type SegmentPortDataSource struct {
	client client.Client
}

type SegmentPortDataSourceModel struct {
	SegmentId   string      `tfsdk:"segment_id"`
	VmName      string      `tfsdk:"vm_name"`
	SegmentPort SegmentPort `tfsdk:"segment_port"`
}

func (d *SegmentPortDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// IMPORTANT: This method is called MULTIPLE times. An initial call might not have configured the Provider yet, so we need
		// to handle this gracefully. It will eventually be called with a configured provider.
		return
	}

	p, ok := req.ProviderData.(*NsxIntervlanRoutingProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Provider Data",
			fmt.Sprintf("Expected *NsxIntervlanRoutingProviderData with initialized client, got: %T", req.ProviderData),
		)
		return
	}

	//nolint:staticcheck // SA4005 This is in line with the terraform example
	d.client = p.Client
}

// Metadata returns the data source type name.
func (d *SegmentPortDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment_port"
}

// Schema defines the schema for the data source.
func (d *SegmentPortDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a segment port by segment_id and vm_name.",
		Attributes: map[string]schema.Attribute{
			"segment_id": schema.StringAttribute{
				Description: "Identifier for this segment.",
				Required:    true,
			},
			"vm_name": schema.StringAttribute{
				Description: "Name of the VM that this segment is associated with.",
				Required:    true,
			},
			"segment_port": schema.SingleNestedAttribute{
				Description:         "The segment port definition.",
				MarkdownDescription: "The segment port definition",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"address_bindings": schema.ListNestedAttribute{
						Description:         "List of IP address bindings. Only required when creating a CHILD port.",
						MarkdownDescription: "List of IP address bindings. Only required when creating a CHILD port.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"ip_address": schema.StringAttribute{
									Description:         "IP address of segment port",
									MarkdownDescription: "IP address of segment port",
									Computed:            true,
								},
								"mac_address": schema.StringAttribute{
									Description:         "MAC address of segment port",
									MarkdownDescription: "MAC address of segment port",
									Computed:            true,
								},
								"vlan_id": schema.StringAttribute{
									Description:         "VLAN ID associated with this segment port",
									MarkdownDescription: "VLAN ID associated with this segment port",
									Computed:            true,
								},
							},
						},
					},
					"admin_state": schema.StringAttribute{
						Description:         "Admin state of the segment port. Can only be UP or DOWN values.",
						MarkdownDescription: "Admin state of the segment port. Can only be UP or DOWN values.",
						Computed:            true,
					},
					"attachment": schema.SingleNestedAttribute{
						Description:         "Attachment object definition",
						MarkdownDescription: "Attachment object definition",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description:         "VIF UUID in NSX. Required if type is PARENT.",
								MarkdownDescription: "VIF UUID in NSX. Required if type is PARENT.",
								Computed:            true,
							},
							"context_id": schema.StringAttribute{
								Description:         "Attachment UUID of the PARENT port. Only required when type is CHILD.",
								MarkdownDescription: "Attachment UUID of the PARENT port. Only required when type is CHILD.",
								Computed:            true,
							},
							"traffic_tag": schema.StringAttribute{
								Description:         "VLAN ID to tag traffic with. Only required when type is CHILD.",
								MarkdownDescription: "VLAN ID to tag traffic with. Only required when type is CHILD.",
								Computed:            true,
							},
							"app_id": schema.StringAttribute{
								Description:         "Application ID associated with this port. Can be the same as the display name. Only required when type is CHILD.",
								MarkdownDescription: "Application ID associated with this port. Can be the same as the display name. Only required when type is CHILD.",
								Computed:            true,
							},
							"type": schema.StringAttribute{
								Description:         "Type of attachment. Case sensitive. Can be either PARENT or CHILD.",
								MarkdownDescription: "Type of attachment. Case sensitive. Can be either PARENT or CHILD.",
								Computed:            true,
							},
						},
					},
					"description": schema.StringAttribute{
						Description:         "Description of segment port",
						MarkdownDescription: "Description of segment port",
						Computed:            true,
					},
					"display_name": schema.StringAttribute{
						Description:         "Display name of segment port",
						MarkdownDescription: "Display name of segment port",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						Description:         "Id of segment port. Can be the same as display_name.",
						MarkdownDescription: "Id of segment port. Can be the same as display_name.",
						Computed:            true,
					},
					"resource_type": schema.StringAttribute{
						Description:         "Resource type of segment port. MUST be set to 'SegmentPort'",
						MarkdownDescription: "Resource type of segment port. Can only be set to 'SegmentPort'",
						Computed:            true,
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *SegmentPortDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read item data source")
	var state SegmentPortDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	portsResponse, err := d.client.ListSegmentPorts(ctx, state.SegmentId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read segment ports for ",
			err.Error(),
		)
		return
	}

	var segmentPorts client.ListSegmentPortsResponse
	if portsResponse.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP error code received for Item",
			portsResponse.Status,
		)
		return
	}

	if err := json.NewDecoder(portsResponse.Body).Decode(&segmentPorts); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received for segment ports",
			err.Error(),
		)
		return
	}

	// Map response body to model
	lowerVmName := strings.ToLower(state.VmName)
	for _, segment := range segmentPorts.Results {
		lowerDisplayName := strings.ToLower(segment.DisplayName)
		if strings.HasPrefix(lowerDisplayName, lowerVmName) {
			state.SegmentPort = SegmentPort{
				AddressBindings: segment.AddressBindings,
				AdminState:      segment.AdminState,
				Attachment:      segment.Attachment,
				Description:     segment.Description,
				DisplayName:     segment.DisplayName,
				Id:              segment.Id,
			}

			// We found the port. no need to keep going
			break
		}
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, "Finished reading segment ports data source", map[string]any{"success": true})
}
