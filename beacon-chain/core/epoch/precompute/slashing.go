package precompute

import (
	"github.com/theQRL/qrysm/v4/beacon-chain/core/helpers"
	"github.com/theQRL/qrysm/v4/beacon-chain/core/time"
	"github.com/theQRL/qrysm/v4/beacon-chain/state"
	"github.com/theQRL/qrysm/v4/config/params"
	"github.com/theQRL/qrysm/v4/consensus-types/primitives"
	"github.com/theQRL/qrysm/v4/math"
	zondpb "github.com/theQRL/qrysm/v4/proto/prysm/v1alpha1"
)

// ProcessSlashingsPrecompute processes the slashed validators during epoch processing.
// This is an optimized version by passing in precomputed total epoch balances.
func ProcessSlashingsPrecompute(s state.BeaconState, pBal *Balance) error {
	currentEpoch := time.CurrentEpoch(s)
	exitLength := params.BeaconConfig().EpochsPerSlashingsVector

	// Compute the sum of state slashings
	slashings := s.Slashings()
	totalSlashing := uint64(0)
	for _, slashing := range slashings {
		totalSlashing += slashing
	}

	minSlashing := math.Min(totalSlashing*params.BeaconConfig().ProportionalSlashingMultiplier, pBal.ActiveCurrentEpoch)
	epochToWithdraw := currentEpoch + exitLength/2

	var hasSlashing bool
	// Iterate through validator list in state, stop until a validator satisfies slashing condition of current epoch.
	err := s.ReadFromEveryValidator(func(idx int, val state.ReadOnlyValidator) error {
		correctEpoch := epochToWithdraw == val.WithdrawableEpoch()
		if val.Slashed() && correctEpoch {
			hasSlashing = true
		}
		return nil
	})
	if err != nil {
		return err
	}
	// Exit early if there's no meaningful slashing to process.
	if !hasSlashing {
		return nil
	}

	increment := params.BeaconConfig().EffectiveBalanceIncrement
	validatorFunc := func(idx int, val *zondpb.Validator) (bool, *zondpb.Validator, error) {
		correctEpoch := epochToWithdraw == val.WithdrawableEpoch
		if val.Slashed && correctEpoch {
			penaltyNumerator := val.EffectiveBalance / increment * minSlashing
			penalty := penaltyNumerator / pBal.ActiveCurrentEpoch * increment
			if err := helpers.DecreaseBalance(s, primitives.ValidatorIndex(idx), penalty); err != nil {
				return false, val, err
			}
			return true, val, nil
		}
		return false, val, nil
	}

	return s.ApplyToEveryValidator(validatorFunc)
}
