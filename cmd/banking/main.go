package main

import (
	"banking/pkg/account"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	uc := account.NewUsecase()
	delivery := account.NewHTTPDelivery(uc)

	r.Post("/", delivery.Register)
	r.Get("/amount/{id}", delivery.Amount)
	r.Post("/transfer", delivery.Transfer)

	log.Fatal(http.ListenAndServe(":8080", r))
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
