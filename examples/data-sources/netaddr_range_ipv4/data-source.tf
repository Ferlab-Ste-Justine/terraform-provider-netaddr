data "netaddr_range_ipv4" "test" {
    key_prefix = "/test/ipv4/"
}

output "data_range_ipv4_test" {
  value = "first_address: ${data.netaddr_range_ipv4.test.first_address}, last_address: ${data.netaddr_range_ipv4.test.last_address}"
}