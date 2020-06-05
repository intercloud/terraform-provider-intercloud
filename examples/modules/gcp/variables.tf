variable "tag_name" {
  description = "The tag name to include in resources"
  type        = string
}

variable "google_region" {
  description = "Google region where to manage resources"
  type = string
}

variable "google_zone" {
  description = "Google zone where to manage resources"
  type = string
}

variable "google_interconnect_attachment_edge_availability_domain" {
  description = "Interconnect attachement edge availability domain"
  type = string
}

variable "google_interconnect_attachment_router" {
  description = "Interconnect attachment router"
  type = string
}


variable "google_project" {
  description = "Google project containing managed resources"
  type = string
  
}



variable "ssh_public_key" {
  description = "The public key used to ssh into the virtual machine"
  type        = string
}
