package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clientType "github.com/cosmos/ibc-go/v3/modules/core/02-client/keeper"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
		hooks    types.ICQHooks

		clientKeeper clientType.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	clientKeeper clientType.Keeper,
) *Keeper {

	return &Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		memKey:       memKey,
		clientKeeper: clientKeeper,
		hooks:        nil,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) SetHooks(icqHooks types.ICQHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set icq hooks twice")
	}

	k.hooks = icqHooks

	return k
}
