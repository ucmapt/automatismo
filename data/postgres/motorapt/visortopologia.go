package motorapt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/ucmapt/automatismo/models"
)

// Cominicación de datos para las estructuras del esquema visor: topologia, alimentador, modo, linea y equipo_seccionamiento para manejo topologico
type VisorTopologyRepo struct {
	db *gorm.DB
}

func NewVisorTopologyRepo(db *gorm.DB) *VisorTopologyRepo {
	return &VisorTopologyRepo{db: db}
}

func (e *VisorTopologyRepo) GetById(id string) (*models.VisorTopologia, error) {

	single := &models.VisorTopologia{}

	err := e.db.Where(&models.VisorTopologia{ID: id}).First(single).Error

	return single, err
}

func (e *VisorTopologyRepo) GetAll() ([]*models.VisorTopologia, error) {

	group := []*models.VisorTopologia{}

	err := e.db.Find(&group).Error

	return group, err
}

func (e *VisorTopologyRepo) GetByNombre(nombre string) (*models.VisorTopologia, error) {

	single := &models.VisorTopologia{}

	err := e.db.Where(&models.VisorTopologia{Nombre: nombre}).First(single).Error

	return single, err
}

func (e *VisorTopologyRepo) GetFeeders() ([]*models.VisorTopologia, error) {
	group := []*models.VisorTopologia{}

	err := e.db.Where("type = 'alimentador'").Find(&group).Error

	return group, err
}

func (e *VisorTopologyRepo) GetNodes() ([]*models.VisorTopologia, error) {
	group := []*models.VisorTopologia{}

	err := e.db.Where("type = 'nodo'").Find(&group).Error

	return group, err
}

func (e *VisorTopologyRepo) GetLines() ([]*models.VisorTopologia, error) {
	group := []*models.VisorTopologia{}

	err := e.db.Where("type IN ('linea', 'linea_asimetrica')").Find(&group).Error

	return group, err
}

func (e *VisorTopologyRepo) GetSwLines() ([]*models.SwLine, error) {
	group := []*models.SwLine{}

	sb := strings.Builder{}
	sb.WriteString("SELECT s.name as switch, s.sw1 as sd1, s.sw2 as sd2, t.* ")
	sb.WriteString("FROM postgis.topology t, ")
	sb.WriteString("(select * FROM postgis.topology where type = 'DISCSWITCH_2') s ")
	sb.WriteString("WHERE t.name = s.node1")

	err := e.db.Raw(sb.String()).Find(&group).Error

	return group, err
}


/*
SELECT t.id, t.nombre_red, t.nombre, t.tipo_elemento, t.nodo1, t.nodo2, 
       l.nombre, l.longitud, 
	   c.r1, c.x1, c.r0, c.x0
FROM visor. topologia t
LEFT JOIN visor.linea l ON t.id = l.id
LEFT JOIN public.catalogo_linea c ON l.id_catalogo = c.id
WHERE t.tipo_elemento in ('linea', 'linea_asimetrica') 
*/

func (e *VisorTopologyRepo) GetFullLines() ([]*models.OldExtendedLine, error) {
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

func (e *VisorTopologyRepo) UpdateViews(offNodes, offLines, failZone string, circuitoscve, circuitos []string) error {

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