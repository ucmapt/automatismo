package main

// Versión 2.0 del Módulo del Automatismo y Procesador de Topología de la UCM-CFE (c) 2021 - 2022

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
	"github.com/ucmapt/automatismo/api"
	"github.com/ucmapt/automatismo/common/events"
	"github.com/ucmapt/automatismo/common/fsm"
	"github.com/ucmapt/automatismo/common/graphs"
	"github.com/ucmapt/automatismo/config"
	"github.com/ucmapt/automatismo/data/postgres/motorapt"
	"github.com/ucmapt/automatismo/ihm"
	"github.com/ucmapt/automatismo/models"
)

/// Manejo de máquina de estados
const (
	Cargando         fsm.StateType = "Cargando"
	Bloqueado        fsm.StateType = "Bloqueado"
	FueraLinea       fsm.StateType = "FueraLinea"
	Operando         fsm.StateType = "Operando"
	SinConfiguracion fsm.StateType = "SinConfiguracion"

	IniciaCarga          fsm.EventType = "IniciaCarga"
	NoHayConfiguracion   fsm.EventType = "NoHayConfiguracion"
	SeDesactivaApt       fsm.EventType = "SeDesactivaApt"
	SeTransfiereDb       fsm.EventType = "SeTransfiereDb"
	SeReactivaApt        fsm.EventType = "SeReactivaApt"
	TerminaTransferencia fsm.EventType = "TerminaTransferencia"
)

// Acciones al cargar
type CargandoAction struct{}

func (a *CargandoAction) Execute(evContext fsm.EventContext) fsm.EventType {
	ihm.Letrero("Cargando configuraciones ...")
	var err error
	if cfgA.Cargar("automatismo.ucm") {
		// Cargar datos ...
		err = initDB()
		if err != nil {
			log.Fatalf("Error al conectar datos: %s", err)
			lanzarEvento("PROBLEMA AL CARGAR CONFIGURACIONES", "AYPT reporta problemas al cargar", "ERROR")
			return NoHayConfiguracion
		}

		bulk, err = recoverData()
		if err != nil {
			log.Fatalf("Error al recuperar datos: %s", err)
			lanzarEvento("PROBLEMA AL CARGAR CONFIGURACIONES", "AYPT reporta problemas al cargar", "ERROR")
			return NoHayConfiguracion
		}

		ihm.Suceso("Configuración cargada")
		lanzarEvento("CONFIGURACION CARGADA", "AYPT termina carga de configuraciones", "AVISO")
		return SeReactivaApt
	}

	return NoHayConfiguracion
}

type OperandoAction struct{}

func (a *OperandoAction) Execute(evContext fsm.EventContext) fsm.EventType {
	ihm.Letrero("Automatismo en operación, esperando eventos ...")
	preparaEventos()

	return fsm.NoAction
}

type BloqueadoAction struct{}

func (a *BloqueadoAction) Execute(evContext fsm.EventContext) fsm.EventType {
	ihm.Suceso("Se bloquea automatismo ...")
	return fsm.NoAction
}

type FueraLineaAction struct{}

func (a *FueraLineaAction) Execute(evContext fsm.EventContext) fsm.EventType {
	ihm.Suceso("Se desactiva automatismo ...")
	return fsm.NoAction
}

type SinConfiguracionAction struct{}

func (a *SinConfiguracionAction) Execute(evContext fsm.EventContext) fsm.EventType {
	ihm.Problema("No se pudo cargar configuracón")
	return fsm.NoAction
}

// MAQUINA FINITA DE ESTADOS DE LA APICACION
// Su fundamneto está en fsm, conbase en rutinas que respondan a cada cambio de estado
// Las rutinas deben estar en el mismo alcance, por eso se han manejado aquí,
// se traajaba en adecuarlo de otra manera

func newAptEngineFSM() *fsm.StateMachine {
	return &(fsm.StateMachine{
		States: fsm.States{
			fsm.Default: fsm.State{
				Events: fsm.Events{
					IniciaCarga: Cargando,
				},
			},
			Cargando: fsm.State{
				Action: &CargandoAction{},
				Events: fsm.Events{
					NoHayConfiguracion: SinConfiguracion,
					SeTransfiereDb:     Bloqueado,
					SeDesactivaApt:     FueraLinea,
					SeReactivaApt:      Operando,
				},
			},
			Operando: fsm.State{
				Action: &OperandoAction{},
				Events: fsm.Events{
					SeTransfiereDb: Bloqueado,
					SeDesactivaApt: FueraLinea,
				},
			},
			Bloqueado: fsm.State{
				Action: &BloqueadoAction{},
				Events: fsm.Events{
					TerminaTransferencia: Cargando,
				},
			},
			FueraLinea: fsm.State{
				Action: &FueraLineaAction{},
				Events: fsm.Events{
					NoHayConfiguracion: SinConfiguracion,
					SeTransfiereDb:     Bloqueado,
					SeReactivaApt:      Operando,
				},
			},
			SinConfiguracion: fsm.State{
				Action: &SinConfiguracionAction{},
				Events: fsm.Events{},
			},
		},
	})
}

const (
	EXCHANGE_BITACORA = "bitacoras"
	EDITOPO_QUEUE     = "editopoqueue"
	AUTOMATISMO_QUEUE = "automatismoqueue"
)

var (
	// Conector para RabbitMQ
	conn *amqp.Connection

	// Manejador de la configuración
	cfgA config.UcmAptConfig

	// Canales de eventos
	chDetener   chan events.DataEvent
	chReiniciar chan events.DataEvent

	// Bus de eventos
	eb = &events.EventBus{
		Subscribers: map[string]events.DataChannelSlice{},
	}

	// Grafo crudo
	bulk *graphs.BulkGraph
)

// Rutina general para manejo de mensajes de error
func manejaError(err error, msg string) {
	if err != nil {
		ihm.Problema(msg)
		log.Fatalf("%s: %s", msg, err)
	}
}

// Rutina de configuración de parámetros de RabbitMQ
func conejoURL() string {
	pc, err := cfgA.ExtraeParam("rmqsrv")
	if err != nil {
		return ""
	}

	pu, err := cfgA.ExtraeParam("rmquser")
	if err != nil {
		return ""
	}

	pp, err := cfgA.ExtraeParam("rmqpwd")
	if err != nil {
		return ""
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", pu.Valor, pp.Valor, pc.Valor)
	return url
}

func actualizaCircuitos() {
	var err error
	grafo := graphs.NewGraph()

	for _, n := range bulk.Nodes {
		err = grafo.AddNode(*n.Name)
		if err != nil {
			fmt.Printf("Error al manejar nodos: %s\n", err)
		}
	}
	for _, f := range bulk.Feeders {
		err = grafo.AddNode(*f.Name)
		if err != nil {
			fmt.Printf("Error al manejar nodos: %s\n", err)
		}
		err = grafo.AddFullFeeder(*f.CircuitoFijo, *f.Name, 23000)
		if err != nil {
			fmt.Printf("Error al manejar alimentadores: %s\n", err)
		}
	}

	for _, l := range bulk.Lines {
		//err = grafo.AddLine(*l.Node1, *l.Node2)
		err = grafo.AddLineFull(*l.Name, *l.Node1, *l.Node2, l.Sw1, l.Sw2, true, l.X0, l.X1, l.C0, l.C1, l.R0, l.R1)
		if err != nil {
			fmt.Printf("Error al manejar lineas: %s\n", err)
		}
	}

	fmt.Println("Comienza coloreo ....")
	printSummary(bulk)
	grafo.Colorize()

	grafo.UpdateViews()

}

// PRoceso principal de preparación de eventos
// TODO - Se estaba trabajndo en adecuarlo como un brocker con funciones más simples, eliminando variables globales
// TODO - Incorporar funciones de configuración
// TODO - Simplificar los llamados a APIs
func preparaEventos() {
	// Preparando eventos gestionados a través de RabbitMQ
	// Falla franca detectada

	conn, err := amqp.Dial(conejoURL())
	log.Println(conejoURL())
	manejaError(err, "Problema al conectar con RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	manejaError(err, "No se ha podido crear canal en RabbitMQ")
	/// 2DO : Integrar esquema de configuracion en colas de Rabbit

	defer ch.Close()
	qFallas, err := ch.QueueDeclare(
		AUTOMATISMO_QUEUE, // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	manejaError(err, "No se pudo gestionar cola de mensajes")

	err = ch.QueueBind(
		qFallas.Name,      // name
		"#",               //key
		EXCHANGE_BITACORA, //exchange
		false,             //noWait
		nil,               //args
	)
	manejaError(err, "No se pudo gestionar manejo de mensajes")

	qEditopo, err := ch.QueueDeclare(
		EDITOPO_QUEUE, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	manejaError(err, "No se pudo gestionar cola de mensajes")

	err = ch.QueueBind(
		qEditopo.Name,     // name
		"#",               //key
		EXCHANGE_BITACORA, //exchange
		false,             //noWait
		nil,               //args
	)
	manejaError(err, "No se pudo gestionar manejo de mensajes")

	messageFallasChannel, err := ch.Consume(
		qEditopo.Name,
		"editopofalla_rmq",
		false,
		false,
		false,
		false,
		nil,
	)
	manejaError(err, "No se ha registrado el consumidor")

	messageEditopoChannel, err := ch.Consume(
		qEditopo.Name,
		"editopo_rmq",
		false,
		false,
		false,
		false,
		nil,
	)
	manejaError(err, "No se ha registrado el consumidor")

	vuelta := make(chan bool)

	go func() {
		for d := range messageFallasChannel {
			buf := bytes.NewReader(d.Body)
			rmqBitacora := &models.RmqBitacoraMensaje{}
			err := binary.Read(buf, binary.LittleEndian, rmqBitacora)

			log.Printf("%+v \n\n", rmqBitacora.Adecuada())
			if err != nil {
				log.Printf("Error reconociendo estructura: %s", err)
			}
		}

		for d := range messageEditopoChannel {
			buf := bytes.NewReader(d.Body)
			rmqBitacora := &models.RmqBitacoraMensaje{}
			err := binary.Read(buf, binary.LittleEndian, rmqBitacora)

			log.Printf("%+v \n\n", rmqBitacora.Adecuada())
			if err != nil {
				log.Printf("Error reconociendo estructura: %s", err)
			}

			//Procesar instrucciones
			if procesarMensaje(rmqBitacora) {
				if err := d.Ack(false); err != nil {
					log.Printf("Error al reconocer mensaje : %s", err)
				}
			} else {
				if err := d.Nack(true, false); err != nil {
					log.Printf("Error al reconocer mensaje : %s", err)
				}
			}

		}
	}()

	log.Printf(" [*] Esperando actividad ... ")
	<-vuelta
}

// Rutina general para procesamiento de mensajes de RabbitMQ, para limpiar
func procesarMensaje(rmq *models.RmqBitacoraMensaje) bool {
	if (rmq.CodigoMsj == 3008) && (rmq.TipoMensaje == 5000) && (rmq.Comando == 7) { //Detener
		//		log.Printf("Bitacora: %+v \n", rmq.Adecuada())
		go procesoDetener()
		return true
	}
	if (rmq.CodigoMsj == 3008) && (rmq.TipoMensaje == 5000) && (rmq.Comando == 6) { //Iniciar
		//		log.Printf("Bitacora: %+v \n", rmq.Adecuada())
		go procesoReinicio()
		return true
	}
	return false
}

// Procesamiento de solicitudes de paro de automatismo

func procesoDetener() {
	var l api.UcmLog
	err := l.Iniciar("rmata", "rmata26") // Se usa esta combinación para usar ciertas APIs, verificar credenciales
	if err != nil {
		log.Fatalf("No jala: %s", err)
	}
	log.Printf("Procesando mensaje de detener \n")
	strMsg := fmt.Sprintf("DETENIENDO AUTOMATISMO")
	s, err := l.EnviarDetenido(cfgA.ExtraeAPI("LOGIN"), cfgA.ExtraeAPI("BITACORA"), cfgA)
	if err == nil {
		log.Printf("[OK] >> %s \n", s)
		ihm.Avisar(strMsg)
		strMsg = fmt.Sprintf("ENVIANDO CONFIRMACION DE AUTOMATISMO DETENIDO")
		ihm.Avisar(strMsg)
	} else {
		log.Printf("Problema enviando mensaje %s \n", s)
		strMsg = fmt.Sprintf("NO SE PUDO DETENER AUTOMATISMO")
		ihm.Problema(strMsg)
	}

}

func procesoReinicio() {
	var l api.UcmLog
	err := l.Iniciar("rmata", "rmata26") // Se usa esta combinación para usar ciertas APIs, verificar credenciales
	if err != nil {
		log.Fatalf("No jala: %s", err)
	}
	log.Printf("Procesando mensaje de detener \n")
	strMsg := fmt.Sprintf("REINICIANDO AUTOMATISMO")
	s, err := l.EnviarDetenido(cfgA.ExtraeAPI("LOGIN"), cfgA.ExtraeAPI("BITACORA"), cfgA)
	if err == nil {
		log.Printf("[OK] >> %s \n", s)
		ihm.Avisar(strMsg)
		strMsg = fmt.Sprintf("ENVIANDO CONFIRMACION DE AUTOMATISMO DETENIDO")
		ihm.Avisar(strMsg)
	} else {
		log.Printf("Problema enviando mensaje %s \n", s)
		strMsg = fmt.Sprintf("NO SE PUDO DETENER AUTOMATISMO")
		ihm.Problema(strMsg)
	}
	strMsg = fmt.Sprintf("REINICIANDO AUTOMATISMO")
	ihm.Avisar(strMsg)
}

func procesarFalla() {
	var l api.UcmLog
	err := l.Iniciar("rmata", "rmata26") // Se usa esta combinación para usar ciertas APIs, verificar credenciales
	if err != nil {
		log.Fatalf("No jala: %s", err)
	}
	log.Printf("Procesando falla \n")
	strMsg := fmt.Sprintf("INICIANDO ATENCION")
	s, err := l.EnviarDetenido(cfgA.ExtraeAPI("LOGIN"), cfgA.ExtraeAPI("BITACORA"), cfgA)
	if err == nil {
		log.Printf("[OK] >> %s \n", s)
		ihm.Avisar(strMsg)
		strMsg = fmt.Sprintf("ENVIANDO NOTIFICACIONES")
		ihm.Avisar(strMsg)
	} else {
		log.Printf("Problema enviando mensaje %s \n", s)
		strMsg = fmt.Sprintf("NO SE PUDO DETENER AUTOMATISMO")
		ihm.Problema(strMsg)
	}
	strMsg = fmt.Sprintf("REINICIANDO AUTOMATISMO")
	ihm.Avisar(strMsg)
}

func lanzarEvento(descripcion string, detalle string, tipo string) error {
	repo := motorapt.NewUcmEventRegRepo(dbApt)
	return repo.InsertEventFromtext(descripcion, detalle, "AYPT", tipo)
}

const executableID = "MOTORAPT"

var (
	Sha1ver    string // clave sha1 del commit usado para compilar el programa
	Branch     string //nombre de la rama usada para compilar el programa
	BuildTime  string // when the executable was built
	flgVersion bool
	dbTop      *gorm.DB
	dbUcm      *gorm.DB
	dbApt      *gorm.DB
)

func initDB() error {
	var err error
	var ambiente config.Ambiente
	dbTop, err = ambiente.GetTopoDb()
	if err != nil {
		return err
	}

	dbUcm, err = ambiente.GetUcmDb()
	if err != nil {
		return err
	}

	dbApt, err = ambiente.GetUcmDb()
	if err != nil {
		return err
	}
	return nil
}

func recoverData() (*graphs.BulkGraph, error) {
	var err error

	oT := motorapt.NewTopologyRepo(dbTop)
	bk := graphs.NewBulkGraph()

	bk.Feeders, err = oT.GetFeeders()
	if err != nil {
		return nil, err
	}
	/*
		fmt.Println("Alimentadores")
		for _, a := range bk.Feeders {
			fmt.Println(a.FeederString())
		}
	*/
	bk.Lines, err = oT.GetFullLines()
	if err != nil {
		return nil, err
	}
	/*
		fmt.Println("Lineas")
		for _, a := range bk.Lines {
			fmt.Println(a.LineString())
		}
	*/
	bk.Nodes, err = oT.GetNodes()
	if err != nil {
		return nil, err
	}
	/*
		fmt.Println("Nodes")
		for _, a := range bk.Lines {
			fmt.Println(a.NodeString())
		}
	*/
	bk.SwLines, err = oT.GetSwLines()
	if err != nil {
		return nil, err
	}
	/*
		fmt.Println("Equipos de comuntación")
		for _, a := range bk.SwLines {
			fmt.Println(a.SwitchString())
		}
	*/
	return bk, nil
}

func printSummary(bulk *graphs.BulkGraph) {
	printColor("Circuitos:", color.FgHiWhite)
	printColor(fmt.Sprintf("%d", len(bulk.Feeders)), color.FgGreen)
	printColor(", Nodos:", color.FgHiWhite)
	printColor(fmt.Sprintf("%d", len(bulk.Nodes)), color.FgGreen)
	printColor(", Líneas:", color.FgHiWhite)
	printColor(fmt.Sprintf("%d", len(bulk.Lines)), color.FgGreen)
	printColor(", Equipos:", color.FgHiWhite)
	printColor(fmt.Sprintf("%d\n", len(bulk.SwLines)), color.FgGreen)
}

func printColor(m string, a color.Attribute) {
	color.Set(a)
	fmt.Printf(m)
	color.Unset()
}

func parseCmdLineFlags() {
	flag.BoolVar(&flgVersion, "version", false, "si true, imprime la versión y termina el programa")
	flag.Parse()

	if flgVersion {
		fmt.Printf("Fecha: %s Rama: %s Commit (sha1): %s\n", BuildTime, Branch, Sha1ver)
		os.Exit(0)
	}
}

func main() {
	// address, err := env.GetMotorAptAddr()
	//parseCmdLineFlags()

	fmt.Println()
	ihm.Letrero("Iniciando motor de automatismo y procesador de topologías")
	ihm.Letrero("UCM-CFE v2.0")

	aptEngine := newAptEngineFSM()
	err := aptEngine.SendEvent(IniciaCarga, nil) //arranca máquina de estados
	if err != nil {
		fmt.Printf("No se pudo inicializar la maquina de estados, err: %v \n", err)
		return
	}

}
