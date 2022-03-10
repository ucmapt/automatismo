package graphs

import (
	"bufio"
	"fmt"
	"os"
)

// Graph structure
type UniGraph struct {
	nodes     []*Node
	lines     []*Line
	feeders   []string
	adjacency [][]bool
	circuits  []*Circuit
	loops     [][]string
}

// NewGraph crea un grafo vacío
func NewGraph() *UniGraph {
	g := UniGraph{
		nodes:     []*Node{},
		lines:     []*Line{},
		feeders:   make([]string, 0),
		adjacency: make([][]bool, 0),
		circuits:  []*Circuit{},
		loops:     make([][]string, 0),
	}
	return &g
}

// Estructura para manejo de incidencias en lugar de
// archivos CSV
type IncidenciaCSV struct {
	nodo   Node
	lineas []Line
}


// Inicializa estructura sustituta para manejo de incidencias
func NewIncidenciaCSV() *IncidenciaCSV {
	ai := IncidenciaCSV{
		nodo:   Node{},
		lineas: []Line{},
	}
	return &ai
}

// Carga arreglo de incidencias desde el Grafo
func LoadIncidenciasCSV(cve string) *IncidenciaCSV {
	return NewIncidenciaCSV()
}

// Circuit structure
type Circuit struct {
	key     string
	keynode string
	nominal float64
	nodes   []string
	lines   []string
}

// Node structure
type Node struct {
	key      string
	circuito string
	adjacent []*Node
}

// Line structure
type Line struct {
	Key      string
	From     string
	To       string
	Sw1      bool
	Sw2      bool
	Status   bool
	Ltype    string
	X0       float64
	X1       float64
	C0       float64
	C1       float64
	R0       float64
	R1       float64
	L        float64
	Circuito string
}

// Add Feeder solo para pruebas iniciales ....
func (g *UniGraph) AddFeederBasic(k string) error {
	if contiene(g.feeders, k) {
		err := fmt.Errorf("Alimentador %v no fue agregado pues ya existe la clave", k)
		return err
	} else if g.getNode(k) == nil {
		err := fmt.Errorf("Alimentador %v no fue agregado no han agregado el nodo correspondiente", k)
		return err
	} else {
		g.feeders = append(g.feeders, k)
	}
	return nil
}

func (g *UniGraph) GetNode(cve string) *Node {
	return g.getNode(cve)
}

func (g *UniGraph) GetLine(from, to string) *Line {
	return g.getLine(from, to)
}

func (g *UniGraph) GetBii(from, to string) *Line {
	return g.getLine(from, to)
}

func (g *UniGraph) GetGii(from, to string) *Line {
	return g.getLine(from, to)
}

func (g *UniGraph) SolveGrafo(from, to string) *Line {
	return g.getLine(from, to)
}

func (g *UniGraph) Simplificar(m interface{}) *Line {
	return nil
}

func (g *UniGraph) GetNodeLines(cve string) []*Line {
	ls := []*Line{}
	n1 := g.getNodeId(cve)
	for n2 := 0; n2 < len(g.adjacency[n1]); n2++ {
		if g.adjacency[n1][n2] {
			l := g.getLine(g.nodes[n1].key, g.nodes[n2].key)
			ls = append(ls, l)
		}
	}
	return ls
}

func (g *UniGraph) AddFullFeeder(k string, keynode string, nominal float64) error {
	if contiene(g.feeders, k) {
		err := fmt.Errorf("Alimentador %v no fue agregado pues ya existe la clave", k)
		return err
	} else if g.getNode(k) == nil {
		err := fmt.Errorf("Alimentador %v no fue agregado no han agregado el nodo correspondiente", k)
		return err
	} else {
		g.feeders = append(g.feeders, k)
	}
	return nil
}

// AddNode agrega un nodo al grafo
func (g *UniGraph) AddNode(k string) error {
	if contains(g.nodes, k) {
		err := fmt.Errorf("Nodo %v no fue agregado pues existe la clave", k)
		// fmt.Println(err.Error())
		return err
	} else {
		g.nodes = append(g.nodes, &Node{key: k})
		l := len(g.nodes)
		mat := make([][]bool, 0)
		for i := 0; i < l-1; i++ {
			row := append(g.adjacency[i], false)
			mat = append(mat, row)
		}
		g.adjacency = append(mat, make([]bool, l))
	}
	return nil
}

// AddLineBasic agrega una línea al grafo solo para pruebas iniciales ....
func (g *UniGraph) AddLineBasic(from, to string) error {
	// get nodes
	fromNode := g.getNode(from)
	toNode := g.getNode(to)
	// check errors
	if fromNode == nil || toNode == nil {
		err := fmt.Errorf("Linea invalida (%v-->%v)", from, to)
		// fmt.Println(err.Error())
		return err
	} else if contains(fromNode.adjacent, to) {
		err := fmt.Errorf("Linea existente (%v-->%v)", from, to)
		// fmt.Println(err.Error())
		return err
	} else {
		// add line
		fromNode.adjacent = append(fromNode.adjacent, toNode)
		toNode.adjacent = append(toNode.adjacent, fromNode)
		g.adjacency[g.getNodeId(from)][g.getNodeId(to)] = true
		g.adjacency[g.getNodeId(to)][g.getNodeId(from)] = true
	}
	return nil
}


// Para realizar las operaciones eléctricas se requiere líneas con la información completa
// Las pruebas de recorridos originalmente podían omitir TODA la información, para modelos más útiles se validará

// AddLineFull agrega una línea al grafo con toda la información eléctrica incluída en la DB
func (g *UniGraph) AddLineFull(key, from, to string, sw1, sw2, status bool, x0, x1, c0, c1, r0, r1 float64) error {
	// get nodes
	fromNode := g.getNode(from)
	toNode := g.getNode(to)
	// check errors
	if fromNode == nil || toNode == nil {
		err := fmt.Errorf("Linea invalida (%v-->%v)", from, to)
		// fmt.Println(err.Error())
		return err
	} else if contains(fromNode.adjacent, to) {
		err := fmt.Errorf("Linea existente (%v-->%v)", from, to)
		// fmt.Println(err.Error())
		return err
	} else {
		// add line
		fromNode.adjacent = append(fromNode.adjacent, toNode)
		toNode.adjacent = append(toNode.adjacent, fromNode)
		g.adjacency[g.getNodeId(from)][g.getNodeId(to)] = true
		g.adjacency[g.getNodeId(to)][g.getNodeId(from)] = true
		lineAux := Line{
			Key:      "",
			From:     from,
			To:       to,
			Sw1:      sw1,
			Sw2:      sw2,
			Status:   status,
			Ltype:    "",
			X0:       x0,
			X1:       x1,
			C0:       c0,
			C1:       c1,
			R0:       r0,
			R1:       r1,
			L:        0,
			Circuito: "",
		}
		g.lines = append(g.lines, &lineAux)
	}
	return nil
}

// Recuperar el nodo cuya clave se suministre 
// si no se encuentra la clave, entonces 
// getNode devuelve nulo
func (g *UniGraph) getNode(k string) *Node {
	for i, n := range g.nodes {
		if n.key == k {
			return g.nodes[i]
		}
	}
	return nil
}

// Recuperar el índice del nodo corresponda la clave se suministra
// si no se encuentra la clave, entonces
// getNodeId devuelve -1
func (g *UniGraph) getNodeId(k string) int {
	for i, n := range g.nodes {
		if n.key == k {
			return i
		}
	}
	return -1
}

// Similar a los nodos, las líneas requieren manipulación directa por línea o por índice
// getLine recupera la línea desde las claves de los nodos
// esta función devueve la primera línea, en caso de redundancia, se ignora ya que eléctricamente no 
func (g *UniGraph) getLine(from, to string) *Line {
	for _, l := range g.lines {
		if l.From == from && l.To == to {
			return l
		}
	}
	return nil
}

// contains
func contains(s []*Node, k string) bool {
	for _, n := range s {
		if k == n.key {
			return true
		}
	}
	return false
}

// Print imprime el contenido de los nodos adyacentes a cada nodo del grafo
func (g UniGraph) Print() {
	for _, n := range g.nodes {
		fmt.Printf("\nNodo %v [%s]: ", n.key, n.circuito)
		for _, a := range n.adjacent {
			fmt.Printf(" %v ", a.key)
		}
	}
	fmt.Printf("\nMatriz de adyacencia\n")
	for _, r := range g.adjacency {
		s := ""
		for _, x := range r {
			s = s + "[" + markBool(x) + "]"
		}
		fmt.Printf("%s\n", s)
	}
	fmt.Printf("Ciclos encontrados: %d \n", len(g.loops))
	for i, l := range g.loops {
		fmt.Printf("Ciclo %d: ", i+1)
		marka := ""
		for _, s := range l {
			fmt.Printf("%s%s", marka, s)
			marka = "-"
		}
		fmt.Println()
	}
}

// hasLine
func (g *UniGraph) hasLine(from, to string) (bool, error) {
	// get nodes
	fromNodeId := g.getNodeId(from)
	toNodeId := g.getNodeId(to)
	// check errors
	if fromNodeId == -1 || toNodeId == -1 {
		err := fmt.Errorf("Linea invalida (%v-->%v)", from, to)
		//fmt.Println(err.Error())
		return false, err
	} else if g.adjacency[fromNodeId][toNodeId] {
		return true, nil
	}
	return false, nil
}

func (g *UniGraph) DepthFirstSearch(cve string) []string {
	var q []string = []string{} //queue
	var explored []string = []string{}

	root := g.getNode(cve)

	if root != nil {
		explored = append(explored, cve)
		q = append(q, root.key)

		for len(q) > 0 { //while q is not empty
			// dequeue
			v := g.getNode(q[0])
			q = q[1:]
			for _, adj := range v.adjacent {
				if !contiene(explored, adj.key) {
					explored = append(explored, adj.key)
					q = append(q, adj.key) //enqueue
				}
			}
		}
	}
	return explored
}

func (g *UniGraph) getDownstreamPaths(from string) [][]string {
	down := [][]string{}
	return down
}

func (g *UniGraph) getDownstreamGraph(from string) *UniGraph {
	return nil
}

func markBool(b bool) string {
	if b {
		return "X"
	}
	return " "
}

type Section struct {
	nodes []string
	lines []string
}

func (g *UniGraph) colorSectionDebug(n *Node, visited []string, c string) {
	// marcar nodo con el número de alimentador
	n.circuito = c
	ignore := false
	// incluir nodo en visitadas
	visited = append(visited, n.key)
	/*
		fmt.Printf("Adyacencia a: %s %d\n", n.key, len(n.adjacent))
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	*/
	for _, a := range n.adjacent {
		//verificar estado de la línea -----FALTA IMPLEMENTAR
		if len(visited) > 1 {
			ignore = visited[len(visited)-2] == a.key // ignorar regresos instantaneos
		}
		if !ignore {
			// si visitadas contiene adyacencia, obtener loop, incluir en loops
			if contiene(visited, a.key) {
				fmt.Printf("Procesando loop entre: %s y %s\n", n.key, a.key)
				fmt.Printf("VISITADOS: %v\n", visited)
				bufio.NewReader(os.Stdin).ReadBytes('\n')
				i := len(visited) - 1
				aux := []string{a.key}
				for i >= 0 {
					if visited[i] != a.key {
						aux = append(aux, visited[i])
					} else {
						aux = append(aux, visited[i])
						if !incluye(g.loops, aux) && !incluye(g.loops, invierte(aux)) {
							g.loops = append(g.loops, aux)
						}
						i = 0
					}
					i = i - 1
				}
			} else { //sino,llamar colorearNodo(adyacencia)
				g.colorSectionDebug(a, visited, c)
			}
		}
	}
}

func (g *UniGraph) colorSection(n *Node, visited []string, c string) {
	// marcar nodo con el número de alimentador
	n.circuito = c
	ultimo := ""
	if len(visited) > 0 {
		ultimo = visited[len(visited)-1]
	}

	// incluir nodo en visitadas
	visited = append(visited, n.key)

	// por cada adyacencia
	for _, a := range n.adjacent {
		// verificar estado de la línea -----FALTA IMPLEMENTAR

		if ultimo != "" && ultimo == a.key { // ignorar regresos instantaneos

		} else {
			// si visitadas contiene adyacencia, obtener loop, incluir en loops
			if contiene(visited, a.key) {
				i := len(visited) - 1
				aux := []string{a.key}
				for i >= 0 {
					if visited[i] != a.key {
						aux = append(aux, visited[i])
					} else {
						aux = append(aux, visited[i])
						if !incluye(g.loops, aux) && !incluye(g.loops, invierte(aux)) {
							g.loops = append(g.loops, aux)
						}
						i = 0
					}
					i--
				}
			} else { //sino,llamar colorearNodo(adyacencia)
				g.colorSection(a, visited, c)
			}
		}
	}
}

func contiene(a []string, c string) bool {
	for _, s := range a {
		if s == c {
			return true
		}
	}
	return false
}

func incluye(a [][]string, c []string) bool {
	f := false
	for _, b := range a {
		if len(b) == len(c) {
			f = true
			for i := range b {
				if b[i] != c[i] {
					f = false
					break
				}
			}
			if f {
				return true
			}
		}
	}
	return f
}

func invierte(ab []string) []string {
	ba := []string{}

	for _, s := range ab {
		ba = append([]string{s}, ba...)
	}
	//fmt.Printf("%v", ab)
	return ba
}

func (g *UniGraph) Colorize() {
	// marca todo OFF
	for _, n := range g.nodes {
		(*n).circuito = "OFF"
	}
	for _, l := range g.lines {
		(*l).Circuito = "OFF"
	}

	// loops
	g.loops = make([][]string, 0)

	for _, f := range g.feeders {
		n := g.getNode(f)
		v := []string{}
		g.colorSection(n, v, f)
	}
	// tramos con ciclo
	// tramos desconectados
	// alimentadores desconectados
}

func (g *UniGraph) UpdateViews() {

	// CONCATENAR COMANDOS

	// update postgis.graphic_data SET resultado = null;

	// Todos los Circuito = OFF
	// update postgis.graphic_line SET desconectado = TRUE where name IN ('24PAR04030');

	// Refrescar resultado DION
	// update postgis.graphic_data SET resultado = 'FALLA PERMANENTE'
	// where name IN ('DC24PAR04030np1878152');

	// Por cada circuito
	// update postgis.graphic_line SET desconectado = FALSE where name IN ('24PAR04030');
	// update postgis.graphic_line SET circuito_dinamico = '24PAR04060' where name IN ();

}
