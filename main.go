package main

import (
	rdx "be/config/radix"
	user_mongo_repo "be/domain/user/repository/mongodb"
	rdx_repo "be/domain/user/repository/radix"
	rest_external "be/rest/external"
	service "be/service"
	"log"

	db "be/config/mongodb"

	"github.com/gin-gonic/gin"
)

func main() {
	dbClient, err := db.New()
	if err != nil {
		log.Println(err)
	}

	mongoRepo := user_mongo_repo.New(dbClient)

	radixInit, err := rdx.New()
	if err != nil {
		log.Println(err)
	}

	rdx_repo_obj := rdx_repo.New(radixInit, mongoRepo)

	r := gin.Default()
	svc := service.New(rdx_repo_obj)
	rest_external.New(svc).Register(r)
	r.Run(":3000") // listen and serve on 0.0.0.0:8080
}
