package main

/*eslint-disable */
import (
	//"encoding/json"
	//"log"

	//"github.com/dgraph-io/dgo"
	//"github.com/dgraph-io/dgo/protos/api"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"

	"github.com/kataras/iris/mvc"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"

	//"ehome/datamodels"
	"Robotic/web/controllers"
	//grpc "google.golang.org/grpc"
)

// dgraph client
// func newClient() *dgo.Dgraph {
// 	// Dial a gRPC connection. The address to dial to can be configured when
// 	// setting up the dgraph cluster.
// 	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return dgo.NewDgraphClient(
// 		api.NewDgraphClient(d),
// 	)
// }

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())
	app.StaticWeb("/static", "./web/static")
	app.RegisterView(iris.HTML("./web/views", ".html"))

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*", "http://localhost:8080"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
	})

	mvc.New(app.Party("/Rotation", crs)).Handle(new(controllers.RotationController))

	// Method : GET
	// resource: @base/root/

	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/ping
	app.Get("/", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello
	app.Run(iris.Addr(":9090"), iris.WithoutServerError(iris.ErrServerClosed))
}
