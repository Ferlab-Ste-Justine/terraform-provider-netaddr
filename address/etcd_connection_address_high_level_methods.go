package address

import (
	"bytes"
	"errors"
	"fmt"
)


func (conn *EtcdConnection) GenerateGeneratedAddressWithValidation(name string, prefixes []string, rangeType string, toleratePresent bool, addrIsGreater AddressIsGreater, incAddr IncrementAddress) (bool, []byte, string, error) {
	addrDetExists, addrDetIsHardcoded, addrDet, addrDetPrefix, detailsErr := conn.FindAddressDetailsInRanges(prefixes, name)
	if detailsErr != nil {
		return false, []byte{}, "", detailsErr
	}

	if addrDetExists {
		addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(addrDetPrefix)
		if addrRangeErr != nil {
			return false, []byte{}, "", addrRangeErr
		}
		if !addrRangeExists {
			return false, []byte{}, "", errors.New(fmt.Sprintf("Error creating address in range with prefix '%s': Range doesn't exist", addrDetPrefix))
		}
		if addrRange.Type != rangeType {
			return false, []byte{}, "", errors.New(fmt.Sprintf("Error creating address in range with prefix '%s': Range type doesn't match the created address type", addrDetPrefix))
		}

		if !toleratePresent {
			return false, []byte{}, "", errors.New(fmt.Sprintf("Error creating address '%s': Address was already present in range with prefix '%s'", name, addrDetPrefix))
		}

		if addrDetIsHardcoded {
			return false, []byte{}, "", errors.New(fmt.Sprintf("Error creating address in range with prefix '%s': An existing address with the same name didn't match the expected hardcoded setting", addrDetPrefix))
		}

		return addrDetExists, addrDet, addrDetPrefix, nil
	}

	for _, prefix := range prefixes {
		addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(prefix)
		if addrRangeErr != nil {
			return false, []byte{}, "", addrRangeErr
		}
		if !addrRangeExists {
			return false, []byte{}, "", errors.New(fmt.Sprintf("Error creating address in range with prefix '%s': Range doesn't exist", prefix))
		}
		if addrRange.Type != rangeType {
			return false, []byte{}, "", errors.New(fmt.Sprintf("Error creating address in range with prefix '%s': Range type doesn't match the created address type", prefix))
		}

		genAddr, full, genErr := conn.createGeneratedAddressWithRetries(prefix, prefixes, name, addrIsGreater, incAddr, conn.Retries)
		if genErr != nil {
			return false, []byte{}, "", genErr
		}

		if full {
			continue
		}

		return addrDetExists, genAddr, prefix, nil
	}

	return false, []byte{}, "", errors.New("Error creating address in range with prefix '%s': Associated ranges are full")
}

func (conn *EtcdConnection) GenerateHardcodedAddressWithValidation(name string, prefixes []string, addr []byte, rangeType string, toleratePresent bool, prettify PrettifyAddr) (bool, string, error) {
	prefix, addrRange, matchFound, err := conn.FindAddressRangeByBoundaries(prefixes, addr)
	if err != nil {
		return false, "", err
	}

	if !matchFound {
		return false, "", errors.New(fmt.Sprintf("Error creating hardcoded address '%s': Address is outside boundaries of the input ranges", name))
	}

	addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(prefix)
	if addrRangeErr != nil {
		return false, "", addrRangeErr
	}
	if !addrRangeExists {
		return false, "", errors.New(fmt.Sprintf("Error creating hardcoded address in range at prefix '%s': Range doesn't exist", prefix))
	}
	if addrRange.Type != rangeType {
		return false, "", errors.New(fmt.Sprintf("Error creating hardcoded address in range at prefix '%s': Range type doesn't match the created address type", prefix))
	}

	addrDetExists, addrDetIsHardcoded, addrDet, detailsErr := conn.GetAddressDetails(prefix, name)
	if detailsErr != nil {
		return false, "", detailsErr
	}

	if addrDetExists {
		if !toleratePresent {
			return false, "", errors.New(fmt.Sprintf("Error creating hardcoded address '%s': Address was already present in range with prefix '%s'", name, prefix))
		}

		if !addrDetIsHardcoded {
			return false, "", errors.New(fmt.Sprintf("Error creating hardcoded address in range with prefix '%s': An existing address with the same name didn't match the expected hardcoded setting", prefix))
		}

		if bytes.Compare(addr, addrDet) != 0 {
			return false, "", errors.New(fmt.Sprintf("Error creating hardcoded address in range with prefix '%s': An existing address with the same name didn't match the expected address value", prefix))
		}

		return addrDetExists, prefix, nil
	}

	return addrDetExists, prefix, conn.CreateHardcodedAddress(prefix, name, addr, prettify)
}

func (conn *EtcdConnection) GetAddressWithValidation(name string, keyPrefix string, rangeType string, tolerateMissing bool) ([]byte, bool, error) {
	addrRange, addrRangeExists, addrRangeErr := conn.GetAddrRange(keyPrefix)
	if addrRangeErr != nil {
		return []byte{}, false, errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': %s", keyPrefix, addrRangeErr.Error()))
	}
	if !addrRangeExists {
		return []byte{}, false, errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range does not exist", keyPrefix))
	}
	if addrRange.Type != rangeType {
		return []byte{}, false, errors.New(fmt.Sprintf("Error retrieving address range at prefix '%s': Range type does not match address type", keyPrefix))
	}	

	addr, found, err := conn.FindAddress(keyPrefix, name)
	if err != nil {
		return []byte{}, false, err
	}

	if !found {
		if !tolerateMissing{
			return []byte{}, false, errors.New(fmt.Sprintf("Error retrieving address '%s' in range at prefix '%s': Address was not found in range", name, keyPrefix))
		}

		return []byte{}, found, nil
	}

	return addr, found, nil
}

func (conn *EtcdConnection) DeleteAddressWithValidation(name string, keyPrefix string, isHardcoded bool, addr []byte, tolerateMissing bool, prettify PrettifyAddr, addrIsLess AddressIsLess) (bool, error) {
	addrDetExists, addrDetIsHardcoded, addrDet, detailsErr := conn.GetAddressDetails(keyPrefix, name)
	if detailsErr != nil {
		return false, detailsErr
	}

	if !addrDetExists {
		if !tolerateMissing{
			return false, errors.New(fmt.Sprintf("Error deleting address '%s' in range at prefix '%s': Address was not found in range", name, keyPrefix))
		}
		
		return addrDetExists, nil
	}

	if addrDetIsHardcoded != isHardcoded {
		return false, errors.New(fmt.Sprintf("Error deleting address '%s' in range at prefix '%s': Address didn't match expected hardcoded setting", name, keyPrefix))
	}

	if bytes.Compare(addr, addrDet) != 0 {
		return false, errors.New(fmt.Sprintf("Error deleting address '%s' in range at prefix '%s': Address didn't have expected value", name, keyPrefix))
	}

	if isHardcoded {
		err := conn.DeleteHardcodedAddress(keyPrefix, name, addr, prettify, addrIsLess)
		if err != nil {
			return addrDetExists, err
		}
	} else {
		err := conn.DeleteGeneratedAddress(keyPrefix, name, addr, prettify)
		if err != nil {
			return addrDetExists, err
		}
	}

	return addrDetExists, nil
}

