package main

import (
	"context"
	"log"
	"os"
	

	"github.com/joho/godotenv"
	"ServiceBookingApp/internal/infrastructure/db"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "ServiceBookingApp/docs"
	
	authService "ServiceBookingApp/internal/auth"
	authHandler "ServiceBookingApp/internal/handlers/auth"
	firebase "firebase.google.com/go/v4"
	
	
	
	"ServiceBookingApp/internal/handlers/services"
	
	"ServiceBookingApp/internal/handlers/providers"
	
	"ServiceBookingApp/internal/handlers/appointments"
	
	"ServiceBookingApp/internal/handlers/users"
	
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize Database
	
	baseRepo, err := db.NewFirestoreRepository()
	
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer baseRepo.Close()

	
	// Initialize Auth Service
	
	var authSvc authService.AuthService
	if os.Getenv("MOCK_AUTH") == "true" {
		log.Println("Using Mock Auth Service")
		authSvc = &authService.MockAuthService{}
	} else {
		// Initialize Firebase Auth
		app, err := firebase.NewApp(context.Background(), &firebase.Config{ProjectID: "turnero-165d4"})
		if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
		}
		authClient, err := app.Auth(context.Background())
		if err != nil {
			log.Fatalf("error getting Auth client: %v\n", err)
		}
		authSvc = &authService.FirebaseAuthService{Client: authClient}
	}
	// Initialize User Handler
	
	userRepo := db.NewUsersRepository(baseRepo.(*db.FirestoreRepository))
	
	userHdl := authHandler.NewUserHandler(authSvc, userRepo, "users")
	
	

	

	// Setup Router
	r := gin.Default()

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	
	// Auth Routes
	authGroup := r.Group("/auth")
	
	authGroup.POST("/login", authService.AuthMiddleware(authSvc), userHdl.Login)
	
	authGroup.GET("/me", authService.AuthMiddleware(authSvc), userHdl.GetMe)
	authGroup.GET("/roles", authService.AuthMiddleware(authSvc), userHdl.GetRoles)
	

	

	
	// Routes for services
	{
		
		repo := db.NewServicesRepository(baseRepo.(*db.FirestoreRepository))
		
		handler := services.NewServicesHandler(repo)

		group := r.Group("/api/services")
		
		
		group.Use(authService.AuthMiddleware(authSvc))
		
		
		group.GET("", handler.List)
		group.GET("/:id", handler.Get)
		group.POST("", handler.Create)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
	}
	
	// Routes for providers
	{
		
		repo := db.NewProvidersRepository(baseRepo.(*db.FirestoreRepository))
		
		handler := providers.NewProvidersHandler(repo)

		group := r.Group("/api/providers")
		
		
		group.Use(authService.AuthMiddleware(authSvc))
		
		
		group.GET("", handler.List)
		group.GET("/:id", handler.Get)
		group.POST("", handler.Create)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
	}
	
	// Routes for appointments
	{
		
		repo := db.NewAppointmentsRepository(baseRepo.(*db.FirestoreRepository))
		
		handler := appointments.NewAppointmentsHandler(repo)

		group := r.Group("/api/appointments")
		
		
		group.Use(authService.AuthMiddleware(authSvc))
		
		
		group.GET("", handler.List)
		group.GET("/:id", handler.Get)
		group.POST("", handler.Create)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
	}
	
	// Routes for users
	{
		
		repo := db.NewUsersRepository(baseRepo.(*db.FirestoreRepository))
		
		handler := users.NewUsersHandler(repo)

		group := r.Group("/api/users")
		
		
		group.Use(authService.AuthMiddleware(authSvc))
		
		
		group.GET("", handler.List)
		group.GET("/:id", handler.Get)
		group.POST("", handler.Create)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
	}
	

	log.Printf("Starting server for project: ServiceBookingApp")
	r.Run(":8080")
}


