package main

import (
	auth "be/auth"
	rdx "be/config/radix"
	user_mongo_repo "be/domain/user/repository/mongodb"
	rdx_repo "be/domain/user/repository/radix"
	rest_external "be/rest/external"
	service "be/service"
	"log"
	"os"

	db "be/config/mongodb"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	envInit()

	dbClient, err := db.New()
	if err != nil {
		log.Println(err)
	}

	port := getEnv("PORT", "")
	if port != "" {
		port = ":" + port
	}

	googleAuth := auth.New(oauth2.Config{
		RedirectURL:  getEnv("HTTP_PROTOCOL", "") + "://" + getEnv("BASE_URL", "") + port + "/auth/google/callback",
		ClientID:     getEnv("GOOGLE_OAUTH_CLIENT_ID", ""),
		ClientSecret: getEnv("GOOGLE_OAUTH_CLIENT_SECRET", ""),
	})

	mongoRepo := user_mongo_repo.New(dbClient)

	radixInit, err := rdx.New()
	if err != nil {
		log.Println(err)
	}

	rdx_repo_obj := rdx_repo.New(radixInit, mongoRepo)

	r := gin.Default()
	svc := service.New(rdx_repo_obj)
	rest_external.New(svc, googleAuth).Register(r)
	r.Run(":" + getEnv("PORT", "")) // listen and serve on 0.0.0.0:8080
}

func envInit() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
