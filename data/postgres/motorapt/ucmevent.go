package motorapt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/ucmapt/automatismo/models"
)

type UcmEventRegRepo struct {
	db *gorm.DB
}

func NewUcmEventRegRepo(db *gorm.DB) *UcmEventRegRepo {
	return &UcmEventRegRepo{db: db}
}

func (e *UcmEventRegRepo) GetById(id string) (*models.UcmEventReg, error) {

	single := &models.UcmEventReg{}

	err := e.db.Where(&models.UcmEventReg{Id: id}).First(single).Error

	return single, err
}


func (e *UcmEventRegRepo) InsertEvent(ev models.UcmEventReg)  error {
	sb := strings.Builder{}
	var err error

	sb.WriteString("INSERT INTO apt_process.apt_eventos (id, momento, detalle, descripcion, origen, tipo) ")
	sb.WriteString("VALUES(gen_random_uuid(), now(), ")
	sb.WriteString(fmt.Sprintf("'%s', '%s', '%s', '%s');", *ev.Detalle, *ev.Descripcion, *ev.Origen,*ev.Tipo))
	err = e.db.Exec(sb.String()).Error
	if err != nil {
		return errors.New(fmt.Sprintf("Problema limpiando mensajes [QRY] %s", sb.String()))
	}
	return err
}

func (e *UcmEventRegRepo) InsertEventFromTexts(descripcion string, detalle string, origen string, tipo string)  error {
	sb := strings.Builder{}
	var err error

	sb.WriteString("INSERT INTO apt_process.apt_eventos (id, momento, detalle, descripcion, origen, tipo) ")
	sb.WriteString("VALUES(gen_random_uuid(), now(), ")
	sb.WriteString(fmt.Sprintf("'%s', '%s', '%s', '%s');", detalle, descripcion, origen, tipo))
	err = e.db.Exec(sb.String()).Error
	if err != nil {
		return errors.New(fmt.Sprintf("Problema insertando eventos [QRY] %s", sb.String()))
	}
	return err
}
