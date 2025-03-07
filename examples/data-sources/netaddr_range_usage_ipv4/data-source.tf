data "netaddr_range_ipv4" "test" {
    key_prefix = "/test/ipv4/"
}

data "netaddr_range_usage_ipv4" "test" {
  range_id = data.netaddr_range_ipv4.test.id
}

output "range_capacity" {
  description = "The range supports the following number of addresses."
  value       = data.netaddr_range_usage_ipv4.test.capacity
}

output "range_used_capacity" {
  description = "The range has allocated the following number of addresses."
  value       = data.netaddr_range_usage_ipv4.test.used_capacity
}

output "range_free_capacity" {
  description = "The range can allocate the following number of addresses before running out of ips."
  value       = data.netaddr_range_usage_ipv4.test.free_capacity
}