# NSX Guest InterVLAN Routing Terraform Provider

[![Latest Release](https://img.shields.io/github/v/tag/technofish-au/terraform-provider-nsx-intervlan-routing?label=latest%20release&style=for-the-badge)](https://github.com/technofish-au/terraform-provider-nsx-intervlan-routing/releases/latest) [![License](https://img.shields.io/github/license/technofish-au/terraform-provider-nsx-intervlan-routing.svg?style=for-the-badge)](LICENSE)

The Terraform/OpenTofu provider for [NSX Guest InterVLAN Routing](https://techdocs.broadcom.com/us/en/vmware-cis/nsx/vmware-nsx/4-2/administration-guide.html) is a plugin for Terraform/OpenTofu that allows you to interact with VMware NSX to create Parent & Child ports for Guest InterVLAN Routing.

Learn More:

- Read the provider [documentation](https://search.opentofu.org/provider/technofish-au/nsx-intervlan-routing/latest)

## Requirements

- [VMware NSX]

  The following table lists the supported product version for this provider.

  | NSX Version | Supported |
  |-------------|-----------|
  | < 4.0       | No        |
  | 4.0.x       | Yes       |
  | 4.1.x       | Yes       |
  | 4.2.x       | Yes       |

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23

## Using the provider

Simply specify the version of the provider required within your terraform versions block and configure the provider block with the secrets required to access your NSX environment.

```hcl
provider "nsx-intervlan-routing" {
    nsx_url = "https://nsx-manager.example.com"
    nsx_username = "admin"
    nsx_password = "password"
}
```

Now start defining parent and child ports for Guest InterVLAN Routing.

## Contributing

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `make install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Support

The provider is currently supported by [Technofish](https://www.technofish.com.au/). For bugs and feature requests, please open a [Github Issue](https://github.com/technofish-au/terraform-provider-nsx-intervlan-routing/issues) and label it appropriately.

Pull Requests are welcome, but must be reviewed and signed off by an approved member of the [Technofish](https://www.technofish.com.au/) team.
