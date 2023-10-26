package provider

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"google.golang.org/grpc/status"

	authproto "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &SchemaResource{}

//var _ resource.ResourceWithImportState = &SchemaResource{}

func NewSchemaResource() resource.Resource {
	return &SchemaResource{}
}

// SchemaResource defines the resource implementation.
type SchemaResource struct {
	client *authzed.Client
}

// SchemaResourceModel describes the resource data model.
type SchemaResourceModel struct {
	Schema types.String `tfsdk:"schema"`
}

func (r *SchemaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

func (r *SchemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "SpiceDB Schema",

		Attributes: map[string]schema.Attribute{
			"schema": schema.StringAttribute{
				MarkdownDescription: "SpiceDB Schema",
				Required:            true,
			},
		},
	}
}

func (r *SchemaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*authzed.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *authzed.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *SchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SchemaResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.WriteSchema(ctx, &authproto.WriteSchemaRequest{
		Schema: data.Schema.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("SpiceDB Schema Error", fmt.Sprintf("Unable to create schema, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Created SpiceDB Schema")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SchemaResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	readSchema, err := ReadSchema(r.client, ctx)
	if err != nil {
		resp.Diagnostics.AddError("SpiceDB Schema Error", fmt.Sprintf("Unable to read schema, got error: %s", status.Code(err)))
		return
	}

	currentSchema := NormaliseString(SortDefinitions(data.Schema.ValueString()))
	newSchema := NormaliseString(readSchema)

	if newSchema != currentSchema {
		data.Schema = types.StringValue(newSchema)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SchemaResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.WriteSchema(ctx, &authproto.WriteSchemaRequest{
		Schema: data.Schema.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("SpiceDB Schema Error", fmt.Sprintf("Unable to update schema, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SchemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SchemaResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.WriteSchema(ctx, &authproto.WriteSchemaRequest{
		Schema: "",
	})
	if err != nil {
		resp.Diagnostics.AddError("SpiceDB Schema Error", fmt.Sprintf("Unable to delete schema, got error: %s", err))
		return
	}
}

//func (r *SchemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
//	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
//}

/* Naive schema normalisation
 * It appears that SpiceDB re-arranges definitions in alphabetical order
 * and cleans up spaces when returning the schema.
 * Blindly checking schema strings for equality will not work here and some
 * naive smartness is required.
 * So, we also try to clean up some spaces and sort definitions alphabetically
 * before we compare the schema values
 */

func NormaliseString(input string) string {
	spaces := regexp.MustCompile(`\s+`)
	colon := regexp.MustCompile(`\s+:`)
	return colon.ReplaceAllString(spaces.ReplaceAllString(input, " "), ":")
}

func SortDefinitions(schema string) string {
	sections := strings.Split(schema, "definition")

	// Remove the empty strings and whitespace.
	var fragments []string
	for _, section := range sections {
		trimmedSection := strings.TrimSpace(section)
		if len(trimmedSection) > 0 {
			fragments = append(fragments, "definition "+trimmedSection)
		}
	}

	// Sort the fragments in alphabetical order.
	sort.Strings(fragments)

	// Create a new string with the sorted fragments.
	return strings.Join(fragments, "\n")
}
