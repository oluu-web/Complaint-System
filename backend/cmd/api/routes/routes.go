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

	authHandler := func(handler http.HandlerFunc) http.HandlerFunc {
		return middleware.Authenticate(handler).ServeHTTP
	}
	router.HandlerFunc(http.MethodPost, "/complaint", authHandler(controllers.NewComplaint))

	return middleware.EnableCORS(router)
}
