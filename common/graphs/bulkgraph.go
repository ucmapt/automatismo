package graphs

import (
	"github.com/ucmapt/automatismo/models"
)

type BulkGraph struct {
	Nodes   []*models.OldTopology
	Lines   []*models.OldExtendedLine
	Feeders []*models.OldTopology
	SwLines []*models.SwLine
}

func NewBulkGraph() *BulkGraph {
	return &BulkGraph{
		Nodes:   []*models.OldTopology{},
		Lines:   []*models.OldExtendedLine{},
		Feeders: []*models.OldTopology{},
		SwLines: []*models.SwLine{},
	}
}
