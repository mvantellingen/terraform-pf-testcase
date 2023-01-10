package hashicups

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dummyRemoteValue struct {
	ID               string
	ComputedValue    int
	MyBoolean        bool
	MyDefaultBoolean bool
}

var remoteValues = map[string]dummyRemoteValue{}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &orderResource{}
	_ resource.ResourceWithConfigure   = &orderResource{}
	_ resource.ResourceWithImportState = &orderResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewOrderResource() resource.Resource {
	return &orderResource{}
}

// orderResource is the resource implementation.
type orderResource struct {
}

// orderResourceModel maps the resource schema data.
type orderResourceModel struct {
	ID               types.String `tfsdk:"id"`
	ComputedValue    types.Int64  `tfsdk:"computed_value"`
	MyBoolean        types.Bool   `tfsdk:"my_boolean"`
	MyDefaultBoolean types.Bool   `tfsdk:"my_default_boolean"`
}

// Metadata returns the data source type name.
func (r *orderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_order"
}

// Schema defines the schema for the data source.
func (r *orderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an order.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier",
				Computed:    true,
			},
			"computed_value": schema.Int64Attribute{
				Description: "Remote value.",
				Computed:    true,
			},
			"my_boolean": schema.BoolAttribute{
				Description: "Remote value.",
				Required:    true,
			},
			"my_default_boolean": schema.BoolAttribute{
				Description: "Remote value.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					BoolDefault(false),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *orderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *orderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan orderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := uuid.GenerateUUID()
	rv := dummyRemoteValue{
		ID:               id,
		ComputedValue:    1,
		MyBoolean:        plan.MyBoolean.ValueBool(),
		MyDefaultBoolean: plan.MyDefaultBoolean.ValueBool(),
	}
	remoteValues[id] = rv
	state := r.transform(id)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *orderResource) transform(id string) orderResourceModel {
	rv := remoteValues[id]
	return orderResourceModel{
		ID:               types.StringValue(rv.ID),
		ComputedValue:    types.Int64Value(int64(rv.ComputedValue)),
		MyBoolean:        types.BoolValue(rv.MyBoolean),
		MyDefaultBoolean: types.BoolValue(rv.MyDefaultBoolean),
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *orderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state orderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state = r.transform(state.ID.ValueString())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *orderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan orderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	rv := remoteValues[id]
	rv.ComputedValue = rv.ComputedValue + 1
	rv.MyBoolean = plan.MyBoolean.ValueBool()
	rv.MyDefaultBoolean = plan.MyDefaultBoolean.ValueBool()
	remoteValues[id] = rv

	state := r.transform(id)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *orderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state orderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	delete(remoteValues, state.ID.ValueString())
}

func (r *orderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
