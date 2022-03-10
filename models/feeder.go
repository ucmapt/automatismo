package models

type Feeder struct{
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