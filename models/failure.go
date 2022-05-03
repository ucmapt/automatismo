package models

import (
	"time"
)

type FallaFranca struct {
	CveComponente string
	CveCircuito   string
	ICC           *float64
	CargaPerdida  *float64
	OperadorId    *string
	FailedAt      *time.Time
	ExposedSince  *time.Time
}

type FailZone struct {
	// Nodos
	// Lineas
}

type FailResult int