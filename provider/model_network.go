package provider

type AddressRange struct {
	First []byte
	Last []byte
}

type Network struct {
	Ipv4 AddressRange
	Ipv6 AddressRange
	Mac AddressRange
}