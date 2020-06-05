variable "tag_name" {
  description = "The tag name to include in resources"
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


variable "ssh_public_key" {
  description = "The public key used to ssh into the virtual machine"
  type        = string
}
