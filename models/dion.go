package models

import (
	"fmt"
	"time"
)

type FallaFrancaPlain struct {
	Id            int
	Elemento      string
	EstampaTiempo time.Time
}

// Descartando timeout y otros 
func (f FallaFrancaPlain) DescartarOtros() error {
	var e error
	e = nil
	if time.Now().Before(f.EstampaTiempo.Add(2 * time.Minute)) {
		e = fmt.Errorf("Verificar elemento %s", f.Elemento)
	}
	return e
}
