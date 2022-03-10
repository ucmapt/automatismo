package models

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/streadway/amqp"
)

type RmqBitacoraMensaje struct {
	CodigoMsj   int32      `json:"codigo_msj"`
	TipoMensaje int32      `json:"tipo_mensaje"`
	IndicePunto int32      `json:"indice_punto"`
	Comando     int32      `json:"comando"`
	Fecha       int64      `json:"fecha"`
	IdInterfaz  int32      `json:"id_interfaz"`
	IdRegistro  [48]byte   `json:"id_registro"`
	Comentario  [128]byte  `json:"comentario"`
	Persistente bool       `json:"persistente"`
	Usuario     [48]byte   `json:"usuario"`
	Consola     [48]byte   `json:"consola"`
	Metadata    [4096]byte `json:"metadata"`
}

type RmqBitacoraMensajeStr struct {
	CodigoMsj   int32  `json:"codigo_msj"`
	TipoMensaje int32  `json:"tipo_mensaje"`
	IndicePunto int32  `json:"indice_punto"`
	Comando     int32  `json:"comando"`
	Fecha       int64  `json:"fecha"`
	IdInterfaz  int32  `json:"id_interfaz"`
	IdRegistro  string `json:"id_registro"`
	Comentario  string `json:"comentario"`
	Persistente bool   `json:"persistente"`
	Usuario     string `json:"usuario"`
	Consola     string `json:"consola"`
	Metadata    string `json:"metadata"`
}

type RmqBitacoraMensajeB struct {
	Id          bson.ObjectId `bson:"_id" json:"_id,omitempty"`
	TipoMensaje int32         `json:"tipo_mensaje"`
	IndicePunto int32         `json:"indice_punto"`
	Comando     int32         `json:"comando"`
	Fecha       int64         `json:"fecha"`
	IdInterfaz  int32         `json:"id_interfaz"`
	IdRegistro  string        `json:"id_registro"`
	Comentario  string        `json:"comentario"`
	Usuario     string        `json:"usuario"`
	Consola     string        `json:"consola"`
	CodigoMsj   int32         `json:"codigo_msj"`
	Metadata    string        `json:"metadata"`
	Persistente bool          `json:"persistente"`
}

type RmqBitacoraComando struct {
	CodigoMsj   int32      `json:"codigo_msj"`
	TipoMensaje int32      `json:"tipo_mensaje"`
	IndicePunto int32      `json:"indice_punto"`
	Comando     int32      `json:"comando"`
	Fecha       int64      `json:"fecha"`
	IdInterfaz  int32      `json:"id_interfaz"`
	IdRegistro  [48]byte   `json:"id_registro"`
	Comentario  [128]byte  `json:"comentario"`
	Persistente bool       `json:"persistente"`
	Usuario     [48]byte   `json:"usuario"`
	Consola     [48]byte   `json:"consola"`
	Metadata    [4096]byte `json:"metadata"`
}

type RmqExchange struct {
	Name          string     `json:"nombre"`
	Kind          string     `json:"tipo"`
	Durable       *bool      `json:"durable"`
	AutoDelete    *bool      `json:"auto_borrar"`
	Internal      bool       `json:"interno"`
	NoWait        bool       `json:"no_esperar"`
	Args          amqp.Table `json:"args"`
	CommonHeaders amqp.Table `json:"cabeceras_comunes"`
}

func RmqMensajeIniciar() RmqBitacoraMensaje {
	r := RmqBitacoraMensaje{
		CodigoMsj:   0,
		TipoMensaje: 5000,
		IndicePunto: -1,
		Comando:     6,
		Fecha:       time.Now().UnixNano() / 1000000,
		IdInterfaz:  8004,
		IdRegistro:  [48]byte{},
		Comentario:  [128]byte{},
		Persistente: true,
		Usuario:     [48]byte{},
		Consola:     [48]byte{},
		Metadata:    [4096]byte{},
	}
	copy(r.Consola[:], StrToBytes("Automatismo", 48)) //se presentará en la consola del automatismo
	copy(r.Comentario[:], StrToBytes("Solicitud de transferencia de datos desde Editor", 128))
	copy(r.IdRegistro[:], StrToBytes("", 48))
	copy(r.Usuario[:], StrToBytes("", 48))
	copy(r.Metadata[:], StrToBytes("", 48)) //no hace falta incluir los mensajes usados en automatismo
	return r
}

func RmqMensajeTerminaOk() RmqBitacoraMensaje {
	r := RmqBitacoraMensaje{
		CodigoMsj:   0,
		TipoMensaje: 5000,
		IndicePunto: -1,
		Comando:     7,
		Fecha:       time.Now().UnixNano() / 1000000,
		IdInterfaz:  8004,
		IdRegistro:  [48]byte{},
		Comentario:  [128]byte{},
		Persistente: true,
		Usuario:     [48]byte{},
		Consola:     [48]byte{},
		Metadata:    [4096]byte{},
	}
	copy(r.Consola[:], StrToBytes("Automatismo", 48)) //se presentará en la consola del automatismo
	copy(r.Comentario[:], StrToBytes("Transferencia de datos desde Editor concluida normal", 128))
	copy(r.IdRegistro[:], StrToBytes("", 48))
	copy(r.Usuario[:], StrToBytes("", 48))
	copy(r.Metadata[:], StrToBytes("", 48)) //no hace falta incluir los mensajes usados en automatismo
	return r
}

func RmqMensajeTerminaMal() RmqBitacoraMensaje {
	r := RmqBitacoraMensaje{
		CodigoMsj:   0,
		TipoMensaje: 5000,
		IndicePunto: -1,
		Comando:     8,
		Fecha:       time.Now().UnixNano() / 1000000,
		IdInterfaz:  8004,
		IdRegistro:  [48]byte{},
		Comentario:  [128]byte{},
		Persistente: true,
		Usuario:     [48]byte{},
		Consola:     [48]byte{},
		Metadata:    [4096]byte{},
	}
	copy(r.Consola[:], StrToBytes("Automatismo", 48)) //se presentará en la consola del automatismo
	copy(r.Comentario[:], StrToBytes("Transferencia de datos desde Editor fallida", 128))
	copy(r.IdRegistro[:], StrToBytes("", 48))
	copy(r.Usuario[:], StrToBytes("", 48))
	copy(r.Metadata[:], StrToBytes("", 48)) //no hace falta incluir los mensajes usados en automatismo
	return r
}

func (msg RmqBitacoraMensaje) EsIniciar() bool {
	return (msg.IdInterfaz == 8004) || (msg.Comando == 6)
}

func AdecuarCadenaEnMapa(m map[string]interface{}) (map[string]interface{}, error) {
	for k := range m {
		if m[k] == nil {
			return nil, errors.New("Problema al convertir, existe valor nulo en el campo: " + k)
		}

		if reflect.TypeOf(m[k]).String() == "string" {
			var array [4096]byte
			slice := []byte(m[k].(string))

			if len(slice) >= len(array) {
				return nil, errors.New("Problema al convertir string a []byte.")
			}

			copy(array[:], slice)
			array[len(slice)] = '\000'
			m[k] = array

		}
	}
	return m, nil
}

func StrToBytes(s string, l int) []byte {
	var array = make([]byte, l)
	slice := []byte(s)
	copy(array[:], slice)
	return array
}

func BytesToString(data []byte) string {
	return string(data[:])
}

func InterfaceToBytes(e interface{}) ([]byte, error) {

	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.LittleEndian, e)
	if err != nil {
		panic(err)
	}

	return buf.Bytes(), nil
}

func (r RmqBitacoraMensaje) Adecuada() RmqBitacoraMensajeStr {
	rx := RmqBitacoraMensajeStr{
		CodigoMsj:   r.CodigoMsj,
		TipoMensaje: r.TipoMensaje,
		IndicePunto: r.IndicePunto,
		Comando:     r.Comando,
		Fecha:       r.Fecha,
		IdInterfaz:  r.IdInterfaz,
		IdRegistro:  BytesToString(r.IdRegistro[:]),
		Comentario:  BytesToString(r.Comentario[:]),
		Persistente: r.Persistente,
		Usuario:     BytesToString(r.Usuario[:]),
		Consola:     BytesToString(r.Consola[:]),
		Metadata:    BytesToString(r.Metadata[:]),
	}
	return rx
}
