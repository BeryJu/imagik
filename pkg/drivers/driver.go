package drivers

import "github.com/gorilla/mux"

type Driver interface {
	Init()
}

type HTTPDriver interface {
	Driver
	InitRoutes(r *mux.Router)
}
