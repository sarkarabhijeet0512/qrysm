package precompute_test

import (
	"context"
	"testing"

	"github.com/theQRL/go-bitfield"
	"github.com/theQRL/qrysm/v4/beacon-chain/core/epoch/precompute"
	"github.com/theQRL/qrysm/v4/beacon-chain/core/helpers"
	"github.com/theQRL/qrysm/v4/config/params"
	zondpb "github.com/theQRL/qrysm/v4/proto/prysm/v1alpha1"
	"github.com/theQRL/qrysm/v4/proto/prysm/v1alpha1/attestation"
	"github.com/theQRL/qrysm/v4/runtime/version"
	"github.com/theQRL/qrysm/v4/testing/assert"
	"github.com/theQRL/qrysm/v4/testing/require"
	"github.com/theQRL/qrysm/v4/testing/util"
)

func TestUpdateValidator_Works(t *testing.T) {
	e := params.BeaconConfig().FarFutureSlot
	vp := []*precompute.Validator{{}, {InclusionSlot: e}, {}, {InclusionSlot: e}, {}, {InclusionSlot: e}}
	record := &precompute.Validator{IsCurrentEpochAttester: true, IsCurrentEpochTargetAttester: true,
		IsPrevEpochAttester: true, IsPrevEpochTargetAttester: true, IsPrevEpochHeadAttester: true}
	a := &zondpb.PendingAttestation{InclusionDelay: 1, ProposerIndex: 2}

	// Indices 1 3 and 5 attested
	vp = precompute.UpdateValidator(vp, record, []uint64{1, 3, 5}, a, 100)

	wanted := &precompute.Validator{IsCurrentEpochAttester: true, IsCurrentEpochTargetAttester: true,
		IsPrevEpochAttester: true, IsPrevEpochTargetAttester: true, IsPrevEpochHeadAttester: true,
		ProposerIndex: 2, InclusionDistance: 1, InclusionSlot: 101}
	wantedVp := []*precompute.Validator{{}, wanted, {}, wanted, {}, wanted}
	assert.DeepEqual(t, wantedVp, vp, "Incorrect attesting validator calculations")
}

func TestUpdateValidator_InclusionOnlyCountsPrevEpoch(t *testing.T) {
	e := params.BeaconConfig().FarFutureSlot
	vp := []*precompute.Validator{{InclusionSlot: e}}
	record := &precompute.Validator{IsCurrentEpochAttester: true, IsCurrentEpochTargetAttester: true}
	a := &zondpb.PendingAttestation{InclusionDelay: 1, ProposerIndex: 2}

	// Verify inclusion info doesn't get updated.
	vp = precompute.UpdateValidator(vp, record, []uint64{0}, a, 100)
	wanted := &precompute.Validator{IsCurrentEpochAttester: true, IsCurrentEpochTargetAttester: true, InclusionSlot: e}
	wantedVp := []*precompute.Validator{wanted}
	assert.DeepEqual(t, wantedVp, vp, "Incorrect attesting validator calculations")
}

func TestUpdateBalance(t *testing.T) {
	vp := []*precompute.Validator{
		{IsCurrentEpochAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsCurrentEpochTargetAttester: true, IsCurrentEpochAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsCurrentEpochTargetAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsPrevEpochAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsPrevEpochAttester: true, IsPrevEpochTargetAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsPrevEpochHeadAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsPrevEpochAttester: true, IsPrevEpochHeadAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsSlashed: true, IsCurrentEpochAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
	}
	wantedPBal := &precompute.Balance{
		ActiveCurrentEpoch:         params.BeaconConfig().EffectiveBalanceIncrement,
		ActivePrevEpoch:            params.BeaconConfig().EffectiveBalanceIncrement,
		CurrentEpochAttested:       200 * params.BeaconConfig().EffectiveBalanceIncrement,
		CurrentEpochTargetAttested: 200 * params.BeaconConfig().EffectiveBalanceIncrement,
		PrevEpochAttested:          300 * params.BeaconConfig().EffectiveBalanceIncrement,
		PrevEpochTargetAttested:    100 * params.BeaconConfig().EffectiveBalanceIncrement,
		PrevEpochHeadAttested:      200 * params.BeaconConfig().EffectiveBalanceIncrement,
	}
	pBal := precompute.UpdateBalance(vp, &precompute.Balance{}, version.Phase0)
	assert.DeepEqual(t, wantedPBal, pBal, "Incorrect balance calculations")
}

func TestUpdateBalanceDifferentVersions(t *testing.T) {
	vp := []*precompute.Validator{
		{IsCurrentEpochAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsCurrentEpochTargetAttester: true, IsCurrentEpochAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsCurrentEpochTargetAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsPrevEpochAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsPrevEpochAttester: true, IsPrevEpochTargetAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsPrevEpochHeadAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsPrevEpochAttester: true, IsPrevEpochHeadAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
		{IsSlashed: true, IsCurrentEpochAttester: true, CurrentEpochEffectiveBalance: 100 * params.BeaconConfig().EffectiveBalanceIncrement},
	}
	wantedPBal := &precompute.Balance{
		ActiveCurrentEpoch:         params.BeaconConfig().EffectiveBalanceIncrement,
		ActivePrevEpoch:            params.BeaconConfig().EffectiveBalanceIncrement,
		CurrentEpochAttested:       200 * params.BeaconConfig().EffectiveBalanceIncrement,
		CurrentEpochTargetAttested: 200 * params.BeaconConfig().EffectiveBalanceIncrement,
		PrevEpochAttested:          params.BeaconConfig().EffectiveBalanceIncrement,
		PrevEpochTargetAttested:    100 * params.BeaconConfig().EffectiveBalanceIncrement,
		PrevEpochHeadAttested:      200 * params.BeaconConfig().EffectiveBalanceIncrement,
	}
	pBal := precompute.UpdateBalance(vp, &precompute.Balance{}, version.Bellatrix)
	assert.DeepEqual(t, wantedPBal, pBal, "Incorrect balance calculations")

	pBal = precompute.UpdateBalance(vp, &precompute.Balance{}, version.Capella)
	assert.DeepEqual(t, wantedPBal, pBal, "Incorrect balance calculations")
}

func TestSameHead(t *testing.T) {
	beaconState, _ := util.DeterministicGenesisState(t, 100)
	require.NoError(t, beaconState.SetSlot(1))
	att := &zondpb.Attestation{Data: &zondpb.AttestationData{
		Target: &zondpb.Checkpoint{Epoch: 0}}}
	r := [32]byte{'A'}
	br := beaconState.BlockRoots()
	br[0] = r[:]
	require.NoError(t, beaconState.SetBlockRoots(br))
	att.Data.BeaconBlockRoot = r[:]
	same, err := precompute.SameHead(beaconState, &zondpb.PendingAttestation{Data: att.Data})
	require.NoError(t, err)
	assert.Equal(t, true, same, "Head in state does not match head in attestation")
	newRoot := [32]byte{'B'}
	att.Data.BeaconBlockRoot = newRoot[:]
	same, err = precompute.SameHead(beaconState, &zondpb.PendingAttestation{Data: att.Data})
	require.NoError(t, err)
	assert.Equal(t, false, same, "Head in state matches head in attestation")
}

func TestSameTarget(t *testing.T) {
	beaconState, _ := util.DeterministicGenesisState(t, 100)
	require.NoError(t, beaconState.SetSlot(1))
	att := &zondpb.Attestation{Data: &zondpb.AttestationData{
		Target: &zondpb.Checkpoint{Epoch: 0}}}
	r := [32]byte{'A'}
	br := beaconState.BlockRoots()
	br[0] = r[:]
	require.NoError(t, beaconState.SetBlockRoots(br))
	att.Data.Target.Root = r[:]
	same, err := precompute.SameTarget(beaconState, &zondpb.PendingAttestation{Data: att.Data}, 0)
	require.NoError(t, err)
	assert.Equal(t, true, same, "Head in state does not match head in attestation")
	newRoot := [32]byte{'B'}
	att.Data.Target.Root = newRoot[:]
	same, err = precompute.SameTarget(beaconState, &zondpb.PendingAttestation{Data: att.Data}, 0)
	require.NoError(t, err)
	assert.Equal(t, false, same, "Head in state matches head in attestation")
}

func TestAttestedPrevEpoch(t *testing.T) {
	beaconState, _ := util.DeterministicGenesisState(t, 100)
	require.NoError(t, beaconState.SetSlot(params.BeaconConfig().SlotsPerEpoch))
	att := &zondpb.Attestation{Data: &zondpb.AttestationData{
		Target: &zondpb.Checkpoint{Epoch: 0}}}
	r := [32]byte{'A'}
	br := beaconState.BlockRoots()
	br[0] = r[:]
	require.NoError(t, beaconState.SetBlockRoots(br))
	att.Data.Target.Root = r[:]
	att.Data.BeaconBlockRoot = r[:]
	votedEpoch, votedTarget, votedHead, err := precompute.AttestedPrevEpoch(beaconState, &zondpb.PendingAttestation{Data: att.Data})
	require.NoError(t, err)
	assert.Equal(t, true, votedEpoch, "Did not vote epoch")
	assert.Equal(t, true, votedTarget, "Did not vote target")
	assert.Equal(t, true, votedHead, "Did not vote head")
}

func TestAttestedCurrentEpoch(t *testing.T) {
	beaconState, _ := util.DeterministicGenesisState(t, 100)
	require.NoError(t, beaconState.SetSlot(params.BeaconConfig().SlotsPerEpoch+1))
	att := &zondpb.Attestation{Data: &zondpb.AttestationData{
		Target: &zondpb.Checkpoint{Epoch: 1}}}
	r := [32]byte{'A'}

	br := beaconState.BlockRoots()
	br[params.BeaconConfig().SlotsPerEpoch] = r[:]
	require.NoError(t, beaconState.SetBlockRoots(br))
	att.Data.Target.Root = r[:]
	att.Data.BeaconBlockRoot = r[:]
	votedEpoch, votedTarget, err := precompute.AttestedCurrentEpoch(beaconState, &zondpb.PendingAttestation{Data: att.Data})
	require.NoError(t, err)
	assert.Equal(t, true, votedEpoch, "Did not vote epoch")
	assert.Equal(t, true, votedTarget, "Did not vote target")
}

func TestProcessAttestations(t *testing.T) {
	params.SetupTestConfigCleanup(t)
	params.OverrideBeaconConfig(params.MinimalSpecConfig())

	validators := uint64(128)
	beaconState, _ := util.DeterministicGenesisState(t, validators)
	require.NoError(t, beaconState.SetSlot(params.BeaconConfig().SlotsPerEpoch))
	c := helpers.SlotCommitteeCount(validators)
	bf := bitfield.NewBitlist(c)
	att1 := &zondpb.Attestation{Data: &zondpb.AttestationData{
		Target: &zondpb.Checkpoint{Epoch: 0}},
		AggregationBits: bf}
	att2 := &zondpb.Attestation{Data: &zondpb.AttestationData{
		Target: &zondpb.Checkpoint{Epoch: 0}},
		AggregationBits: bf}
	rt := [32]byte{'A'}
	att1.Data.Target.Root = rt[:]
	att1.Data.BeaconBlockRoot = rt[:]
	br := beaconState.BlockRoots()
	newRt := [32]byte{'B'}
	br[0] = newRt[:]
	require.NoError(t, beaconState.SetBlockRoots(br))
	att2.Data.Target.Root = newRt[:]
	att2.Data.BeaconBlockRoot = newRt[:]
	err := beaconState.AppendPreviousEpochAttestations(&zondpb.PendingAttestation{Data: att1.Data, AggregationBits: bf, InclusionDelay: 1})
	require.NoError(t, err)
	err = beaconState.AppendCurrentEpochAttestations(&zondpb.PendingAttestation{Data: att2.Data, AggregationBits: bf, InclusionDelay: 1})
	require.NoError(t, err)

	pVals := make([]*precompute.Validator, validators)
	for i := 0; i < len(pVals); i++ {
		pVals[i] = &precompute.Validator{CurrentEpochEffectiveBalance: 100}
	}
	pVals, _, err = precompute.ProcessAttestations(context.Background(), beaconState, pVals, &precompute.Balance{})
	require.NoError(t, err)

	committee, err := helpers.BeaconCommitteeFromState(context.Background(), beaconState, att1.Data.Slot, att1.Data.CommitteeIndex)
	require.NoError(t, err)
	indices, err := attestation.AttestingIndices(att1.AggregationBits, committee)
	require.NoError(t, err)
	for _, i := range indices {
		if !pVals[i].IsPrevEpochAttester {
			t.Error("Not a prev epoch attester")
		}
	}
	committee, err = helpers.BeaconCommitteeFromState(context.Background(), beaconState, att2.Data.Slot, att2.Data.CommitteeIndex)
	require.NoError(t, err)
	indices, err = attestation.AttestingIndices(att2.AggregationBits, committee)
	require.NoError(t, err)
	for _, i := range indices {
		assert.Equal(t, true, pVals[i].IsPrevEpochAttester, "Not a prev epoch attester")
		assert.Equal(t, true, pVals[i].IsPrevEpochTargetAttester, "Not a prev epoch target attester")
		assert.Equal(t, true, pVals[i].IsPrevEpochHeadAttester, "Not a prev epoch head attester")
	}
}

func TestEnsureBalancesLowerBound(t *testing.T) {
	b := &precompute.Balance{}
	b = precompute.EnsureBalancesLowerBound(b)
	balanceIncrement := params.BeaconConfig().EffectiveBalanceIncrement
	assert.Equal(t, balanceIncrement, b.ActiveCurrentEpoch, "Did not get wanted active current balance")
	assert.Equal(t, balanceIncrement, b.ActivePrevEpoch, "Did not get wanted active previous balance")
	assert.Equal(t, balanceIncrement, b.CurrentEpochAttested, "Did not get wanted current attested balance")
	assert.Equal(t, balanceIncrement, b.CurrentEpochTargetAttested, "Did not get wanted target attested balance")
	assert.Equal(t, balanceIncrement, b.PrevEpochAttested, "Did not get wanted prev attested balance")
	assert.Equal(t, balanceIncrement, b.PrevEpochTargetAttested, "Did not get wanted prev target attested balance")
	assert.Equal(t, balanceIncrement, b.PrevEpochHeadAttested, "Did not get wanted prev head attested balance")
}
