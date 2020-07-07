provider "intercloud" {
  api_endpoint    = var.intercloud_api_endpoint
  organization_id = var.intercloud_organization_id
}

provider "aws" {
  version = "=2.58"
  alias   = "ic"
  region  = var.aws_region
}

provider "azurerm" {
  alias           = "ic"
  version         = "=2.6"
  subscription_id = var.azure_subscription_id
  features {}
}