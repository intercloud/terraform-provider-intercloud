
variable "tag_name" {
  description = "The tag name to include in resources to easily identify them"
  type        = string
}

variable "ssh_public_key" {
  description = "The public key used to ssh into the virtual machine"
  type        = string
}
