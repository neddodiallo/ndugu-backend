package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"ndugu-backend/internal/common"
	"ndugu-backend/internal/repository"
	"ndugu-backend/internal/services"
)

func main() {
	// Initialiser le logger
	logger := common.NewSugarLogger()

	// Initialiser les repositories
	userRepo := repository.NewMockUserRepository() // TODO: Remplacer par une vraie implémentation
	oryClient := repository.NewOryClient(logger)

	// Initialiser les services
	authService := services.NewAuthService(userRepo, oryClient, logger)

	// Créer le serveur gRPC
	grpcServer := NewGRPCServer(authService, logger)

	// Démarrer le serveur gRPC
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			logger.Error("Erreur lors de l'écoute gRPC: %v", err)
			os.Exit(1)
		}

		logger.Info("Serveur gRPC démarré sur le port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("Erreur lors du démarrage du serveur gRPC: %v", err)
			os.Exit(1)
		}
	}()

	logger.Info("🚀 Serveur Ndugu Backend démarré")
	logger.Info("📡 gRPC Server: localhost:50051")
	logger.Info("")
	logger.Info("🔗 Endpoints gRPC disponibles:")
	logger.Info("    - ndugu.v1.AuthService/CreateUser - Créer un utilisateur")
	logger.Info("    - ndugu.v1.AuthService/GetUser - Récupérer un utilisateur")
	logger.Info("    - ndugu.v1.AuthService/ValidateSession - Valider une session")
	logger.Info("    - ndugu.v1.AuthService/CreateOAuth2Client - Créer un client OAuth2")
	logger.Info("    - ndugu.v1.AuthService/CreatePermission - Créer une permission")
	logger.Info("    - ndugu.v1.AuthService/CheckPermission - Vérifier une permission")
	logger.Info("    - ndugu.v1.CustomerService/CreateCustomer - Créer un client")
	logger.Info("")
	logger.Info("🔧 Services Ory:")
	logger.Info("  - Kratos: http://localhost:4433 (public), http://localhost:4434 (admin)")
	logger.Info("  - Hydra: http://localhost:4444 (public), http://localhost:4445 (admin)")
	logger.Info("  - Keto: http://localhost:4466 (read), http://localhost:4467 (write)")
	logger.Info("")
	logger.Info("🌍 API Gateway (APISIX): http://localhost:9080")

	// Attendre un signal d'arrêt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 Arrêt du serveur...")
	// TODO: Implémenter l'arrêt gracieux
}
