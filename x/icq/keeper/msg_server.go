package keeper

import (
	"context"
	"fmt"

	ics23 "github.com/confio/ics23/go"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	commitmenttypes "github.com/cosmos/ibc-go/v3/modules/core/23-commitment/types"

	"github.com/simplyvc/interchainqueries/x/icq/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// SubmitICQResult submits a result for a given ICQ instance request
func (k msgServer) SubmitICQResult(
	goCtx context.Context,
	msg *types.MsgSubmitICQResults,
) (*types.MsgSubmitICQResultsResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Retrieve the given ICQ instance that we want to service to validate its existence
	icqInstance, found := k.GetPendingICQInstance(ctx, msg.QueryId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrICQNotFound, "(%d) not found", msg.QueryId)
	}

	// Retrieve the periodic ICQ for the full details of the query
	periodicICQ, found := k.GetPeriodicICQs(ctx, msg.PeriodicId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrPeriodicICQsNotFound, "(%d) not found", msg.PeriodicId)
	}

	if icqInstance.PeriodicId != msg.PeriodicId {
		return nil, sdkerrors.Wrapf(types.ErrPeriodicIdNoMatch, "stored: (%d) doesn't match submitted (%d)",
			icqInstance.PeriodicId, msg.PeriodicId)
	}

	// Retrieve the consensus state to validate the proofs
	consensusState, found := k.clientKeeper.GetClientConsensusState(ctx, periodicICQ.ClientId, msg.Height)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrClientConsensusNotFound, "(%s) not found at height (%d)",
			periodicICQ.ClientId, msg.Height.RevisionHeight)
	}

	individualRequests := periodicICQ.IndividualICQs
	individualICQs := msg.IndividualResults

	if len(individualICQs) != len(individualRequests) {
		return nil, sdkerrors.Wrapf(types.ErrIavlRootVerification, "elements are diff lengths")
	}

	for i, element := range individualICQs {
		if element.Id != individualRequests[element.Id].Id {
			// TODO: There is a missconfiguration in the way the list was setup
			return nil, sdkerrors.Wrapf(types.ErrIavlRootVerification, "elements do not match up, not in order")
		}

		merkleProof, err := commitmenttypes.ConvertProofs(element.Proof)
		if err != nil {
			return nil, err
		}

		// NOTE: merkleProof.Proofs[0] is the IavlSpec Proof
		// 		 merkleProof.Proofs[1] is the TendermintSpec Proof
		iavlRoot, err := merkleProof.Proofs[0].Calculate()
		if err != nil {
			return nil, err
		}

		// Verify if the IavlSpec Root exists in the TendermintSpec Proof.
		rootVerified := ics23.VerifyMembership(
			ics23.TendermintSpec,
			consensusState.GetRoot().GetHash(),
			merkleProof.Proofs[1],
			merkleProof.Proofs[1].GetExist().Key,
			iavlRoot,
		)
		if !rootVerified {
			return nil, sdkerrors.Wrapf(types.ErrIavlRootVerification, "root not verified")
		}

		// Verify that the value exists for the stored key using the verified iavlRoot
		kvVerified := ics23.VerifyMembership(
			ics23.IavlSpec,
			iavlRoot,
			merkleProof.Proofs[0],
			individualRequests[i].QueryParameters,
			element.Result,
		)
		if !kvVerified {
			return nil, sdkerrors.Wrapf(types.ErrKVVerification, "value: (%s) not verified for key (%s)",
				element.Result, individualRequests[i].QueryParameters,
			)
		}

		// Verification passed we call the hook
		k.AfterDataIsValidated(ctx, *element)
		fmt.Print(element.Result)
	}

	k.RemovePendingICQInstance(ctx, msg.QueryId)

	return &types.MsgSubmitICQResultsResponse{}, nil
}
