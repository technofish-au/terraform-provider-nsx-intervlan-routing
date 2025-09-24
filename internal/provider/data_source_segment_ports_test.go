// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccSegmentPortDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExampleDataSourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.nsx-intervlan-routing_segment_ports.example",
						tfjsonpath.New("segment_id"),
						knownvalue.StringExact("4d4c0f0a-6c5 0-420b-90f1-68fb7585cda4"),
					),
				},
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
data "nsx-intervlan-routing_segment_ports" "example" {
  segment_id    = "4d4c0f0a-6c5 0-420b-90f1-68fb7585cda4"
}
`

const testAccExampleDataSourceResult = `
{
  "results": [
    {
      "resource_type": "SegmentPort",
	  "id": "a274ac51-88f5-491f-a46f-840d409ce82f",
	  "display_name": "a274ac51-88f5-491f-a46f-840d409ce82f",
	  "path": "/infra/segments/production-t1-seg/ports/a274ac51-88f5-491f-a46f-840d409ce82f",
	  "relative_path": "a274ac51-88f5-491f-a46f-840d409ce82f",
	  "parent_path": "/infra/segments/production-t1-seg",
	  "marked_for_delete": false,
	  "_create_user": "system",
	  "_create_time": 1544503100539,
	  "_last_modified_user": "system",
	  "_last_modified_time": 1544503100539,
	  "_system_owned": true,
	  "_protection": "NOT_PROTECTED",
	  "_revision": 0
    }
  ],
  "result_count": 1,
  "sort_by": "display_name",
  "sort_ascending": true,
}
`
