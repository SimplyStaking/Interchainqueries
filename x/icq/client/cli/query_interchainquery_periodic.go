package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

func CmdListPeriodicICQs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-periodic-icqs",
		Short: "list all periodic interchainqueries",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllPeriodicICQsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.PeriodicICQsAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowPeriodicICQs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-periodic-icq [id]",
		Short: "shows a periodic interchainquery for the specified id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetPeriodicICQsRequest{
				Id: id,
			}

			res, err := queryClient.PeriodicICQs(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
