package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hxzqlh/tm-tools/convert"
	cmn "github.com/tendermint/tmlibs/common"
)

var oldData = flag.String("old", os.ExpandEnv("$HOME/.tendermint"), "old tendermint dir")
var newData = flag.String("new", os.ExpandEnv("$HOME/.tendermint.new"), "new tendermint dir")
var privData = flag.String("priv", "", "other priv_validator.json configs dir")
var startHeight = flag.Int("s", 1, "start from height")

func main() {
	if len(os.Args) <= 6 {
		fmt.Printf("Usage: %s -old tmroot -new tmroot -priv priv_dir [-s startHeight]\n", os.Args[0])
		os.Exit(0)
	}

	flag.Parse()

	//pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	oTmRoot, _ := filepath.Abs(*oldData)
	nTmRoot, _ := filepath.Abs(*newData)
	privDir, _ := filepath.Abs(*privData)
	cmn.EnsureDir(nTmRoot, 0755)

	convert.OnStart(oTmRoot, nTmRoot)

	// gen config.toml
	convert.OnConfigToml(nTmRoot + "/config.toml")

	// genesis
	convert.OnGenesisJSON(oTmRoot+"/genesis.json", nTmRoot+"/genesis.json")

	// me priv
	convert.OnPrivValidatorJSON(oTmRoot+"/priv_validator.json", nTmRoot+"/priv_validator.json")

	// old priv validatores
	convert.LoadPrivValidators(privDir)

	convert.TotalHeight()

	convert.OnBlockStore(*startHeight)

	convert.OnStop()
}
