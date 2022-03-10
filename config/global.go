package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// Interfaces es un listado de API disponibles
var Interfaces = []string{
	"CONFIGURACIONES",
	"DESPLIEGUE",
	"TOPOLOGIA",
	"HISTORIAL",
	"VALORES",
	"PUNTOS",
	"LICENCIAS",
	"LOGIN",
	"BITACORA",
}

// UcmAptParamConfig estructura para manejo de parámetros de configuración de APT
type UcmAptParamConfig struct {
	Clave string `json:"clave"`
	Valor string `json:"valor"`
}

// UcmAptAPIConfig estructura para manejo de parámetros de configuración de APT
type UcmAptAPIConfig struct {
	API   string `json:"api"`
	Descr string `json:"descr"`
	URL   string `json:"url"`
	Tipo  string `json:"tipo"`
}

// UcmAptConfig estructura para manejo de parámetros de configuración de APT
type UcmAptConfig struct {
	Ccd        string              `json:"ccd"`
	Version    string              `json:"version"`
	Parametros []UcmAptParamConfig `json:"parametros"`
	ListaAPI   []UcmAptAPIConfig   `json:"api_lista"`
}

// Cargar carga la configuración a partir del archivo
func (u *UcmAptConfig) Cargar(archivo string) bool {
	c := true
	jsonFile, err := os.Open(archivo)
	if err != nil {
		log.Fatalln(err)
		c = false
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, u)
	if err != nil {
		log.Fatalln(err)
		c = false
	}
	return c
}

// ExtraeError estructura para manejo propio de errores al extraer configuraciones
type ExtraeError struct {
	problema string
}

// Error cumple la interface error
func (e ExtraeError) Error() string {
	return e.problema
}

// ExtraeParam adquiere un parámetro específico por índice o clave
func (u UcmAptConfig) ExtraeParam(cual interface{}) (UcmAptParamConfig, error) {
	var par UcmAptParamConfig

	switch cual.(type) {
	case int:
		// Extraer por número
		if cual.(int) >= len(u.Parametros) {
			return par, errors.New("Fuera de rango")
		}
		par = u.Parametros[cual.(int)]

	default:
		// Extraer por clave
		for _, a := range u.Parametros {
			if a.Clave == cual.(string) {
				par = a
				break
			}
		}
	}

	return par, nil
}

// ExtraeAPI adquiere una API específica por índice o clave
func (u UcmAptConfig) ExtraeAPI(cual interface{}) UcmAptAPIConfig {
	var api UcmAptAPIConfig
	var apiServer string

	apiServer = ""
	svr, err := u.ExtraeParam("apisrv")

	if err == nil {
		apiServer = svr.Valor
	}

	switch cual.(type) {
	case int:
		// Extraer por número
		api = u.ListaAPI[cual.(int)]
		api.URL = apiServer + api.URL
	default:
		// Extraer por clave
		for _, a := range u.ListaAPI {
			if a.API == cual.(string) {
				api = a
				api.URL = apiServer + api.URL
				break
			}
		}
	}
	return api
}