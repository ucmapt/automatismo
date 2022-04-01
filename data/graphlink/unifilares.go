package graphlink

import "time"

type GraphLinkTopology struct {
	ID                 string      `json:"_id"`
	Nombre             string      `json:"nombre"`
	Clave              string      `json:"clave"`
	Descripcion        string      `json:"descripcion"`
	Padre              string      `json:"padre"`
	Siguiente          string      `json:"siguiente"`
	Anterior           string      `json:"anterior"`
	Modelo             Modelo      `json:"modelo"`
	Capas              []Capas     `json:"capas"`
	Aor                string      `json:"aor"`
	Estado             string      `json:"estado"`
	Propiedades        Propiedades `json:"propiedades"`
	Version            string      `json:"version"`
	FechaCreacion      time.Time   `json:"fecha_creacion"`
	FechaActualizacion time.Time   `json:"fecha_actualizacion"`
}
type LinkDataArray struct {
	EsLink    bool        `json:"esLink"`
	From      string      `json:"from"`
	FromPort  string      `json:"fromPort"`
	ID        interface{} `json:"id"`
	IDVista   string      `json:"idVista"`
	Key       int         `json:"key"`
	Layer     string      `json:"layer"`
	Points    []int       `json:"points"`
	TipoLinea string      `json:"tipoLinea"`
	To        string      `json:"to"`
	ToPort    string      `json:"toPort"`
}
type Variable struct {
	Clave string `json:"clave"`
	ID    string `json:"id"`
}
type Dei struct {
	Dispositivo string   `json:"dispositivo"`
	Equipo      string   `json:"equipo"`
	Valor       int      `json:"valor"`
	Variable    Variable `json:"variable"`
}
type VerticalAlignment struct {
	OffsetX int     `json:"offsetX"`
	OffsetY int     `json:"offsetY"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
}
type NodeDataArray struct {
	Angle             int               `json:"angle"`
	Category          string            `json:"category"`
	ID                string            `json:"id,omitempty"`
	IDVista           string            `json:"idVista,omitempty"`
	Key               string            `json:"key"`
	Layer             string            `json:"layer"`
	Loc               string            `json:"loc"`
	MainCategory      string            `json:"mainCategory,omitempty"`
	Nombre            string            `json:"nombre,omitempty"`
	ZOrder            int               `json:"zOrder"`
	Dei               Dei               `json:"dei,omitempty"`
	Background        string            `json:"background,omitempty"`
	EsTexto           bool              `json:"esTexto,omitempty"`
	Font              string            `json:"font,omitempty"`
	Size              string            `json:"size,omitempty"`
	Stroke            string            `json:"stroke,omitempty"`
	Text              string            `json:"text,omitempty"`
	Through           bool              `json:"through,omitempty"`
	Underline         bool              `json:"underline,omitempty"`
	VerticalAlignment VerticalAlignment `json:"verticalAlignment,omitempty"`
	EsConector        bool              `json:"esConector,omitempty"`
}
type Modelo struct {
	Class                  string          `json:"class"`
	CopiesArrayObjects     bool            `json:"copiesArrayObjects"`
	CopiesArrays           bool            `json:"copiesArrays"`
	LinkDataArray          []LinkDataArray `json:"linkDataArray"`
	LinkFromPortIDProperty string          `json:"linkFromPortIdProperty"`
	LinkKeyProperty        string          `json:"linkKeyProperty"`
	LinkToPortIDProperty   string          `json:"linkToPortIdProperty"`
	NodeDataArray          []NodeDataArray `json:"nodeDataArray"`
}
type Capas struct {
	ID        string `json:"_id"`
	Bloqueada bool   `json:"bloqueada"`
	Clave     string `json:"clave"`
	Editable  bool   `json:"editable"`
	Nombre    string `json:"nombre"`
	Orden     int    `json:"orden"`
	Visible   bool   `json:"visible"`
}
type Historial struct {
	Accion        string `json:"accion"`
	Aor           string `json:"aor"`
	Categoria     string `json:"categoria"`
	Fecha         string `json:"fecha"`
	IDNodo        int    `json:"id_Nodo"`
	Propiedad     string `json:"propiedad"`
	ValorAnterior string `json:"valorAnterior"`
	ValorNuevo    string `json:"valorNuevo"`
}
type Propiedades struct {
	EsActual   bool          `json:"esActual"`
	Grid       string        `json:"grid"`
	Historial  []Historial   `json:"historial"`
	Plantillas []interface{} `json:"plantillas"`
}