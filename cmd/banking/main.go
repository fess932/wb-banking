package main

import (
	"banking/pkg/account"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", rest()))
}

func rest() http.Handler {
	r := chi.NewRouter()
	uc := account.NewUsecase()
	delivery := account.NewHTTPDelivery(uc)

	r.Post("/", delivery.Register)
	r.Get("/amount/{id}", delivery.Amount)
	r.Post("/transfer", delivery.Transfer)

	return r
}

// create account
// post
// id - app
// email  - user
// balance - user

// balance
// get /id

// transfer
// post
// {
// 	from: id
//  to: id
//  amount: 123
// }
