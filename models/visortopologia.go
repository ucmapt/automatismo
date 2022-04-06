package models

import (
	//"github.com/ucmapt/automatismo/config"
)

// VisorTopologia estructura principal para manejo de registros de topología
type VisorTopologia struct {
	ID               string `json:"id"`
	NombreRed        string `json:"nombre_red"`
	Nombre           string `json:"nombre"`
	TipoElemento     string `json:"tipo_elemento"`
	Nodo1            string `json:"nodo1"`
	Nodo2            string `json:"nodo2"`
	Sw1              bool   `json:"sw1"`
	Sw2              bool   `json:"sw2"`
	Division         string `json:"division"`
	CircuitoFijo     string `json:"circuito_fijo"`
	CircuitoDinamico string `json:"circuito_dinamico"`
}

// METODO VisorTopologia.Tablename - para cumplir con GORM
func (a *VisorTopologia) TableName() string {
	return "visor.topologia"
}

// VisorGrafico estructura complementaria para manejo de la representación gráfica de cada elemento del modelo georreferenciado
type VisorGrafico struct {
	ID                   string `json:"id"`
	NombreRed            string `json:"nombre_red"`
	Nombre               string `json:"nombre"`
	Diagrama             string `json:"diagrama"`
	CoordenadasSrid      string `json:"coordenadas_srid"`
	CircuitoFijo         string `json:"circuito_fijo"`
	TipoElemento         string `json:"tipo_elemento"`
	Estado               bool   `json:"estado"`
	EsAereo              bool   `json:"es_aereo"`
	EsParticular         bool   `json:"es_particular"`
	TipoEstilo           string `json:"tipo_estilo"`
	CoordenadasFijasSrid string `json:"coordenadas_fijas_srid"`
	EstadoUcm            int    `json:"estado_ucm"`
}

// METODO VisorGrafico.Tablename - para cumplir con GORM
func (a *VisorGrafico) TableName() string {
	return "visor.grafico"
}

// VisorTransformador estructura complementaria para manejo de información de transformadores
type VisorTransformador struct {
	ID            string  `json:"id"`
	NombreRed     string  `json:"nombre_red"`
	Info          string  `json:"info"`
	Nombre        string  `json:"nombre"`
	Fases         string  `json:"fases"`
	Potencia      int     `json:"potencia"`
	PotenciaFaseA int     `json:"potencia_fase_a"`
	PotenciaFaseB int     `json:"potencia_fase_b"`
	PotenciaFaseC int     `json:"potencia_fase_c"`
	TipoConexion  string  `json:"tipo_conexion"`
	EsParticular  bool    `json:"es_particular"`
	EsAereo       bool    `json:"es_aereo"`
	P             float64 `json:"p"`
	Q             float64 `json:"q"`
	NombreUsuario string  `json:"nombre_usuario"`
	RpuMedidor    string  `json:"rpu_medidor"`
	Dato1         string  `json:"dato1"`
	Dato2         string  `json:"dato2"`
	Dato3         string  `json:"dato3"`
}

// METODO VisorTransformador.Tablename - para cumplir con GORM
func (a *VisorTransformador) TableName() string {
	return "visor.transformador"
}

// VisorFusible estructura complementaria para manejo de información de fusibles
type VisorFusible struct {
	ID              string      `json:"id"`
	NombreRed       string      `json:"nombre_red"`
	Info            string      `json:"info"`
	Nombre          string      `json:"nombre"`
	NumeroEconomico string      `json:"numero_economico"`
	Fases           string      `json:"fases"`
	Capacidad       int         `json:"capacidad"`
	Tipo            string      `json:"tipo"`
	Dato1           string      `json:"dato1"`
	Dato2           string      `json:"dato2"`
	Dato3           string      `json:"dato3"`
}
// METODO VisorFusible.Tablename - para cumplir con GORM
func (a *VisorFusible) TableName() string {
	return "visor.fusible"
}

// VisorLinea estructura complementaria para manejo de información de fusibles
type VisorLinea struct {
	ID         string  `json:"id"`
	NombreRed  string  `json:"nombre_red"`
	Info       string  `json:"info"`
	Nombre     string  `json:"nombre"`
	Longitud   float64 `json:"longitud"`
	EsAerea    bool    `json:"es_aerea"`
	IDCatalogo string  `json:"id_catalogo"`
	IDMaterial string  `json:"id_material"`
	Fases      string  `json:"fases"`
	EsEnlace   bool    `json:"es_enlace"`
	Dato1      string  `json:"dato1"`
	Dato2      string  `json:"dato2"`
	Dato3      string  `json:"dato3"`
}

// METODO VisorLinea.Tablename - para cumplir con GORM
func (a *VisorLinea) TableName() string {
	return "visor.linea"
}

// VisorAlimentador estructura complementaria para manejo de información de alimentadores
type VisorAlimentador struct {
	ID        string `json:"id"`
	NombreRed string `json:"nombre_red"`
	Info      string `json:"info"`
	Nombre    string `json:"nombre"`
	Dato1     string `json:"dato1"`
	Dato2     string `json:"dato2"`
	Dato3     string `json:"dato3"`
}

// METODO VisorAlimentador.Tablename - para cumplir con GORM
func (a *VisorAlimentador) TableName() string {
	return "visor.alimentador"
}

// VisorNodo estructura complementaria para manejo de información de nodos
type VisorNodo struct {
	ID        string  `json:"id"`
	NombreRed string  `json:"nombre_red"`
	Info      string  `json:"info"`
	Nombre    string  `json:"nombre"`
	Voltaje   float64 `json:"voltaje"`
	EsAereo   bool    `json:"es_aereo"`
	Dato1     string  `json:"dato1"`
	Dato2     string  `json:"dato2"`
	Dato3     string  `json:"dato3"`
}

// METODO VisorNodo.Tablename - para cumplir con GORM
func (a *VisorNodo) TableName() string {
	return "visor.nodo"
}

// VisorRestaurador estructura complementaria para manejo de información de restauradores
type VisorRestaurador []struct {
	ID               string `json:"id"`
	NombreRed        string `json:"nombre_red"`
	Info             string `json:"info"`
	Nombre           string `json:"nombre"`
	Dato1            string `json:"dato1"`
	Dato2            string `json:"dato2"`
	Dato3            string `json:"dato3"`
	EsTelecontrolado bool   `json:"es_telecontrolado"`
}

// METODO VisorRestaurador.Tablename - para cumplir con GORM
func (a *VisorRestaurador) TableName() string {
	return "visor.restaurador"
}

// VisorEquipoSw estructura complementaria para manejo de información de equipos de seccionamiento
type VisorEquipoSw struct {
	ID               string  `json:"id"`
	NombreRed        string  `json:"nombre_red"`
	Info             string  `json:"info"`
	Nombre           string  `json:"nombre"`
	Fases            string  `json:"fases"`
	Ir               float64 `json:"ir"`
	Tipo             string  `json:"tipo"`
	Dato1            string  `json:"dato1"`
	Dato2            string  `json:"dato2"`
	Dato3            string  `json:"dato3"`
	EsTelecontrolado bool    `json:"es_telecontrolado"`
}

// METODO VisorEquipoSw.Tablename - para cumplir con GORM
func (a *VisorEquipoSw) TableName() string {
	return "visor.equipo_seccionamiento"
}

// VisorCapacitor estructura complementaria para manejo de información de capacitores
type VisorCapacitor struct {
	ID           string  `json:"id"`
	NombreRed    string  `json:"nombre_red"`
	Info         string  `json:"info"`
	Nombre       string  `json:"nombre"`
	P1           int     `json:"p1"`
	Q1           float64 `json:"q1"`
	P0           int     `json:"p0"`
	Q0           int     `json:"q0"`
	Voltaje      float64 `json:"voltaje"`
	EsFijo       bool    `json:"es_fijo"`
	Fases        string  `json:"fases"`
	TipoConexion string  `json:"tipo_conexion"`
	Dato1        string  `json:"dato1"`
	Dato2        string  `json:"dato2"`
	Dato3        string  `json:"dato3"`
}

// METODO VisorCapacitor.Tablename - para cumplir con GORM
func (a *VisorCapacitor) TableName() string {
	return "visor.capacitor"
}
