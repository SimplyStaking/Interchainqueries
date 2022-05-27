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

func (k Keeper) PeriodicLastDataPointIdAll(
	c context.Context,
	req *types.QueryAllPeriodicLastDataPointIdRequest,
) (*types.QueryAllPeriodicLastDataPointIdResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var ids []string
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	idStore := prefix.NewStore(store, types.KeyPrefix(types.PeriodicLastDataPointIdKey))

	pageRes, err := query.Paginate(idStore, req.Pagination, func(_ []byte, value []byte) error {
		ids = append(ids, string(value))
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPeriodicLastDataPointIdResponse{LastDataPointId: ids, Pagination: pageRes}, nil
}

func (k Keeper) PeriodicLastDataPointId(
	c context.Context,
	req *types.QueryGetPeriodicLastDataPointIdRequest,
) (*types.QueryGetPeriodicLastDataPointIdResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	id, found := k.GetPeriodicLastDataPointId(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetPeriodicLastDataPointIdResponse{LastDataPointId: id}, nil
}

func (k Keeper) DataPointAll(
	c context.Context,
	req *types.QueryAllDataPointRequest,
) (*types.QueryAllDataPointResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var dataPoints []types.DataPoint
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	dataPointStore := prefix.NewStore(store, types.KeyPrefix(types.DataPointKey))

	pageRes, err := query.Paginate(dataPointStore, req.Pagination, func(_ []byte, value []byte) error {
		var dataPoint types.DataPoint
		if err := k.cdc.Unmarshal(value, &dataPoint); err != nil {
			return err
		}

		dataPoints = append(dataPoints, dataPoint)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDataPointResponse{DataPoint: dataPoints, Pagination: pageRes}, nil
}

func (k Keeper) DataPoint(
	c context.Context,
	req *types.QueryGetDataPointRequest,
) (*types.QueryGetDataPointResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	dataPoint, found := k.GetDataPoint(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetDataPointResponse{DataPoint: dataPoint}, nil
}

func (k Keeper) AllDataPointsForPeriodic(
	c context.Context,
	req *types.QueryAllDataPointsForPeriodicRequest,
) (*types.QueryAllDataPointsForPeriodicResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	// TODO add pagination?
	dataPoints := k.GetAllDataPointsForPeriodic(ctx, req.Id)

	return &types.QueryAllDataPointsForPeriodicResponse{DataPoint: dataPoints}, nil
}
