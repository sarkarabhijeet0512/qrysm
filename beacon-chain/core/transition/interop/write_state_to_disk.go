package interop

import (
	"fmt"
	"os"
	"path"

	"github.com/theQRL/qrysm/v4/config/features"
	"github.com/theQRL/qrysm/v4/io/file"
	"github.com/theQRL/qrysm/v4/beacon-chain/state"
)

// WriteStateToDisk as a state ssz. Writes to temp directory. Debug!
func WriteStateToDisk(state state.ReadOnlyBeaconState) {
	if !features.Get().WriteSSZStateTransitions {
		return
	}
	fp := path.Join(os.TempDir(), fmt.Sprintf("beacon_state_%d.ssz", state.Slot()))
	log.Warnf("Writing state to disk at %s", fp)
	enc, err := state.MarshalSSZ()
	if err != nil {
		log.WithError(err).Error("Failed to ssz encode state")
		return
	}
	if err := file.WriteFile(fp, enc); err != nil {
		log.WithError(err).Error("Failed to write to disk")
	}
}
