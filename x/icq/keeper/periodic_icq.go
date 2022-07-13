package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

// GetPeriodicICQsCount gets the total number of periodic interchainqueries
func (k Keeper) GetPeriodicICQsCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.PeriodicICQsCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetPeriodicICQs(ctx sdk.Context, interchainquery types.PeriodicICQs) {
	// TODO When setting a PeriodicICQs we need to clear all the Datapoints for it as well as an ICQ instances
	// this needs to happen in case of data changes.
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicICQsKey))
	b := k.cdc.MustMarshal(&interchainquery)
	store.Set(GetPeriodicICQsIDBytes(interchainquery.Id), b)
}

// setPeriodicICQsCount sets the total number of Periodic ICQs, to be only used internally
func (k Keeper) setPeriodicICQsCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.PeriodicICQsCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendPeriodicICQs appends a periodic ICQ to the store and increases count
func (k Keeper) AppendPeriodicICQs(
	ctx sdk.Context,
	interchainquery types.PeriodicICQs,
) uint64 {
	// Create the interchainquery
	count := k.GetPeriodicICQsCount(ctx)

	// Set the ID of the appended value
	interchainquery.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicICQsKey))
	appendedValue := k.cdc.MustMarshal(&interchainquery)
	store.Set(GetPeriodicICQsIDBytes(interchainquery.Id), appendedValue)

	// Update pending interchainquery count
	k.setPeriodicICQsCount(ctx, count+1)

	return count
}

// GetPeriodicICQs returns a periodic ICQ from its id
func (k Keeper) GetPeriodicICQs(ctx sdk.Context, id uint64) (val types.PeriodicICQs, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicICQsKey))
	b := store.Get(GetPeriodicICQsIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePeriodicICQs removes a periodic ICQ from the store as well as all ICQ instances and DataPoints
func (k Keeper) RemovePeriodicICQs(ctx sdk.Context, periodicId uint64) {

	// Clean up all data dependant on the Periodic Query
	k.RemoveAllICQInstancesForPeriodic(ctx, periodicId)
	k.RemoveICQTimeouts(ctx, periodicId)

	// Finally remove the periodic query itself
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicICQsKey))
	store.Delete(GetPeriodicICQsIDBytes(periodicId))
}

// GetAllPeriodicICQs returns all periodic ICQs in store
func (k Keeper) GetAllPeriodicICQs(ctx sdk.Context) (list []types.PeriodicICQs) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicICQsKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PeriodicICQs
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetPeriodicICQsIDBytes returns the byte representation of the ID
func GetPeriodicICQsIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetPeriodicICQsIDFromBytes returns ID in uint64 format from a byte array
func GetPeriodicICQsIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

// ExecutePeriodicQueries iterates through all the periodic queries,
// if they meet the criteria to be executed then they are added to the queries list.
func (k Keeper) ExecutePeriodicQueries(ctx sdk.Context) {
	currHeight := uint64(ctx.BlockHeight())

	for _, query := range k.GetAllPeriodicICQs(ctx) {
		if query.LastHeightExecuted+query.BlockRepeat == currHeight {
			k.AppendPendingICQInstance(ctx, query, currHeight)
			query.LastHeightExecuted = currHeight
			k.SetPeriodicICQs(ctx, query)
		}
	}
}
