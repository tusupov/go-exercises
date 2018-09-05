package handle

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/tusupov/go-exercises/bank-account"
)

// Error codes returned by failures to parse an expression.
var (
	errCreated     = errors.New("Account is created")
	errFailCreated = errors.New("Failed to create account, a negative initial deposit")
	errNotCreated  = errors.New("Account is not created")
	errNotEnough   = errors.New("Not enough money")
	errClosed      = errors.New("Account is closed")
)

type handle struct {
	account *account.Account
}

// New handle
func New() *handle {
	return &handle{}
}

// AccountOpen handle function
func (h *handle) AccountOpen(w http.ResponseWriter, r *http.Request) {
	// Account opened
	if h.account != nil {
		JSONErrorResponse(w, http.StatusBadRequest, errCreated.Error())
		return
	}

	// Read body
	var accountRequest account.Request
	err := json.NewDecoder(r.Body).Decode(&accountRequest)
	if err != nil {
		JSONErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Open account
	h.account = account.Open(accountRequest.Balance)
	if h.account == nil {
		JSONErrorResponse(w, http.StatusBadRequest, errFailCreated.Error())
		return
	}

	JSONSuccessResponse(w, account.Response{Balance: accountRequest.Balance})
}

// AccountBalance handle function
func (h *handle) AccountBalance(w http.ResponseWriter, r *http.Request) {
	// Account not created
	if h.account == nil {
		JSONErrorResponse(w, http.StatusBadRequest, errNotCreated.Error())
		return
	}

	// Get account balance
	balance, ok := h.account.Balance()
	if !ok {
		JSONErrorResponse(w, http.StatusBadRequest, errClosed.Error())
		return
	}

	JSONSuccessResponse(w, account.Response{Balance: balance})
}

// AccountDeposit handle function
func (h *handle) AccountDeposit(w http.ResponseWriter, r *http.Request) {
	// Account not created
	if h.account == nil {
		JSONErrorResponse(w, http.StatusBadRequest, errNotCreated.Error())
		return
	}

	// Check to closed
	if _, ok := h.account.Balance(); !ok {
		JSONErrorResponse(w, http.StatusBadRequest, errClosed.Error())
		return
	}

	// Read from body
	var accountResponse account.Response
	err := json.NewDecoder(r.Body).Decode(&accountResponse)
	if err != nil {
		JSONErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Deposit
	newBalance, ok := h.account.Deposit(accountResponse.Balance)
	if !ok {
		JSONErrorResponse(w, http.StatusBadRequest, errNotEnough.Error())
		return
	}

	JSONSuccessResponse(w, account.Response{Balance: newBalance})
}

// AccountClose handle function
func (h *handle) AccountClose(w http.ResponseWriter, r *http.Request) {
	// Account not created
	if h.account == nil {
		JSONErrorResponse(w, http.StatusBadRequest, errNotCreated.Error())
		return
	}

	// Close account
	payout, ok := h.account.Close()
	if !ok {
		JSONErrorResponse(w, http.StatusBadRequest, errClosed.Error())
		return
	}

	JSONSuccessResponse(w, account.Response{Balance: payout})
}

func JSONSuccessResponse(w http.ResponseWriter, resp interface{}) {
	JSONResponse(w, http.StatusOK, resp)
}

// Write
func JSONErrorResponse(w http.ResponseWriter, errorCode int, err string) {
	type ApiError struct {
		Error string `json:"error"`
	}
	JSONResponse(w, errorCode, ApiError{
		Error: err,
	})
}

// Write JSON data to http writer
// Encode resp to json and write to w
func JSONResponse(w http.ResponseWriter, statusCode int, resp interface{}) {

	respJson, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(respJson); err != nil {
		log.Println(err)
	}

}
