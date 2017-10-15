package convert

import (
	"fmt"

	"github.com/hxzqlh/tm-tools/util"
)

func TotalHeight() {
	blockStore := util.LoadOldBlockStoreStateJSON(oBlockDb)
	totalHeight = blockStore.Height

	fmt.Println("total height", totalHeight)
}
