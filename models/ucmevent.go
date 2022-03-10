package models

import (
	"errors"
)

type UcmEventReg struct {
	Id          string  `json:"id"`
	Momento     string  `json:"momento"`
	Detalle     *string `json:"detalle"`
	Descripcion *string `json:"descripcion"`
	Origen      *string `json:"origen"`
	Tipo        *string `json:"tipo"`
}

func (a *UcmEventReg) TableName() string {
	return "apt_process.apt_eventos"
}

func (a *UcmEventReg) Validate() error {

	if a.Descripcion == nil {
		return errors.New("el campo descripcion no puede ser nulo")
	}

	if a.Origen == nil {
		return errors.New("el campo origen no puede ser nulo")
	}

	if a.Tipo == nil {
		return errors.New("el campo tipo no puede ser nulo")
	}
	return nil
}

func (a *UcmEventReg) ValidateUpdate() error {

	err := a.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (a *UcmEventReg) SetDefaults() UcmEventReg {

	nowObject := UcmEventReg{
		Id:          a.Id,
		Momento:     a.Momento,
		Detalle:     a.Detalle,
		Descripcion: a.Descripcion,
		Origen:      a.Origen,
		Tipo:        a.Tipo,
	}

	return nowObject
}
