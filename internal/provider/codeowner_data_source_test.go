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
					resource.TestCheckResourceAttr("data.codeowners_codeowner.world", "path", "hello/world"),
					resource.TestCheckResourceAttr("data.codeowners_codeowner.world", "owners.0", "lexton"),
					resource.TestCheckResourceAttr("data.codeowners_codeowner.world", "owners.#", "1"),
					resource.TestCheckResourceAttr("data.codeowners_codeowner.there1", "owners.0", "lexton2"),
					resource.TestCheckResourceAttr("data.codeowners_codeowner.there1", "owners.#", "1"),
					resource.TestCheckResourceAttr("data.codeowners_codeowner.others", "owners.0", "other"),
					resource.TestCheckResourceAttr("data.codeowners_codeowner.others", "owners.#", "1"),
					resource.TestCheckResourceAttr("data.codeowners_codeowner.multiple", "owners.0", "lexton"),
					resource.TestCheckResourceAttr("data.codeowners_codeowner.multiple", "owners.1", "lexton@example.com"),
					resource.TestCheckResourceAttr("data.codeowners_codeowner.multiple", "owners.#", "2"),
				),
			},
		},
	})
}

const testAccCodeownerDataSourceConfig = `

provider "codeowners" {
	codeowner_path = "TEST_CODEOWNERS"
}

data "codeowners_codeowner" "world" {
  path = "hello/world"
}

data "codeowners_codeowner" "there1" {
  path = "hello/there/world"
}

data "codeowners_codeowner" "others" {
  path = "foo/bar"
}

data "codeowners_codeowner" "multiple" {
  path = "prefix/hello/foo"
}
`
