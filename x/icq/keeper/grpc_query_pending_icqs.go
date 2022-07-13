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

func (k Keeper) PendingICQsRequestAll(
	c context.Context,
	req *types.QueryAllPendingICQsRequest,
) (*types.QueryAllPendingICQsRequestResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pendingICQRequests []types.PendingICQsRequest
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	pendingICQStore := prefix.NewStore(store, types.KeyPrefix(types.PendingICQInstanceKey))

	pageRes, err := query.Paginate(pendingICQStore, req.Pagination, func(key []byte, value []byte) error {
		var pendingICQ types.PendingICQInstance

		if err := k.cdc.Unmarshal(value, &pendingICQ); err != nil {
			return err
		}

		periodicICQ, found := k.GetPeriodicICQs(ctx, pendingICQ.PeriodicId)
		if !found {
			return status.Error(codes.Internal, "Periodic ICQ doesn't exist for this ICQ Request!")
		}
		pendingICQRequests = append(
			pendingICQRequests,
			NewPendingICQsRequest(pendingICQ, periodicICQ),
		)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPendingICQsRequestResponse{PendingICQsRequest: pendingICQRequests, Pagination: pageRes}, nil
}

func (k Keeper) PendingICQsRequest(
	c context.Context,
	req *types.QueryGetPendingICQsRequest,
) (*types.QueryGetPendingICQsRequestResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	pendingICQ, found := k.GetPendingICQInstance(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}
	periodicICQ, found := k.GetPeriodicICQs(ctx, pendingICQ.PeriodicId)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetPendingICQsRequestResponse{
		PendingICQsRequest: NewPendingICQsRequest(pendingICQ, periodicICQ)}, nil
}
