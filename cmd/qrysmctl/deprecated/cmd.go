package deprecated

import (
	"github.com/theQRL/qrysm/v4/cmd/qrysmctl/deprecated/checkpoint"
	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{}

func init() {
	Commands = append(Commands, checkpoint.Commands...)
}
