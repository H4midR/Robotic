package controllers

import (
	"Robotic/datamodels"
	"Robotic/helperfunctions"
	"encoding/json"
	"math"

	"github.com/kataras/iris"
)

// NcController : controller for managing nc command
type RotationController struct {
	key string
}

//Get : Get index form
func (c *RotationController) Get(ctx iris.Context) {

	//------------------------------------					sample of definition of a Rotation matrix and calculate it's reverse and it's product
	// R := datamodels.RMatrix{}
	// R.Value[0] = [3]float64{0.8660254, -0.5, 0}
	// R.Value[1] = [3]float64{0.5, 0.8660254, 0}
	// R.Value[2] = [3]float64{0, 0, 1}
	// R2 := R.Inverse()
	// q, _ := json.Marshal(R.Multiply(&R2))
	// ctx.Write(q)

	// test the data for the first case        home work , step 1&2
	R := datamodels.Rotation{}
	pi := math.Pi
	R.FixedAngles(pi/3, pi/4, pi/6)
	ctx.WriteString("\n\n<h1> calculete Rotation Matix step 1 , Fixed Angles</h1>")
	q, _ := json.Marshal(R)
	ctx.Write(q)
	R2 := datamodels.Rotation{}
	R2.RMatrix = R.RMatrix
	R2.InitFixedAnglesInverse()
	ctx.WriteString("\n\n<h1> calculete angles from Rotation Matix step 2 Inverse Fixed Angles</h1>")
	q, _ = json.Marshal(R2)
	ctx.Write(q)

}
func (c *RotationController) GetHomework(ctx iris.Context) {
	var req map[string]string
	ctx.ReadForm(&req)
	q, _ := json.Marshal(req)
	ctx.WriteString("\n\n<h1> request was :</h1>")
	ctx.Write(q)

	R := datamodels.Rotation{}
	//reading angles
	alpha := helperfunctions.Rad2Deg(helperfunctions.Str2Num(req["alpha"]))
	betta := helperfunctions.Rad2Deg(helperfunctions.Str2Num(req["betta"]))
	gamma := helperfunctions.Rad2Deg(helperfunctions.Str2Num(req["gamma"]))
	//seting angles and calculate the rotation matirx
	R.FixedAngles(gamma, betta, alpha)
	//writing data
	ctx.WriteString("\n\n<h1> calculete Rotation Matix step 1 , Fixed Angles</h1>")
	q, _ = json.Marshal(R)
	ctx.Write(q)
	//create new rotation
	R2 := datamodels.Rotation{}
	//setting rotation matirx
	R2.RMatrix = R.RMatrix
	//calculate the angles from fixed angles inverse way
	R2.InitFixedAnglesInverse()
	//writing data
	ctx.WriteString("\n\n<h1> calculete angles from Rotation Matix step 2 Inverse Fixed Angles</h1>")
	q, _ = json.Marshal(R2)
	ctx.Write(q)
}
