package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"ndugu-backend/internal/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server encapsule le serveur gRPC et HTTP
type Server struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	oryClient  *auth.OryClient
}

// NewServer crée une nouvelle instance du serveur
func NewServer() *Server {
	// Créer le client Ory
	oryClient := auth.NewOryClient()

	// Créer le serveur gRPC
	grpcServer := grpc.NewServer()

	// Enregistrer les services
	authService := NewAuthServiceServer(oryClient)
	RegisterAuthServiceServer(grpcServer, authService)

	// Activer la réflexion gRPC pour le débogage
	reflection.Register(grpcServer)

	return &Server{
		grpcServer: grpcServer,
		oryClient:  oryClient,
	}
}

// Start démarre les serveurs gRPC et HTTP
func (s *Server) Start(grpcPort, httpPort string) error {
	// Démarrer le serveur gRPC
	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("Erreur lors de l'écoute gRPC: %v", err)
		}

		log.Printf("Serveur gRPC démarré sur le port %s", grpcPort)
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Fatalf("Erreur lors du démarrage du serveur gRPC: %v", err)
		}
	}()

	// Démarrer le serveur HTTP
	s.setupHTTPServer()
	go func() {
		log.Printf("Serveur HTTP démarré sur le port %s", httpPort)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur lors du démarrage du serveur HTTP: %v", err)
		}
	}()

	return nil
}

// Stop arrête les serveurs
func (s *Server) Stop(ctx context.Context) error {
	// Arrêter le serveur gRPC
	s.grpcServer.GracefulStop()

	// Arrêter le serveur HTTP
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}

	return nil
}

// setupHTTPServer configure le serveur HTTP avec les endpoints REST
func (s *Server) setupHTTPServer() {
	mux := http.NewServeMux()

	// Endpoints REST pour les utilisateurs
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			s.handleCreateUserHTTP(w, r)
		case "GET":
			s.handleGetUserHTTP(w, r)
		default:
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		}
	})

	// Endpoint de santé
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Services Ory opérationnels")
	})

	s.httpServer = &http.Server{
		Handler: mux,
	}
}

// handleCreateUserHTTP gère la création d'utilisateur via HTTP
func (s *Server) handleCreateUserHTTP(w http.ResponseWriter, r *http.Request) {
	// Parser les données JSON
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erreur de parsing JSON", http.StatusBadRequest)
		return
	}

	// Créer l'utilisateur
	user, err := s.oryClient.CreateUser(r.Context(), req.Email, req.FirstName, req.LastName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retourner la réponse
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// handleGetUserHTTP gère la récupération d'utilisateur via HTTP
func (s *Server) handleGetUserHTTP(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "ID utilisateur requis", http.StatusBadRequest)
		return
	}

	// Récupérer l'utilisateur
	user, err := s.oryClient.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retourner la réponse
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
