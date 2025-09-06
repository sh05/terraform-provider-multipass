package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure MultipassProvider satisfies various provider interfaces.
var _ provider.Provider = &MultipassProvider{}
var _ provider.ProviderWithFunctions = &MultipassProvider{}

// MultipassProvider defines the provider implementation.
type MultipassProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance testing.
	version string
}

// MultipassProviderModel describes the provider data model.
type MultipassProviderModel struct {
	BinaryPath types.String `tfsdk:"binary_path"`
}

func (p *MultipassProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "multipass"
	resp.Version = p.version
}

func (p *MultipassProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Multipass provider enables Terraform to manage Ubuntu virtual machines using Canonical Multipass.",
		Attributes: map[string]schema.Attribute{
			"binary_path": schema.StringAttribute{
				MarkdownDescription: "Path to the multipass binary. Defaults to 'multipass' if not specified.",
				Optional:            true,
			},
		},
	}
}

func (p *MultipassProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data MultipassProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	binaryPath := data.BinaryPath.ValueString()

	// Create the Multipass client
	client := NewMultipassClient(binaryPath)

	// Make the client available during resource operations
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *MultipassProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewInstanceResource,
	}
}

func (p *MultipassProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewInstanceDataSource,
	}
}

func (p *MultipassProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// No functions defined yet
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MultipassProvider{
			version: version,
		}
	}
}
