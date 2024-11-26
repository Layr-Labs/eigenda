package v2

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
)

type BlobVersionParameterMap = common.ReadOnlyMap[BlobVersion, *core.BlobVersionParameters]

func NewBlobVersionParameterMap(params map[BlobVersion]*core.BlobVersionParameters) *BlobVersionParameterMap {
	return common.NewReadOnlyMap(params)
}
