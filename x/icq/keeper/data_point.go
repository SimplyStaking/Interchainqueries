package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

// SetPeriodicLastDataPointId set a specific string to that points to the last data point in the store
func (k Keeper) SetPeriodicLastDataPointId(
	ctx sdk.Context,
	periodicId uint64,
	lastResultId string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicLastDataPointIdKey))
	store.Set(GetPeriodicLastDataPointIdIDBytes(periodicId), GetDataPointIDBytes(lastResultId))
}

// GetPeriodicLastDataPointId returns a last data point id
func (k Keeper) GetPeriodicLastDataPointId(ctx sdk.Context, id uint64) (val string, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicLastDataPointIdKey))
	b := store.Get(GetPeriodicLastDataPointIdIDBytes(id))
	if b == nil {
		return val, false
	}
	return string(b), true
}

// RemovePeriodicLastDataPointId removes a last data point id from the store
func (k Keeper) RemovePeriodicLastDataPointId(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicLastDataPointIdKey))
	store.Delete(GetPeriodicLastDataPointIdIDBytes(id))
}

// GetAllPeriodicLastDataPointId returns all last data point ids
func (k Keeper) GetAllPeriodicLastDataPointId(ctx sdk.Context) (list []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PeriodicLastDataPointIdKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, string(iterator.Value()))
	}

	return
}

// ----------------- Functions for DataPoints Keeper --------------------------- //

func (k Keeper) SetDataPoint(ctx sdk.Context, dataPoint types.DataPoint) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DataPointKey))
	b := k.cdc.MustMarshal(&dataPoint)
	store.Set(GetDataPointIDBytes(dataPoint.Id), b)
}

func (k Keeper) GetDataPoint(ctx sdk.Context, id string) (val types.DataPoint, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DataPointKey))
	b := store.Get(GetDataPointIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveDataPoint(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DataPointKey))
	store.Delete(GetDataPointIDBytes(id))
}

func (k Keeper) GetAllDataPoints(ctx sdk.Context) (list []types.DataPoint) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DataPointKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DataPoint
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllDataPointsForPeriodic helper function to get all data points for a given periodic query
func (k Keeper) GetAllDataPointsForPeriodic(ctx sdk.Context, periodicId uint64) (list []types.DataPoint) {
	lastResultId, found := k.GetPeriodicLastDataPointId(ctx, periodicId)
	if !found {
		return
	}

	breakPointId := lastResultId

	// Infinite loop to get all DataPoints until they are all retrieved
	for {
		prevDataPoint, found := k.GetDataPoint(ctx, lastResultId)
		if !found {
			return
		}

		lastResultId = prevDataPoint.PrevDataPointId
		// This check is here to ensure we do not double add if we do not exceed max results
		if prevDataPoint.Id == prevDataPoint.PrevDataPointId {
			return
		}

		list = append(list, prevDataPoint)
		if breakPointId == lastResultId {
			return
		}
	}
}

// RemoveAllDataPoints helper function to remove all data points for a given Periodic ICQ
func (k Keeper) RemoveAllDataPoints(ctx sdk.Context, periodicId uint64) {
	lastResultId, found := k.GetPeriodicLastDataPointId(ctx, periodicId)
	if !found {
		return
	}

	// Infinite loop to delete all DataPoints until they are no longer found
	for {
		prevDataPoint, found := k.GetDataPoint(ctx, lastResultId)
		if !found {
			return
		}
		lastResultId = prevDataPoint.PrevDataPointId
		k.RemoveDataPoint(ctx, prevDataPoint.Id)
	}
}

func GetDataPointIDBytes(id string) []byte {
	return []byte(id)
}

// GetPeriodicLastDataPointIdIDBytes returns the byte representation of the ID
func GetPeriodicLastDataPointIdIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetPeriodicLastDataPointIdIDFromBytes returns ID in uint64 format from a byte array
func GetPeriodicLastDataPointIdIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
