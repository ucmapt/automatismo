package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Ambiente struct {
}

// Rutina auxiliar para extraer el parámetro p del objeto de recolecta la configuración u
func extraeParametro(u UcmAptConfig, p string) (string, error) {
	par , err  := u.ExtraeParam(p)
	if err != nil {
		return "", err
	}
	return par.Valor, nil
}

// Obtener connectividad con la BD de datos georreferenciados
func (a Ambiente) GetGeoDb(u UcmAptConfig) (*gorm.DB, error) {
	user, err := extraeParametro(u, "geodbuser")
	if err != nil {
		return nil, fmt.Errorf("Error al extraer parámetros")
	}
	pass, err := extraeParametro(u, "geodbpwd")
	if err != nil {
		return nil, fmt.Errorf("Error al extraer parámetros")
	}
	port, err := extraeParametro(u, "geodbport")
	if err != nil {
		return nil, fmt.Errorf("Error al extraer parámetros")
	}
	host, err := extraeParametro(u, "geodbsvr")
	if err != nil {
		return nil, fmt.Errorf("Error al extraer parámetros")
	}
	dbName := "UCM_Views"

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbName, pass)

	db, err := gorm.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Obtener conectividad con la BD principal de la UCM-CFE
func (a Ambiente) GetUcmDb(u UcmAptConfig) (*gorm.DB, error) {
	user, err := extraeParametro(u, "ucmdbuser")
	if err != nil {
		return nil, fmt.Errorf("Error al extraer parámetro: %s", "Usuario - UCM")
	}
	pass, err := extraeParametro(u, "ucmdbpwd")
	if err != nil {
		return nil, fmt.Errorf("Error al extraer parámetro: %s", "Contraseña- UCM")
	}
	port, err := extraeParametro(u, "ucmdbport")
	if err != nil {
		return nil, fmt.Errorf("Error al extraer parámetro: %s", "Puerto - UCM")
	}
	host, err := extraeParametro(u, "ucmdbsvr")
	if err != nil {
		return nil, fmt.Errorf("Error al extraer parámetro: %s", "Servidor - UCM")
	}
	dbName := "p_ucm_cfe"

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbName, pass)

	db, err := gorm.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (a Ambiente) GetAptDb(u UcmAptConfig) (*gorm.DB, error) {
	user, err := extraeParametro(u, "aptdbuser")
	if err != nil {
		return nil, fmt.Errorf("Error al extraer parámetro: %s", "Usuario - UCM")
	}
	pass, err := extraeParametro(u, "aptdbpwd")
	if err != nil {
		return nil, fmt.Errorf("Error al extraer parámetro: %s", "Contraseña- UCM")
	}
	port, err := extraeParametro(u, "aptdbport")
	if err != nil {
		return nil, fmt.Errorf("Error al extraer parámetro: %s", "Puerto - UCM")
	}
	host, err := extraeParametro(u, "aptdbsvr")
	if err != nil {
		return nil, fmt.Errorf("Error al extraer parámetro: %s", "Servidor - UCM")
	}
	dbName := "p_ucm_cfe"

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbName, pass)

	db, err := gorm.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return db, nil
}
