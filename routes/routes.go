// Package routes consist of router path used for handling incoming request //
package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rzknugraha/zorro-mark/controllers"
)

// Route is
type Route struct{}

// Init is
func (r *Route) Init() *mux.Router {
	// Initialize controller //
	healthCheckController := controllers.InitHealthCheckController()
	playerController := controllers.InitPlayerController()
	userController := controllers.InitUserController()
	uploadController := controllers.InitUploadController()
	documentController := controllers.InitDocumentController()
	esignController := controllers.InitEsignController()

	// Initialize router //
	router := mux.NewRouter().StrictSlash(false)
	v1 := router.PathPrefix("/v1").Subrouter()

	v1.HandleFunc("/cors", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "OPTIONS" {
			w.Write([]byte("allowed"))
			return
		}

		w.Write([]byte("hello"))
	}).Methods("GET")
	v1.HandleFunc("/healthcheck", healthCheckController.HealthCheck).Methods("GET")
	v1.HandleFunc("/player", playerController.StorePlayer).Methods("POST")
	v1.HandleFunc("/login", userController.Login).Methods("POST")

	ClientAuth := v1.PathPrefix("/client").Subrouter()
	ClientAuth.Use(JWTAuthMiddleware)

	ClientAuth.HandleFunc("/profile", userController.Profile).Methods(http.MethodGet)
	ClientAuth.HandleFunc("/profile/update/file", userController.UploadProfile).Methods(http.MethodPost)
	ClientAuth.HandleFunc("/file/upload", uploadController.Upload).Methods("POST")
	ClientAuth.HandleFunc("/file/get", uploadController.GetFile).Methods("POST")
	ClientAuth.HandleFunc("/document/get", documentController.GetDocuments).Methods("GET")
	ClientAuth.HandleFunc("/document/get/{IDDoc}", documentController.GetSingleDocument).Methods("GET")
	ClientAuth.HandleFunc("/document/update", documentController.UpdateDocument).Methods("POST")

	ClientAuth.HandleFunc("/document/activity/get/{IDDoc}", documentController.GetDocActivity).Methods("GET")

	//esign
	ClientAuth.HandleFunc("/sign/doc", esignController.SignDoc).Methods("POST")

	//Users
	ClientAuth.HandleFunc("/users/get", userController.GetAll).Methods("GET")

	return v1
}
