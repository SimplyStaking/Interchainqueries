package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

// GetPeriodicICQCount gets the total number of periodic interchainqueries
func (k Keeper) GetPeriodicICQCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.PeriodicICQCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetPeriodicICQ(ctx sdk.Context, interchainquery types.PeriodicICQ) {
	// TODO When setting a PeriodicICQ we need to clear all the Datapoints for it as well as an ICQ instances
	// this needs to happen in case of data changes.
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicICQKey))
	b := k.cdc.MustMarshal(&interchainquery)
	store.Set(GetPeriodicICQIDBytes(interchainquery.Id), b)
}

// setPeriodicICQCount sets the total number of Periodic ICQs, to be only used internally
func (k Keeper) setPeriodicICQCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.PeriodicICQCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendPeriodicICQ appends a periodic ICQ to the store and increases count
func (k Keeper) AppendPeriodicICQ(
	ctx sdk.Context,
	interchainquery types.PeriodicICQ,
) uint64 {
	// Create the interchainquery
	count := k.GetPeriodicICQCount(ctx)

	// Set the ID of the appended value
	interchainquery.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicICQKey))
	appendedValue := k.cdc.MustMarshal(&interchainquery)
	store.Set(GetPeriodicICQIDBytes(interchainquery.Id), appendedValue)

	// Update pending interchainquery count
	k.setPeriodicICQCount(ctx, count+1)

	return count
}

// GetPeriodicICQ returns a periodic ICQ from its id
func (k Keeper) GetPeriodicICQ(ctx sdk.Context, id uint64) (val types.PeriodicICQ, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicICQKey))
	b := store.Get(GetPeriodicICQIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePeriodicICQ removes a periodic ICQ from the store as well as all ICQ instances and DataPoints
func (k Keeper) RemovePeriodicICQ(ctx sdk.Context, periodicId uint64) {

	// Clean up all data dependant on the Periodic Query
	k.RemoveAllICQInstancesForPeriodic(ctx, periodicId)
	k.RemoveAllDataPointResults(ctx, periodicId)
	k.RemoveICQResult(ctx, periodicId)
	k.RemoveICQTimeouts(ctx, periodicId)

	// Finally remove the periodic query itself
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicICQKey))
	store.Delete(GetPeriodicICQIDBytes(periodicId))
}

// GetAllPeriodicICQ returns all periodic ICQs in store
func (k Keeper) GetAllPeriodicICQ(ctx sdk.Context) (list []types.PeriodicICQ) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicICQKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PeriodicICQ
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetPeriodicICQIDBytes returns the byte representation of the ID
func GetPeriodicICQIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetPeriodicICQIDFromBytes returns ID in uint64 format from a byte array
func GetPeriodicICQIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

// ExecutePeriodicQueries iterates through all the periodic queries,
// if they meet the criteria to be executed then they are added to the queries list.
func (k Keeper) ExecutePeriodicQueries(ctx sdk.Context) {
	currHeight := uint64(ctx.BlockHeight())

	for _, query := range k.GetAllPeriodicICQ(ctx) {
		if query.LastHeightExecuted+query.BlockRepeat == currHeight {
			k.AppendPendingICQInstance(ctx, query, currHeight)
			query.LastHeightExecuted = currHeight
			k.SetPeriodicICQ(ctx, query)
		}
	}
}
