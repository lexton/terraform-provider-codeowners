// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hmarr/codeowners"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CodeownerDataSource{}

func NewCodeownerDataSource() datasource.DataSource {
	return &CodeownerDataSource{}
}

// CodeownerDataSource defines the data source implementation.
type CodeownerDataSource struct {
	ruleset codeowners.Ruleset
}

// CodeownerDataSourceModel describes the data source data model.
type CodeownerDataSourceModel struct {
	Id     types.String `tfsdk:"id"`
	Path   types.String `tfsdk:"path"`
	Owners types.List   `tfsdk:"owners"`
}

func (d *CodeownerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_codeowner"
}

func (d *CodeownerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Codeowners will find the codeowners for the file at a path",
		Attributes: map[string]schema.Attribute{
			// Required for testing data resources
			"id": schema.StringAttribute{
				Computed: true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Path to determine codeownership, this path is only measured relative to the root of the directory.",
				Required:            true,
			},
			"owners": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of Codeowners",
				Computed:            true,
			},
		},
	}
}

func (d *CodeownerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	ruleset, ok := req.ProviderData.(codeowners.Ruleset)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *codeowners.Ruleset, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.ruleset = ruleset
}

func (d *CodeownerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CodeownerDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(data.Path.ValueString())

	rule, err := d.ruleset.Match(data.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to match %s", data.Path.String()), err.Error())
		return
	}

	if rule == nil {
		resp.Diagnostics.AddError(fmt.Sprintf("No CODEOWNERS matched path %s", data.Path.String()), "")
		return
	}

	tflog.Trace(ctx, "Matched Rule", map[string]any{
		"line":    rule.LineNumber,
		"comment": rule.Comment,
		"owners":  rule.Owners,
	})
	ownerList := []string{}
	for _, o := range rule.Owners {
		ownerList = append(ownerList, o.Value)
	}

	owners, diags := types.ListValueFrom(ctx, types.StringType, ownerList)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	data.Owners = owners

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
