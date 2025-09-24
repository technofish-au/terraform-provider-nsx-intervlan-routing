// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

//
//import (
//	"testing"
//
//	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
//	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
//	"github.com/hashicorp/terraform-plugin-testing/statecheck"
//	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
//)
//
//func TestAccSegmentPortParentResource(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { testAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			// Create and Read testing
//			{
//				Config: testAccSegmentPortResourceParentConfig,
//				ConfigStateChecks: []statecheck.StateCheck{
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_id"),
//						knownvalue.StringExact("4d4c0f0a-6c5 0-420b-90f1-68fb7585cda4"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("port_id"),
//						knownvalue.StringExact("a274ac51-88f5-491f-a46f-840d409ce82f"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.admin_state"),
//						knownvalue.StringExact("UP"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.attachment.id"),
//						knownvalue.StringExact("9765bf41-9725-4714-977e-7f7395920de2"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.attachment.traffic_tag"),
//						knownvalue.StringExact("1000"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.attachment.type"),
//						knownvalue.StringExact("PARENT"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.description"),
//						knownvalue.StringExact("GCVE-PA-VM-ESX-2 Parent Port"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.display_name"),
//						knownvalue.StringExact("GCVE-PA-VM-ESX-2.vmx@060af2c2-e9ff-4686-866c-c0daab1748d6"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.id"),
//						knownvalue.StringExact("060af2c2-e9ff-4686-866c-c0daab1748d6"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.resource_type"),
//						knownvalue.StringExact("SegmentPort"),
//					),
//				},
//			},
//		},
//	})
//}
//
//func TestAccSegmentPortChildResource(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:                 func() { testAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithEcho,
//		Steps: []resource.TestStep{
//			// Create and Read testing
//			{
//				Config: testAccSegmentPortResourceChildConfig,
//				ConfigStateChecks: []statecheck.StateCheck{
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_id"),
//						knownvalue.StringExact("2bfe8abf-4161-4788-9cbe-c444e9bf7454"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("port_id"),
//						knownvalue.StringExact("a274ac51-88f5-491f-a46f-840d409ce82f"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.address_bindings.ip_address"),
//						knownvalue.StringExact("169.254.254.169"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.address_bindings.mac_address"),
//						knownvalue.StringExact("00:50:56:ad:5e:64"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.address_bindings.vlan_id"),
//						knownvalue.StringExact("1001"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.admin_state"),
//						knownvalue.StringExact("UP"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.attachment.context_id"),
//						knownvalue.StringExact("9765bf41-9725-4714-977e-7f7395920de2"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.attachment.traffic_tag"),
//						knownvalue.StringExact("1001"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.attachment.app_id"),
//						knownvalue.StringExact("Segment1001"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.attachment.type"),
//						knownvalue.StringExact("CHILD"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.description"),
//						knownvalue.StringExact("GCVE-PA-VM-ESX-2 Child Port 1001"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.display_name"),
//						knownvalue.StringExact("GCVE-PA-VM-ESX-2.vmx@a274ac51-88f5-491f-a46f-840d409ce82f"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.id"),
//						knownvalue.StringExact("a274ac51-88f5-491f-a46f-840d409ce82f"),
//					),
//					statecheck.ExpectKnownValue(
//						"nsx_intervlan_routing_segment_port.parent_example",
//						tfjsonpath.New("segment_port.resource_type"),
//						knownvalue.StringExact("SegmentPort"),
//					),
//				},
//			},
//		},
//	})
//}
//
//const testAccSegmentPortResourceParentConfig = `
//resource "nsx_intervlan_routing_segment_port" "parent_example" {
//  segment_id = "4d4c0f0a-6c5 0-420b-90f1-68fb7585cda4"
//  port_id    = "a274ac51-88f5-491f-a46f-840d409ce82f"
//  segment_port = {
//    admin_state = "UP"
//    attachment = {
//      id          = "9765bf41-9725-4714-977e-7f7395920de2"
//      traffic_tag = "1000"
//      type        = "PARENT"
//    }
//    description   = "GCVE-PA-VM-ESX-2 Parent Port"
//    display_name  = "GCVE-PA-VM-ESX-2.vmx@060af2c2-e9ff-4686-866c-c0daab1748d6"
//    id            = "060af2c2-e9ff-4686-866c-c0daab1748d6"
//    resource_type = "SegmentPort"
//  }
//}
//`
//
//const testAccSegmentPortResourceChildConfig = `
//resource "nsx_intervlan_routing_segment_port" "child_example" {
//  segment_id = "2bfe8abf-4161-4788-9cbe-c444e9bf7454"
//  port_id    = "a274ac51-88f5-491f-a46f-840d409ce82f"
//  segment_port = {
//    address_bindings = [
//      {
//        ip_address  = "169.254.254.169"
//        mac_address = "00:50:56:ad:5e:64"
//        vlan_id     = "1001"
//      },
//    ]
//    admin_state = "UP"
//    attachment = {
//      context_id  = "9765bf41-9725-4714-977e-7f7395920de2"
//      traffic_tag = "1001"
//      app_id      = "Segment1001"
//      type        = "CHILD"
//    }
//    description   = "GCVE-PA-VM-ESX-2 Child Port 1001"
//    display_name  = "GCVE-PA-VM-ESX-2.vmx@a274ac51-88f5-491f-a46f-840d409ce82f"
//    id            = "a274ac51-88f5-491f-a46f-840d409ce82f"
//    resource_type = "SegmentPort"
//  }
//}
//`
