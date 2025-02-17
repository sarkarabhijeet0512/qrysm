package beacon_api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/theQRL/qrysm/v4/api/gateway/apimiddleware"
	zondpb "github.com/theQRL/qrysm/v4/proto/prysm/v1alpha1"
	"github.com/theQRL/qrysm/v4/testing/assert"
	"github.com/theQRL/qrysm/v4/validator/client/beacon-api/mock"
)

func TestProposeBeaconBlock_Error(t *testing.T) {
	testSuites := []struct {
		name                 string
		expectedErrorMessage string
		expectedHttpError    *apimiddleware.DefaultErrorJson
	}{
		{
			name:                 "error 202",
			expectedErrorMessage: "block was successfully broadcasted but failed validation",
			expectedHttpError: &apimiddleware.DefaultErrorJson{
				Code:    http.StatusAccepted,
				Message: "202 error",
			},
		},
		{
			name:                 "request failed",
			expectedErrorMessage: "failed to send POST data to REST endpoint",
			expectedHttpError:    nil,
		},
	}

	testCases := []struct {
		name             string
		consensusVersion string
		endpoint         string
		block            *zondpb.GenericSignedBeaconBlock
	}{
		{
			name:             "phase0",
			consensusVersion: "phase0",
			endpoint:         "/zond/v1/beacon/blocks",
			block: &zondpb.GenericSignedBeaconBlock{
				Block: generateSignedPhase0Block(),
			},
		},
		{
			name:             "altair",
			consensusVersion: "altair",
			endpoint:         "/zond/v1/beacon/blocks",
			block: &zondpb.GenericSignedBeaconBlock{
				Block: generateSignedAltairBlock(),
			},
		},
		{
			name:             "bellatrix",
			consensusVersion: "bellatrix",
			endpoint:         "/zond/v1/beacon/blocks",
			block: &zondpb.GenericSignedBeaconBlock{
				Block: generateSignedBellatrixBlock(),
			},
		},
		{
			name:             "blinded bellatrix",
			consensusVersion: "bellatrix",
			endpoint:         "/zond/v1/beacon/blinded_blocks",
			block: &zondpb.GenericSignedBeaconBlock{
				Block: generateSignedBlindedBellatrixBlock(),
			},
		},
		{
			name:             "blinded capella",
			consensusVersion: "capella",
			endpoint:         "/zond/v1/beacon/blinded_blocks",
			block: &zondpb.GenericSignedBeaconBlock{
				Block: generateSignedBlindedCapellaBlock(),
			},
		},
	}

	for _, testSuite := range testSuites {
		for _, testCase := range testCases {
			t.Run(testSuite.name+"/"+testCase.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				ctx := context.Background()
				jsonRestHandler := mock.NewMockjsonRestHandler(ctrl)

				headers := map[string]string{"Eth-Consensus-Version": testCase.consensusVersion}
				jsonRestHandler.EXPECT().PostRestJson(
					ctx,
					testCase.endpoint,
					headers,
					gomock.Any(),
					nil,
				).Return(
					testSuite.expectedHttpError,
					errors.New("foo error"),
				).Times(1)

				validatorClient := &beaconApiValidatorClient{jsonRestHandler: jsonRestHandler}
				_, err := validatorClient.proposeBeaconBlock(ctx, testCase.block)
				assert.ErrorContains(t, testSuite.expectedErrorMessage, err)
				assert.ErrorContains(t, "foo error", err)
			})
		}
	}
}

func TestProposeBeaconBlock_UnsupportedBlockType(t *testing.T) {
	validatorClient := &beaconApiValidatorClient{}
	_, err := validatorClient.proposeBeaconBlock(context.Background(), &zondpb.GenericSignedBeaconBlock{})
	assert.ErrorContains(t, "unsupported block type", err)
}
