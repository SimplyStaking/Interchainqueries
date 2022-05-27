package cli

import (
	"fmt"
	// "strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(_ string) *cobra.Command {
	// Group icq queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdListPendingICQRequests())
	cmd.AddCommand(CmdShowPendingICQRequest())

	cmd.AddCommand(CmdListPeriodicICQLastDataPointIds())
	cmd.AddCommand(CmdShowPeriodicICQLastDataPointId())

	cmd.AddCommand(CmdListPendingICQRequestsTimeouts())
	cmd.AddCommand(CmdShowPendingICQRequestTimeouts())

	cmd.AddCommand(CmdListPeriodicICQs())
	cmd.AddCommand(CmdShowPeriodicICQ())

	cmd.AddCommand(CmdListDataPoints())
	cmd.AddCommand(CmdShowDataPoint())

	return cmd
}
