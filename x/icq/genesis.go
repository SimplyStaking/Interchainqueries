package icq

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/simplyvc/interchainqueries/x/icq/keeper"
	"github.com/simplyvc/interchainqueries/x/icq/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(_ sdk.Context, _ keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
