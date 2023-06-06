// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hmarr/codeowners"
)

// Ensure CodeownerProvider satisfies various provider interfaces.
var _ provider.Provider = &CodeownerProvider{}

// CodeownerProvider defines the provider implementation.
type CodeownerProvider struct {
	version string
}

// CodeownerProviderModel describes the provider data model.
type CodeownerProviderModel struct {
	CodeownerPath types.String `tfsdk:"codeowner_path"`
}

func (p *CodeownerProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "codeowners"
	resp.Version = p.version
}

func (p *CodeownerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"codeowner_path": schema.StringAttribute{
				MarkdownDescription: "Path to the CODEOWNERS file",
				Required:            true,
			},
		},
	}
}

func (p *CodeownerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CodeownerProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		resp.Diagnostics.AddError("Failed to get current working directory", err.Error())
		return
	}

	path := filepath.Join(wd, data.CodeownerPath.ValueString())

	tflog.Info(ctx, "CODEOWNERS path", map[string]any{
		"path": path,
	})

	file, err := os.Open(path)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to open %s", path), err.Error())
		return
	}
	defer file.Close()

	ruleset, err := codeowners.ParseFile(file)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to parse %s", data.CodeownerPath.String()), err.Error())
		return
	}

	resp.DataSourceData = ruleset
	resp.ResourceData = ruleset
}

func (p *CodeownerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *CodeownerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCodeownerDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CodeownerProvider{
			version: version,
		}
	}
}
