output "aws_destinations" {
  value = data.intercloud_destinations.dest_aws
}

output "gcp_destinations" {
  value = data.intercloud_destinations.dest_gcp
}