package routers

import (
	"backend/config"
	"backend/controllers"
	"backend/core/authentication"
	"github.com/aerogo/aero"
	"net/http"
)

func SetUserRoutes(c *config.ConfigService, app *aero.Application) *aero.Application {
	app.Post(`/api/` + c.Config.Version + `/register`, controllers.Register)
	app.Get(`/api/` + c.Config.Version + `/profile/:user`, controllers.Profile)
	app.Post(`/api/` + c.Config.Version + `/profile`, controllers.Profile)
	app.Get(`/api/` + c.Config.Version + `/settings/:user`, controllers.Settings)
	app.Post(`/api/` + c.Config.Version + `/settings`, controllers.Settings)
	app.Get(`/api/` + c.Config.Version + `/leads/:user`, controllers.Leads)
	app.Post(`/api/` + c.Config.Version + `/leads`, controllers.Leads)
	app.Post(`/api/` + c.Config.Version + `/team`, controllers.Team)
	app.Post(`/api/` + c.Config.Version + `/calendar`, controllers.Calendar)
	app.Router().Add(http.MethodOptions, `/api/` + c.Config.Version + `/users/all`, authentication.RequireAdminRole(controllers.Users))
	app.Post(`/api/` + c.Config.Version + `/users/all`, authentication.RequireAdminRole(controllers.Users))
	return app
}
