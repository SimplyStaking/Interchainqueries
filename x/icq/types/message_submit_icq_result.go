package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/tendermint/tendermint/proto/tendermint/crypto"
)

// interchainquery message types
const (
	TypeMsgSubmitICQResult = "submiticqresult"
)

var _ sdk.Msg = &MsgSubmitICQResult{}

// NewMsgSubmitICQResult - construct a msg to fulfil query request.
//nolint:interfacer
func NewMsgSubmitICQResult(chain_id string, result []byte, from_address sdk.Address, proof *crypto.ProofOps,
	periodic_id, query_id uint64, height *ibcclienttypes.Height) *MsgSubmitICQResult {
	return &MsgSubmitICQResult{
		ChainId:     chain_id,
		QueryId:     query_id,
		Result:      result,
		Height:      height,
		FromAddress: from_address.String(),
		Proof:       proof,
		PeriodicId:  periodic_id,
	}
}

// Route Implements Msg.
func (msg MsgSubmitICQResult) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSubmitICQResult) Type() string { return TypeMsgSubmitICQResult }

// ValidateBasic Implements Msg.
func (msg MsgSubmitICQResult) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSubmitICQResult) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgSubmitICQResult) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}
