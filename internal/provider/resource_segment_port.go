// Copyright (c) Technofish Consulting Pty Ltd.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"terraform-provider-nsx-intervlan-routing/client"
	"terraform-provider-nsx-intervlan-routing/helpers"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.ResourceWithConfigure   = &SegmentPortResource{}
	_ resource.Resource                = &SegmentPortResource{}
	_ resource.ResourceWithImportState = &SegmentPortResource{}
)

func NewSegmentPortResource() resource.Resource {
	return &SegmentPortResource{}
}

type SegmentPortResource struct {
	client client.Client
}

type SegmentPortResourceModel struct {
	SegmentId   types.String         `tfsdk:"segment_id"`
	PortId      types.String         `tfsdk:"port_id"`
	SegmentPort *helpers.SegmentPort `tfsdk:"segment_port"`
}

func (r *SegmentPortResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = p.Client
}

// Metadata returns the resource type name.
func (r *SegmentPortResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment_port"
}

func (r *SegmentPortResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a segment port.",
		Attributes: map[string]schema.Attribute{
			"segment_id": schema.StringAttribute{
				Description:         "Identifier for this segment.",
				MarkdownDescription: "Identifier for this segment.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"port_id": schema.StringAttribute{
				Description:         "Identifier for this port.",
				MarkdownDescription: "Identifier for this port.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"segment_port": schema.SingleNestedAttribute{
				Description:         "The segment port definition.",
				MarkdownDescription: "The segment port definition",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"address_bindings": schema.ListNestedAttribute{
						Description:         "List of IP address bindings. Only required when creating a CHILD port.",
						MarkdownDescription: "List of IP address bindings. Only required when creating a CHILD port.",
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"ip_address": schema.StringAttribute{
									Description:         "IP address of segment port",
									MarkdownDescription: "IP address of segment port",
									Required:            true,
								},
								"mac_address": schema.StringAttribute{
									Description:         "MAC address of segment port",
									MarkdownDescription: "MAC address of segment port",
									Required:            true,
								},
								"vlan_id": schema.StringAttribute{
									Description:         "VLAN ID associated with this segment port",
									MarkdownDescription: "VLAN ID associated with this segment port",
									Required:            true,
								},
							},
						},
					},
					"admin_state": schema.StringAttribute{
						Description:         "Admin state of the segment port. Can only be UP or DOWN values.",
						MarkdownDescription: "Admin state of the segment port. Can only be UP or DOWN values.",
						Required:            true,
					},
					"attachment": schema.SingleNestedAttribute{
						Description:         "Attachment object definition",
						MarkdownDescription: "Attachment object definition",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Description:         "VIF UUID in NSX. Required if type is PARENT.",
								MarkdownDescription: "VIF UUID in NSX. Required if type is PARENT.",
								Optional:            true,
							},
							"context_id": schema.StringAttribute{
								Description:         "Attachment UUID of the PARENT port. Only required when type is CHILD.",
								MarkdownDescription: "Attachment UUID of the PARENT port. Only required when type is CHILD.",
								Optional:            true,
							},
							"traffic_tag": schema.Int32Attribute{
								Description:         "VLAN ID to tag traffic with. Only required when type is CHILD.",
								MarkdownDescription: "VLAN ID to tag traffic with. Only required when type is CHILD.",
								Optional:            true,
							},
							"allocate_addresses": schema.StringAttribute{
								Description:         "Indicate how IP will be allocated for the port. Enum: IP_POOL, MAC_POOL, BOTH, DHCP, DHCPV6, SLAAC, NONE",
								MarkdownDescription: "Indicate how IP will be allocated for the port. Enum: IP_POOL, MAC_POOL, BOTH, DHCP, DHCPV6, SLAAC, NONE",
								Optional:            true,
							},
							"app_id": schema.StringAttribute{
								Description:         "Application ID associated with this port. Can be the same as the display name. Only required when type is CHILD.",
								MarkdownDescription: "Application ID associated with this port. Can be the same as the display name. Only required when type is CHILD.",
								Optional:            true,
							},
							"type": schema.StringAttribute{
								Description:         "Type of attachment. Case sensitive. Can be either PARENT or CHILD.",
								MarkdownDescription: "Type of attachment. Case sensitive. Can be either PARENT or CHILD.",
								Required:            true,
							},
						},
					},
					"description": schema.StringAttribute{
						Description:         "Description of segment port",
						MarkdownDescription: "Description of segment port",
						Optional:            true,
					},
					"display_name": schema.StringAttribute{
						Description:         "Display name of segment port",
						MarkdownDescription: "Display name of segment port",
						Required:            true,
					},
					"id": schema.StringAttribute{
						Description:         "Id of segment port. Can be the same as display_name.",
						MarkdownDescription: "Id of segment port. Can be the same as display_name.",
						Required:            true,
					},
					"parent_path": schema.StringAttribute{
						Description:         "Parent path of segment port",
						MarkdownDescription: "Parent path of segment port",
						Computed:            true,
					},
					"path": schema.StringAttribute{
						Description:         "Path of segment port",
						MarkdownDescription: "Path of segment port",
						Computed:            true,
					},
					"relative_path": schema.StringAttribute{
						Description:         "Relative path of segment port",
						MarkdownDescription: "Relative path of segment port",
						Computed:            true,
					},
					"resource_type": schema.StringAttribute{
						Description:         "Resource type of segment port. MUST be set to 'SegmentPort'",
						MarkdownDescription: "Resource type of segment port. Can only be set to 'SegmentPort'",
						Required:            true,
					},
				},
			},
		},
	}
}

// Create a new resource.
func (r *SegmentPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Preparing to create segment port resource")
	// Retrieve values from plan
	var plan SegmentPortResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	segmentId := plan.SegmentId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Segment ID: %s", segmentId))
	portId := plan.PortId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Port ID: %s", portId))

	segmentPort := helpers.ConvertTFToSegmentPort(*plan.SegmentPort)
	patchRequest := helpers.PatchSegmentPortRequest{
		SegmentId:      segmentId,
		PortId:         portId,
		ApiSegmentPort: segmentPort,
	}

	// Create new item
	spResponse, err := r.client.PatchSegmentPort(ctx, patchRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Segment Port",
			err.Error(),
		)
		return
	}

	if spResponse.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"An invalid response was received. Code: "+strconv.Itoa(spResponse.StatusCode),
			spResponse.Status,
		)
		return
	}

	var newSegmentPort helpers.ApiSegmentPort
	if err := json.NewDecoder(spResponse.Body).Decode(&newSegmentPort); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received for Item",
			err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Created segment port resource", map[string]any{"segment_port": newSegmentPort})

	// This should contain the computed values as well.
	tfSegmentPort := helpers.ConvertSegmentPortToTF(newSegmentPort)
	plan.SegmentPort = &tfSegmentPort
	tflog.Debug(ctx, "COMPUTED SEGMENT PORT", map[string]any{"segment_port": tfSegmentPort})

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Created segment port resource", map[string]any{"success": true})
}

// Read resource information.
func (r *SegmentPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read segment port resource")
	// Get current state
	var state SegmentPortResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	spResponse, err := r.client.GetSegmentPort(ctx, state.SegmentId.ValueString(), state.PortId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Segment Port configuration",
			err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Read segment port resource", map[string]any{"segment_port": spResponse})

	// Treat HTTP 404 Not Found status as a signal to remove/recreate resource
	if spResponse.StatusCode == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	if spResponse.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Unexpected HTTP error code received for segment port",
			spResponse.Status,
		)
		return
	}

	var newSegmentPort helpers.ApiSegmentPort
	if err := json.NewDecoder(spResponse.Body).Decode(&newSegmentPort); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received for segment port",
			err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Read segment port resource", map[string]any{"segment_port": newSegmentPort})

	// Map response body to model
	convertedSegment := helpers.ConvertSegmentPortToTF(newSegmentPort)
	state = SegmentPortResourceModel{
		SegmentId:   state.SegmentId,
		PortId:      state.PortId,
		SegmentPort: &convertedSegment,
	}
	tflog.Debug(ctx, "Conversion complete", map[string]any{"segment_port": convertedSegment})

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Finished reading segment port resource", map[string]any{"success": true})
}

func (r *SegmentPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Preparing to update segment port resource")
	// Retrieve values from plan
	var plan SegmentPortResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Segment ID to update: %s", plan.SegmentId.ValueString()))
	segmentId := plan.SegmentId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Port ID to update: %s", plan.PortId.ValueString()))
	portId := plan.PortId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Segment Port details: %+v", &plan.SegmentPort))
	segmentPort := helpers.ConvertTFToSegmentPort(*plan.SegmentPort)

	patchRequest := helpers.PatchSegmentPortRequest{
		SegmentId:      segmentId,
		PortId:         portId,
		ApiSegmentPort: segmentPort,
	}
	tflog.Debug(ctx, fmt.Sprintf("Updating segment port with request %+v", req))

	// Create new item
	spResponse, err := r.client.PatchSegmentPort(ctx, patchRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Segment Port",
			err.Error(),
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("PatchSegmentPort response: %+v", spResponse))

	if spResponse.StatusCode != 200 {
		resp.Diagnostics.AddError(
			fmt.Sprintf("An invalid response was received. Code: %d", spResponse.StatusCode),
			spResponse.Status,
		)
		return
	}

	var updatedSegmentPort helpers.ApiSegmentPort
	if err := json.NewDecoder(spResponse.Body).Decode(&updatedSegmentPort); err != nil {
		resp.Diagnostics.AddError(
			"Invalid format received for segment ports",
			err.Error(),
		)
		return
	}
	convertedSegment := helpers.ConvertSegmentPortToTF(updatedSegmentPort)
	tflog.Debug(ctx, fmt.Sprintf("Converted segment port to TF: %+v", convertedSegment))
	plan.SegmentPort = &convertedSegment

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, fmt.Sprintf("Error encountered setting state: %s", resp.Diagnostics.Errors()))
		return
	}
	tflog.Debug(ctx, "Updated segment port resource", map[string]any{"success": true})
}

func (r *SegmentPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Preparing to delete segment port resource")
	// Retrieve values from state
	var state SegmentPortResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// delete item
	_, err := r.client.DeleteSegmentPort(ctx, state.SegmentId.ValueString(), state.PortId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Item",
			err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Deleted segment port resource", map[string]any{"success": true})
}

func (r *SegmentPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	// If our ID was a string then we could do this
	resource.ImportStatePassthroughID(ctx, path.Root("port_id"), req, resp)

	//id, err := strconv.ParseInt(req.ID, 10, 64)
	//
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Error importing item",
	//		"Could not import item, unexpected error (ID should be an integer): "+err.Error(),
	//	)
	//	return
	//}

	//resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
