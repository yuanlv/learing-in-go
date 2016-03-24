package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
	"github.com/yuanlv/learning-in-go/rest_email"
	"log"
)

func main() {
	email := rest_email.Email{}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/", email.GetServiceList),
		rest.Get("/email", email.GetEmailListAndContent),
		rest.Post("/email", email.SendEmailContent),
		rest.Delete("/email", email.DeleteEmail),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Printf("email api server start...")
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
