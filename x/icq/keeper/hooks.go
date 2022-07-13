package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/simplyvc/interchainqueries/x/icq/types"
)

// Implements ICQHooks interface
var _ types.ICQHooks = Keeper{}

// AfterValidatorCreated - call hook if registered
func (k Keeper) AfterDataIsValidated(ctx sdk.Context, result types.IndividualResult) {
	if k.hooks != nil {
		k.hooks.AfterDataIsValidated(ctx, result)
	}
}
