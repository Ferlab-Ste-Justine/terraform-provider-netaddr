//Invalid address range, key prefix already taken
/*resource "netaddr_range_ipv4" "invalid" {
    key_prefix = "/test/ipv4/"
    first_address = "192.169.0.1"
    last_address = "192.169.2.254"
}*/

//Invalid mac address, uses network of wrong type
/*resource "netaddr_address_mac" "invalid" {
    range_id = netaddr_range_ipv4.test.id
    name = "test6"
}*/

//Invalid mac address, name already taken
/*resource "netaddr_address_mac" "invalid2" {
    range_id = netaddr_range_mac.test.id
    name = "test5"
}*/

//Invalid mac address, hardcoded address already taken by generated address
/*resource "netaddr_address_mac" "invalid3" {
    range_id = netaddr_range_mac.test.id
    name = "test6"
    hardcoded_address = "52:54:00:00:00:01"
}*/

//Invalid mac address, hardcoded address already taken by another hardcoded address
/*resource "netaddr_address_mac" "invalid3" {
    range_id = netaddr_range_mac.test.id
    name = "test6"
    hardcoded_address = "52:54:00:00:00:01"
}*/

//Invalid, range will run out of addresses
/*resource "netaddr_range_ipv4" "runout" {
    key_prefix = "/test/runout/ipv4/"
    first_address = "192.169.0.1"
    last_address = "192.169.0.2"
}

resource "netaddr_address_ipv4" "runout1" {
    range_id = netaddr_range_ipv4.runout.id
    name = "runout1"
    hardcoded_address = "192.169.0.1"
}

resource "netaddr_address_ipv4" "runout2" {
    range_id = netaddr_range_ipv4.runout.id
    name = "runout2"
}*/

/*resource "netaddr_address_ipv4" "runout3" {
    range_id = netaddr_range_ipv4.runout.id
    name = "runout_this_is_it"
}*/