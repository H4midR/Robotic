package datamodels

//RMatrix : Rotation Matrix
type RMatrix struct {
	Value [3][3]float64 `json:"Value"`
}

//Rotation : Rotation
type Rotation struct {
	Alpha   float64 `json:"Alpha,omitempty"`   // alpha angle around x axis
	Betta   float64 `json:"Betta,omitempty"`   // betta angle around y axis
	Gamma   float64 `json:"Gamma,omitempty"`   // gamma angle around x axis
	Axis    Vector  `json:"Axis,omitempty"`    // axis of rotation
	Phi     float64 `json:"Phi,omitempty"`     // angle of rotation around the axis
	RMatrix RMatrix `json:"RMatrix,omitempty"` // rotation matrix
}

//Inverse : inverse matrix
func (m *RMatrix) Inverse() RMatrix {
	return m.Transpose()
}

//Transpose : transpose matrix
func (m *RMatrix) Transpose() RMatrix {
	res := RMatrix{}
	for i, element := range m.Value {
		for j, e := range element {
			res.Value[j][i] = e
		}
	}
	return res
}

//Multiply : Multiply to Matix A.Multiply(B) ~ A x B
func (m *RMatrix) Multiply(m2 *RMatrix) RMatrix {
	res := RMatrix{}
	for i, element := range res.Value {
		for j := range element {
			//res[i][j]=
			var value float64
			value = 0
			lenght := len(element)
			for k := 0; k < lenght; k++ {
				value = value + m.Value[i][k]*m2.Value[k][j]
			}
			res.Value[i][j] = value
		}
	}
	return res
}
