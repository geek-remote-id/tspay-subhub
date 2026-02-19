package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/geek-remote-id/tspay-subhub/docs"
	"github.com/geek-remote-id/tspay-subhub/handlers"
	"github.com/geek-remote-id/tspay-subhub/utils"

	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Tspay Subhub API
// @version         1.0
// @description     API Server for Tspay Subhub.
// @BasePath        /api

// @Summary Health Check
// @Description Get the status of the API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Response{
		Status:  "success",
		Message: "API Running",
	})
}

func main() {
	// load .env using viper
	viper.SetConfigFile(".env")
	viper.AutomaticEnv() // read value from system env too

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using system environment variables")
	} else {
		log.Println(".env file loaded successfully")
	}

	portStr := viper.GetString("PORT")
	if portStr == "" {
		portStr = "8080"
	}

	// Dynamic Swagger Host
	appHost := viper.GetString("APP_HOST")
	if appHost != "" {
		docs.SwaggerInfo.Host = appHost
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
	} else {
		docs.SwaggerInfo.Host = "localhost:" + portStr
		docs.SwaggerInfo.Schemes = []string{"http"}
	}
	log.Println("Swagger Host set to:", docs.SwaggerInfo.Host)

	// Routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/incoming", handlers.GenerateIncomingHandler())

	// Swagger
	http.HandleFunc("/", httpSwagger.WrapHandler)

	fmt.Println("Server running on http://localhost:" + portStr)
	err := http.ListenAndServe(":"+portStr, nil)
	if err != nil {
		fmt.Println("Error running server:", err)
	}
}
