package core

import (
	"log"
	"net/http"

	"../../internal/core/app/handler"
	"../../internal/core/service"
)

func main() {
	// Initialize the user service
	userService := service.NewUserService()

	// Initialize the handler with the user service
	userHandler := handler.NewUserHandler(userService)

	// Setup the HTTP route for adding a new user
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			userHandler.HandleAddUser(w, r)
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
