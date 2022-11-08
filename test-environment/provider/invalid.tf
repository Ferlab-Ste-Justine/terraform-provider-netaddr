//Invalid address range, key prefix already taken
//Note: Will pass if strict is true
/*resource "netaddr_range_ipv4" "invalid" {
    key_prefix = "/test/basic-ipv4/"
    first_address = "192.169.0.1"
    last_address = "192.169.2.254"
}*/

//Invalid mac address, uses network of wrong type
/*resource "netaddr_address_mac" "invalid_wrong_type" {
    range_id = netaddr_range_ipv4.basic_ipv4.id
    name = "invalid-addr1"
}*/

//Invalid mac address, name already taken
//Note: Will pass if strict is true
/*resource "netaddr_address_mac" "invalid_name_taken" {
    range_id = netaddr_range_mac.basic_mac.id
    name = "addr1"
}*/

//Invalid mac address, hardcoded address already taken by generated address
/*resource "netaddr_address_mac" "invalid_addr_taken" {
    range_id = netaddr_range_mac.basic_mac.id
    name = "invalid-addr2"
    hardcoded_address = "52:54:00:00:00:01"
}*/

//Invalid mac address, hardcoded address already taken by another hardcoded address
/*resource "netaddr_address_mac" "invalid_addr_taken_take2" {
    range_id = netaddr_range_mac.basic_mac.id
    name = "invalid-addr3"
    hardcoded_address = "52:54:00:00:00:02"
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