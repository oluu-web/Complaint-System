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
	router.HandlerFunc(http.MethodGet, "/complaint/:id", authHandler(controllers.GetComplaintByObjectID))
	router.HandlerFunc(http.MethodGet, "/complaints/:id", authHandler(controllers.GetComplaintsByStudentID))
	router.HandlerFunc(http.MethodGet, "/courses/:id", authHandler(controllers.GetCoursesByStudentID))
	router.HandlerFunc(http.MethodGet, "/staff-complaints/:id", authHandler(controllers.GetComplaintsByStaffID))
	router.HandlerFunc(http.MethodPut, "/approved-by-lecturer/:id", authHandler(controllers.ChangeComplaintStatusByLecturer))
	router.HandlerFunc(http.MethodPut, "/approved-by-advisor/:id", authHandler(controllers.ChangeComplaintStatusByAdvisor))
	router.HandlerFunc(http.MethodPut, "/approved-by-hod/:id", authHandler(controllers.ChangeComplaintStatusByHOD))
	router.HandlerFunc(http.MethodPut, "/approved-by-senate/:id", authHandler(controllers.ChangeComplaintStatusBySenate))
	return middleware.EnableCORS(router)
}
