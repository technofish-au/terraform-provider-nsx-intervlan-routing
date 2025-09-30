// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"strings"

	"terraform-provider-nsx-intervlan-routing/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &SegmentPortDataSource{}
	_ datasource.DataSourceWithConfigure = &SegmentPortDataSource{}
)

func NewSegmentPortDataSource() datasource.DataSource {
	return &SegmentPortDataSource{}
}

type SegmentPortDataSource struct {
	client *client.Client
}

type SegmentPortDataSourceModel struct {
	SegmentId   string `tfsdk:"segment_id"`
	VmName      string `tfsdk:"vm_name"`
	SegmentPort SegmentPort
}

func (d SegmentPortDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		tflog.Error(ctx, "Unable to prepare client")
		return
	}

	//nolint:staticcheck // SA4005 This is in line with the terraform example
	d.client = client
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
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *SegmentPortDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read item data source")
	var state SegmentPortDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

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
	state = SegmentPortDataSourceModel{}
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
