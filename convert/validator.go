package convert

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

var privValidators *PrivValidators
var priv_files []string

type PrivValidators struct {
	Validators []*types.PrivValidator
}

func (arr *PrivValidators) Len() int {
	return len(arr.Validators)
}

func (arr *PrivValidators) Less(i, j int) bool {
	a := hex.EncodeToString(arr.Validators[i].Address)
	b := hex.EncodeToString(arr.Validators[j].Address)
	return strings.Compare(a, b) < 1
}

func (arr *PrivValidators) Swap(i, j int) {
	arr.Validators[i], arr.Validators[j] = arr.Validators[j], arr.Validators[i]
}

func LoadPrivValidators(folder string) {
	err := filepath.Walk(folder, walkFunc)
	if err != nil {
		panic(err)
	}

	privValidators = &PrivValidators{
		Validators: []*types.PrivValidator{},
	}
	for _, path := range priv_files {
		priv := NewPrivValidator(path)
		// add other priv_validators
		privValidators.Validators = append(privValidators.Validators, priv)
	}

	// sort priv_validators
	sort.Sort(privValidators)
}

func walkFunc(path string, f os.FileInfo, err error) error {
	if !cmn.FileExists(path) || f.IsDir() {
		return nil
	}

	fmt.Println("priv: ", path)
	if strings.HasPrefix(filepath.Base(path), "priv_validator") {
		priv_files = append(priv_files, path)
	}

	return nil
}
