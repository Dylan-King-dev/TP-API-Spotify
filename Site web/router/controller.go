package router

import (
	"TPSpotify/controller"
	"net/http"
)

func SetupRoutes() {
	http.HandleFunc("/", controller.DaHome)
	http.HandleFunc("/damso", controller.DamsoPage)
	http.HandleFunc("/laylow", controller.LaylowPage)

}
