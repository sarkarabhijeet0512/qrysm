package beacon

import (
	"github.com/theQRL/qrysm/v4/beacon-chain/rpc/eth/shared"
)

type BlockRootResponse struct {
	Data *struct {
		Root string `json:"root"`
	} `json:"data"`
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
}

type GetCommitteesResponse struct {
	Data                []*shared.Committee `json:"data"`
	ExecutionOptimistic bool                `json:"execution_optimistic"`
	Finalized           bool                `json:"finalized"`
}

type DepositContractResponse struct {
	Data *struct {
		ChainId string `json:"chain_id"`
		Address string `json:"address"`
	} `json:"data"`
}

type ListAttestationsResponse struct {
	Data []*shared.Attestation `json:"data"`
}

type SubmitAttestationsRequest struct {
	Data []*shared.Attestation `json:"data"`
}

type ListVoluntaryExitsResponse struct {
	Data []*shared.SignedVoluntaryExit `json:"data"`
}

type SubmitSyncCommitteeSignaturesRequest struct {
	Data []*shared.SyncCommitteeMessage `json:"data"`
}

type GetStateForkResponse struct {
	Data                *shared.Fork `json:"data"`
	ExecutionOptimistic bool         `json:"execution_optimistic"`
	Finalized           bool         `json:"finalized"`
}

type GetFinalityCheckpointsResponse struct {
	ExecutionOptimistic bool                 `json:"execution_optimistic"`
	Finalized           bool                 `json:"finalized"`
	Data                *FinalityCheckpoints `json:"data"`
}

type FinalityCheckpoints struct {
	PreviousJustified *shared.Checkpoint `json:"previous_justified"`
	CurrentJustified  *shared.Checkpoint `json:"current_justified"`
	Finalized         *shared.Checkpoint `json:"finalized"`
}

type GetGenesisResponse struct {
	Data *Genesis `json:"data"`
}

type Genesis struct {
	GenesisTime           string `json:"genesis_time"`
	GenesisValidatorsRoot string `json:"genesis_validators_root"`
	GenesisForkVersion    string `json:"genesis_fork_version"`
}

type GetBlockHeadersResponse struct {
	Data                []*shared.SignedBeaconBlockHeaderContainer `json:"data"`
	ExecutionOptimistic bool                                       `json:"execution_optimistic"`
	Finalized           bool                                       `json:"finalized"`
}

type GetBlockHeaderResponse struct {
	ExecutionOptimistic bool                                     `json:"execution_optimistic"`
	Finalized           bool                                     `json:"finalized"`
	Data                *shared.SignedBeaconBlockHeaderContainer `json:"data"`
}

type GetValidatorsResponse struct {
	ExecutionOptimistic bool                  `json:"execution_optimistic"`
	Finalized           bool                  `json:"finalized"`
	Data                []*ValidatorContainer `json:"data"`
}

type GetValidatorResponse struct {
	ExecutionOptimistic bool                `json:"execution_optimistic"`
	Finalized           bool                `json:"finalized"`
	Data                *ValidatorContainer `json:"data"`
}

type GetValidatorBalancesResponse struct {
	ExecutionOptimistic bool                `json:"execution_optimistic"`
	Finalized           bool                `json:"finalized"`
	Data                []*ValidatorBalance `json:"data"`
}

type ValidatorContainer struct {
	Index     string     `json:"index"`
	Balance   string     `json:"balance"`
	Status    string     `json:"status"`
	Validator *Validator `json:"validator"`
}

type Validator struct {
	Pubkey                     string `json:"pubkey"`
	WithdrawalCredentials      string `json:"withdrawal_credentials"`
	EffectiveBalance           string `json:"effective_balance"`
	Slashed                    bool   `json:"slashed"`
	ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
	ActivationEpoch            string `json:"activation_epoch"`
	ExitEpoch                  string `json:"exit_epoch"`
	WithdrawableEpoch          string `json:"withdrawable_epoch"`
}

type ValidatorBalance struct {
	Index   string `json:"index"`
	Balance string `json:"balance"`
}
