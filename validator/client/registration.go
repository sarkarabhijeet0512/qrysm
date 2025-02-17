package client

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/theQRL/go-zond/common/hexutil"
	"github.com/theQRL/qrysm/v4/beacon-chain/builder"
	"github.com/theQRL/qrysm/v4/beacon-chain/core/signing"
	"github.com/theQRL/qrysm/v4/config/params"
	"github.com/theQRL/qrysm/v4/encoding/bytesutil"
	zondpb "github.com/theQRL/qrysm/v4/proto/prysm/v1alpha1"
	validatorpb "github.com/theQRL/qrysm/v4/proto/prysm/v1alpha1/validator-client"
	"github.com/theQRL/qrysm/v4/validator/client/iface"
	"go.opencensus.io/trace"
)

// SubmitValidatorRegistrations signs validator registration objects and submits it to the beacon node.
func SubmitValidatorRegistrations(
	ctx context.Context,
	validatorClient iface.ValidatorClient,
	signedRegs []*zondpb.SignedValidatorRegistrationV1,
) error {
	ctx, span := trace.StartSpan(ctx, "validator.SubmitValidatorRegistrations")
	defer span.End()

	if len(signedRegs) == 0 {
		return nil
	}

	if _, err := validatorClient.SubmitValidatorRegistrations(ctx, &zondpb.SignedValidatorRegistrationsV1{
		Messages: signedRegs,
	}); err != nil {
		if strings.Contains(err.Error(), builder.ErrNoBuilder.Error()) {
			log.Warnln("Beacon node does not utilize a custom builder via the --http-mev-relay flag. Validator registration skipped.")
			return nil
		}
		return errors.Wrap(err, "could not submit signed registrations to beacon node")
	}
	log.Infoln("Submitted builder validator registration settings for custom builders")
	return nil
}

// Sings validator registration obj with the proposer domain and private key.
func signValidatorRegistration(ctx context.Context, signer iface.SigningFunc, reg *zondpb.ValidatorRegistrationV1) ([]byte, error) {
	// Per spec, we want the fork version and genesis validator to be nil.
	// Which is genesis value and zero by default.
	d, err := signing.ComputeDomain(
		params.BeaconConfig().DomainApplicationBuilder,
		nil, /* fork version */
		nil /* genesis val root */)
	if err != nil {
		return nil, err
	}

	r, err := signing.ComputeSigningRoot(reg, d)
	if err != nil {
		return nil, errors.Wrap(err, signingRootErr)
	}

	sig, err := signer(ctx, &validatorpb.SignRequest{
		PublicKey:       reg.Pubkey,
		SigningRoot:     r[:],
		SignatureDomain: d,
		Object:          &validatorpb.SignRequest_Registration{Registration: reg},
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not sign validator registration")
	}
	return sig.Marshal(), nil
}

// SignValidatorRegistrationRequest compares and returns either the cached validator registration request or signs a new one.
func (v *validator) SignValidatorRegistrationRequest(ctx context.Context, signer iface.SigningFunc, newValidatorRegistration *zondpb.ValidatorRegistrationV1) (*zondpb.SignedValidatorRegistrationV1, error) {
	signedReg, ok := v.signedValidatorRegistrations[bytesutil.ToBytes2592(newValidatorRegistration.Pubkey)]
	if ok && isValidatorRegistrationSame(signedReg.Message, newValidatorRegistration) {
		return signedReg, nil
	} else {
		sig, err := signValidatorRegistration(ctx, signer, newValidatorRegistration)
		if err != nil {
			return nil, err
		}
		newRequest := &zondpb.SignedValidatorRegistrationV1{
			Message:   newValidatorRegistration,
			Signature: sig,
		}
		v.signedValidatorRegistrations[bytesutil.ToBytes2592(newValidatorRegistration.Pubkey)] = newRequest
		return newRequest, nil
	}
}

func isValidatorRegistrationSame(cachedVR *zondpb.ValidatorRegistrationV1, newVR *zondpb.ValidatorRegistrationV1) bool {
	isSame := true
	if cachedVR.GasLimit != newVR.GasLimit {
		isSame = false
	}
	if hexutil.Encode(cachedVR.FeeRecipient) != hexutil.Encode(newVR.FeeRecipient) {
		isSame = false
	}
	return isSame
}
