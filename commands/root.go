package commands

import (
	"github.com/spf13/cobra"

	//"github.com/tendermint/tendermint/libs/log"
	cfg "github.com/triasteam/go-streamnet/config"
	//"os"
)

var (
	config = cfg.DefaultConfig()
	//logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)


// RootCmd is the root command for StreamNet core.
var RootCmd = &cobra.Command{
	Use:   "sng",
	Short: "StreamNet in Go",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
				return nil
	},
}

