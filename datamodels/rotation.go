package datamodels

import "math"

//RMatrix : Rotation Matrix
type RMatrix struct {
	Value [3][3]float64 `json:"Value"`
}

//Rotation : Rotation
type Rotation struct {
	Alpha   float64 `json:"Alpha,omitempty"`   // alpha angle around z axis
	Betta   float64 `json:"Betta,omitempty"`   // betta angle around y axis
	Gamma   float64 `json:"Gamma,omitempty"`   // gamma angle around x axis
	Axis    Vector  `json:"Axis,omitempty"`    // axis of rotation
	Phi     float64 `json:"Phi,omitempty"`     // angle of rotation around the axis
	RMatrix RMatrix `json:"RMatrix,omitempty"` // rotation matrix
}

//--------------------------------------------------------------------------------------------------#
//										RMatrix														||
//__________________________________________________________________________________________________#

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

//Multiply : Multiply to Matix A.Multiply(&B) ~ A x B
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

//--------------------------------------------------------------------------------------------------#
//										Rotation														||
//__________________________________________________________________________________________________#

//FixedAngles : calculate the Rotation Matrix from alpha,betta,gamma in the fixed-angles way
func (r *Rotation) FixedAngles(gamma float64, betta float64, alpha float64) Rotation {
	r.Alpha = alpha
	r.Betta = betta
	r.Gamma = gamma
	r.InitFixedAngles()
	return *r
}

//InitFixedAngles : calculate the Rotation Matrix from alpha,betta,gamma in the fixed-angles way , by it's preset parameters
func (r *Rotation) InitFixedAngles() RMatrix {
	//Flag := "XYZ"
	res := RMatrix{}
	c := math.Cos
	s := math.Sin

	res.Value[0][0] = c(r.Alpha) * c(r.Betta)
	res.Value[0][1] = c(r.Alpha)*s(r.Betta)*s(r.Gamma) - s(r.Alpha)*c(r.Gamma)
	res.Value[0][2] = c(r.Alpha)*s(r.Betta)*c(r.Gamma) + s(r.Alpha)*s(r.Gamma)

	res.Value[1][0] = s(r.Alpha) * c(r.Betta)
	res.Value[1][1] = s(r.Alpha)*s(r.Betta)*s(r.Gamma) + c(r.Alpha)*c(r.Gamma)
	res.Value[1][2] = s(r.Alpha)*s(r.Betta)*c(r.Gamma) - c(r.Alpha)*s(r.Gamma)

	res.Value[2][0] = -s(r.Betta)
	res.Value[2][1] = c(r.Betta) * s(r.Gamma)
	res.Value[2][2] = c(r.Betta) * c(r.Gamma)

	r.RMatrix = res
	return res
}

//FixedAnglesInverse : calculate the alpha , betta and gamma from the rotation Matrix in the inverse fixed angles way
func (r *Rotation) FixedAnglesInverse(m *RMatrix) Rotation {

	return *r
}

//InitFixedAnglesInverse : calculate the alpha , betta and gamma from the rotation Matrix in the inverse fixed angles way , by it's preset parameters
func (r *Rotation) InitFixedAnglesInverse() (gamma float64, betta float64, alpha float64) {
	r.Betta = math.Atan2(-r.RMatrix.Value[2][0], math.Sqrt(math.Pow(r.RMatrix.Value[0][0], 2)+math.Pow(r.RMatrix.Value[1][0], 2)))
	cB := math.Cos(r.Betta)
	r.Alpha = math.Atan2(r.RMatrix.Value[1][0]/cB, r.RMatrix.Value[0][0]/cB)
	r.Gamma = math.Atan2(r.RMatrix.Value[2][1]/cB, r.RMatrix.Value[2][2]/cB)
	return r.Gamma, r.Betta, r.Alpha
}
