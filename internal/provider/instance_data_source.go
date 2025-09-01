package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sh05/terraform-provider-multipass/internal/common"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &InstanceDataSource{}

func NewInstanceDataSource() datasource.DataSource {
	return &InstanceDataSource{}
}

// InstanceDataSource defines the data source implementation.
type InstanceDataSource struct {
	client *MultipassClient
}

// InstanceDataSourceModel describes the data source data model.
type InstanceDataSourceModel struct {
	Id        types.String        `tfsdk:"id"`
	Name      types.String        `tfsdk:"name"`
	Instance  *InstanceDataModel  `tfsdk:"instance"`
	Instances []InstanceDataModel `tfsdk:"instances"`
}

// InstanceDataModel represents an instance in the data source
type InstanceDataModel struct {
	Name      types.String `tfsdk:"name"`
	State     types.String `tfsdk:"state"`
	IPv4      types.List   `tfsdk:"ipv4"`
	Release   types.String `tfsdk:"release"`
	ImageHash types.String `tfsdk:"image_hash"`
}

func (d *InstanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

func (d *InstanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Multipass instance data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Data source identifier",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the instance to retrieve. If not specified, all instances will be returned.",
				Optional:            true,
			},
			"instance": schema.SingleNestedAttribute{
				MarkdownDescription: "Instance details (when querying by name)",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "Instance name",
						Computed:            true,
					},
					"state": schema.StringAttribute{
						MarkdownDescription: "Instance state",
						Computed:            true,
					},
					"ipv4": schema.ListAttribute{
						MarkdownDescription: "IPv4 addresses",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"release": schema.StringAttribute{
						MarkdownDescription: "Ubuntu release",
						Computed:            true,
					},
					"image_hash": schema.StringAttribute{
						MarkdownDescription: "Image hash",
						Computed:            true,
					},
				},
			},
			"instances": schema.ListNestedAttribute{
				MarkdownDescription: "List of all instances (when not querying by name)",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Instance name",
							Computed:            true,
						},
						"state": schema.StringAttribute{
							MarkdownDescription: "Instance state",
							Computed:            true,
						},
						"ipv4": schema.ListAttribute{
							MarkdownDescription: "IPv4 addresses",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"release": schema.StringAttribute{
							MarkdownDescription: "Ubuntu release",
							Computed:            true,
						},
						"image_hash": schema.StringAttribute{
							MarkdownDescription: "Image hash",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *InstanceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*MultipassClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *MultipassClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *InstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InstanceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		// Query specific instance by name
		instanceName := data.Name.ValueString()
		tflog.Trace(ctx, "reading multipass instance", map[string]interface{}{"name": instanceName})

		instance, err := d.client.GetInstance(instanceName)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read instance, got error: %s", err))
			return
		}

		// Convert to data model
		instanceData := d.convertToDataModel(instance)
		data.Instance = &instanceData
		data.Id = types.StringValue(instanceName)
	} else {
		// List all instances
		tflog.Trace(ctx, "listing all multipass instances")

		instances, err := d.client.ListInstances()
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list instances, got error: %s", err))
			return
		}

		// Convert to data models
		instanceDataList := make([]InstanceDataModel, len(instances))
		for i, instance := range instances {
			instanceDataList[i] = d.convertToDataModel(&instance)
		}

		data.Instances = instanceDataList
		data.Id = types.StringValue("all-instances")
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// convertToDataModel converts a common.MultipassInstance to InstanceDataModel
func (d *InstanceDataSource) convertToDataModel(instance *common.MultipassInstance) InstanceDataModel {
	model := InstanceDataModel{
		Name:      types.StringValue(instance.Name),
		State:     types.StringValue(instance.State),
		Release:   types.StringValue(instance.Release),
		ImageHash: types.StringValue(instance.ImageHash),
	}

	// Convert IPv4 addresses to list
	if len(instance.IPv4) > 0 {
		ipv4Values := make([]attr.Value, len(instance.IPv4))
		for i, ip := range instance.IPv4 {
			ipv4Values[i] = types.StringValue(ip)
		}
		model.IPv4 = types.ListValueMust(types.StringType, ipv4Values)
	} else {
		model.IPv4 = types.ListValueMust(types.StringType, []attr.Value{})
	}

	return model
}
