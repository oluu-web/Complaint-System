package routes

import (
	"complaints/cmd/api/controllers"
	"complaints/cmd/api/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func InitRoutes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, "/register", controllers.Register)
	router.HandlerFunc(http.MethodPost, "/login", controllers.Login)

	return middleware.EnableCORS(router)
}
