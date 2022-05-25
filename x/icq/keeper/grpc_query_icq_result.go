package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

func (k Keeper) ICQResultAll(
	c context.Context,
	req *types.QueryAllICQResultRequest,
) (*types.QueryAllICQResultResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var icqResults []types.ICQResult
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	icqResultStore := prefix.NewStore(store, types.KeyPrefix(types.ICQResultKey))

	pageRes, err := query.Paginate(icqResultStore, req.Pagination, func(key []byte, value []byte) error {
		var icqResult types.ICQResult
		if err := k.cdc.Unmarshal(value, &icqResult); err != nil {
			return err
		}

		icqResults = append(icqResults, icqResult)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllICQResultResponse{ICQResult: icqResults, Pagination: pageRes}, nil
}

func (k Keeper) ICQResult(
	c context.Context,
	req *types.QueryGetICQResultRequest,
) (*types.QueryGetICQResultResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	icqResult, found := k.GetICQResult(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetICQResultResponse{ICQResult: icqResult}, nil
}

func (k Keeper) DataPointResultAll(
	c context.Context,
	req *types.QueryAllDataPointResultRequest,
) (*types.QueryAllDataPointResultResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var dataPointResults []types.DataPointResult
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	dataPointResultStore := prefix.NewStore(store, types.KeyPrefix(types.DataPointResultKey))

	pageRes, err := query.Paginate(dataPointResultStore, req.Pagination, func(key []byte, value []byte) error {
		var dataPointResult types.DataPointResult
		if err := k.cdc.Unmarshal(value, &dataPointResult); err != nil {
			return err
		}

		dataPointResults = append(dataPointResults, dataPointResult)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDataPointResultResponse{DataPointResult: dataPointResults, Pagination: pageRes}, nil
}

func (k Keeper) DataPointResult(
	c context.Context,
	req *types.QueryGetDataPointResultRequest,
) (*types.QueryGetDataPointResultResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	dataPointResult, found := k.GetDataPointResult(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetDataPointResultResponse{DataPointResult: dataPointResult}, nil
}
