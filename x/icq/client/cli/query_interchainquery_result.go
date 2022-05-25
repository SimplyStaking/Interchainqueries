package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

func CmdListPendingICQRequestsResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-interchainquery-result",
		Short: "list all interchainquery results for all periodic queries",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllICQResultRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ICQResultAll(context.Background(), params)
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

func CmdShowPendingICQRequestResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-interchainquery-result [id]",
		Short: "shows an interchainquery result for a prieodic query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetICQResultRequest{
				Id: id,
			}

			res, err := queryClient.ICQResult(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListDataPointResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-datapoint-results",
		Short: "list all datapoint results for all periodic queries",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllDataPointResultRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.DataPointResultAll(context.Background(), params)
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

func CmdShowDataPointResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-datapoint-result [id]",
		Short: "shows a datapoint result for a specific identifier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetDataPointResultRequest{
				Id: args[0],
			}

			res, err := queryClient.DataPointResult(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
