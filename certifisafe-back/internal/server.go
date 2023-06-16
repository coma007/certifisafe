package internal

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

type Server interface {
	ListenAndServe()
}

type DefaultServer struct {
	app *DefaultAppFactory
}

func NewDefaultRouter(app *DefaultAppFactory) *DefaultServer {
	return &DefaultServer{
		app: app,
	}
}

func (server DefaultServer) ListenAndServe() {
	fmt.Println("http server runs on :8080")
	router := server.initRoutes()
	handler := server.handleCORS(router)
	http.ListenAndServe(":8080", handler)
}

func (server DefaultServer) handleCORS(router *mux.Router) http.Handler {
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
	}).Handler(router)
	return handler
}

func (server *DefaultServer) initRoutes() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/certificate/{id}", server.middleware(server.app.Controllers.CertificateController.GetCertificate)).Methods("GET")
	router.HandleFunc("/api/certificate", server.middleware(server.app.Controllers.CertificateController.GetCertificates)).Methods("GET")
	router.HandleFunc("/api/certificate/{id}/download", server.middleware(server.app.Controllers.CertificateController.DownloadCertificate)).Methods("GET")
	router.HandleFunc("/api/certificate/{id}/withdraw", server.middleware(server.app.Controllers.CertificateController.WithdrawCertificate)).Methods("PATCH")
	router.HandleFunc("/api/certificate/{id}/valid", server.middleware(server.app.Controllers.CertificateController.IsValid)).Methods("GET")
	router.HandleFunc("/api/certificate/valid", server.middleware(server.app.Controllers.CertificateController.IsValidFile)).Methods("POST")

	router.HandleFunc("/api/request", server.middleware(server.app.Controllers.RequestController.CreateRequest)).Methods("POST")
	router.HandleFunc("/api/request/{id}", server.middleware(server.app.Controllers.RequestController.GetRequest)).Methods("GET")
	router.HandleFunc("/api/request/signing", server.middleware(server.app.Controllers.RequestController.GetAllRequestsByUserSigning)).Methods("GET")
	router.HandleFunc("/api/request/user", server.middleware(server.app.Controllers.RequestController.GetAllRequestsByUser)).Methods("GET")
	router.HandleFunc("/api/request/accept/{id}", server.middleware(server.app.Controllers.RequestController.AcceptRequest)).Methods("PATCH")
	router.HandleFunc("/api/request/decline/{id}", server.middleware(server.app.Controllers.RequestController.DeclineRequest)).Methods("PATCH")
	router.HandleFunc("/api/request/delete/{id}", server.middleware(server.app.Controllers.RequestController.DeleteRequest)).Methods("PATCH")
	router.HandleFunc("/api/certificate/generate", server.middleware(server.app.Controllers.RequestController.GenerateCertificates)).Methods("PATCH")

	router.HandleFunc("/api/login", server.app.Controllers.AuthController.Login).Methods("POST")
	router.HandleFunc("/api/register", server.app.Controllers.AuthController.Register).Methods("POST")
	router.HandleFunc("/api/verify-email/{verificationCode}", server.app.Controllers.AuthController.VerifyEmail).Methods("GET")
	router.HandleFunc("/api/password-recovery-request", server.app.Controllers.AuthController.PasswordRecoveryRequest).Methods("POST")
	router.HandleFunc("/api/password-recovery", server.app.Controllers.AuthController.PasswordRecovery).Methods("POST")

	router.HandleFunc("/api/oauth", server.app.Controllers.OauthController.Oauth).Methods("GET")
	router.HandleFunc("/api/oauth/callback", server.app.Controllers.OauthController.OauthCallback).Methods("GET")

	return router
}

func (server *DefaultServer) middleware(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		tokenValid, err := server.app.Services.AuthService.ValidateToken(token)
		if err != nil || !tokenValid {
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		f(w, r)
	}
}
