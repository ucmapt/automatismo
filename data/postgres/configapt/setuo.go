package configapt

import (
	"github.com/jinzhu/gorm"
	"github.com/ucmapt/automatismo/models"
)

type AptSetUpRepo struct {
	db *gorm.DB
}

func NewAptSetUpRepo(db *gorm.DB) *AptSetUpRepo {
	return &AptSetUpRepo{db: db}
}

func (e *AptSetUpRepo) GetFirst() (*models.AptSetUp, error) {

	single := &models.AptSetUp{}

	err := e.db.Where(&models.UcmEventReg{}).First(single).Error

	return single, err
}
