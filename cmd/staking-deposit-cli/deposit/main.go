package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/theQRL/qrysm/v4/cmd/staking-deposit-cli/deposit/existingseed"
	"github.com/theQRL/qrysm/v4/cmd/staking-deposit-cli/deposit/generatedilithiumtoexecutionchange"
	"github.com/theQRL/qrysm/v4/cmd/staking-deposit-cli/deposit/newseed"
	"github.com/urfave/cli/v2"
)

var depositCommands []*cli.Command

func main() {
	app := &cli.App{
		Commands: depositCommands,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	depositCommands = append(depositCommands, existingseed.Commands...)
	depositCommands = append(depositCommands, newseed.Commands...)
	depositCommands = append(depositCommands, generatedilithiumtoexecutionchange.Commands...)
}
