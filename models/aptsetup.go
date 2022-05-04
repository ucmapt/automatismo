package models

import (
)

type AptSetUp struct {
	Id                                 string   `json:"id" gorm:"type:uuid;primary_key"`
	EstadoGeneral                      *bool    `json:"estado_general"`
	ToleranciaPerimetro                *float32 `json:"tolerancia_perimetro"`
	LimiteCorrientes                   *float32 `json:"limite_corrientes"`
	TiempoPostLicencia                 *int32   `json:"tiempo_post_licencia"`
	DistinguirSeccionamientoAutomatico *bool    `json:"distinguir_seccionamiento_automatico"`
	SugerenciasOperacion               *bool    `json:"sugerencias_operacion"`
}

func (c *AptSetUp) TableName() string {
	return "aut_configuracion_general"
}

func (c *AptSetUp) Validate() error {

	return nil
}

func (c *AptSetUp) ValidateUpdate() error {

	err := c.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (c *AptSetUp) SetDefaults() AptSetUp {

	nowObject := AptSetUp{
		Id:                                 c.Id,
		EstadoGeneral:                      c.EstadoGeneral,
		ToleranciaPerimetro:                c.ToleranciaPerimetro,
		LimiteCorrientes:                   c.LimiteCorrientes,
		TiempoPostLicencia:                 c.TiempoPostLicencia,
		DistinguirSeccionamientoAutomatico: c.DistinguirSeccionamientoAutomatico,
		SugerenciasOperacion:               c.SugerenciasOperacion,
	}

	return nowObject
}
