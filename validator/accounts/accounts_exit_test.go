package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	dilithium2 "github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/qrysm/v4/beacon-chain/rpc/apimiddleware"
	"github.com/theQRL/qrysm/v4/build/bazel"
	"github.com/theQRL/qrysm/v4/encoding/bytesutil"
	"github.com/theQRL/qrysm/v4/io/file"
	zond "github.com/theQRL/qrysm/v4/proto/prysm/v1alpha1"
	"github.com/theQRL/qrysm/v4/testing/assert"
	"github.com/theQRL/qrysm/v4/testing/require"
)

func TestDisplayExitInfo(t *testing.T) {
	logHook := test.NewGlobal()
	key := []byte("0x123456")
	displayExitInfo([][]byte{key}, []string{string(key)})
	assert.LogsContain(t, logHook, "https://beaconcha.in/validator/3078313233343536")
}

func TestDisplayExitInfo_NoKeys(t *testing.T) {
	logHook := test.NewGlobal()
	displayExitInfo([][]byte{}, []string{})
	assert.LogsContain(t, logHook, "No successful voluntary exits")
}

func TestPrepareAllKeys(t *testing.T) {
	key1 := bytesutil.ToBytes2592([]byte("key1"))
	key2 := bytesutil.ToBytes2592([]byte("key2"))
	raw, formatted := prepareAllKeys([][dilithium2.CryptoPublicKeyBytes]byte{key1, key2})
	require.Equal(t, 2, len(raw))
	require.Equal(t, 2, len(formatted))
	assert.DeepEqual(t, bytesutil.ToBytes2592([]byte{107, 101, 121, 49}), bytesutil.ToBytes2592(raw[0]))
	assert.DeepEqual(t, bytesutil.ToBytes2592([]byte{107, 101, 121, 50}), bytesutil.ToBytes2592(raw[1]))
	assert.Equal(t, "0x6b6579310000", formatted[0])
	assert.Equal(t, "0x6b6579320000", formatted[1])
}

func TestWriteSignedVoluntaryExitJSON(t *testing.T) {
	sve := &zond.SignedVoluntaryExit{
		Exit: &zond.VoluntaryExit{
			Epoch:          5,
			ValidatorIndex: 300,
		},
		Signature: []byte{0x01, 0x02},
	}

	output := path.Join(bazel.TestTmpDir(), "TestWriteSignedVoluntaryExitJSON")
	require.NoError(t, writeSignedVoluntaryExitJSON(context.Background(), sve, output))

	b, err := file.ReadFileAsBytes(path.Join(output, "validator-exit-300.json"))
	require.NoError(t, err)

	svej := &apimiddleware.SignedVoluntaryExitJson{}
	require.NoError(t, json.Unmarshal(b, svej))

	require.Equal(t, fmt.Sprintf("%d", sve.Exit.Epoch), svej.Exit.Epoch)
	require.Equal(t, fmt.Sprintf("%d", sve.Exit.ValidatorIndex), svej.Exit.ValidatorIndex)
	require.Equal(t, "0x0102", svej.Signature)
}
