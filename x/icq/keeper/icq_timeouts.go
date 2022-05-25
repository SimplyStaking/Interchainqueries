package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

func (k Keeper) SetICQTimeouts(ctx sdk.Context, timedoutICQs types.ICQTimeouts) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ICQTimeoutsKey))
	b := k.cdc.MustMarshal(&timedoutICQs)
	store.Set(GetICQTimeoutsIDBytes(timedoutICQs.PeriodicId), b)
}

func (k Keeper) GetICQTimeouts(ctx sdk.Context, id uint64) (val types.ICQTimeouts, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ICQTimeoutsKey))
	b := store.Get(GetICQTimeoutsIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// AppendICQTimeouts appends a ICQTimeouts in the store with a new id and update the count
func (k Keeper) AppendICQTimeouts(
	ctx sdk.Context,
	icq types.PendingICQInstance,
) uint64 {

	timeouts := uint64(1)
	// Get the Timeout data for periodic query
	q, found := k.GetICQTimeouts(ctx, icq.PeriodicId)
	if found {
		timeouts = timeouts + q.Timeouts
	}

	timedoutICQ := types.ICQTimeouts{
		PeriodicId:        icq.PeriodicId,
		Timeouts:          timeouts,
		LastTimeoutHeight: icq.TimeoutHeight,
	}

	k.SetICQTimeouts(ctx, timedoutICQ)

	return timeouts
}

// RemoveICQTimeouts removes a ICQTimeouts from the store
func (k Keeper) RemoveICQTimeouts(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ICQTimeoutsKey))
	store.Delete(GetICQTimeoutsIDBytes(id))
}

// GetAllICQTimeouts returns all ICQTimeouts
func (k Keeper) GetAllICQTimeouts(ctx sdk.Context) (list []types.ICQTimeouts) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ICQTimeoutsKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ICQTimeouts
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetICQTimeoutsIDBytes returns the byte representation of the ID
func GetICQTimeoutsIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetICQTimeoutsIDFromBytes returns ID in uint64 format from a byte array
func GetICQTimeoutsIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
