package address

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

func IncAddressBy1(toInc []byte) []byte {
	cp := make([]byte, len(toInc))
	copy(cp, toInc)
	idx := len(toInc) - 1
	for idx >= 0 {
		val := uint8(cp[idx])
		if val == 255 {
			val = 0
			cp[idx] = byte(val)
		} else {
			val += 1
			cp[idx] = byte(val)
			return cp
		}

		idx -= 1
	}

	return nil
}

func Ipv4StringToBytes(ipv4 string) ([]byte, error) {
	byteRepr := net.ParseIP(ipv4)
	if byteRepr == nil || byteRepr.To4() == nil {
		return []byte{}, errors.New(fmt.Sprintf("%s is not a valid ipv4 address", ipv4))
	}

	return []byte(byteRepr), nil
}

func Ipv4BytesTo4(ipv4 []byte) []byte {
	return net.IP(ipv4).To4()
}

func Ipv4BytesToString(ipv4 []byte) string {
	return net.IP(ipv4).String()
}

func MacStringToBytes(mac string) ([]byte, error) {
	byteRepr, err := net.ParseMAC(mac)
	if err != nil {
		return []byte{}, err
	}
	return []byte(byteRepr), nil
}

func MacBytesToString(mac []byte) string {
	return net.HardwareAddr(mac).String()
}

func AddressWithinBoundaries(addr []byte, lower []byte, higher []byte) bool {
	lowerRangeCmp := bytes.Compare(addr, lower)
	upperRangeCmp := bytes.Compare(addr, higher)
	return (lowerRangeCmp == 0 || lowerRangeCmp == 1) && (upperRangeCmp == 0 || upperRangeCmp == -1)
}

func AddressLessThan(addr []byte, otherAddress []byte) bool {
	cmp := bytes.Compare(addr, otherAddress)
	return cmp == -1
}

func AddressGreaterThan(addr []byte, otherAddress []byte) bool {
	cmp := bytes.Compare(addr, otherAddress)
	return cmp == 1
}

func Ipv4RangeAddressCount(firstAddr []byte, lastAddr []byte) int64 {
	return int64(binary.BigEndian.Uint32(Ipv4BytesTo4(lastAddr))) - int64(binary.BigEndian.Uint32(Ipv4BytesTo4(firstAddr))) + int64(1)
}

func Ipv6StringToBytes(ipv6 string) ([]byte, error) {
	byteRepr := net.ParseIP(ipv6)
	if byteRepr == nil || byteRepr.To16() == nil {
		return []byte{}, errors.New(fmt.Sprintf("%s is not a valid ipv6 address", ipv6))
	}

	return []byte(byteRepr), nil
}

func Ipv6BytesToString(ipv6 []byte) string {
	return net.IP(ipv6).String()
}