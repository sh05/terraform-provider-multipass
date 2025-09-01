package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sh05/terraform-provider-multipass/internal/common"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &InstanceResource{}
var _ resource.ResourceWithImportState = &InstanceResource{}

func NewInstanceResource() resource.Resource {
	return &InstanceResource{}
}

// InstanceResource defines the resource implementation.
type InstanceResource struct {
	client *MultipassClient
}

// InstanceResourceModel describes the resource data model.
type InstanceResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Image     types.String `tfsdk:"image"`
	CPU       types.String `tfsdk:"cpu"`
	Memory    types.String `tfsdk:"memory"`
	Disk      types.String `tfsdk:"disk"`
	CloudInit types.String `tfsdk:"cloud_init"`
	State     types.String `tfsdk:"state"`
	IPv4      types.List   `tfsdk:"ipv4"`
}

func (r *InstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

func (r *InstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Multipass instance resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Instance identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Instance name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image": schema.StringAttribute{
				MarkdownDescription: "Ubuntu image to use (e.g. 'ubuntu', '22.04', '20.04')",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cpu": schema.StringAttribute{
				MarkdownDescription: "Number of CPUs to allocate",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"memory": schema.StringAttribute{
				MarkdownDescription: "Amount of memory to allocate (e.g. '1G', '512M')",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disk": schema.StringAttribute{
				MarkdownDescription: "Disk space to allocate (e.g. '5G', '10G')",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloud_init": schema.StringAttribute{
				MarkdownDescription: "Path to cloud-init configuration file",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "Current state of the instance",
				Computed:            true,
			},
			"ipv4": schema.ListAttribute{
				MarkdownDescription: "IPv4 addresses assigned to the instance",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *InstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*MultipassClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *MultipassClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *InstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InstanceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create launch options
	opts := &common.LaunchOptions{
		Name:      data.Name.ValueString(),
		Image:     data.Image.ValueString(),
		CPU:       data.CPU.ValueString(),
		Memory:    data.Memory.ValueString(),
		Disk:      data.Disk.ValueString(),
		CloudInit: data.CloudInit.ValueString(),
	}

	// Launch the instance
	tflog.Trace(ctx, "launching multipass instance", map[string]interface{}{"name": opts.Name})

	err := r.client.Launch(opts)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create instance, got error: %s", err))
		return
	}

	// Set the ID
	data.Id = data.Name

	// Read the instance to get current state
	instance, err := r.client.GetInstance(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read instance after creation, got error: %s", err))
		return
	}

	// Update the model with instance data
	r.updateModelFromInstance(&data, instance)

	tflog.Trace(ctx, "created multipass instance")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InstanceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the instance
	instance, err := r.client.GetInstance(data.Name.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			// Instance doesn't exist, remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read instance, got error: %s", err))
		return
	}

	// Update the model with instance data
	r.updateModelFromInstance(&data, instance)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InstanceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// For now, most updates require replacement since Multipass doesn't support
	// changing instance configuration after creation. The schema marks most
	// attributes as requiring replacement.

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InstanceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the instance
	tflog.Trace(ctx, "deleting multipass instance", map[string]interface{}{"name": data.Name.ValueString()})

	err := r.client.DeleteInstance(data.Name.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete instance, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "deleted multipass instance")
}

func (r *InstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// updateModelFromInstance updates the resource model with data from a Multipass instance
func (r *InstanceResource) updateModelFromInstance(data *InstanceResourceModel, instance *common.MultipassInstance) {
	data.State = types.StringValue(instance.State)

	// Convert IPv4 addresses to list
	if len(instance.IPv4) > 0 {
		ipv4Values := make([]attr.Value, len(instance.IPv4))
		for i, ip := range instance.IPv4 {
			ipv4Values[i] = types.StringValue(ip)
		}
		data.IPv4 = types.ListValueMust(types.StringType, ipv4Values)
	}
}
