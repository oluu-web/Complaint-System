package routes

import (
	"complaints/cmd/api/controllers"
	"complaints/cmd/api/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func InitRoutes() http.Handler {
	router := httprouter.New()

	//serve static files
	router.ServeFiles("/uploads/*filepath", http.Dir("./uploads"))
	router.HandlerFunc(http.MethodPost, "/register", controllers.Register)
	router.HandlerFunc(http.MethodPost, "/login", controllers.Login)

	authHandler := func(handler http.HandlerFunc) http.HandlerFunc {
		return middleware.Authenticate(handler).ServeHTTP
	}
	router.HandlerFunc(http.MethodPost, "/complaint", authHandler(controllers.NewComplaint))
	router.HandlerFunc(http.MethodGet, "/complaint/:id", authHandler(controllers.GetComplaintByObjectID))
	router.HandlerFunc(http.MethodGet, "/complaints/:id", authHandler(controllers.GetComplaintsByStudentID))
	router.HandlerFunc(http.MethodGet, "/courses/:id", authHandler(controllers.GetCoursesByStudentID))
	router.HandlerFunc(http.MethodGet, "/lecturer-courses/:id", authHandler(controllers.GetCoursesByStaffID))
	router.HandlerFunc(http.MethodGet, "/staff-complaints/:id", authHandler(controllers.GetComplaintsByStaffID))
	router.HandlerFunc(http.MethodGet, "/hod-complaints", authHandler(controllers.GetComplaintsForHOD))
	router.HandlerFunc(http.MethodGet, "/senate-complaints", authHandler(controllers.GetComplaintsForSenate))
	router.HandlerFunc(http.MethodGet, "/lecturer-complaints/:id", authHandler(controllers.GetComplaintsByCourseCode))
	router.HandlerFunc(http.MethodPut, "/approved-by-lecturer/:id", authHandler(controllers.ChangeComplaintStatusByLecturer))
	router.HandlerFunc(http.MethodPut, "/approved-by-advisor/:id", authHandler(controllers.ChangeComplaintStatusByAdvisor))
	router.HandlerFunc(http.MethodPut, "/approved-by-hod/:id", authHandler(controllers.ChangeComplaintStatusByHOD))
	router.HandlerFunc(http.MethodPut, "/approved-by-senate/:id", authHandler(controllers.ChangeComplaintStatusBySenate))
	router.HandlerFunc(http.MethodPut, "/decline/:id", authHandler((controllers.DeclineRequest)))

	//serve static files
	// router.Handler(http.MethodGet, "/uploads/*filepath", http.StripPrefix("/uploads", http.FileServer(http.Dir("uploads"))))

	return middleware.EnableCORS(router)
}
