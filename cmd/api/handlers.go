package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Config) GetReccurringPayments(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")

	balance, err := app.Models.RecurringPayment.GetReccurringPayments(u)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//TODO: Add logging
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Grabbed recurring payments for user %s", u),
		Data:    balance,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetAllReccurringPayments(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")

	balance, err := app.Models.RecurringPayment.GetAllReccurringPayments()
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//TODO: Add logging
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Grabbed recurring payments for user %s", u),
		Data:    balance,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) AddReccurringPayment(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	var requestPayload struct {
		PaymentAmount      float32 `json:"amount"`
		PaymentName        string  `json:"paymentName"`
		PaymentDescription string  `json:"paymentDescription"`
		PaymentDate        string  `json:"paymentDate"`
	}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	balance, err := app.Models.RecurringPayment.AddReccurringPayment(u, requestPayload.PaymentAmount, requestPayload.PaymentName, requestPayload.PaymentDescription, requestPayload.PaymentDate)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Added recurring payment for user %s", u),
		Data:    balance,
	}
	app.writeJSON(w, http.StatusAccepted, payload)
}
func (app *Config) GetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "user")
	var requestPayload struct {
		PaymentID int `json:"paymentID"`
	}

	transactions, err := app.Models.PaymentHistory.GetPaymentHistory(requestPayload.PaymentID)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Retrieved payment history for user %s", u),
		Data:    transactions,
	}
	app.writeJSON(w, http.StatusAccepted, payload)

}
