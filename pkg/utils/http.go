package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labstack/echo/v5"
)

func GetHTTPVars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

func HttpMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// If it's a preflight OPTIONS request, stop here
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func EchoMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("Access-Control-Allow-Origin", "*")
		c.Request().Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH")
		c.Request().Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request().Method == "OPTIONS" {
			c.JSON(http.StatusOK, "")
			return nil
		}

		return next(c)
	}
}

func FormatResponse(w http.ResponseWriter, response interface{}, code int) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if response != nil {
		if err := json.NewEncoder(w).Encode(response); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

type GenericHttpResponse struct {
	Data    interface{} `json:"data"`
	Error   string      `json:"error_message"`
	Message string      `json:"message"`
}
