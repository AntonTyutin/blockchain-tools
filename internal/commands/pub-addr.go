package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AntonTyutin/blockchain-tools/internal/utils"
	"github.com/alecthomas/kingpin/v2"
	"github.com/go-chain/go-hdwallet"
)

func ConfigurePubAddrCommand(app *kingpin.Application) {
	c := &PubAddrOptions{}
	pubAddr := app.Command("pub-addr", "Get public address").Action(c.run)
	pubAddr.Flag("seed-phrase-source", "Seed Phrase Source").Required().Short('f').StringVar(&c.SeedPhraseSource)
	pubAddr.Arg("derivation-paths", "Derivation Paths").Required().StringsVar(&c.DerivationPaths)
}

type PubAddrOptions struct {
	SeedPhraseSource string
	DerivationPaths  []string
}

func (opts *PubAddrOptions) run(c *kingpin.ParseContext) error {
	seedPhrase, err := getSeedPhrase(&opts.SeedPhraseSource)
	if err != nil {
		return err
	}

	seed, err := hdwallet.NewSeed(string(seedPhrase), "", hdwallet.English)
	if err != nil {
		return err
	}
	masterKey, err := hdwallet.NewKey(hdwallet.Seed(seed))
	if err != nil {
		return err
	}

	for _, derivationPath := range opts.DerivationPaths {
		childKey, err := getChild(masterKey, derivationPath)
		if err != nil {
			return err
		}
		address, err := childKey.GetAddress()
		if err != nil {
			return err
		}
		fmt.Println(address)
	}
	return nil
}

func getChild(key *hdwallet.Key, derivationPath string) (hdwallet.Wallet, error) {
	path, err := utils.ParseDerivationPath(derivationPath)
	if err != nil {
		return nil, err
	}

	walletOptions := func(o *hdwallet.Options) {
		o.Purpose = path.Purpose
		o.CoinType = path.CoinType
		o.Account = path.Account
		o.Change = path.Change
		o.AddressIndex = path.AddressIndex
	}

	return key.GetWallet(walletOptions)
}

func getSeedPhrase(seedPhraseSource *string) ([]byte, error) {

	var seedPhraseSourceFileStream *os.File

	if *seedPhraseSource == "" {
		seedPhraseSourceFileStream = os.Stdin
	} else {
		var err error
		seedPhraseSourceFileStream, err = os.Open(*seedPhraseSource)
		if err != nil {
			return []byte(""), fmt.Errorf("Unable to open file %s. %w", *seedPhraseSource, err)
		}
	}

	seedPhrase, err := bufio.NewReader(seedPhraseSourceFileStream).ReadBytes('\n')

	if err != nil {
		return []byte(""), fmt.Errorf("Unable to read seed phrase from given source. %s", err)
	}

	return seedPhrase[:len(seedPhrase)-1], nil
}
