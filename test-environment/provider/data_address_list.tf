data "netaddr_address_list_ipv4" "basic_ipv4" {
  range_id = netaddr_range_ipv4.basic_ipv4.id
}

data "netaddr_address_list_mac" "basic_mac" {
  range_id = netaddr_range_mac.basic_mac.id
}

output "data_address_list_ipv4" {
  value = data.netaddr_address_list_ipv4.basic_ipv4
}

output "data_address_list_mac" {
  value = data.netaddr_address_list_mac.basic_mac
}