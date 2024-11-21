package main

import (
	"os"

	"github.com/AntonTyutin/blockchain-tools/internal/commands"
	"github.com/alecthomas/kingpin/v2"
)

func main() {
	app := kingpin.New("blockchain-tools", "Utilities for doing blockchain-related cryptographic things")
	commands.ConfigurePubAddrCommand(app)
	commands.ConfigureBase58Command(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
