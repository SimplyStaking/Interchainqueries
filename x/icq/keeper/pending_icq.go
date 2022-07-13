package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

// GetPendingICQInstanceCount get the total number of pending icq instances
func (k Keeper) GetPendingICQInstanceCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.PendingICQInstanceCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// setPendingICQInstanceCount set the total number of icq instances count
func (k Keeper) setPendingICQInstanceCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.PendingICQInstanceCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendPendingICQInstance appends an icq instance in the store with a new id and update the count
func (k Keeper) AppendPendingICQInstance(
	ctx sdk.Context,
	query types.PeriodicICQs,
	currHeight uint64,
) uint64 {

	pendingICQInstance := types.PendingICQInstance{
		TimeoutHeight: currHeight + query.TimeoutHeightPadding,
		TargetHeight:  query.TargetHeight,
		PeriodicId:    query.Id,
	}

	count := k.GetPendingICQInstanceCount(ctx)

	// Set the ID of the appended value
	pendingICQInstance.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingICQInstanceKey))
	appendedValue := k.cdc.MustMarshal(&pendingICQInstance)
	store.Set(GetPendingICQIDBytes(pendingICQInstance.Id), appendedValue)

	// Update pending icq count
	k.setPendingICQInstanceCount(ctx, count+1)

	return count
}

// GetPendingICQInstance returns a pending icq instance
func (k Keeper) GetPendingICQInstance(ctx sdk.Context, id uint64) (val types.PendingICQInstance, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingICQInstanceKey))
	b := store.Get(GetPendingICQIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePendingICQInstance removes a pending icq instance from the store
func (k Keeper) RemovePendingICQInstance(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingICQInstanceKey))
	store.Delete(GetPendingICQIDBytes(id))
}

// GetAllPendingICQInstance returns all pending interchainquery
func (k Keeper) GetAllPendingICQInstance(ctx sdk.Context) (list []types.PendingICQInstance) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingICQInstanceKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PendingICQInstance
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// ProcessBeginBlockPendingQueries iterates through all the pending queries,
// if a query is found to be expired it is appeneded to the Timedout list
func (k Keeper) ProcessBeginBlockPendingQueries(ctx sdk.Context) {
	currHeight := uint64(ctx.BlockHeight())
	for _, query := range k.GetAllPendingICQInstance(ctx) {
		if currHeight >= query.TimeoutHeight {
			k.AppendICQTimeouts(ctx, query)
			k.RemovePendingICQInstance(ctx, query.Id)
		}
	}
}

// RemoveAllICQInstancesForPeriodic remove all ICQ instance for the periodic ID
func (k Keeper) RemoveAllICQInstancesForPeriodic(ctx sdk.Context, periodicId uint64) {
	pendingICQInstances := k.GetAllPendingICQInstance(ctx)
	for _, icqInstance := range pendingICQInstances {
		if icqInstance.PeriodicId == periodicId {
			k.RemovePendingICQInstance(ctx, icqInstance.Id)
		}
	}
}

func NewPendingICQsRequest(
	pendingICQ types.PendingICQInstance,
	periodicICQ types.PeriodicICQs,
) types.PendingICQsRequest {
	return types.PendingICQsRequest{
		Id:             pendingICQ.Id,
		IndividualICQs: periodicICQ.IndividualICQs,
		TimeoutHeight:  pendingICQ.TimeoutHeight,
		TargetHeight:   pendingICQ.TargetHeight,
		ClientId:       periodicICQ.ClientId,
		PeriodicId:     periodicICQ.Id,
	}
}

// GetPendingICQIDBytes returns the byte representation of the ID
func GetPendingICQIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}
