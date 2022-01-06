data "netaddr_range_mac" "test" {
    key_prefix = "/test/mac/"
}

data "netaddr_address_mac" "test" {
    range_id = data.netaddr_range_mac.test.id
    name = "test"
}

output "data_mac_test" {
  value = data.netaddr_address_mac.test.address
}