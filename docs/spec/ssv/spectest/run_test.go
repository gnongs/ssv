package spectest

import (
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/ssv/docs/spec/qbft"
	tests2 "github.com/bloxapp/ssv/docs/spec/ssv/spectest/tests"
	"github.com/bloxapp/ssv/docs/spec/types/testingutils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAll(t *testing.T) {
	for _, test := range AllTests {
		t.Run(test.Name, func(t *testing.T) {
			runTest(t, test)
		})
	}
}

func TestJson(t *testing.T) {
	basedir, _ := os.Getwd()
	path := filepath.Join(basedir, "generate")
	fileName := "tests.json"
	tests := map[string]*tests2.SpecTest{}
	byteValue, err := ioutil.ReadFile(path + "/" + fileName)
	require.NoError(t, err)

	if err := json.Unmarshal(byteValue, &tests); err != nil {
		require.NoError(t, err)
	}

	for _, test := range tests {

		// a little trick we do to instantiate all the internal controller params
		byts, err := test.Runner.QBFTController.Encode()
		require.NoError(t, err)
		newContr := qbft.NewController(
			[]byte{1, 2, 3, 4},
			testingutils.TestingShare,
			testingutils.TestingConfig.Domain,
			testingutils.TestingConfig.Signer,
			testingutils.TestingConfig.ValueCheck,
			testingutils.TestingConfig.Storage,
			testingutils.TestingConfig.Network,
		)
		require.NoError(t, newContr.Decode(byts))
		test.Runner.QBFTController = newContr

		for idx, i := range test.Runner.QBFTController.StoredInstances {
			if i == nil {
				continue
			}
			fixedInst := fixQBFTInstanceForRun(t, i)
			test.Runner.QBFTController.StoredInstances[idx] = fixedInst

			if test.Runner.State != nil &&
				test.Runner.State.RunningInstance != nil &&
				test.Runner.State.RunningInstance.GetHeight() == fixedInst.GetHeight() {
				test.Runner.State.RunningInstance = fixedInst
			}
		}
		t.Run(test.Name, func(t *testing.T) {
			runTest(t, test)
		})
	}
}

func runTest(t *testing.T, test *tests2.SpecTest) {
	v := testingutils.BaseValidator()
	v.DutyRunners[test.Runner.BeaconRoleType] = test.Runner

	lastErr := v.StartDuty(test.Duty)
	for _, msg := range test.Messages {
		err := v.ProcessMessage(msg)
		if err != nil {
			lastErr = err
		}
	}

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, lastErr, test.ExpectedError)
	} else {
		require.NoError(t, lastErr)
	}

	postRoot, err := test.Runner.State.GetRoot()
	require.NoError(t, err)

	require.EqualValues(t, test.PostDutyRunnerStateRoot, hex.EncodeToString(postRoot))
}

func fixQBFTInstanceForRun(t *testing.T, i *qbft.Instance) *qbft.Instance {
	// a little trick we do to instantiate all the internal instance params
	if i == nil {
		return nil
	}
	byts, _ := i.Encode()
	newInst := qbft.NewInstance(testingutils.TestingConfig, i.State.Share, i.State.ID)
	require.NoError(t, newInst.Decode(byts))
	return newInst
}
