data "netaddr_range_ipv4" "range1" {
    key_prefix = "/test/ipv4/"
}

data "netaddr_range_ipv4" "range2" {
    key_prefix = "/test/ipv4-extras/"
}

data "netaddr_address_ipv4_v2" "test" {
    range_ids = [data.netaddr_range_ipv4.range1.id, data.netaddr_range_ipv4.range2.id]
    name = "test"
}

output "data_ipv4_test" {
  value = data.netaddr_address_ipv4.test.address
}