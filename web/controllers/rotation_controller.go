package controllers

import (
	"Robotic/datamodels"
	"encoding/json"

	"github.com/kataras/iris"
)

// NcController : controller for managing nc command
type RotationController struct {
	key string
}

//Get : Get index form
func (c *RotationController) Get(ctx iris.Context) {
	R := datamodels.RMatrix{}
	R.Value[0] = [3]float64{0.8660254, -0.5, 0}
	R.Value[1] = [3]float64{0.5, 0.8660254, 0}
	R.Value[2] = [3]float64{0, 0, 1}
	R2 := R.Inverse()
	q, _ := json.Marshal(R.Multiply(&R2))

	ctx.Write(q)

}
