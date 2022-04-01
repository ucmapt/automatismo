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

type BulkVisor struct {
	Nodes   []*models.VisorTopology
	Lines   []*models.OldExtendedLine
	Feeders []*models.OldTopology
	SwLines []*models.SwLine
}

func NewBulkVisor() *BulkVisor {
	return &BulkVisor{
		Nodes:   []*models.VisorTopology{},
		Lines:   []*models.OldExtendedLine{},
		Feeders: []*models.OldTopology{},
		SwLines: []*models.SwLine{},
	}
}
