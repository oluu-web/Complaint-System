package routes

import (
	"complaints/cmd/api/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func InitRoutes() http.Handler {
	router := httprouter.New()

	return middleware.EnableCORS(router)
}
