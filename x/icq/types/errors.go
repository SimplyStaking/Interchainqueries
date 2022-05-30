package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/icq module sentinel errors
var (
	ErrICQNotFound                = sdkerrors.Register(ModuleName, 1, "icq not found")
	ErrPeriodicICQNotFound        = sdkerrors.Register(ModuleName, 2, "periodic icq not found")
	ErrPeriodicIdNoMatch          = sdkerrors.Register(ModuleName, 3, "periodic ids do not match")
	ErrIavlRootVerification       = sdkerrors.Register(ModuleName, 4, "failed to verify iavl root in ibc client")
	ErrKVVerification             = sdkerrors.Register(ModuleName, 5, "failed to verify the kv from iavl proof")
	ErrDuplicateHeightSubmissions = sdkerrors.Register(ModuleName, 6, "result for the target height already submitted")
	ErrClientConsensusNotFound    = sdkerrors.Register(ModuleName, 7, "client consensus not found")
	ErrDataPointKeyBroken         = sdkerrors.Register(ModuleName, 8, "data point id key broken")
	ErrResultNotSubmitted         = sdkerrors.Register(ModuleName, 9, "result is nil")
)
