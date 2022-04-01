package dion

import (
	"fmt"
	"strings"

	"gonum.org/v1/gonum/mat"
)

// Reemplaza la función de MATLAB para crear matrices bidimensionales con ceros
func zeros(x, y int) *mat.Dense {
	z := make([]float64, x*y)
	for i := range z {
		z[i] = 0
	}
	aux := mat.NewDense(x, y, z)
	return aux
}

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

// FUNCION   ArregloIncidencia
// PARAMETRO n1 - Matriz columna con los indices de nodo registrados al recuperar de la BD
// PARAMETRO n2 - Matriz columna con los indices de nodo registrados al recuperar de la BD
// NOTA      Se asumen n1 y n2 como matrices 2d tipo columna con múltiples filas y una sola columna
// NOTA      Solo para pruebas preliminares, se sustituye por CargaGrafos
func ArregloIncidencia(n1, n2 mat.Dense) (AuxIncidencias, error) {
	//formar una lista de nodos que incluya los n1 y n2 sin repetir

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
		ns2= append(ns2, aux)
	}

	t := len(ns)
	for i:= 0;i<t;i++ {
		adjs := []int{}

		if !contiene(adjs, ns[j]){
			
		}
	}
	}

	return AuxIncidencias{}, nil
}

// Función auxiliar
func contiene(l []int, valor int) bool {
	for _, e := range l {
		if e == valor {
			return true
		}
	}
	return false
}

// Solo para pruebas preliminares, eliminar en versión final
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
