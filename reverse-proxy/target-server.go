package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/younesabouali/formal-challenges/reverse-proxy/utils"
)

func startTargetServer(targetUrl string) {
	router := chi.NewRouter()

	type message struct {
		Password string
		Email    string
	}
	response := message{Email: "ab@mail.com", Password: "123"}
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, 200, response)
	})

	router.Post("/*", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, 200, message{Password: "Heelo"})
	})

	router.Put("/*", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, 200, response)
	})

	router.Patch("/*", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, 200, response)
	})

	router.Delete("/*", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithJSON(w, 200, response)
	})
	port := strings.Split(targetUrl, ":")[2]
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	fmt.Println("server runing on PORT : ", port)
	srv.ListenAndServe()

}
