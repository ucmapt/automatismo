package dion

import (
	"fmt"

	"github.com/ucmapt/automatismo/common/graphs"
	"github.com/ucmapt/automatismo/models"

	"gonum.org/v1/gonum/mat"
)

var (
	Gbus mat.VecDense
	Bbus mat.VecDense
)

type ybus struct {
}


func admitancia(g *graphs.UniGraph, a string) (mat.VecDense, mat.VecDense) {
	var G, B mat.VecDense
	G = *mat.NewVecDense(100, make([]float64, 0))
	B = *mat.NewVecDense(100, make([]float64, 0))
	return G,B

}

func clasificacionCircuito(g *graphs.UniGraph, a string) error{
	// obtener incidencias 
	// OJO en las rutinas de MATLAB se calculan a partir de dos archivos csv
	// tomar como base el grafo y el alimentador

	u := mat.NewVecDense(3, []float64{1, 2, 3})
	e := u.AtVec(1)
	if e != 0 {
		return fmt.Errorf("No manejado")
	}
	return nil 
}

// Para el cálculo de las corrientes de cortocircuito en cada nodo del circuito, se requirie estimar 
// los valores complejos de corriente de cortocircuito en cada nodo del circuito, aplicando un cociente complejo
// sobre los datos del flujo de cargas  
func cocienteComplejo(a, b complex128) complex128{
	den := real(b) *real (b) + imag(b) * imag(b)
	return complex(real(a)*real(b)+imag(a)*imag(b)/den,imag(a)*real(b)-real(a)*imag(b)/den)
}

// Manejo de
func formacionYbusZbus(){

}

// Manejo de arrgelo de incidencias a partir de datos del grafo
func arregloIncidencia(n1, n2 graphs.Node) []int {

	// en el algoritmo original, se reconstruía el archivo de incidencias.csv
	// aquí se contruye una estructura dinámica
	return nil 
}

// CASO DE USO: Cálculo de polígono de falla
// FEATURE: El módulo de automatismo y procesador de topologías podrá calcular las posiciones más probables 
// del origen de una falla franca con base en un análisis a partir de la corriente de corto circuito reportada durante el suceso.
// R017 - R021
func AnalizarCortoCircuito(b graphs.BulkGraph, g *graphs.UniGraph, cve string, iff float64) error {
	// Resolver punto de falla
	// Buscar el dispositivo en el grafo para determinar el circuito a usar
	// Ubicar falla
	var swd *models.SwLine
	swd = nil
	for i := 0; i < len(b.SwLines); i++ {
		if b.SwLines[i].Switch == cve {
			swd = b.SwLines[i]
		}
	}
	aux := *swd.Circuito

	// Limitar datos a usar (circuito, nodos, líneas)
	bii := g.GetBii(aux, aux)
	gii := g.GetGii(aux, aux)

	// Construir estructuras auxiliares correspondiente
	// Levantar evento “Calculando perímetro de falla”
	if bii.Status {
		return fmt.Errorf("Probema al gestionar cálculos ")
	}

	if gii.Status {
		return fmt.Errorf("Probema al gestionar cálculos ")
	}
	// Aplicar clasificación por circuito Ybus = Gbus - jBbus
	clasificacionCircuito(g, aux)

	// Formar Ybus y calcular Zbus conforme a grafo a actualizar

	// Proceso iterativo para resolver por Newton Raphson
	// Aplicar corrección voltaje nodo oscilatorio
	// Cálculo de Corto Circuito
	// Reflejar resultados

/*
	H02_formacion_Ybus_calculo_Zbus_datos_iniciales
	k_it = 0;
	norm_b = 100;
	while  k_it < 6 & norm_b > 10
		k_it = k_it + 1;
		H04_Newton_Raphson_proceso_iterativo
	end
	%% 
	
	F05_correccion_voltaje_nodo_oscilatorio_y_calculo_de_CC
	H06_grafo_vs_distribucion_CC
*/

	return nil
}






