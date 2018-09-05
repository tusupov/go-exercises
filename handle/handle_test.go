package handle

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("An error occurred. %v", err)
	}
}

func TestHandle_Account(t *testing.T) {

	h := New()

	var testCase = []struct {
		id           int
		handleFunc   func(http.ResponseWriter, *http.Request)
		body         string
		responseCode int
	}{

		// not active
		{
			1,
			h.AccountBalance,
			``,
			400,
		},
		{
			2,
			h.AccountDeposit,
			`{"amount":10}`,
			400,
		},
		{
			3,
			h.AccountClose,
			``,
			400,
		},

		// active
		{
			4,
			h.AccountOpen,
			`{initialAmount": -100}`,
			500,
		},
		{
			5,
			h.AccountOpen,
			`{"initialAmount": -100}`,
			400,
		},
		{
			6,
			h.AccountOpen,
			`{"initialAmount": 100}`,
			200,
		},

		// balance
		{
			7,
			h.AccountBalance,
			``,
			200,
		},

		// deposit
		{
			8,
			h.AccountDeposit,
			`{amount": 10}`,
			500,
		},
		{
			9,
			h.AccountDeposit,
			`{"amount": 10}`,
			200,
		},
		{
			10,
			h.AccountDeposit,
			`{"amount": -200}`,
			400,
		},

		// close
		{
			11,
			h.AccountClose,
			``,
			200,
		},

		// closed
		{
			12,
			h.AccountBalance,
			``,
			400,
		},
		{
			13,
			h.AccountDeposit,
			`{"amount":10}`,
			400,
		},
		{
			14,
			h.AccountClose,
			``,
			400,
		},
	}

	for i, test := range testCase {

		req, err := http.NewRequest("PUT", "/account", strings.NewReader(test.body))
		checkError(t, err)

		recorder := httptest.NewRecorder()

		http.
			HandlerFunc(test.handleFunc).
			ServeHTTP(recorder, req)

		//Confirm the response has the right status code
		if status := recorder.Code; status != test.responseCode {
			t.Fatalf("Status code differs. Expected %d.\n Got %d instead\n%d, #%v", test.responseCode, status, i, test)
		}

	}

}
