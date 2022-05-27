package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

func CmdListPeriodicICQLastDataPointIds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-periodic-last-data-point-ids",
		Short: "list all the last data point ids for all the periodic queries",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllPeriodicLastDataPointIdRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.PeriodicLastDataPointIdAll(context.Background(), params)
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

func CmdShowPeriodicICQLastDataPointId() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-periodic-last-data-point [id]",
		Short: "shows the last data point id for a periodic query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetPeriodicLastDataPointIdRequest{
				Id: id,
			}

			res, err := queryClient.PeriodicLastDataPointId(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdListDataPoints() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-datapoints",
		Short: "list all datapoints for all periodic queries",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllDataPointRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.DataPointAll(context.Background(), params)
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

func CmdListDataPointsForPeriodic() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-datapoints-for-periodic",
		Short: "list all datapoints for all periodic queries",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllDataPointsForPeriodicRequest{
				Id:         id,
				Pagination: pageReq,
			}

			res, err := queryClient.AllDataPointsForPeriodic(context.Background(), params)
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

func CmdShowDataPoint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-datapoint [id]",
		Short: "shows a datapoint for a specific identifier",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetDataPointRequest{
				Id: args[0],
			}

			res, err := queryClient.DataPoint(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
