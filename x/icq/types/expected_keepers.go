package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Event Hooks
// These can be utilized to communicate between the ICQ keeper and another keeper.

// ICQHooks event hooks
type ICQHooks interface {
	AfterDataIsValidated(ctx sdk.Context, result IndividualResult)
}
