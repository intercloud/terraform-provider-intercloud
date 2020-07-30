# ------------------------------------------------
# terraform configuration
# ------------------------------------------------
terraform {
  required_version = ">= 0.12.0"
  required_providers {
    # Build has to be named terraform-provider-intercloud_v0.0.0
    #intercloud = ">= 0.0.0"
    tls    = "~> 2.1"
    aws    = "~> 2.58"
    random = "~> 2.2"
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

module "ic_gcp" {
  source = "../modules/gcp"
  providers = {
    google = google.ic
  }
  # module variables
  tag_name                                                = random_pet.tag.id
  google_project                                          = var.google_project
  google_region                                           = var.google_region
  google_zone                                             = var.google_zone
  google_interconnect_attachment_edge_availability_domain = var.google_interconnect_attachment_edge_availability_domain
  google_interconnect_attachment_router                   = var.google_interconnect_attachment_router
  ssh_public_key                                          = tls_private_key.common.public_key_openssh
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

data "intercloud_destinations" "dest_gcp" {
  family = "gcp"
}


data "intercloud_destination" "aws_destination" {
  location = "telehouse tsh london"
  family   = "aws"
}

data "intercloud_destination" "gcp_destination" {
  location = "equinix ld5 london"
  family   = "gcp"
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

resource "intercloud_connector" "connector_gcp" {
  name           = "conn-gcp-${random_pet.tag.id}"
  description    = "descript"
  group_id       = intercloud_group.group_1.id
  destination_id = data.intercloud_destination.gcp_destination.id
  gcp {
    pairing_key = module.ic_gcp.google_compute_interconnect_attachment_pairing_key
    med         = "low"
  }
}


resource "intercloud_link" "link_gcp_aws" {
  name        = "link-${random_pet.tag.id}"
  description = "descript"
  from        = intercloud_connector.connector_gcp.id
  to          = intercloud_connector.connector_aws.id
}