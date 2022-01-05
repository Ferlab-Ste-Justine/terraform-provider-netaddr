resource "netaddr_network" "test" {
    key_prefix = "/test-network/"
    ipv4 = {
        first = "192.168.0.1"
        last = "192.168.2.254"
    }
    ipv6 = {
        first = "fe80::54d2:baff:fe57:e92"
        last = "fe80::54d2:baff:fe57:e92"
    }
    mac = {
        first = "52:54:00:00:00:00"
        last = "52:54:00:ff:ff:ff"
    }
}

resource "netaddr_ipv4" "test" {
    network_id = netaddr_network.test.id
    name = "test"
}

resource "netaddr_ipv4" "test2" {
    network_id = netaddr_network.test.id
    name = "test2"
}

output "test" {
  value = netaddr_ipv4.test.address
}

output "test2" {
  value = netaddr_ipv4.test2.address
}