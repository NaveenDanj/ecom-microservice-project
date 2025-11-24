package routes

import (
	"net/http"
	"user-service/controllers"
	"user-service/middleware"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/auth/register", controllers.RegisterHandler).Methods("POST")
	r.HandleFunc("/auth/login", controllers.LoginHandler).Methods("POST")
	r.Handle("/auth/current-user", middleware.UserAuthMiddleware(http.HandlerFunc(controllers.CurrentUserHandler))).Methods("GET")
	// 	r.HandleFunc("/auth/login", controllers.LoginHandler).Methods("POST")
	// 	r.HandleFunc("/auth/profile", controllers.ProfileHandler).Methods("GET")
	// 	r.HandleFunc("/auth/update-profile", controllers.UpdateHandler).Methods("PUT")
}

// func AdminRoutes(r *mux.Router) {
// 	r.HandleFunc("/admin/users", controllers.GetAllUsersHandler).Methods("GET")
// 	r.HandleFunc("/admin/delete-user", controllers.DeleteUserHandler).Methods("DELETE")
// }
