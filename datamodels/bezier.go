package datamodels

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/kataras/iris"
)

const ZeroLimit = 0.025
const ConfPath = ".conf"

var lock sync.Mutex
var Marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
var Unmarshal = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

//Config : structures and parameters configuration
type Config struct {
	XPitch              float32 `json:"XPitch"`
	XStepPerRound       uint32  `json:"XStepPerRound"`
	YPitch              float32 `json:"YPitch"`
	YStepPerRound       uint32  `json:"YStepPerRound"`
	ZPitch              float32 `json:"ZPitch"`
	ZStepPerRound       uint32  `json:"ZStepPerRound"`
	ErrFactor           float64 `json:"ErrFactor"`
	LengthCalResolotion uint32  `json:"LengthCalResolotion"`
}

//Save : save the config
func (con *Config) Save() error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(ConfPath)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := Marshal(con)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

func (con *Config) Load() error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Open(ConfPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return Unmarshal(f, con)
}

// CPoints : cpoints
type CPoints struct {
	Curve      string  `json:"curve,omitempty"`
	Req        string  `json:"req,omitempty"`
	Points     []Point `json:"points,omitempty"`
	Resolotion int     `json:"resolotion,omitempty"`
	Rapidity   float64 `json:"rapidity,omitempty"`
}

// Bezier : bezier struct
type Bezier struct {
	CP     []Point               `json:"cp"`
	QP     []Point               `json:"qp,omitempty"`
	BF     []BernsteinPolynomial `json:"basicFunctions"`
	BDifF  []BernsteinPolynomial `json:"basicDiffFunctions"`
	U      float64               `json:"u"`
	N      uint64                `json:"n"`
	Length float64               `json:"length"`
	Conf   Config                `json:"config"`
	//ErrCof float64               `json:"errCof"`
}

//BernsteinPolynomial : bezier basic functions
type BernsteinPolynomial struct {
	Pu  Polynomial `json:"p(u)"`
	Pui Polynomial `json:"p(u-1)"`
	N   uint64     `json:"n"`
	I   uint64     `json:"i"`
}

///----------------------------------------------------------------helper function
// JSON : return a string json format of polynominal
func (b *Bezier) JSON() []byte {
	res, _ := json.Marshal(b)
	return res
}

//FactorialMemoization
func Factorial(n uint64) (res uint64) {

	if n > 0 {
		res = n * Factorial(n-1)
		return res
	}

	return 1
}

//check : cheack error state
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//delayCalculation : calculate the delay time with limitation
func delayCalculation(val float64) (int64, bool) {
	//return 0, 2500 / val
	if val > -1*ZeroLimit && val < ZeroLimit {
		return 100000000, false
	} else {
		return int64(2500000 / val), true
	}
}

///----------------------------------------------------------------Bernstein Polinominal
// Init : initialize the Bernstein Polynomial
func (bp *BernsteinPolynomial) Init() {
	var PuCof = make([]float64, bp.N+1)
	var PuiCof = make([]float64, bp.N+1)
	PuCof[bp.I] = float64(Factorial(bp.N) / (Factorial(bp.I) * Factorial(bp.N-bp.I)))
	PuiCof[bp.N-bp.I] = 1

	bp.Pu = Polynomial{Coefficient: PuCof}
	bp.Pui = Polynomial{Coefficient: PuiCof}
}

// Cal : calculate the value of Bu,i(u)
func (bp *BernsteinPolynomial) Cal(u float64) float64 {
	return bp.Pu.Cal(u) * bp.Pui.Cal(1-u)
}

///----------------------------------------------------------------Bezier
//Init : Init bezier
func (b *Bezier) Init(ctx iris.Context) error {
	b.N = uint64(len(b.CP) - 1)
	if b.N <= 0 {
		ctx.Writef("{\"Message\":\"control points must be more than 1\",\"status\":\"error\"}")
		return errors.New("control points must be more than 1")
	}
	b.BF = make([]BernsteinPolynomial, b.N+1)
	b.QP = make([]Point, b.N)
	b.BDifF = make([]BernsteinPolynomial, b.N)
	var i uint64
	for i = 0; i <= b.N; i++ {
		b.BF[i] = BernsteinPolynomial{N: b.N, I: i}
		b.BF[i].Init()

		//		calculate Q[i]
		if i < b.N {
			tmp1Point := b.CP[i].Minus()
			tmp2Point := b.CP[i+1].Add(&tmp1Point)
			tmp2Point.SMultiplication(float64(b.N))
			b.QP[i] = tmp2Point.Clone()
			b.BDifF[i] = BernsteinPolynomial{N: b.N - 1, I: i}
			b.BDifF[i].Init()
		}

	}

	b.Conf = Config{}
	b.Conf.Load()
	l := b.LengthCal(b.Conf.LengthCalResolotion) // what is the best resulution for calculation of curve length? maybe it's better to be dependent on it's value or curveture 1:1000000 , 2:1000
	log.Printf("B length is %f", l)
	//u := 0.1
	//log.Printf("ds/du %f @ %f & lenght is %f", b.DsPDu(u, 1000), u, l)

	return nil
}

//JSON : string json
func (bp *BernsteinPolynomial) JSON() []byte {
	res, _ := json.Marshal(bp)
	return res
}

//Cal : calculate the curve point at u
func (b *Bezier) Cal(u float64) Point {
	var i uint64
	res := Point{}
	for i = 0; i <= b.N; i++ {
		p := b.CP[i].Multiplication(b.BF[i].Cal(u))
		res.SAdd(&p)
	}
	return res
}

//DiffCal : calculate the curve tangent vector at u
func (b *Bezier) DiffCal(u float64) Vector {
	var i uint64
	res := Point{}
	for i = 0; i < b.N; i++ {
		p := b.QP[i].Multiplication(b.BDifF[i].Cal(u))
		res.SAdd(&p)
	}
	return res.Vector()
}

//Length : calculate the length
func (b *Bezier) LengthCal(n uint32) float64 {
	var res float64
	var i uint32
	for i = 0; i <= n; i++ {
		p2 := b.Cal(float64(i+1) / float64(n+1))
		p1 := b.Cal(float64(i) / float64(n+1))
		res += math.Sqrt(math.Pow(p2.X-p1.X, 2) + math.Pow(p2.Y-p1.Y, 2) + math.Pow(p2.Z-p1.Z, 2))
	}
	b.Length = res
	return res
}

//DsPDu : calculation of Ds/Du , @ u float64 with the calculation Resolution of n uint32 *************************
func (b *Bezier) DsPDu(u float64, n uint32) float64 {
	p1 := b.Cal(u - 1/float64(2*n))
	p2 := b.Cal(u + 1/float64(2*n))
	return math.Sqrt(math.Pow(p2.X-p1.X, 2)+math.Pow(p2.Y-p1.Y, 2)+math.Pow(p2.Z-p1.Z, 2)) * float64(n)

}

//errCF : error coeficent calculation factor for Velocity
func (b *Bezier) errCF(val float64, err float64) float64 {
	return val * (1 + b.Conf.ErrFactor*err*err*err)
}

//Go : Go path
func (b *Bezier) Go(baseFeed float64, resolotion int, ctx iris.Context) {
	start := time.Now()
	Tstr := start.Format("2006_Jan__2_15_04")
	Tlable := fmt.Sprintf("Samples/%s", Tstr)
	os.MkdirAll(Tlable, 0755)
	operationDone := make(chan bool)
	var OD bool
	//var Feed float64

	var U float64
	var V Vector
	Pos := Point{
		X: 0,
		Y: 0,
		Z: 0,
	}
	tPos := Point{
		X: 0,
		Y: 0,
		Z: 0,
	}
	Err := Point{
		X: 0,
		Y: 0,
		Z: 0,
	}
	U = b.U
	//Feed = baseFeed
	OD = false
	V = b.DiffCal(U)
	V.Unic()
	V.SMultiplication(b.Feed(U, baseFeed))

	log.Printf("the curve length is : %f", b.Length)
	//X axis
	go func() {
		var myt time.Time
		var delay int64
		var mypass bool
		filex := fmt.Sprintf("%s/Filex.txt", Tlable)
		Filex, err := os.Create(filex)
		check(err)
		defer Filex.Close()
		var str string
		Filex.WriteString("time(s) \t Pos(mm) \t Speed(mm/s) \t u \t Err\n")
		for OD == false {
			st := time.Now()
			//log.Printf("X %f , %f", V.X, U)
			delay, mypass = delayCalculation(b.errCF(V.X, Err.X))

			if mypass == true {
				myt = time.Now()
				elapsed := myt.Sub(start)
				str = fmt.Sprintf("%f \t %f \t %f \t %f \t %f\n", elapsed.Seconds(), Pos.X, V.X, U, Err.X)
				Filex.WriteString(str)
				if V.X > 0 {
					Pos.X = Pos.X + 0.0025
				} else if V.X < 0 {
					Pos.X = Pos.X - 0.0025
				}
			}
			dt := time.Now().Sub(st).Nanoseconds()
			time.Sleep(time.Duration(delay-dt) * time.Nanosecond)
		}
		operationDone <- true
	}()
	//Y axis
	go func() {
		var myt time.Time
		var delay int64
		var mypass bool
		filey := fmt.Sprintf("%s/Filey.txt", Tlable)
		Filey, err := os.Create(filey)
		check(err)
		defer Filey.Close()
		var str string
		Filey.WriteString("time(s) , Pos(mm) , Speed(mm/s) , u , Err\n")
		for OD == false {
			st := time.Now()
			//log.Printf("Y %f , %f", V.Y, U)
			delay, mypass = delayCalculation(b.errCF(V.Y, Err.Y))
			if mypass == true {
				myt = time.Now()
				elapsed := myt.Sub(start)
				str = fmt.Sprintf("%f , %f , %f , %f , %f\n", elapsed.Seconds(), Pos.Y, V.Y, U, Err.Y)
				Filey.WriteString(str)
				if V.Y > 0 {
					Pos.Y = Pos.Y + 0.0025
				} else if V.Y < 0 {
					Pos.Y = Pos.Y - 0.0025
				}
			}
			dt := time.Now().Sub(st).Nanoseconds()
			time.Sleep(time.Duration(delay-dt) * time.Nanosecond)
		}
		operationDone <- true
	}()
	//Z axis
	go func() {
		var myt time.Time
		var delay int64
		var mypass bool
		filez := fmt.Sprintf("%s/Filez.txt", Tlable)
		Filez, err := os.Create(filez)
		check(err)
		defer Filez.Close()
		var str string
		Filez.WriteString("time(s) , Pos(mm) , Speed(mm/s) , u , Err\n")
		for OD == false {
			st := time.Now()
			//log.Printf("Y %f , %f", V.Y, U)
			delay, mypass = delayCalculation(b.errCF(V.Z, Err.Z))
			if mypass == true {
				myt = time.Now()
				elapsed := myt.Sub(start)
				str = fmt.Sprintf("%f , %f , %f , %f , %f\n", elapsed.Seconds(), Pos.Z, V.Z, U, Err.Z)
				Filez.WriteString(str)
				if V.Z > 0 {
					Pos.Z = Pos.Z + 0.0025
				} else if V.Z < 0 {
					Pos.Z = Pos.Z - 0.0025
				}
			}
			dt := time.Now().Sub(st).Nanoseconds()
			time.Sleep(time.Duration(delay-dt) * time.Nanosecond)
		}
		operationDone <- true
	}()

	// Master Thread
	go func() {
		U = 0
		for U <= 1 {
			st := time.Now()
			U = U + (b.Feed(U, baseFeed)*0.001)/b.DsPDu(U, 10000)
			V = b.DiffCal(U)
			V.Unic()
			V.SMultiplication(b.Feed(U, baseFeed))
			tPos = b.Cal(U)
			tmppos := Pos.Minus()
			Err = tPos.Add(&tmppos)
			dt := time.Now().Sub(st).Nanoseconds()
			time.Sleep(time.Duration(1000000-dt) * time.Nanosecond)
		}
		// when u == 1
		OD = true
	}()
	// Loger thread
	go func() {
		t := int((b.Length * 2000) / (baseFeed * float64(resolotion)))
		// type data struct {
		// 	Points []Point `json:"points"`
		// }
		// CP := data{}
		ctx.WriteString("{\"points\" : [")
		ctx.Writef("[%f,%f,%f,%f]", U, Pos.X, Pos.Y, Pos.Z)
		//CP.Points[i] = Pos.Clone()
		time.Sleep(time.Duration(t) * time.Millisecond)
		// i := 0
		for OD == false {
			ctx.Writef(",[%f,%f,%f,%f]", U, Pos.X, Pos.Y, Pos.Z)
			//CP.Points[i] = Pos.Clone()
			time.Sleep(time.Duration(t) * time.Millisecond)

		}
		ctx.WriteString("]}")
		// str, _ := json.Marshal(CP)
		// ctx.Write(str)
		operationDone <- true
	}()
	<-operationDone
	<-operationDone
	<-operationDone
	<-operationDone

	t := time.Now()
	elapsed := t.Sub(start)
	//bcFile := fmt.Sprintf("%s/BezierCurve", Tlable)
	bcFile, _ := os.Create(fmt.Sprintf("%s/BezierCurve", Tlable))
	bcstr, _ := json.Marshal(b)
	bcFile.Write(bcstr)
	defer bcFile.Close()

	log.Printf("all jobs Done it takes %s", elapsed.String())
	log.Printf("the last Point is : ( %f , %f , %f )", Pos.X, Pos.Y, Pos.Z)
}

func (b *Bezier) Gopos(baseFeed float64, resolotion int, ctx iris.Context) {
	start := time.Now()
	Tstr := start.Format("2006_Jan__2_15_04")
	Tlable := fmt.Sprintf("Samples/%s", Tstr)
	os.MkdirAll(Tlable, 0755)
	operationDone := make(chan bool)
	LoggingDone := make(chan bool)
	var OD bool
	//var Feed float64

	var U float64
	var V Vector
	Pos := Point{
		X: 0,
		Y: 0,
		Z: 0,
	}
	/*tPos := Point{
		X: 0,
		Y: 0,
		Z: 0,
	}
	Err := Point{
		X: 0,
		Y: 0,
		Z: 0,
	}*/
	U = b.U
	//Feed = baseFeed
	OD = false
	V = b.DiffCal(U)
	V.Unic()
	V.SMultiplication(b.Feed(U, baseFeed))

	log.Printf("the curve length is : %f", b.Length)

	// Master Thread
	go func() {
		time.Sleep(time.Duration(300) * time.Millisecond)
	}()
	// Loger thread
	go func() {
		t := int((b.Length * 2000) / (baseFeed * float64(resolotion)))
		// type data struct {
		// 	Points []Point `json:"points"`
		// }
		// CP := data{}
		ctx.WriteString("{\"points\" : [")
		ctx.Writef("[%f,%f,%f,%f]", U, Pos.X, Pos.Y, Pos.Z)
		//CP.Points[i] = Pos.Clone()
		time.Sleep(time.Duration(t) * time.Millisecond)
		// i := 0
		for OD == false {
			ctx.Writef(",[%f,%f,%f,%f]", U, Pos.X, Pos.Y, Pos.Z)
			//CP.Points[i] = Pos.Clone()
			time.Sleep(time.Duration(t) * time.Millisecond)

		}
		ctx.WriteString("]}")
		// str, _ := json.Marshal(CP)
		// ctx.Write(str)
		LoggingDone <- true
	}()
	<-operationDone
	<-operationDone
	<-operationDone

	<-LoggingDone

	t := time.Now()
	elapsed := t.Sub(start)
	//bcFile := fmt.Sprintf("%s/BezierCurve", Tlable)
	bcFile, _ := os.Create(fmt.Sprintf("%s/BezierCurve", Tlable))
	bcstr, _ := json.Marshal(b)
	bcFile.Write(bcstr)
	defer bcFile.Close()

	log.Printf("all jobs Done it takes %s", elapsed.String())
	log.Printf("the last Point is : ( %f , %f , %f )", Pos.X, Pos.Y, Pos.Z)

}

func (b *Bezier) Feed(U float64, baseFeed float64) float64 {
	return baseFeed

}
