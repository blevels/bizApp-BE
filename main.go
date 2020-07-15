package main

import (
	"backend/routers"
	"backend/settings"
	"os"
)

func main() {
	os.Setenv("GO_ENV", "preproduction")
	settings.Init()
	routers.InitAero()
}