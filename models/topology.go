package models

import (
	"errors"
	"fmt"
	//"github.com/ucmapt/automatismo/config"
)

// SELECT id, idalterno, netname, "name", "type", node1, node2, sw1, sw2, tipo_conexion, zona, subestacion, circuito, circuito_fijo, division, circuito_dinamico
// FROM postgis.topology;

type ShortTopology struct {
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

type OldTopology struct {
	Id               string  `json:"id"`
	Idalterno        int     `json:"idalterno"`
	Netname          *string `json:"netname"`
	Name             *string `json:"name"`
	Type             *string `json:"type"`
	Node1            *string `json:"node1"`
	Node2            *string `json:"node2"`
	Sw1              bool    `json:"sw1"`
	Sw2              bool    `json:"sw2"`
	TipoConexion     *string `json:"tipo_conexion"`
	Zona             *string `json:"zona"`
	Subestacion      *string `json:"subestacion"`
	Circuito         *string `json:"circuito"`
	CircuitoFijo     *string `json:"circuito_fijo"`
	Division         *string `json:"division"`
	CircuitoDinamico *string `json:"circuito_dinamico"`
}

type OldExtendedLine struct {
	ID               string  `json:"id"`
	Idalterno        int     `json:"idalterno"`
	Netname          *string `json:"netname"`
	Name             *string `json:"name"`
	Type             *string `json:"type"`
	Node1            *string `json:"node1"`
	Node2            *string `json:"node2"`
	Sw1              bool    `json:"sw1"`
	Sw2              bool    `json:"sw2"`
	TipoConexion     *string `json:"tipo_conexion"`
	Zona             *string `json:"zona"`
	Subestacion      *string `json:"subestacion"`
	Circuito         *string `json:"circuito"`
	CircuitoFijo     *string `json:"circuito_fijo"`
	Division         *string `json:"division"`
	CircuitoDinamico *string `json:"circuito_dinamico"`
	R1               float64 `json:"r1"`
	X1               float64 `json:"x1"`
	C1               float64 `json:"c1"`
	R0               float64 `json:"r0"`
	X0               float64 `json:"x0"`
	C0               float64 `json:"c0"`
	Length           float64 `json:"length"`
}

func (a *OldTopology) TableName() string {
	return "postgis.topology"
}

func (a *OldTopology) Validate() error {

	if a.Netname == nil {
		return errors.New("el campo nombre_red no puede ser nulo")
	}

	if a.Name == nil {
		return errors.New("el campo nombre no puede ser nulo")
	}

	return nil
}

func (a *OldTopology) ValidateUpdate() error {

	err := a.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (a *OldTopology) SetDefaults() OldTopology {

	nowObject := OldTopology{
		Id:               a.Id,
		Idalterno:        a.Idalterno,
		Netname:          a.Netname,
		Name:             a.Name,
		Type:             a.Type,
		Node1:            a.Node1,
		Node2:            a.Node2,
		Sw1:              a.Sw1,
		Sw2:              a.Sw2,
		TipoConexion:     a.TipoConexion,
		Zona:             a.Zona,
		Subestacion:      a.Subestacion,
		Circuito:         a.Circuito,
		CircuitoFijo:     a.CircuitoFijo,
		Division:         a.Division,
		CircuitoDinamico: a.CircuitoDinamico,
	}

	return nowObject
}

func (a *OldTopology) FeederString() string {
	mark := " "
	if a.Sw1 == true {
		mark = "✓"
	}
	return fmt.Sprintf("Circuito[%s][%s], Fijo[%s], Dinámico[%s], Nodo[%s]", *a.Netname, mark, *a.CircuitoFijo, *a.CircuitoDinamico, *a.Node1)
}

func (a *OldTopology) LineString() string {
	mark1 := " "
	mark2 := " "
	if a.Sw1 == true {
		mark1 = "✓"
	}
	if a.Sw2 == true {
		mark2 = "✓"
	}
	return fmt.Sprintf("Linea[%s][%s-%s], Tipo[%s], Fijo[%s], Dinámico[%s], Nodos[%s - %s]", *a.Name, mark1, mark2, *a.Type, *a.CircuitoFijo, *a.CircuitoDinamico, *a.Node1, *a.Node2)
}

func (a *OldTopology) NodeString() string {
	return fmt.Sprintf("Nodo[%s], Fijo[%s], Dinámico[%s]", *a.Name, *a.CircuitoFijo, *a.CircuitoDinamico)
}

type SwLine struct {
	Switch           string  `json:"switch"`
	Sd1              bool    `json:"sd1"`
	Sd2              bool    `json:"sd2"`
	Id               string  `json:"id"`
	Idalterno        int     `json:"idalterno"`
	Netname          *string `json:"netname"`
	Name             *string `json:"name"`
	Type             *string `json:"type"`
	Node1            *string `json:"node1"`
	Node2            *string `json:"node2"`
	Sw1              bool    `json:"sw1"`
	Sw2              bool    `json:"sw2"`
	TipoConexion     *string `json:"tipo_conexion"`
	Zona             *string `json:"zona"`
	Subestacion      *string `json:"subestacion"`
	Circuito         *string `json:"circuito"`
	CircuitoFijo     *string `json:"circuito_fijo"`
	Division         *string `json:"division"`
	CircuitoDinamico *string `json:"circuito_dinamico"`
}

func (a *SwLine) SetDefaults() SwLine {

	nowObject := SwLine{
		Switch:           a.Switch,
		Sd1:              a.Sw1,
		Sd2:              a.Sw2,
		Id:               a.Id,
		Idalterno:        a.Idalterno,
		Netname:          a.Netname,
		Name:             a.Name,
		Type:             a.Type,
		Node1:            a.Node1,
		Node2:            a.Node2,
		Sw1:              a.Sw1,
		Sw2:              a.Sw2,
		TipoConexion:     a.TipoConexion,
		Zona:             a.Zona,
		Subestacion:      a.Subestacion,
		Circuito:         a.Circuito,
		CircuitoFijo:     a.CircuitoFijo,
		Division:         a.Division,
		CircuitoDinamico: a.CircuitoDinamico,
	}

	return nowObject
}

func (a *SwLine) SwitchString() string {
	mark1 := " "
	mark2 := " "
	mark3 := " "
	mark4 := " "
	if a.Sd1 == true {
		mark1 = "✓"
	}
	if a.Sd2 == true {
		mark2 = "✓"
	}
	if a.Sw1 == true {
		mark3 = "✓"
	}
	if a.Sw2 == true {
		mark4 = "✓"
	}
	return fmt.Sprintf("Equipo[%s][%s-%s], Linea[%s][%s-%s]", *&a.Switch, mark1, mark2, *a.Name, mark3, mark4)
}
