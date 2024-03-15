package core

import (
	"log"
	"net/http"

	"../../internal/core/app/handler"
	"../../internal/core/repository"
	"../../internal/core/service"
)

// ServiceSetup holds all the services and their handlers for the application.
type ServiceSetup struct {
	UserServiceHandler *handler.UserHandler
	// Add other service handlers here as you expand
}

// initializeServices sets up all the services, repositories, and handlers.
func initializeServices() ServiceSetup {
	// Initialize all repositories
	userRepo := repository.NewMockUserRepository()

	// Initialize all services with their respective repositories
	userService := service.NewUserService(userRepo)

	// Initialize all handlers with their respective services
	userServiceHandler := handler.NewUserHandler(userService)

	// Return a struct containing all handlers for easy access
	return ServiceSetup{
		UserServiceHandler: userServiceHandler,
	}
}

func main() {
	services := initializeServices()

	// Setup the HTTP route for adding a new user
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			services.UserServiceHandler.HandleAddUser(w, r)
			return
		}

		// Method Not Allowed for non-POST requests
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	// Start the HTTP server
	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
