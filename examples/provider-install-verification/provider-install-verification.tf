terraform {
  required_providers {
    nsx-intervlan-routing = {
      source = "technofish-au/nsx-intervlan-routing"
    }
  }
}

provider "nsx-intervlan-routing" {
  host           = "127.0.0.1"
  insecure       = true
  username       = "admin"
  password       = "password"
}

data "nsx-intervlan-routing_segment_ports" "example" {
  segment_id    = "4d4c0f0a-6c5 0-420b-90f1-68fb7585cda4"
}
