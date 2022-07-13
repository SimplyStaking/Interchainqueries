package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// interchainquery message types
const (
	TypeMsgSubmitICQResults = "submiticqresult"
)

var _ sdk.Msg = &MsgSubmitICQResults{}

// TODO fix this
// // NewMsgSubmitICQResults - construct a msg to fulfil query request.
// //nolint:interfacer
// func NewMsgSubmitICQResults(result []byte, from_address sdk.Address, proof *crypto.ProofOps,
// 	periodic_id, query_id uint64, height *ibcclienttypes.Height) *MsgSubmitICQResults {
// 	return &MsgSubmitICQResults{
// 		QueryId:     query_id,
// 		Result:      result,
// 		Height:      height,
// 		FromAddress: from_address.String(),
// 		Proof:       proof,
// 		PeriodicId:  periodic_id,
// 	}
// }

// Route Implements Msg.
func (msg MsgSubmitICQResults) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSubmitICQResults) Type() string {
	return TypeMsgSubmitICQResults
}

// ValidateBasic Implements Msg.
func (msg MsgSubmitICQResults) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSubmitICQResults) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgSubmitICQResults) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{fromAddress}
}
