package netcat

import (
	"github.com/feliux/netdbg/internal/logger"
	"github.com/spf13/cobra"
)

func ExecuteCommand(cmd *cobra.Command, args []string) {
	listen, _ := cmd.Flags().GetBool("listen")
	address, _ := cmd.Flags().GetString("address")
	port, _ := cmd.Flags().GetInt("port")
	protocol, _ := cmd.Flags().GetString("protocol")

	connector := NewConnector(protocol)

	if listen {
		err := connector.Listen(address, port)
		if err != nil {
			logger.Error("Error starting server", "err", err)
		}
	} else {
		zero, _ := cmd.Flags().GetBool("zero")
		err := connector.Connect(address, port, zero)
		if err != nil {
			logger.Error("Error connecting to server", "err", err)
		}
	}
}
