package routers

import (
	"github.com/aerogo/aero"
	"net/http"

	"backend/config"
	"backend/controllers"
	"backend/core/authentication"
)

func SetAuthenticationRoutes(c *config.ConfigService, app *aero.Application) *aero.Application {
	app.Post(`/api/` + c.Config.Version + `/login`, controllers.Login)
	app.Post(`/api/` + c.Config.Version + `/refresh`, controllers.RefreshToken)
	app.Get(`/api/` + c.Config.Version + `/logout`, controllers.Logout)
	app.Router().Add(http.MethodOptions, `/api/` + c.Config.Version + `/logout`, controllers.Logout)

	// Register middleware
	app.Use(authentication.Cors)
	app.Use(authentication.RequireTokenAuthentication)
	return app
}