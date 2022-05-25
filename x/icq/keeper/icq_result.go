package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

// GetICQResultCount get the total number of interchainqueryResult
func (k Keeper) GetICQResultCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.ICQResultCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// setICQResultCount set the total number of interchainqueryResult
func (k Keeper) setICQResultCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.ICQResultCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendICQResult appends a interchainqueryResult in the store with a new id and update the count
func (k Keeper) AppendICQResult(
	ctx sdk.Context,
	interchainqueryResult types.ICQResult,
) uint64 {
	// Create the interchainqueryResult
	count := k.GetICQResultCount(ctx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ICQResultKey))
	appendedValue := k.cdc.MustMarshal(&interchainqueryResult)
	store.Set(GetICQResultIDBytes(interchainqueryResult.PeriodicId), appendedValue)

	// Update interchainqueryResult count
	k.setICQResultCount(ctx, count+1)

	return count
}

// SetICQResult set a specific interchainqueryResult in the store
func (k Keeper) SetICQResult(ctx sdk.Context, interchainqueryResult types.ICQResult) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ICQResultKey))
	b := k.cdc.MustMarshal(&interchainqueryResult)
	store.Set(GetICQResultIDBytes(interchainqueryResult.PeriodicId), b)
}

// GetICQResult returns a interchainqueryResult from its id
func (k Keeper) GetICQResult(ctx sdk.Context, id uint64) (val types.ICQResult, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ICQResultKey))
	b := store.Get(GetICQResultIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveICQResult removes a interchainqueryResult from the store
func (k Keeper) RemoveICQResult(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ICQResultKey))
	store.Delete(GetICQResultIDBytes(id))
}

// GetAllICQResult returns all interchainqueryResult
func (k Keeper) GetAllICQResult(ctx sdk.Context) (list []types.ICQResult) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ICQResultKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ICQResult
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// ----------------- Functions for DataPoint Results Keeper --------------------------- //

func (k Keeper) SetDataPointResult(ctx sdk.Context, dataPointResult types.DataPointResult) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DataPointResultKey))
	b := k.cdc.MustMarshal(&dataPointResult)
	store.Set(GetDataPointResultIDBytes(dataPointResult.Id), b)
}

func (k Keeper) GetDataPointResult(ctx sdk.Context, id string) (val types.DataPointResult, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DataPointResultKey))
	b := store.Get(GetDataPointResultIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveDataPointResult(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DataPointResultKey))
	store.Delete(GetDataPointResultIDBytes(id))
}

func (k Keeper) GetAllDataPointResults(ctx sdk.Context) (list []types.DataPointResult) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DataPointResultKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DataPointResult
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// RemoveAllDataPointResults helper function to remove all data points for a given Periodic ICQ
func (k Keeper) RemoveAllDataPointResults(ctx sdk.Context, periodicId uint64) {
	icqResult, found := k.GetICQResult(ctx, periodicId)
	if !found {
		return
	}

	lastResultId := icqResult.LastResultId

	// Infinite loop to delete all DataPoints until they are no longer found
	for {
		lastDataPoint, found := k.GetDataPointResult(ctx, lastResultId)
		if !found {
			return
		}
		lastResultId = lastDataPoint.LastDataPointId
		k.RemoveDataPointResult(ctx, lastDataPoint.Id)
	}
}

func GetDataPointResultIDBytes(id string) []byte {
	return []byte(id)
}

// GetICQResultIDBytes returns the byte representation of the ID
func GetICQResultIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetICQResultIDFromBytes returns ID in uint64 format from a byte array
func GetICQResultIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
