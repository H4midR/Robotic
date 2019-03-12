package datamodels

//Frame : frame
type Frame struct {
	Name           string         `json:"Name,omitempty"`           //Name of Frame
	Origin         Point          `json:"Origin,omitempty"`         //Origin Of frame
	X              Vector         `json:"X,omitempty"`              //X axis Of frame in Ref Frame
	Y              Vector         `json:"Y,omitempty"`              //Y axis of frame in Ref Frame
	Z              Vector         `json:"Z,omitempty"`              //Z axis of frame in Ref Frame
	Ref            interface{}    `json:"Ref,omitempty"`            //Refrence Frame
	Rotation       Rotation       `json:"Rotation,omitempty"`       //Rotation by Ref Frame
	Transformation Transformation `json:"Transformation,omitempty"` //Transformation by Ref Frame
}
