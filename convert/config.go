package convert

import (
	"io/ioutil"

	"github.com/hxzqlh/tm-tools/old"
	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

func OnConfigToml(configFilePath string) {
	var configTmpl = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

proxy_app = "trade"
moniker = "anonymous"
node_laddr = "tcp://0.0.0.0:46656"
seeds = ""
fast_sync = true
db_backend = "leveldb"
log_level = "info"
rpc_laddr = "tcp://0.0.0.0:46657"
`
	cmn.WriteFile(configFilePath, []byte(configTmpl), 0644)
}

func OnGenesisJSON(oPath, nPath string) {
	jsonBytes, err := ioutil.ReadFile(oPath)
	if err != nil {
		panic(err)
	}

	oGen := old.GenesisDocFromJSON(jsonBytes)
	nGen := NewGenesisDoc(oGen)
	nGen.SaveAs(nPath)
}

func NewPrivValidator(oPath string) *types.PrivValidator {
	privVali := &types.PrivValidator{}
	old := old.LoadPrivValidator(oPath)
	privVali.Address = old.Address
	privVali.LastHeight = old.LastHeight
	privVali.LastRound = old.LastRound
	privVali.LastSignature = old.LastSignature
	privVali.LastSignBytes = old.LastSignBytes
	privVali.LastStep = old.LastStep
	privVali.PrivKey = old.PrivKey
	privVali.PubKey = old.PubKey
	return privVali
}

func OnPrivValidatorJSON(oPath, nPath string) {
	privVali := NewPrivValidator(oPath)
	privVali.SetFile(nPath)
	privVali.Save()

	// add me to priv_validators
	privValidators = &PrivValidators{
		Validators: []*types.PrivValidator{},
	}

	privValidators.Validators = append(privValidators.Validators, privVali)
}
