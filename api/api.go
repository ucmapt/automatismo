package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ucmapt/automatismo/config"
)

const maxUint = ^uint(0)
const minUint = 0
const maxInt = int(maxUint >> 1)
const minInt = -maxInt - 1

// MANEJO DE CONFIGURACIONES DE IHM

// IhmConfiguraciones estructura principal para manejo de configuracion de la IHM
type IhmConfiguraciones struct {
	Configuraciones []IhmConfiguracion `json:"configuraciones"`
}

// IhmConfiguracion estructura auxiliar para manejo de configuracion de la IHM
type IhmConfiguracion struct {
	ID             string         `json:"_id"`
	Tipo           string         `json:"tipo"`
	Usuario        string         `json:"usuario"`
	Data           IhmConfData    `json:"data"`
	Menu           NullableString `json:"menu"`
	MenuPlantillas NullableString `json:"menu_plantillas"`
	MenuVistas     NullableString `json:"menu_vistas"`
}

// NullableString tipo de dato para validar datos nulos
type NullableString string

// MarshalJSON implementa la serializaci�n de NullableString
func (c NullableString) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if len(string(c)) == 0 {
		buf.WriteString(`null`)
	} else {
		buf.WriteString(`"` + string(c) + `"`) // add double quation mark as json format required
	}
	return buf.Bytes(), nil
}

// UnmarshalJSON implementa la deserialización de NullableString
func (c *NullableString) UnmarshalJSON(in []byte) error {
	str := string(in)
	if str == `null` {
		*c = ""
		return nil
	}
	res := NullableString(str)
	if len(res) >= 2 {
		res = res[1 : len(res)-1] // remove the wrapped qutation
	}
	*c = res
	return nil
}

// NullableNodo tipo de dato para validar datos nulos
type NullableNodo int

// MarshalJSON implementa la serialización de NullableNodo
func (c NullableNodo) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	if len(strconv.Itoa(int(c))) == 0 {
		buf.WriteString(`null`)
	} else if int(c) == minInt {
		buf.WriteString(`- - -`)
	} else {
		buf.WriteString(strconv.Itoa(int(c)))
	}
	return buf.Bytes(), nil
}

// UnmarshalJSON implementa la deserialización de NullableNodo
func (c *NullableNodo) UnmarshalJSON(in []byte) error {
	str := string(in)
	nodo, err := strconv.Atoi(str)
	if err != nil {
		*c = NullableNodo(minInt) // N�mero negativo m�s peque�o
		return nil
	}
	*c = NullableNodo(nodo)
	return nil
}

// IhmConfData estructura auxiliar para manejo de configuración de la IHM
type IhmConfData struct {
	MenuDespliegues []IhmConfMenu     `json:"menu_despliegues,omitempty"`
	MenuPlantillas  []IhmConfMenu     `json:"menu_plantillas,omitempty"`
	MenuTopologias  []IhmConfMenu     `json:"menu_topologias,omitempty"`
	MenuVistas      []IhmConfMenu     `json:"menu_vistas,omitempty"`
	Shortcuts       []IhmConfShortcut `json:"shortcuts,omitempty"`
}

// IhmConfMenu estructura auxiliar para manejo de configuración de la IHM
type IhmConfMenu struct {
	ID       string         `json:"id"`
	ParentID NullableString `json:"parentid"`
	Text     string         `json:"text"`
	Value    string         `json:"value"`
}

// IhmConfShortcut estructura auxiliar para manejo de configuración de la IHM
type IhmConfShortcut struct {
	Clave      int               `json:"clave"`
	Comando    string            `json:"comando"`
	KeyBinding IhmConfKeyBinding `json:"keybinding"`
	Modulo     string            `json:"modulo"`
}

// IhmConfKeyBinding estructura auxiliar para manejo de configuración de la IHM
type IhmConfKeyBinding struct {
	AltKey   bool   `json:"altKey"`
	CtrlKey  bool   `json:"ctrlKey"`
	Key      string `json:"key"`
	KeyCode  int    `json:"keyCode"`
	ShiftKey bool   `json:"shiftKey"`
}

// CargaConfig carga datos de las configuraciones de IHM a partir de la configuración
func (a *IhmConfiguraciones) CargaConfig(conf config.UcmAptAPIConfig) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(conf.URL)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
	strBody := string(body)
	err = json.Unmarshal([]byte(strBody), a)
	if err != nil {
		log.Fatalln(err)
	}
}

// CargaArchivo carga datos de las configuraciones de IHM a partir de un archivo
func (a *IhmConfiguraciones) CargaArchivo(archivo string) {
	jsonFile, err := os.Open(archivo)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Archivo de configuración abierto sin complicacicones")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, a)
	if err != nil {
		fmt.Println(err)
	}
}

// MANEJO DE DESPLIEGUES

// IhmDespliegue estructura principal para manejo de vistas de despliegue
type IhmDespliegue struct {
	ID                 string          `json:"_id"`
	Nombre             string          `json:"nombre"`
	Clave              string          `json:"clave"`
	Descripcion        string          `json:"descripcion"`
	Padre              string          `json:"padre"`
	Siguiente          string          `json:"siguiente"`
	Anterior           string          `json:"anterior"`
	Modelo             DespModelo      `json:"modelo"`
	Capas              []DespCapas     `json:"capas"`
	Aor                string          `json:"aor"`
	Estado             string          `json:"estado"`
	Propiedades        DespPropiedades `json:"propiedades"`
	Version            string          `json:"version"`
	FechaCreacion      time.Time       `json:"fecha_creacion"`
	FechaActualizacion time.Time       `json:"fecha_actualizacion"`
}

// DespLinkDataArray estructura auxiliar para manejo de vistas de despliegue
type DespLinkDataArray struct {
	EsLink    bool        `json:"esLink"`
	From      int         `json:"from"`
	FromPort  string      `json:"fromPort"`
	ID        interface{} `json:"id"`
	IDVista   string      `json:"idVista"`
	Key       int         `json:"key"`
	Layer     string      `json:"layer"`
	Points    []float32   `json:"points"`
	TipoLinea string      `json:"tipoLinea"`
	To        int         `json:"to"`
	ToPort    string      `json:"toPort"`
	ZOrder    int         `json:"zOrder"`
}

// DespDei estructura auxiliar para manejo de vistas de despliegue
type DespDei struct {
	Clave string `json:"clave"`
	Grupo string `json:"grupo"`
	ID    string `json:"id"`
	Valor int    `json:"valor"`
}

// DespNodeDataArray estructura auxiliar para manejo de vistas de despliegue
type DespNodeDataArray struct {
	Angle        int     `json:"angle"`
	Category     string  `json:"category"`
	Dei          DespDei `json:"dei,omitempty"`
	ID           string  `json:"id,omitempty"`
	IDVista      string  `json:"idVista,omitempty"`
	Key          int     `json:"key"`
	Layer        string  `json:"layer"`
	Loc          string  `json:"loc"`
	MainCategory string  `json:"mainCategory,omitempty"`
	Nombre       string  `json:"nombre,omitempty"`
	ZOrder       int     `json:"zOrder"`
	EsFigura     bool    `json:"esFigura,omitempty"`
	Fill         string  `json:"fill,omitempty"`
	Size         string  `json:"size,omitempty"`
	Stroke       string  `json:"stroke,omitempty"`
	StrokeWidth  int     `json:"strokeWidth,omitempty"`
	TipoLinea    string  `json:"tipoLinea,omitempty"`
}

// DespModelo estructura auxiliar para manejo de vistas de despliegue
type DespModelo struct {
	Class                  string              `json:"class"`
	CopiesArrayObjects     bool                `json:"copiesArrayObjects"`
	CopiesArrays           bool                `json:"copiesArrays"`
	LinkDataArray          []DespLinkDataArray `json:"linkDataArray"`
	LinkFromPortIDProperty string              `json:"linkFromPortIdProperty"`
	LinkKeyProperty        string              `json:"linkKeyProperty"`
	LinkToPortIDProperty   string              `json:"linkToPortIdProperty"`
	NodeDataArray          []DespNodeDataArray `json:"nodeDataArray"`
}

// DespCapas estructura auxiliar para manejo de vistas de despliegue
type DespCapas struct {
	ID        string `json:"_id"`
	Bloqueada bool   `json:"bloqueada"`
	Clave     string `json:"clave"`
	Editable  bool   `json:"editable"`
	Nombre    string `json:"nombre"`
	Orden     int    `json:"orden"`
	OrdenSup  int    `json:"ordenSup"`
	Visible   bool   `json:"visible"`
}

// DespHistorial estructura auxiliar para manejo de vistas de despliegue
type DespHistorial struct {
	Accion        string       `json:"accion"`
	Aor           string       `json:"aor"`
	Categoria     string       `json:"categoria"`
	Fecha         string       `json:"fecha"`
	IDNodo        NullableNodo `json:"id_Nodo"`
	Propiedad     string       `json:"propiedad"`
	ValorAnterior string       `json:"valorAnterior"`
	ValorNuevo    string       `json:"valorNuevo"`
}

// DespPropiedades estructura auxiliar para manejo de vistas de despliegue
type DespPropiedades struct {
	Grid      string          `json:"grid"`
	Historial []DespHistorial `json:"historial"`
}

// CargaConfig carga datos de un despliegue a partir de la configuración
func (a *IhmDespliegue) CargaConfig(conf config.UcmAptAPIConfig, params ...string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var fullURL strings.Builder
	fullURL.WriteString(conf.URL)
	for _, p := range params {
		fullURL.WriteString(p)
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(fullURL.String())
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
	strBody := string(body)
	err = json.Unmarshal([]byte(strBody), a)
	if err != nil {
		log.Fatalln(err)
	}
}

// CargaArchivo carga datos de un despliegue a partir de un archivo
func (a *IhmDespliegue) CargaArchivo(archivo string) {
	jsonFile, err := os.Open(archivo)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Archivo de configuración abierto sin complicacicones")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, a)
	if err != nil {
		fmt.Println(err)
	}
}

// MANEJO DE DESPLIEGUES CON TOPOLOGIA

// IhmTopologia estructura principal para manejo de vistas de despliegue con topolog�a
type IhmTopologia struct {
	ID                 string          `json:"_id"`
	Nombre             string          `json:"nombre"`
	Clave              string          `json:"clave"`
	Descripcion        string          `json:"descripcion"`
	Padre              string          `json:"padre"`
	Siguiente          string          `json:"siguiente"`
	Anterior           string          `json:"anterior"`
	Modelo             TopoModelo      `json:"modelo"`
	Capas              []TopoCapas     `json:"capas"`
	Aor                string          `json:"aor"`
	Estado             string          `json:"estado"`
	Propiedades        TopoPropiedades `json:"propiedades"`
	Version            string          `json:"version"`
	FechaCreacion      time.Time       `json:"fecha_creacion"`
	FechaActualizacion time.Time       `json:"fecha_actualizacion"`
}

// TopoLinkDataArray estructura auxiliar para manejo de vistas de despliegue con topologia
type TopoLinkDataArray struct {
	EsLink    bool      `json:"esLink"`
	From      string    `json:"from"`
	FromPort  string    `json:"fromPort"`
	Key       string    `json:"key"`
	Layer     string    `json:"layer"`
	Points    []float64 `json:"points"`
	TipoLinea string    `json:"tipoLinea"`
	To        string    `json:"to"`
	ToPort    string    `json:"toPort"`
}

// TopoVariable estructura auxiliar para manejo de vistas de despliegue con topologia
type TopoVariable struct {
	Clave string `json:"clave"`
	ID    string `json:"id"`
}

// TopoDei estructura auxiliar para manejo de vistas de despliegue con topologia
type TopoDei struct {
	Dispositivo string       `json:"dispositivo"`
	Equipo      string       `json:"equipo"`
	Valor       string       `json:"valor"`
	Variable    TopoVariable `json:"variable"`
}

// TopoVerticalAlignment estructura auxiliar para manejo de vistas de despliegue con topologia
type TopoVerticalAlignment struct {
	Class   string  `json:"class"`
	OffsetX int     `json:"offsetX"`
	OffsetY int     `json:"offsetY"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
}

// TopoNodeDataArray estructura auxiliar para manejo de vistas de despliegue con topologia
type TopoNodeDataArray struct {
	Angle             int                   `json:"angle"`
	Category          string                `json:"category"`
	Dei               TopoDei               `json:"dei,omitempty"`
	ID                string                `json:"id"`
	IDVista           string                `json:"idVista"`
	Key               string                `json:"key"`
	Layer             string                `json:"layer"`
	Loc               string                `json:"loc"`
	MainCategory      string                `json:"mainCategory,omitempty"`
	Nombre            string                `json:"nombre,omitempty"`
	ZOrder            int                   `json:"zOrder"`
	Background        string                `json:"background,omitempty"`
	EsTexto           bool                  `json:"esTexto,omitempty"`
	Font              string                `json:"font,omitempty"`
	Size              string                `json:"size,omitempty"`
	Stroke            string                `json:"stroke,omitempty"`
	Text              string                `json:"text,omitempty"`
	Through           bool                  `json:"through,omitempty"`
	Underline         bool                  `json:"underline,omitempty"`
	VerticalAlignment TopoVerticalAlignment `json:"verticalAlignment,omitempty"`
	TextAlign         string                `json:"textAlign,omitempty"`
	Vertical          string                `json:"vertical,omitempty"`
}

// TopoModelo estructura auxiliar para manejo de vistas de despliegue con topologia
type TopoModelo struct {
	Class                  string              `json:"class"`
	CopiesArrayObjects     bool                `json:"copiesArrayObjects"`
	CopiesArrays           bool                `json:"copiesArrays"`
	LinkDataArray          []TopoLinkDataArray `json:"linkDataArray"`
	LinkFromPortIDProperty string              `json:"linkFromPortIdProperty"`
	LinkKeyProperty        string              `json:"linkKeyProperty"`
	LinkToPortIDProperty   string              `json:"linkToPortIdProperty"`
	NodeDataArray          []TopoNodeDataArray `json:"nodeDataArray"`
}

// TopoCapas estructura auxiliar para manejo de vistas de despliegue con topologia
type TopoCapas struct {
	ID        string `json:"_id"`
	Bloqueada bool   `json:"bloqueada"`
	Clave     string `json:"clave"`
	Editable  bool   `json:"editable"`
	Nombre    string `json:"nombre"`
	Orden     int    `json:"orden"`
	OrdenSup  int    `json:"ordenSup"`
	Visible   bool   `json:"visible"`
}

// TopoHistorial estructura auxiliar para manejo de vistas de despliegue con topologia
type TopoHistorial struct {
	Accion        string       `json:"accion"`
	Aor           string       `json:"aor"`
	Categoria     string       `json:"categoria"`
	Fecha         string       `json:"fecha"`
	IDNodo        NullableNodo `json:"id_Nodo"`
	Propiedad     string       `json:"propiedad"`
	ValorAnterior string       `json:"valorAnterior"`
	ValorNuevo    string       `json:"valorNuevo"`
}

// TopoSize estructura auxiliar para manejo de vistas de despliegue con topologia
type TopoSize struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

// TopoTopologia estructura auxiliar para manejo de vistas de despliegue con topologia
type TopoTopologia struct {
	A  string `json:"a"`
	De string `json:"de"`
}

// TopoPropiedades estructura auxiliar para manejo de vistas de despliegue con topologia
type TopoPropiedades struct {
	Grid      string          `json:"grid"`
	Historial []TopoHistorial `json:"historial"`
	Size      TopoSize        `json:"size"`
	Topologia []TopoTopologia `json:"topologia"`
}

// CargaConfig carga datos de un despliegue con topologia a partir de la configuración
func (a *IhmTopologia) CargaConfig(conf config.UcmAptAPIConfig, params ...string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var fullURL strings.Builder
	fullURL.WriteString(conf.URL)
	for _, p := range params {
		fullURL.WriteString(p)
	}
	client := &http.Client{Transport: tr}
	//client.Post()
	resp, err := client.Get(fullURL.String())
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
	strBody := string(body)
	err = json.Unmarshal([]byte(strBody), a)
	if err != nil {
		log.Fatalln(err)
	}
}

// CargaArchivo carga datos de un despliegue a partir de un archivo
func (a *IhmTopologia) CargaArchivo(archivo string) {
	jsonFile, err := os.Open(archivo)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Archivo de configuración abierto sin complicacicones")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, a)
	if err != nil {
		fmt.Println(err)
	}
}

// MANEJO DE CONSULTAS AL HISTORIAL DE MEDICIONES

// HistorialMediciones estructura principal para manejo de mediciones
type HistorialMediciones struct {
	Total    int            `json:"total"`
	Muestras []HistMuestras `json:"muestras"`
}

// HistMuestras estructura auxiliar para manejo de mediciones
type HistMuestras struct {
	ID             string    `json:"_id"`
	AppID          string    `json:"app_id"`
	BanderasDnp    int       `json:"banderas_dnp"`
	BanderasKernel int       `json:"banderas_kernel"`
	Consola        string    `json:"consola"`
	EsEventoDnp    bool      `json:"es_evento_dnp"`
	EsEventoIhm    bool      `json:"es_evento_ihm"`
	FechaLocal     int64     `json:"fecha_local"`
	FechaRemota    int       `json:"fecha_remota"`
	IDAlarma       int       `json:"id_alarma"`
	IDEvento       int       `json:"id_evento"`
	IndicePunto    int       `json:"indice_punto"`
	Silenciado     bool      `json:"silenciado"`
	Timestamp      time.Time `json:"timestamp"`
	Tipo           string    `json:"tipo"`
	TipoMensaje    int       `json:"tipo_mensaje"`
	Usuario        string    `json:"usuario"`
	Valor          float64   `json:"valor"`
	ValorCrudo     int       `json:"valor_crudo"`
}

// CargaConfig carga datos del hitorial de mediciones a partir de la configuración
func (a *HistorialMediciones) CargaConfig(conf config.UcmAptAPIConfig, params ...string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var fullURL strings.Builder
	fullURL.WriteString(conf.URL)
	for _, p := range params {
		fullURL.WriteString(p)
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(fullURL.String())
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
	strBody := string(body)
	err = json.Unmarshal([]byte(strBody), a)
	if err != nil {
		log.Fatalln(err)
	}
}

// CargaArchivo carga datos del historial de mediciones a partir de un archivo
func (a *HistorialMediciones) CargaArchivo(archivo string) {
	jsonFile, err := os.Open(archivo)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Archivo de configuraci�n abierto sin complicacicones")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, a)
	if err != nil {
		fmt.Println(err)
	}
}

// VALORES ACTUALES

// ValoresActuales estructura para manejo de los valores actuales
type ValoresActuales []struct {
	IndicePunto         string `json:"indice_punto"`
	Nombre              string `json:"nombre"`
	Valor               string `json:"valor"`
	BanderasDnp         string `json:"banderas_dnp"`
	BanderasKernel      string `json:"banderas_kernel"`
	InhibirPunto        string `json:"inhibir_punto"`
	Libranza            string `json:"libranza"`
	ModoManual          string `json:"modo_manual"`
	ControlSeleccionado string `json:"control_seleccionado"`
	InhibirAlarmas      string `json:"inhibir_alarmas"`
	IDEvento            string `json:"id_evento"`
}

// CargaConfig carga valores actuales a partir de la configuraci�n
func (a *ValoresActuales) CargaConfig(conf config.UcmAptAPIConfig, params ...string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var fullURL strings.Builder
	fullURL.WriteString(conf.URL)
	for _, p := range params {
		fullURL.WriteString(p)
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(fullURL.String())
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
	strBody := string(body)
	err = json.Unmarshal([]byte(strBody), a)
	if err != nil {
		log.Fatalln(err)
	}
}

// CargaArchivo carga datos del historial de mediciones a partir de un archivo
func (a *ValoresActuales) CargaArchivo(archivo string) {
	jsonFile, err := os.Open(archivo)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Archivo de configuraci�n abierto sin complicacicones")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, a)
	if err != nil {
		fmt.Println(err)
	}
}

// PUNTOS

// PuntosUcm estructura principal para manejar puntos detallados y estados de la UCM
type PuntosUcm []struct {
	ID                      string          `json:"id"`
	IDTipoPuntoUcm          string          `json:"id_tipo_punto_ucm"`
	IDAdr                   string          `json:"id_adr"`
	Clave                   string          `json:"clave"`
	Nombre                  string          `json:"nombre"`
	Descripcion             string          `json:"descripcion"`
	IndicePunto             int             `json:"indice_punto"`
	Identificador           int             `json:"identificador"`
	IdentificadorSecundario interface{}     `json:"identificador_secundario"`
	Formula                 string          `json:"formula"`
	EsNombreEditado         interface{}     `json:"es_nombre_editado"`
	IDOrigenPunto           string          `json:"id_origen_punto"`
	ValorInicial            int             `json:"valor_inicial"`
	Adr                     PtosAdr         `json:"adr,omitempty"`
	TipoPunto               PtosTipoPunto   `json:"tipo_punto"`
	OrigenPunto             PtosOrigenPunto `json:"origen_punto"`
}

// PtosAdr estructura auxiliar para manejar puntos detallados y estados de la UCM
type PtosAdr struct {
	ID          string    `json:"id"`
	Clave       string    `json:"clave"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	FechaAlta   time.Time `json:"fecha_alta"`
	Activo      bool      `json:"activo"`
}

// PtosCatalogoTipoPunto estructura auxiliar para manejar puntos detallados y estados de la UCM
type PtosCatalogoTipoPunto struct {
	ID        string    `json:"id"`
	Clave     string    `json:"clave"`
	Nombre    string    `json:"nombre"`
	FechaAlta time.Time `json:"fecha_alta"`
}

// PtosTipoPunto estructura auxiliar para manejar puntos detallados y estados de la UCM
type PtosTipoPunto struct {
	ID                  string                `json:"id"`
	IDCatalogoTipoPunto string                `json:"id_catalogo_tipo_punto"`
	Clave               string                `json:"clave"`
	Nombre              string                `json:"nombre"`
	Descripcion         string                `json:"descripcion"`
	NombreControl       interface{}           `json:"nombre_control"`
	FechaAlta           time.Time             `json:"fecha_alta"`
	CatalogoTipoPunto   PtosCatalogoTipoPunto `json:"catalogo_tipo_punto"`
}

// PtosOrigenPunto estructura auxiliar para manejar puntos detallados y estados de la UCM
type PtosOrigenPunto struct {
	ID          string `json:"id"`
	Tipo        string `json:"tipo"`
	Origen      string `json:"origen"`
	TipoPunto   string `json:"tipo_punto"`
	Valor       int    `json:"valor"`
	Descripcion string `json:"descripcion"`
}

// CargaConfig carga datos de los puntos detallados y estados a partir de la configuración
func (a *PuntosUcm) CargaConfig(conf config.UcmAptAPIConfig, params ...string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var fullURL strings.Builder
	fullURL.WriteString(conf.URL)
	for _, p := range params {
		fullURL.WriteString(p)
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(fullURL.String())
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
	strBody := string(body)
	err = json.Unmarshal([]byte(strBody), a)
	if err != nil {
		log.Fatalln(err)
	}
}

// CargaArchivo carga datos de los puntos detallados y estados a partir de un archivo
func (a *PuntosUcm) CargaArchivo(archivo string) {
	jsonFile, err := os.Open(archivo)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Archivo de configuración abierto sin complicacicones")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, a)
	if err != nil {
		fmt.Println(err)
	}
}

// LICENCIAS

// Licencias estructura principal para manejo de licencias
type Licencias []struct {
	Punto     string        `json:"Punto"`
	Licencias LicsLicencias `json:"Licencias"`
}

// LicsLicencias estructura auxiliar para manejo de licencias
type LicsLicencias struct {
	Comentario     string `json:"comentario"`
	Estampa        string `json:"estampa"`
	FechaFin       string `json:"fechaFin"`
	FechaIni       string `json:"fechaIni"`
	NumeroLicencia string `json:"numero_licencia"`
	Tipo           string `json:"tipo"`
	Usuario        string `json:"usuario"`
}

// CargaConfig carga datos de las licencias a partir de la configuraci�n
func (a *Licencias) CargaConfig(conf config.UcmAptAPIConfig, params ...string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var fullURL strings.Builder
	fullURL.WriteString(conf.URL)
	for _, p := range params {
		fullURL.WriteString(p)
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 2}
	resp, err := client.Get(fullURL.String())
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
	strBody := string(body)
	err = json.Unmarshal([]byte(strBody), a)
	if err != nil {
		log.Fatalln(err)
	}
}

// CargaArchivo carga datos de las licencias a partir de un archivo
func (a *Licencias) CargaArchivo(archivo string) {
	jsonFile, err := os.Open(archivo)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Archivo de configuración abierto sin complicacicones")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, a)
	if err != nil {
		fmt.Println(err)
	}
}

// LOGIN

// UcmLog estructura para manejo de sesiones
type UcmLog struct {
	Usuario LogRequest
	LogInfo LogResponse
	Expira  time.Time
}

// LogRequest elemento para manejo de sesiones
type LogRequest struct {
	Usuario  string `json:"username"`
	Consigna string `json:"password"`
}

// LogResponse elemento para manejo de sesiones
type LogResponse struct {
	Code   int    `json:"code"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

// Iniciar constructor de UcmLog
func (l *UcmLog) Iniciar(usuario string, consigna string) error {
	l.Usuario = LogRequest{
		Usuario:  usuario,
		Consigna: consigna,
	}
	l.LogInfo = LogResponse{}
	return nil
}

// ObtenerToken verifica la sesión, solicita un token en caso necesario
func (l UcmLog) ObtenerToken(apiAuth config.UcmAptAPIConfig) (bool, error) {
	var valido = false
	var ok = false

	// Ubicar caracteristicas del token actual
	if (l.LogInfo != LogResponse{}) {
		te := time.Now().Add(time.Minute * time.Duration(2))
		if l.Expira.Before(te) {
			valido = true
		}
	}

	if !valido {
		requestBody, err := json.Marshal(map[string]string{
			"username": "rmata",
			"password": "rmata26",
		})
		if err != nil {
			log.Fatalln(err)
		}

		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client := &http.Client{Transport: customTransport}

		request, err := http.NewRequest("POST", apiAuth.URL, bytes.NewBuffer(requestBody))
		request.Header.Set("Content-type", "application/json; charset=utf-8")
		request.Header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
		request.Header.Set("Access-Control-Allow-Origin", "*")

		if err != nil {
			log.Fatalln(err)
		}

		resp, err := client.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var jResp LogResponse
		err = json.Unmarshal(body, &jResp)
		if err != nil {
			log.Fatalln(err)
		}
		l.LogInfo = jResp

		l.Expira = time.Now().Add(time.Minute * time.Duration(30))
		fmt.Println(l.LogInfo)
		ok = true
	} else {
		ok = true
	}

	return ok, nil
}

func MsgBitacoraAction(cargaPerd int64, equipos string) string {
	var byte_buf bytes.Buffer

	byte_buf.WriteString("\"{'id_solicitud':'001',")
	byte_buf.WriteString("'timestamp':")
	byte_buf.WriteString(fmt.Sprintf("%10d", time.Now().Unix()))
	byte_buf.WriteString(",'tipofalla': 'falla permanente',")
	byte_buf.WriteString("'ubicacion': 'Zona Camargo',")
	byte_buf.WriteString("'cargaperdida':'")
	byte_buf.WriteString(fmt.Sprintf("%8d", cargaPerd))
	byte_buf.WriteString(" VA',")
	byte_buf.WriteString("'equipos':'")
	byte_buf.WriteString(equipos)
	byte_buf.WriteString("'}'\"")

	return byte_buf.String()
}

type AptMsgBitacora struct {
	CodigoMsj   int    `json:"codigo_msj"`
	TipoMensaje int    `json:"tipo_mensaje"`
	Comando     int    `json:"comando"`
	Comentario  string `json:"comentario"`
	Consola     string `json:"consola"`
	Fecha       int64  `json:"fecha"`
	IDInterfaz  int    `json:"id_interfaz"`
	IDRegistro  string `json:"id_registro"`
	IndicePunto int    `json:"indice_punto"`
	Metadata    string `json:"metadata"`
	Persistente bool   `json:"persistente"`
	Usuario     string `json:"usuario"`
}

func (l UcmLog) AtenderFalla(apiAuthx config.UcmAptAPIConfig, apiBitacora config.UcmAptAPIConfig, cfg config.UcmAptConfig) (string, error) {
	var detenido string = ""
	usr, err := cfg.ExtraeParam("authuser")
	pwd, err := cfg.ExtraeParam("authpwd")

	requestBody, err := json.Marshal(map[string]string{
		"username": usr.Valor,
		"password": pwd.Valor,
	})
	if err != nil {
		log.Fatalln(err)
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}
	request, err := http.NewRequest("POST", apiAuthx.URL, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-type", "application/json; charset=utf-8")
	request.Header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
	request.Header.Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Autorizado")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var jResp LogResponse
	err = json.Unmarshal(body, &jResp)
	if err != nil {
		log.Fatalln(err)
	}
	l.LogInfo = jResp

	l.Expira = time.Now().Add(time.Minute * time.Duration(30))

	//	log.Println(l.LogInfo.Token)
	if err == nil {
		msgbit := AptMsgBitacora{
			CodigoMsj:   3008,
			TipoMensaje: 5000,
			Comando:     8,
			Comentario:  "Aviso de automatismo detenido",
			Consola:     "Automatismo",
			Fecha:       time.Now().Unix(),
			IDInterfaz:  8004,
			IDRegistro:  "",
			IndicePunto: -1,
			Metadata:    "{}",
			Persistente: true,
			Usuario:     "rmata",
		}

		requestBody, err := json.Marshal(msgbit)
		if err != nil {
			log.Fatalln(err)
		}

		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client := &http.Client{Transport: customTransport}
		var bearer = "Bearer " + l.LogInfo.Token

		request, err := http.NewRequest("POST", apiBitacora.URL, bytes.NewBuffer(requestBody))
		request.Header.Set("Content-type", "application/json; charset=utf-8")
		request.Header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
		request.Header.Set("Access-Control-Allow-Origin", "*")
		request.Header.Add("Authorization", bearer)
		if err != nil {
			log.Fatalln(err)
		}

		resp, err := client.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		/*
			resp, err := client.Post(conf.URL, "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				log.Fatalln(err)
			}*/

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		detenido = string(body)
		log.Println("Enviado")
	}
	return detenido, nil
}


func (l UcmLog) EnviarDetenido(apiAuthx config.UcmAptAPIConfig, apiBitacora config.UcmAptAPIConfig, cfg config.UcmAptConfig) (string, error) {
	var detenido string = ""
	usr, err := cfg.ExtraeParam("authuser")
	pwd, err := cfg.ExtraeParam("authpwd")

	requestBody, err := json.Marshal(map[string]string{
		"username": usr.Valor,
		"password": pwd.Valor,
	})
	if err != nil {
		log.Fatalln(err)
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}
	request, err := http.NewRequest("POST", apiAuthx.URL, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-type", "application/json; charset=utf-8")
	request.Header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
	request.Header.Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Autorizado")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var jResp LogResponse
	err = json.Unmarshal(body, &jResp)
	if err != nil {
		log.Fatalln(err)
	}
	l.LogInfo = jResp

	l.Expira = time.Now().Add(time.Minute * time.Duration(30))

	//	log.Println(l.LogInfo.Token)
	if err == nil {
		msgbit := AptMsgBitacora{
			CodigoMsj:   3008,
			TipoMensaje: 5000,
			Comando:     8,
			Comentario:  "Aviso de automatismo detenido",
			Consola:     "Automatismo",
			Fecha:       time.Now().Unix(),
			IDInterfaz:  8004,
			IDRegistro:  "",
			IndicePunto: -1,
			Metadata:    "{}",
			Persistente: true,
			Usuario:     "rmata",
		}

		requestBody, err := json.Marshal(msgbit)
		if err != nil {
			log.Fatalln(err)
		}

		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client := &http.Client{Transport: customTransport}
		var bearer = "Bearer " + l.LogInfo.Token

		request, err := http.NewRequest("POST", apiBitacora.URL, bytes.NewBuffer(requestBody))
		request.Header.Set("Content-type", "application/json; charset=utf-8")
		request.Header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
		request.Header.Set("Access-Control-Allow-Origin", "*")
		request.Header.Add("Authorization", bearer)
		if err != nil {
			log.Fatalln(err)
		}

		resp, err := client.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		/*
			resp, err := client.Post(conf.URL, "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				log.Fatalln(err)
			}*/

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		detenido = string(body)
		log.Println("Enviado")
	}
	return detenido, nil
}

func (l UcmLog) EnviarReinicio(apiAuthx config.UcmAptAPIConfig, apiBitacora config.UcmAptAPIConfig, cfg config.UcmAptConfig) (string, error) {
	var detenido string = ""
	usr, err := cfg.ExtraeParam("authuser")
	pwd, err := cfg.ExtraeParam("authpwd")

	requestBody, err := json.Marshal(map[string]string{
		"username": usr.Valor,
		"password": pwd.Valor,
	})
	if err != nil {
		log.Fatalln(err)
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}
	request, err := http.NewRequest("POST", apiAuthx.URL, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-type", "application/json; charset=utf-8")
	request.Header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
	request.Header.Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Autorizado")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var jResp LogResponse
	err = json.Unmarshal(body, &jResp)
	if err != nil {
		log.Fatalln(err)
	}
	l.LogInfo = jResp

	l.Expira = time.Now().Add(time.Minute * time.Duration(30))

	//	log.Println(l.LogInfo.Token)
	if err == nil {
		msgbit := AptMsgBitacora{
			CodigoMsj:   3008,
			TipoMensaje: 5000,
			Comando:     8,
			Comentario:  "Aviso de automatismo detenido",
			Consola:     "Automatismo",
			Fecha:       time.Now().Unix(),
			IDInterfaz:  8004,
			IDRegistro:  "",
			IndicePunto: -1,
			Metadata:    "{}",
			Persistente: true,
			Usuario:     "rmata",
		}

		requestBody, err := json.Marshal(msgbit)
		if err != nil {
			log.Fatalln(err)
		}

		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client := &http.Client{Transport: customTransport}
		var bearer = "Bearer " + l.LogInfo.Token

		request, err := http.NewRequest("POST", apiBitacora.URL, bytes.NewBuffer(requestBody))
		request.Header.Set("Content-type", "application/json; charset=utf-8")
		request.Header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
		request.Header.Set("Access-Control-Allow-Origin", "*")
		request.Header.Add("Authorization", bearer)
		if err != nil {
			log.Fatalln(err)
		}

		resp, err := client.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		/*
			resp, err := client.Post(conf.URL, "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				log.Fatalln(err)
			}*/

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		detenido = string(body)
		log.Println("Enviado")
	}
	return detenido, nil
}

/*
	"CONFIGURACIONES", *
	"DESPLIEGUE",      *
	"TOPOLOGIA",       *
	"HISTORIAL",       *
	"VALORES",         *
	"PUNTOS",          *
	"LICENCIAS",       *
	"LOGIN",
	"BITACORA",
*/
