output "aws_destinations" {
  value = data.intercloud_destinations.dest_aws
}

output "azure_destinations" {
  value = data.intercloud_destinations.dest_azure
}