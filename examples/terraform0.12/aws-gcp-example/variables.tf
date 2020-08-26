variable "intercloud_api_endpoint" {
  description = "InterCloud provider API endpoint to interact with"
  type        = string
  default     = "https://api-console-lab.intercloud.io"
}

variable "intercloud_organization_id" {
  description = "InterCloud provider organization id to manage"
  type        = string
}

variable "aws_region" {
  description = "AWS region where to manage resources"
  type        = string
}

variable "google_project" {
  description = "Google project where to manage resources"
  type        = string
}

variable "google_region" {
  description = "Google region where to manage resources"
  type        = string
}

variable "google_zone" {
  description = "Google zone where to manage resources"
  type        = string
}

variable "google_interconnect_attachment_edge_availability_domain" {
  description = "Interconnect attachement edge availability domain"
  type        = string
}

variable "google_interconnect_attachment_router" {
  description = "Interconnect attachment router"
  type        = string
}
