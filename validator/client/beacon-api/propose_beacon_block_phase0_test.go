package beacon_api

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/theQRL/go-zond/common/hexutil"
	"github.com/theQRL/qrysm/v4/beacon-chain/rpc/apimiddleware"
	zondpb "github.com/theQRL/qrysm/v4/proto/prysm/v1alpha1"
	"github.com/theQRL/qrysm/v4/testing/assert"
	"github.com/theQRL/qrysm/v4/testing/require"
	"github.com/theQRL/qrysm/v4/validator/client/beacon-api/mock"
	test_helpers "github.com/theQRL/qrysm/v4/validator/client/beacon-api/test-helpers"
)

func TestProposeBeaconBlock_Phase0(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	jsonRestHandler := mock.NewMockjsonRestHandler(ctrl)

	phase0Block := generateSignedPhase0Block()

	genericSignedBlock := &zondpb.GenericSignedBeaconBlock{}
	genericSignedBlock.Block = phase0Block

	jsonPhase0Block := &apimiddleware.SignedBeaconBlockJson{
		Signature: hexutil.Encode(phase0Block.Phase0.Signature),
		Message: &apimiddleware.BeaconBlockJson{
			ParentRoot:    hexutil.Encode(phase0Block.Phase0.Block.ParentRoot),
			ProposerIndex: uint64ToString(phase0Block.Phase0.Block.ProposerIndex),
			Slot:          uint64ToString(phase0Block.Phase0.Block.Slot),
			StateRoot:     hexutil.Encode(phase0Block.Phase0.Block.StateRoot),
			Body: &apimiddleware.BeaconBlockBodyJson{
				Attestations:      jsonifyAttestations(phase0Block.Phase0.Block.Body.Attestations),
				AttesterSlashings: jsonifyAttesterSlashings(phase0Block.Phase0.Block.Body.AttesterSlashings),
				Deposits:          jsonifyDeposits(phase0Block.Phase0.Block.Body.Deposits),
				Eth1Data:          jsonifyEth1Data(phase0Block.Phase0.Block.Body.Eth1Data),
				Graffiti:          hexutil.Encode(phase0Block.Phase0.Block.Body.Graffiti),
				ProposerSlashings: jsonifyProposerSlashings(phase0Block.Phase0.Block.Body.ProposerSlashings),
				RandaoReveal:      hexutil.Encode(phase0Block.Phase0.Block.Body.RandaoReveal),
				VoluntaryExits:    JsonifySignedVoluntaryExits(phase0Block.Phase0.Block.Body.VoluntaryExits),
			},
		},
	}

	marshalledBlock, err := json.Marshal(jsonPhase0Block)
	require.NoError(t, err)

	ctx := context.Background()

	// Make sure that what we send in the POST body is the marshalled version of the protobuf block
	headers := map[string]string{"Eth-Consensus-Version": "phase0"}
	jsonRestHandler.EXPECT().PostRestJson(
		ctx,
		"/zond/v1/beacon/blocks",
		headers,
		bytes.NewBuffer(marshalledBlock),
		nil,
	)

	validatorClient := &beaconApiValidatorClient{jsonRestHandler: jsonRestHandler}
	proposeResponse, err := validatorClient.proposeBeaconBlock(ctx, genericSignedBlock)
	assert.NoError(t, err)
	require.NotNil(t, proposeResponse)

	expectedBlockRoot, err := phase0Block.Phase0.Block.HashTreeRoot()
	require.NoError(t, err)

	// Make sure that the block root is set
	assert.DeepEqual(t, expectedBlockRoot[:], proposeResponse.BlockRoot)
}

func generateSignedPhase0Block() *zondpb.GenericSignedBeaconBlock_Phase0 {
	return &zondpb.GenericSignedBeaconBlock_Phase0{
		Phase0: &zondpb.SignedBeaconBlock{
			Block:     test_helpers.GenerateProtoPhase0BeaconBlock(),
			Signature: test_helpers.FillByteSlice(96, 110),
		},
	}
}
