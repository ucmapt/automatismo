package models

import (
	//"github.com/ucmapt/automatismo/config"
)

// SELECT id, idalterno, netname, "name", "type", node1, node2, sw1, sw2, tipo_conexion, zona, subestacion, circuito, circuito_fijo, division, circuito_dinamico
// FROM postgis.topology;

type VisorTopology struct {
	Id               string  `json:"id"`
	NombreRed        *string `json:"nombre_red"`
	Nombre           *string `json:"nombre"`
	TipoElemento     *string `json:"tipo_elemento"`
	Nodo1            *string `json:"nodo1"`
	Nodo2            *string `json:"nodo2"`
	Sw1              bool    `json:"sw1"`
	Sw2              bool    `json:"sw2"`
	Division         *string `json:"division"`
	CircuitoFijo     *string `json:"circuito_fijo"`
	CircuitoDinamico *string `json:"circuito_dinamico"`
}

func (a *VisorTopology) TableName() string {
	return "visor.topology"
}

func (a *VisorTopology) SetDefaults() VisorTopology {
	nowObject := VisorTopology{
		Id:               a.Id,
		NombreRed:        a.NombreRed,
		Nombre:           a.Nombre,
		TipoElemento:     a.TipoElemento,
		Nodo1:            a.Nodo1,
		Nodo2:            a.Nodo2,
		Sw1:              a.Sw1,
		Sw2:              a.Sw2,
		Division:         a.Division,
		CircuitoFijo:     a.CircuitoFijo,
		CircuitoDinamico: a.CircuitoDinamico,
	}

	return nowObject
}


