data "netaddr_range_ipv4" "test" {
    key_prefix = "/test/ipv4/"
}

data "netaddr_address_ipv4" "test" {
    range_id = data.netaddr_range_ipv4.test.id
    name = "test"
}

output "data_ipv4_test" {
  value = data.netaddr_address_ipv4.test.address
}