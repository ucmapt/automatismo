package motorapt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/ucmapt/automatismo/models"
)

type OldTopologyRepo struct {
	db *gorm.DB
}

func NewTopologyRepo(db *gorm.DB) *OldTopologyRepo {
	return &OldTopologyRepo{db: db}
}

func (e *OldTopologyRepo) GetById(id string) (*models.OldTopology, error) {

	single := &models.OldTopology{}

	err := e.db.Where(&models.OldTopology{Id: id}).First(single).Error

	return single, err
}

func (e *OldTopologyRepo) GetAll() ([]*models.OldTopology, error) {

	group := []*models.OldTopology{}

	err := e.db.Find(&group).Error

	return group, err
}

func (e *OldTopologyRepo) GetByNombre(nombre string) (*models.OldTopology, error) {

	single := &models.OldTopology{}

	err := e.db.Where(&models.OldTopology{Name: &nombre}).First(single).Error

	return single, err
}

func (e *OldTopologyRepo) Create(single *models.OldTopology) (*models.OldTopology, error) {

	err := e.db.Create(&single).Error

	return single, err
}

func (e *OldTopologyRepo) Update(single *models.OldTopology) (*models.OldTopology, error) {

	err := e.db.Save(&single).Error

	return single, err
}

func (e *OldTopologyRepo) DeleteById(id string) (int64, error) {
	single := &models.OldTopology{}

	db := e.db.Where("id = ?", id).Delete(single)
	rows := db.RowsAffected
	err := db.Error

	return rows, err
}

func (e *OldTopologyRepo) GetFeeders() ([]*models.OldTopology, error) {
	group := []*models.OldTopology{}

	err := e.db.Where("type = 'FEEDER'").Find(&group).Error

	return group, err
}

func (e *OldTopologyRepo) GetNodes() ([]*models.OldTopology, error) {
	group := []*models.OldTopology{}

	err := e.db.Where("type = 'BUSBAR-NODE'").Find(&group).Error

	return group, err
}

func (e *OldTopologyRepo) GetLines() ([]*models.OldTopology, error) {
	group := []*models.OldTopology{}

	err := e.db.Where("type IN ('LINE', 'ASY_LINE')").Find(&group).Error

	return group, err
}

func (e *OldTopologyRepo) GetSwLines() ([]*models.SwLine, error) {
	group := []*models.SwLine{}

	sb := strings.Builder{}
	sb.WriteString("SELECT s.name as switch, s.sw1 as sd1, s.sw2 as sd2, t.* ")
	sb.WriteString("FROM postgis.topology t, ")
	sb.WriteString("(select * FROM postgis.topology where type = 'DISCSWITCH_2') s ")
	sb.WriteString("WHERE t.name = s.node1")

	err := e.db.Raw(sb.String()).Find(&group).Error

	return group, err
}


func (e *OldTopologyRepo) GetFullLines() ([]*models.OldExtendedLine, error) {
	group := []*models.OldExtendedLine{}

	sb := strings.Builder{}
	sb.WriteString("SELECT t.id, t.idalterno, t.netname, t.name, t.type, t.node1, t.node2, t.sw1, t.sw2, ")
	sb.WriteString("t.tipo_conexion, t.zona, t.subestacion, t.circuito, t.circuito_fijo, t.division, t.circuito_dinamico, ")
	sb.WriteString("l.r1, l.x1, l.c1, l.r0, l.x0, l.c0, l.length ")
    sb.WriteString("FROM postgis.topology t inner join postgis.line l on t.id = l.id ")
    sb.WriteString("WHERE t.type in ('LINE', 'ASY_LINE')")

	err := e.db.Raw(sb.String()).Find(&group).Error

	return group, err
}

func (e *OldTopologyRepo) UpdateViews(offNodes, offLines, failZone string, circuitoscve, circuitos []string) error {

	sb := strings.Builder{}
	var err error

	sb.WriteString("SELECT * FROM ")
	sb.WriteString(".truncate_tables();")
	err = e.db.Exec(sb.String()).Error
	if err != nil {
		return errors.New(fmt.Sprintf("Problema limpiando mensajes %s", sb.String()))
	}
	sb.Reset()
	return err
}


