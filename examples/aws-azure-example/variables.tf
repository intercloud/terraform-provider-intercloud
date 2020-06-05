variable "intercloud_api_endpoint"  {
  description = "InterCloud provider API endpoint to interact with"
  type        = string
  default     = "https://api-console-lab.intercloud.io"
}

variable "intercloud_organization_id"  {
  description = "InterCloud provider organization id to manage"
  type        = string
}

variable "aws_region" {
  description = "AWS region where to manage resources"
  type        = string
}

variable "azure_express_route_circuit_name" {
  description = "Azure express circuit name to provision"
  type = string
}

variable "azure_resource_group_name" {
  description = "Azure resource group containing managed resources"
  type = string
}

variable "azure_subscription_id" {
  description = "Azure subscription id"
  type        = string
}
