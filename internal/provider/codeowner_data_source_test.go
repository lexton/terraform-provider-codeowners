// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCodeownersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCodeownerDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.codeowners_codeowner.test", "path", "hello/world"),
				),
			},
		},
	})
}

const testAccCodeownerDataSourceConfig = `

provider "codeowners" {
	codeowner_path = "TEST_CODEOWNERS"
}

data "codeowners_codeowner" "test" {
  path = "hello/world"
}
`
