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

func (k Keeper) PeriodicICQsAll(
	c context.Context,
	req *types.QueryAllPeriodicICQsRequest,
) (*types.QueryAllPeriodicICQsResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var periodicICQs []types.PeriodicICQs
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	periodicICQStore := prefix.NewStore(store, types.KeyPrefix(types.PeriodicICQsKey))

	pageRes, err := query.Paginate(periodicICQStore, req.Pagination, func(key []byte, value []byte) error {
		var periodicICQ types.PeriodicICQs
		if err := k.cdc.Unmarshal(value, &periodicICQ); err != nil {
			return err
		}

		periodicICQs = append(periodicICQs, periodicICQ)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPeriodicICQsResponse{PeriodicICQs: periodicICQs, Pagination: pageRes}, nil
}

func (k Keeper) PeriodicICQs(
	c context.Context,
	req *types.QueryGetPeriodicICQsRequest,
) (*types.QueryGetPeriodicICQsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	interchainqueryResult, found := k.GetPeriodicICQs(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetPeriodicICQsResponse{PeriodicICQs: interchainqueryResult}, nil
}
