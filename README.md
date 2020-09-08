# Terraform Provider for InterCloud

- Website: <https://www.terraform.io>
- Documentation: See [run documentation website section](#documentation-website)
  
<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

## Maintainers

This provider plugin is maintained by the InterCloud team at [InterCloud](https://intercloud.com)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

------------------------------

## Using the provider from releases binaries

Go to [releases download page](https://github.com/intercloud/terraform-provider-intercloud/releases)

Download the right archive from the assets depending on your OS and Arch

:warning: Under MacOS Catalina, errors can appear. Please see this Github issue for more informations and solutions: [macOS Catalina error](https://github.com/hashicorp/terraform/issues/23033)

### Terraform 0.12 installation

Unzip and move the binary in :

| OS                | Location                        |
|-------------------|:--------------------------------|
| Windows           | %APPDATA%\terraform.d\plugins   |
| All other systems | ~/.terraform.d/plugins          |

### Terraform 0.13 installation

Terraform 0.13 uses `source` to make the path to the binary

------------------------------

## Using the provider from source

### Terraform 0.12

```sh
make release-snapshot
## Move/Copy the right binary depending of your {OS}_{Arch} from ./dist to the root of your hcl files
## IE with darwin_amd64
mv ./dist/terraform-provider-intercloud_darwin_amd64/terraform-provider-intercloud_v1.1.0-SNAPSHOT-783c762 ~/my-tf-plan-ie/
cd ~/my-tf-plan-ie
terraform init
terraform plan
terraform apply
```

Windows OS User:

```powershell
# Execute those commands only the first time
go install github.com/goreleaser/goreleaser
go mod tidy

# Then
goreleaser build --rm-dist --snapshot
# mv ./dist/terraform-provider-intercloud_windows_amd64/terraform-provider-intercloud_v1.1.0-SNAPSHOT-783c762.exe E:\my-tf-plan-ie\terraform-provider-intercloud_v1.1.0-SNAPSHOT-783c762.exe
Move-Item -Path ./dist/terraform-provider-intercloud_windows_amd64/terraform-provider-intercloud_v1.1.0-SNAPSHOT-783c762.exe -Destination E:\my-tf-plan-ie\terraform-provider-intercloud_v1.1.0-SNAPSHOT-783c762.exe
E:\my-tf-plan-ie
terraform init
terraform plan
terraform apply

```

The provider can also be moved to one of those directories:

| OS                | Location                        |
|-------------------|:--------------------------------|
| Windows           | %APPDATA%\terraform.d\plugins   |
| All other systems | ~/.terraform.d/plugins          |

### Terraform 0.13

- *Coming soon*

------------------------------

## Building the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.14+ is *required*).
You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-intercloud
...
```

------------------------------

## Examples

Some terraform 0.12 projects are available at [./examples/terraform0.12](?/../examples/terraform0.12)

*Examples with terraform 0.13 are coming soon.*
