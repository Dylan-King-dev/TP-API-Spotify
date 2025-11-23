package router

import (
	"TPSpotify/controller"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/", controller.DaHome)
}
