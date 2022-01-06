data "netaddr_range_mac" "test" {
    key_prefix = "/test/mac/"
}

output "data_range_mac_test" {
  value = "first_address: ${data.netaddr_range_mac.test.first_address}, last_address: ${data.netaddr_range_mac.test.last_address}"
}