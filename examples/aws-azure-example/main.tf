# ------------------------------------------------
# terraform configuration
# ------------------------------------------------
terraform {
  required_version = ">= 0.12.0"
  required_providers {
    # Build has to be named terraform-provider-intercloud_v0.0.0
    #intercloud = ">= 0.0.0"
    tls     = "~> 2.1"
    aws     = "~> 2.58"
    random  = "~> 2.2"
    azurerm = "~> 2.6"
  }
}

# ------------------------------------------------
# aws, azure, google modules
# ------------------------------------------------
module "ic_aws" {
  source = "../modules/aws"
  providers = {
    aws = aws.ic
  }
  # module variables
  tag_name       = random_pet.tag.id
  ssh_public_key = tls_private_key.common.public_key_openssh
}

module "ic_azure" {
  source = "../modules/azure"
  providers = {
    azurerm = azurerm.ic
  }
  # module variables
  tag_name                         = random_pet.tag.id
  azure_express_route_circuit_name = var.azure_express_route_circuit_name
  azure_resource_group_name        = var.azure_resource_group_name
  ssh_public_key                   = tls_private_key.common.public_key_openssh
}

# ------------------------------------------------
# resources
# ------------------------------------------------
resource "random_pet" "tag" {}

resource "tls_private_key" "common" {
  algorithm = "RSA"
  rsa_bits  = 2048
}


data "intercloud_destinations" "dest_aws" {
  family = "aws"
}

data "intercloud_destinations" "dest_azure" {
  family = "azure"
}


data "intercloud_destination" "aws_destination" {
  location = "telehouse tsh london"
  family   = "aws"
}

data "intercloud_destination" "azure_destination" {
  location = "equinix ld5 london"
  family   = "azure"
}

resource "intercloud_group" "group_1" {
  name = "group-${random_pet.tag.id}"
}

resource "intercloud_connector" "connector_aws" {
  name           = "conn-aws-${random_pet.tag.id}"
  description    = "descript"
  group_id       = intercloud_group.group_1.id
  destination_id = data.intercloud_destination.aws_destination.id
  aws {
    aws_account_id = module.ic_aws.current_user.account_id
    aws_bgp_asn    = module.ic_aws.amazon_side_asn
  }
}

resource "intercloud_connector" "connector_azure" {
  name           = "conn_azure-${random_pet.tag.id}"
  description    = "descript"
  group_id       = intercloud_group.group_1.id
  destination_id = data.intercloud_destination.azure_destination.id
  azure {
    service_key = module.ic_azure.express_route_service_key
  }
}

resource "intercloud_link" "link_azure_aws" {
  name        = "link-${random_pet.tag.id}"
  description = "descript"
  from        = intercloud_connector.connector_azure.id
  to          = intercloud_connector.connector_aws.id
}