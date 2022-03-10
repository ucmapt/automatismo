package matrixes

type Matrix2d struct {
	Dimension [2]int
	CrudeData []float64
}

// constructor matriz 2d
func NewMatrix2d(rows int, cols int) Matrix2d {
	var aux Matrix2d
	aux = Matrix2d{
		Dimension: [2]int{rows, cols},
		CrudeData: make([]float64, rows*cols),
	}
	return aux
}


// eqivalente a matriz[row][col] basada en arreglo
func (m Matrix2d) Data(row int, col int) float64 { //zero base location
	var base int = m.Dimension[0]
	return m.CrudeData[row*base+col]
}

// copia clonada de matriz 2d
func (m Matrix2d) Clone() Matrix2d {
	var aux Matrix2d
	aux = Matrix2d{
		Dimension: [2]int{m.Dimension[0], m.Dimension[1]},
		CrudeData: make([]float64, len(m.CrudeData)),
	}
	copy(aux.CrudeData, m.CrudeData)
	return aux
}

// suma de matrices 2d
func (m Matrix2d) Sum(m2 Matrix2d) Matrix2d {
	var aux Matrix2d
	// validate same dimensions
	if m.Dimension[0] == m2.Dimension[0] && m.Dimension[1] == m2.Dimension[1] {
		aux = Matrix2d{
			Dimension: [2]int{m.Dimension[0], m.Dimension[1]},
			CrudeData: make([]float64, len(m.CrudeData)),
		}
		for i:= range m.CrudeData {
			aux.CrudeData[i] = m.CrudeData[i] + m2.CrudeData[i]
		}
	} else{
		aux = Matrix2d{
			Dimension: [2]int{0, 0},
			CrudeData: make([]float64, 0),
		}
	}

	return aux 
}


// resta de matrices 2d
func (m Matrix2d) Minus(m2 Matrix2d) Matrix2d {
	var aux Matrix2d
	// validate same dimensions
	if m.Dimension[0] == m2.Dimension[0] && m.Dimension[1] == m2.Dimension[1] {
		aux = Matrix2d{
			Dimension: [2]int{m.Dimension[0], m.Dimension[1]},
			CrudeData: make([]float64, len(m.CrudeData)),
		}
		for i:= range m.CrudeData {
			aux.CrudeData[i] = m.CrudeData[i] - m2.CrudeData[i]
		}
	} else{
		aux = Matrix2d{
			Dimension: [2]int{0, 0},
			CrudeData: make([]float64, 0),
		}
	}

	return aux 
}
