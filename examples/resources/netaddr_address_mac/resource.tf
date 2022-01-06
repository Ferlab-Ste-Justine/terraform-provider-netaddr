resource "netaddr_range_mac" "test" {
    key_prefix = "/test/mac/"
    first_address = "52:54:00:00:00:00"
    last_address = "52:54:00:ff:ff:ff"
}

resource "netaddr_address_mac" "test" {
    range_id = netaddr_range_mac.test.id
    name = "test"
    hardcoded_address = "52:54:00:00:00:02"
}

resource "netaddr_address_mac" "test2" {
    range_id = netaddr_range_mac.test.id
    name = "test2"
}

output "test_addr" {
  value = netaddr_address_mac.test.address
}

output "test2_addr" {
  value = netaddr_address_mac.test2.address
}