package dion

import (
	"fmt"
	"strings"

	"github.com/ucmapt/automatismo/common/graphs"
	"github.com/ucmapt/automatismo/models"
	"gonum.org/v1/gonum/mat"
)

// FUNCION   zeros
// PARAMETRO rows - Número de columnas a considerar
// PARAMETRO cols - Número de columnas a considerar
// zeros reemplaza la función de MATLAB para crear matrices bidimensionales con ceros
// RETORNO   Matriz
// NOTA      Se asumen n1 y n2 como matrices 2d tipo columna con múltiples filas y una sola columna
// NOTA      Solo para pruebas preliminares, se sustituye por CargaGrafos
func zeros(rows, cols int) *mat.Dense {
	z := make([]float64, rows*cols)
	for i := range z {
		z[i] = 0
	}
	aux := mat.NewDense(rows, cols, z)
	return aux
}

// FUNCION   Admitancia
// PARAMETRO r - Arreglo (slice) de enteros
// PARAMETRO x - Matriz columna con los indices de nodo registrados al recuperar de la BD
// RETORNO   booleano indicando si existe el valor en el arreglo
// NOTA      se sustituye por CargaGrafos
func Admitancia(r, x mat.Dense, n int) (*mat.Dense, *mat.Dense) {
	b := zeros(n, 1)
	g := zeros(n, 1)

	for k := 0; k < n; k++ {
		d := r.At(k, 0)*r.At(k, 0) + x.At(k, 0)*x.At(k, 0)
		g.Set(k, 0, r.At(k, 0)/d)
		b.Set(k, 0, -1*x.At(k, 0)/d)
	}
	return g, b
}

//Estructuras auxiliares para pruebas preliminares de ArregloIncidencia, se sustituyen al manejar GrafoVisor
type AuxIncidencias struct {
	Nodos []NodoIncidencias
}

type NodoIncidencias struct {
	NodoId      int
	Adyacencias []int
}

// FUNCION   arregloIncidencia
// PARAMETRO n1 - Matriz columna con los indices de nodo registrados al recuperar de la BD
// PARAMETRO n2 - Matriz columna con los indices de nodo registrados al recuperar de la BD
// RETORNO   AuxInidencias objeto conteniendo los nodos con sus índices de incidencia
// RETORNO   error objeto de error
func arregloIncidencia(n1, n2 mat.Dense) (AuxIncidencias, error) {
	//formar una lista de nodos que incluya los n1 y n2 sin repetir
	ai := AuxIncidencias{
		Nodos: []NodoIncidencias{},
	}
	ns := []int{}
	ns1 := []int{}
	ns2 := []int{}

	rows, cols := n1.Dims()
	if cols != 1 {
		e := fmt.Errorf("Lista (nodo1) no cumple con formato")
		return AuxIncidencias{}, e
	}

	for i := 0; i < rows; i++ {
		ns = append(ns, int(n1.At(i, 0)))
	}
	ns1 = append(ns1, ns...)

	rows, cols = n2.Dims()
	if cols != 1 {
		e := fmt.Errorf("Lista (nodo2) no cumple con formato")
		return AuxIncidencias{}, e
	}

	for i := 0; i < rows; i++ {
		aux := int(n2.At(i, 0))
		if !contiene(ns, aux) {
			ns = append(ns, aux)
		}
		ns2 = append(ns2, aux)
	}

	t := len(ns)
	for i := 0; i < t; i++ {
		adjs := []int{}
		for j := 0; j < len(ns1); j++ {
			if (ns[i] == ns1[j]) || (ns[i] == ns2[j]) {
				adjs = append(adjs, j)
			}
		}
		nodAux := NodoIncidencias{
			NodoId:      i,
			Adyacencias: adjs,
		}
		ai.Nodos = append(ai.Nodos, nodAux)
	}

	return ai, nil
}

// FUNCION   cocienteComplejo
// PARAMETRO a - Matriz renglón de dos columnas representando un dato complejo
// PARAMETRO b - Matriz renglón de dos columnas representando un dato complejo
// RETORNO   Matriz renglón de dos columnas representando un dato complejo con el cociente complejo de los parámetros de entrada
func cocienteComplejo(a, b mat.Dense) mat.Dense {
	c := *zeros(1, 2)
	den := (b.At(0, 0) * b.At(0, 0)) + (b.At(0, 1) * b.At(0, 1))
	c.Set(0, 0, (a.At(0, 0)*b.At(0, 0))+(a.At(0, 1)*b.At(0, 1))/den)
	c.Set(0, 1, (a.At(0, 1)*b.At(0, 0))+(a.At(0, 0)*b.At(0, 1))/den)
	return c
}

// FUNCION   cocienteComplex
// PARAMETRO a - Valor complejo
// PARAMETRO b - Valor complejo
// RETORNO   valor complejo con el cociente complejo(A,B)
func cocienteComplex(a, b complex128) complex128 {
	den := real(b)*real(b) + imag(b)*imag(b)
	c := complex((real(a)*real(b)+imag(a)*imag(b))/den, (imag(a)*real(b)-real(a)*imag(b))/den)
	return c
}

// FUNCION   contine
// PARAMETRO l - Arreglo (slice) de enteros
// PARAMETRO valor - Matriz columna con los indices de nodo registrados al recuperar de la BD
// RETORNO   booleano indicando si existe el valor en el arreglo
func contiene(l []int, valor int) bool {
	for _, e := range l {
		if e == valor {
			return true
		}
	}
	return false
}

func extraeMatrizBanda() {

}

// FUNCION   ShowMatrix
// PARAMETRO x - Matriz de flotantes
// RETORNO   cadena mostrando la matriz formateada
func ShowMatrix(x *mat.Dense) string {
	var sb strings.Builder
	r, c := x.Dims()
	for i := 0; i < r; i++ {
		sb.WriteString("[")
		for j := 0; j < c; j++ {
			if j != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%f", x.At(i, j)))
		}
		sb.WriteString(fmt.Sprintf("]\n"))
	}
	return sb.String()
}

// Lectura y acondicionamiento inicial de datos
// La versión del Dr. Dionisio partía de recuperar los circuitos en formato matricial,
// aquí se acopla desde un grafo que representa un circuito dinámico extraido

// num_linea,nom_linea,num_circ,longitud,r1,x1,r0,x0
func preparacionInicial(g graphs.UniGraph) (*mat.Dense, *mat.Dense) {
	Gbus := zeros(0, 0)
	Bbus := zeros(0, 0)

	return Gbus, Bbus
}

// Formación de matrices Gbus y Bbus

// AnalisisCortoCircuito recibe Grafo, ICC y elemento qque lo siente
func AnalisisCortoCircuito(g graphs.UniGraph, iff float64, elem string) (models.FailResult, models.FailZone, error) {
	Gbus, Bbus := preparacionInicial(g)
	if Gbus == nil {

	}
	if Bbus == nil {

	}
	return 0, models.FailZone{}, nil
}

/*
Lectura y acondicionamiento inicial de datos
Formación de matrices Gbus y Bbus
Formación de matrices Gbus y Bbus de admitancia
Elaboración y búsqueda en grafo
Reordenamiento de matrices Gbus y Bbus y extracción de matriz banda
Factorización compleja LU datos en banda
Inversa compleja para cálculo Rbus y Xbus
Diagonales de Rbus y Xbus
Elaboración de matrices G y B de admitancia
Elaboración de matrices Gcga y Bcga
Factorización compleja LU cálculo VoBus
Solución de VoBus
Newton Raphson
Corrección de cálculo de nodo oscilatorio y equivalencia Icc
Termina proceso






H01_Clasificacion_por_circuito
H02_formacion_Ybus_calculo_Zbus_datos_iniciales
%%
k_it = 0;
norm_b = 100;
while  k_it < 6 & norm_b > 10
    k_it = k_it + 1;
    H04_Newton_Raphson_proceso_iterativo
end
%%

F05_correccion_voltaje_nodo_oscilatorio_y_calculo_de_CC

DT = [DT;dT];
H06_grafo_vs_distribucion_CC

feval('print',['localizacion_cc_circuito_',num2str(num_foto),'_Iff2250.jpg'],'-djpeg')

end
*/
