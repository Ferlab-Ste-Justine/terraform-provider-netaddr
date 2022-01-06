data "netaddr_range_ipv4" "test" {
    key_prefix = netaddr_range_ipv4.test.id
}

data "netaddr_range_mac" "test" {
    key_prefix = netaddr_range_mac.test.id
}

data "netaddr_address_mac" "test3" {
    range_id = data.netaddr_range_mac.test.id
    name = "test3"
}

data "netaddr_address_mac" "test4" {
    range_id = data.netaddr_range_mac.test.id
    name = "test4"
}


output "data_range_ipv4_test" {
  value = "first_address: ${data.netaddr_range_ipv4.test.first_address}, last_address: ${data.netaddr_range_ipv4.test.last_address}"
}

output "data_range_mac_test" {
  value = "first_address: ${data.netaddr_range_mac.test.first_address}, last_address: ${data.netaddr_range_mac.test.last_address}"
}

output "data_address_mac_test3" {
  value = data.netaddr_address_mac.test3.address
}

output "data_address_mac_test4" {
  value = data.netaddr_address_mac.test4.address
}