package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Ambiente struct {
}

func (a Ambiente) GetTopoDb() (*gorm.DB, error) {
	user := "ucmfs"
	pass := "ucmfs"
	port := 5432
	host := "10.0.113.224"
	dbName := "UCM_Views"

	uri := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbName, pass)

	db, err := gorm.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (a Ambiente) GetUcmDb() (*gorm.DB, error) {
	user := "ucmfs"
	pass := "ucmfs"
	port := 5432
	host := "10.0.113.224"
	dbName := "p_ucm_cfe"

	uri := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbName, pass)

	db, err := gorm.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (a Ambiente) GetAptDb() (*gorm.DB, error) {
	user := "ucmfs"
	pass := "ucmfs"
	port := 5432
	host := "10.0.113.224"
	dbName := "p_ucm_cfe"

	uri := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbName, pass)

	db, err := gorm.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return db, nil
}
