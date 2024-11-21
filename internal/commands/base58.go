package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/alecthomas/kingpin/v2"
	"github.com/btcsuite/btcutil/base58"
)

func ConfigureBase58Command(app *kingpin.Application) {
	c := &Base58Options{}
	pubAddr := app.Command("base58", "Base58 convertions").Action(c.run)
	pubAddr.Flag("decode", "Decode Base58 to hexadecimal").Short('d').BoolVar(&c.decode)
	pubAddr.Arg("address", "Address to encode/decode").Required().StringsVar(&c.address)
}

type Base58Options struct {
	decode  bool
	address []string
}

func (opts *Base58Options) run(c *kingpin.ParseContext) error {
	if !opts.decode {
		for _, hexAddress := range opts.address {
			encoded, err := base58encode(hexAddress)
			if err != nil {
				return err
			}
			fmt.Println(encoded)
		}
	} else {
		for _, address := range opts.address {
			decoded, err := base58decode(address)
			if err != nil {
				return err
			}
			fmt.Println(decoded)
		}
	}

	return nil
}

func base58encode(hexAddress string) (string, error) {
	addr, err := hex.DecodeString(hexAddress)
	if err != nil {
		return "", err
	}

	return base58.CheckEncode(addr[1:], addr[0]), nil
}

func base58decode(address string) (string, error) {
	bytes, ver, err := base58.CheckDecode(address)
	addr := make([]byte, 1, len(bytes)+1)
	addr[0] = ver
	addr = append(addr, bytes...)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(addr), nil
}
