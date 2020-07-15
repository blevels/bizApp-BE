// TODO Add route to ascertain api version
// TODO https://itnext.io/learning-go-mongodb-crud-with-grpc-98e425aeaae6

package routers

import (
	"github.com/aerogo/aero"

	"backend/config"
)

var c = config.CreateConfigService()

func InitAero() {
	app := aero.New()
	app = InitRoutes(app)
	app.Run()
}

// Router is exported and used in main.go

func InitRoutes(app *aero.Application) *aero.Application {
	app = SetAuthenticationRoutes(c, app)
	app = SetUserRoutes(c, app)
	return app
}