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

func (k Keeper) PeriodicICQAll(
	c context.Context,
	req *types.QueryAllPeriodicICQRequest,
) (*types.QueryAllPeriodicICQResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var periodicICQs []types.PeriodicICQ
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	periodicICQStore := prefix.NewStore(store, types.KeyPrefix(types.PeriodicICQKey))

	pageRes, err := query.Paginate(periodicICQStore, req.Pagination, func(key []byte, value []byte) error {
		var periodicICQ types.PeriodicICQ
		if err := k.cdc.Unmarshal(value, &periodicICQ); err != nil {
			return err
		}

		periodicICQs = append(periodicICQs, periodicICQ)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPeriodicICQResponse{PeriodicICQ: periodicICQs, Pagination: pageRes}, nil
}

func (k Keeper) PeriodicICQ(
	c context.Context,
	req *types.QueryGetPeriodicICQRequest,
) (*types.QueryGetPeriodicICQResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	interchainqueryResult, found := k.GetPeriodicICQ(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetPeriodicICQResponse{PeriodicICQ: interchainqueryResult}, nil
}
