output "current_user" {
  value = data.aws_caller_identity.current
}

output "amazon_side_asn" {
  value = aws_vpn_gateway.vpn_gw.amazon_side_asn
}
