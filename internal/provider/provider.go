// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"terraform-provider-nsx-intervlan-routing/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure NsxtIntervlanRoutingProvider satisfies various provider interfaces.
var _ provider.Provider = &NsxtIntervlanRoutingProvider{}

// NsxtIntervlanRoutingProvider defines the provider implementation.
type NsxtIntervlanRoutingProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type NsxtIntervlanRoutingProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

func (p *NsxtIntervlanRoutingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nsx-intervlan-routing"
	resp.Version = p.version
}

func (p *NsxtIntervlanRoutingProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "Hostname or IP address of the NSX endpoint",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username of the NSX endpoint",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password of the NSX endpoint",
				Optional:            true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Whether or not the NSX endpoint is insecure",
				Optional:            true,
			},
		},
	}
}

func (p *NsxtIntervlanRoutingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data NsxtIntervlanRoutingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.Host.IsNull() {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("host"),
			"Missing NSX Manager API Hostname (using default value: 127.0.0.1)",
			"The provider is using a default value as there is a missing or empty value for the NSX-T Manager API hostname. "+
				"Set the host value in the configuration. "+
				"If either is already set, ensure the value is not empty.",
		)
		data.Host = types.StringValue("127.0.0.1")
	}
	if data.Username.IsNull() {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("username"),
			"Missing NSX API username (using default value: admin)",
			"The provider is using a default value as there is a missing or empty value for the NSX-T API username. "+
				"Set the username value in the configuration. "+
				"If either is already set, ensure the value is not empty.",
		)
		data.Username = types.StringValue("admin")
	}
	if data.Password.IsNull() {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("password"),
			"Missing NSX API port (using default value: password)",
			"The provider is using a default value as there is a missing or empty value for the NSX-T API password. "+
				"Set the password value in the configuration. "+
				"If either is already set, ensure the value is not empty.",
		)
		data.Password = types.StringValue("password")
	}
	if data.Insecure.IsUnknown() {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("insecure"),
			"Missing NSX Manager API Insecure (using default value: false)",
			"The provider is using a default value as there is a missing or empty value for the NSX-T Manager API insecure. "+
				"Set the insecure value in the configuration. "+
				"If either is already set, ensure the value is not empty.",
		)
		data.Insecure = types.BoolValue(false)
	}

	// Example client configuration for data sources and resources
	c, err := client.NewClient(data.Host.String(), data.Username.String(), data.Password.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create an instance of the API Client.",
			"The provider has failed to instantiate the API Client at line 117 of the provider.tf."+
				"Please check the instantiation of the client to ensure the params are correct.")
		panic("Failed to create an instance of the API Client. Error is: " + err.Error())
	}
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *NsxtIntervlanRoutingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSegmentPortResource,
	}
}

func (p *NsxtIntervlanRoutingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSegmentPortDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NsxtIntervlanRoutingProvider{
			version: version,
		}
	}
}
