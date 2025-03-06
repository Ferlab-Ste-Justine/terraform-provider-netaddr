package provider

import (
	"testing"
)

func TestIpv4RangeAddressCount(t *testing.T) {
	addr1, addr1Err := Ipv4StringToBytes("10.128.60.190")
	if addr1Err != nil {
		t.Errorf("Address range count test failed getting address 1: %s", addr1Err.Error())
	}

	addr2, addr2Err := Ipv4StringToBytes("10.128.60.254")
	if addr2Err != nil {
		t.Errorf("Address range count test failed getting address 1: %s", addr2Err.Error())
	}

	range1Count := Ipv4RangeAddressCount(addr1, addr2)
	if range1Count != int64(65) {
		t.Errorf("Expected range count between address 1 and address 2 to be 65 and it was %d", range1Count)
	}

	addr3, addr3Err := Ipv4StringToBytes("10.128.50.254")
	if addr3Err != nil {
		t.Errorf("Address range count test failed getting address 1: %s", addr3Err.Error())
	}
	addr4, addr4Err := Ipv4StringToBytes("10.128.60.190")
	if addr4Err != nil {
		t.Errorf("Address range count test failed getting address 1: %s", addr4Err.Error())
	}

	range2Count := Ipv4RangeAddressCount(addr3, addr4)
	expectedRange2Count := int64(2) + int64(191) + int64(9*256)
	if range2Count != expectedRange2Count {
		t.Errorf("Expected range count between address 3 and address 4 to be %d and it was %d", expectedRange2Count, range2Count)
	}
}