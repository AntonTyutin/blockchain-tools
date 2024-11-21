package utils

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-chain/go-hdwallet"
)

const hardenedShift = hdwallet.ZeroQuote

type PathSegments struct {
	Purpose      uint32
	CoinType     uint32
	Account      uint32
	Change       uint32
	AddressIndex uint32
}

func ParseDerivationPath(path string) (*PathSegments, error) {
	r := regexp.MustCompile(`m/(?<Purpose>\d+)(?<PurposeHardened>')?/(?<CoinType>\d+)(?<CoinTypeHardened>')?/(?<Account>\d+)(?<AccountHardened>')?/(?<Change>\d+)(?<ChangeHardened>')?/(?<AddressIndex>\d+)(?<AddressIndexHardened>')?`)
	matches := r.FindStringSubmatch(path)

	if matches == nil {
		return nil, fmt.Errorf("Derivation Path must be complied with https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki")
	}

	purpose, err := fetchSegmentValue(matches, r.SubexpIndex("Purpose"), r.SubexpIndex("PurposeHardened"))
	if err != nil {
		return nil, err
	}
	coinType, err := fetchSegmentValue(matches, r.SubexpIndex("CoinType"), r.SubexpIndex("CoinTypeHardened"))
	if err != nil {
		return nil, err
	}
	account, err := fetchSegmentValue(matches, r.SubexpIndex("Account"), r.SubexpIndex("AccountHardened"))
	if err != nil {
		return nil, err
	}
	change, err := fetchSegmentValue(matches, r.SubexpIndex("Change"), r.SubexpIndex("ChangeHardened"))
	if err != nil {
		return nil, err
	}
	addressIndex, err := fetchSegmentValue(matches, r.SubexpIndex("AddressIndex"), r.SubexpIndex("AddressIndexHardened"))
	if err != nil {
		return nil, err
	}
	
	parsedPath := &PathSegments{purpose, coinType, account, change, addressIndex}
	return parsedPath, nil
}

func fetchSegmentValue(matches []string, valueIndex, hardenedIndex int) (uint32, error) {
	idx, err := strconv.Atoi(matches[valueIndex])
	if err != nil {
		return 0, fmt.Errorf("Derivation Path must be complied with https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki. %w", err)
	}

	value := uint32(idx)
	if matches[hardenedIndex] != "" {
		value = value + hardenedShift
	}

	return value, nil
}
