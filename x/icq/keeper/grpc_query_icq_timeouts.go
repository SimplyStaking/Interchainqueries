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

func (k Keeper) ICQTimeoutsAll(
	c context.Context,
	req *types.QueryAllICQTimeoutsRequest,
) (*types.QueryAllICQTimeoutsResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var icqTimeouts []types.ICQTimeouts
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	icqTimeoutsStore := prefix.NewStore(store, types.KeyPrefix(types.ICQTimeoutsKey))

	pageRes, err := query.Paginate(icqTimeoutsStore, req.Pagination, func(key []byte, value []byte) error {
		var icqTimeout types.ICQTimeouts
		if err := k.cdc.Unmarshal(value, &icqTimeout); err != nil {
			return err
		}

		icqTimeouts = append(icqTimeouts, icqTimeout)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllICQTimeoutsResponse{ICQTimeouts: icqTimeouts, Pagination: pageRes}, nil
}

func (k Keeper) ICQTimeouts(
	c context.Context,
	req *types.QueryGetICQTimeoutsRequest,
) (*types.QueryGetICQTimeoutsResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	icqTimeouts, found := k.GetICQTimeouts(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetICQTimeoutsResponse{ICQTimeouts: icqTimeouts}, nil
}
