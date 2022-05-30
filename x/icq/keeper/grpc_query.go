package keeper

import (
	"github.com/simplyvc/interchainqueries/x/icq/types"
)

var _ types.QueryServer = Keeper{}
