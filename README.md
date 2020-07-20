Terraform Provider for InterCloud
==================

- Website: <https://www.terraform.io>
- Documentation: See [run documentation website section](#documentation-website)
<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by the InterCloud team at [InterCloud](https://intercloud.com)

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12+

Documentation website
----------------------

The documentation is currently not availbale on the public Terraform documentation, this will be done in the future.

The documentation is available by locally running the documentation website.

```bash
cd scripts/
./run-website.sh
```

After starting the server, the documentation is available at <http://localhost:4567/docs/providers/intercloud>

Using the provider
----------------------

## 3rd Party plugin installation

For using the `terraform-provider-intercloud` plugin, the plugin versions need to be present into the terraform third-party plugins directory.

| OS                | Location                        |
|-------------------|:--------------------------------|
| Windows           | %APPDATA%\terraform.d\plugins   |
| All other systems | ~/.terraform.d/plugins          |

More information about this can be found [here](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

The `terraform-provider-intercloud` plugin supports multiple platforms (`linux/amd64`, `darwin/amd64`, `windows/amd64` and so on).

If you already downloaded the plugin binaries (one binary per platform), just copy the binaries it into the plugins directory.

If you want to build and install the plugin, follow the next steps.

After installing the `terraform-provider-intercloud` plugin into the terraform third-party plugins directory will look like this :

```
.terraform.d
└── plugins
    ├── darwin_386
    │   ├── terraform-provider-intercloud_v0.0.1-beta
    ├── darwin_amd64
    │   ├── terraform-provider-intercloud_v0.0.1-beta
    .
    ├── linux_386
    │   ├── terraform-provider-intercloud_v0.0.1-beta
    ├── linux_amd64
    │   ├── terraform-provider-intercloud_v0.0.1-beta
    .
    ├── windows_386
    │   ├── terraform-provider-intercloud_v0.0.1-beta.exe
    └── windows_amd64
        ├── terraform-provider-intercloud_v0.0.1-beta.exe
```

Building the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.14+ is *required*).
You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-intercloud
...
```

Installing the Provider
---------------------------

If you wish to build and automatically install the `terraform-provider-intercloud` into third party plugins from a source code version.

```sh
# replace <VERSION> with the plugin version, e.g. "0.0.1"
$ PROVIDER_VERSION=<VERSION> make install
 ```
