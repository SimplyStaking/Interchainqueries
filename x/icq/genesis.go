package icq

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/simplyvc/interchainqueries/x/icq/keeper"
	"github.com/simplyvc/interchainqueries/x/icq/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/export
	path := "/store/bank/key"
	clientId := "07-tendermint-0"
	chainId := "target-chain"

	add1, _ := sdk.AccAddressFromBech32("icq17dtl0mjt3t77kpuhg2edqzjpszulwhgzp4faq3")
	add2, _ := sdk.AccAddressFromBech32("icq1xprmyxmxhk3j7vl929famszskac8z8l84ljeak")

	individualICQs := []*types.IndividualICQ{
		&types.IndividualICQ{
			Id:              uint64(0),
			Path:            path,
			QueryParameters: append(banktypes.CreateAccountBalancesPrefix(add1), []byte("stake")...),
		}, &types.IndividualICQ{
			Id:              uint64(1),
			Path:            path,
			QueryParameters: append(banktypes.CreateAccountBalancesPrefix(add2), []byte("stake")...),
		}}

	icqPeriodic := types.PeriodicICQs{
		IndividualICQs:       individualICQs,
		TimeoutHeightPadding: uint64(20),
		TargetHeight:         uint64(0),
		ClientId:             clientId,
		ChainId:              chainId,
		BlockRepeat:          uint64(10),
		LastHeightExecuted:   uint64(ctx.BlockHeight()),
		MaxResults:           uint64(10),
	}
	k.AppendPeriodicICQs(ctx, icqPeriodic)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(_ sdk.Context, _ keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	return genesis
}
