//Mac Validation
resource "netaddr_range_mac" "basic_mac" {
    key_prefix = "/test/basic-mac/"
    first_address = "52:54:00:00:00:00"
    last_address = "52:54:00:ff:ff:ff"
}

resource "netaddr_address_mac" "basic_mac_addr1" {
    range_id = netaddr_range_mac.basic_mac.id
    name = "addr1"
    hardcoded_address = "52:54:00:00:00:02"
}

resource "netaddr_address_mac" "basic_mac_addr2" {
    range_id = netaddr_range_mac.basic_mac.id
    name = "addr2"
    depends_on = [netaddr_address_mac.basic_mac_addr1]
}

resource "netaddr_address_mac" "basic_mac_addr3" {
    range_id = netaddr_range_mac.basic_mac.id
    name = "addr3"
    depends_on = [netaddr_address_mac.basic_mac_addr1]
}

resource "netaddr_address_mac" "basic_mac_addr4" {
    range_id = netaddr_range_mac.basic_mac.id
    name = "addr4"
    depends_on = [netaddr_address_mac.basic_mac_addr1]
}

resource "netaddr_address_mac" "basic_mac_addr5" {
    range_id = netaddr_range_mac.basic_mac.id
    name = "addr5"
    depends_on = [netaddr_address_mac.basic_mac_addr1]
}

output "basic_mac_addr1" {
  value = netaddr_address_mac.basic_mac_addr1.address
}

output "basic_mac_addr2" {
  value = netaddr_address_mac.basic_mac_addr2.address
}

output "basic_mac_addr3" {
  value = netaddr_address_mac.basic_mac_addr3.address
}

output "basic_mac_addr4" {
  value = netaddr_address_mac.basic_mac_addr4.address
}

output "basic_mac_addr5" {
  value = netaddr_address_mac.basic_mac_addr5.address
}