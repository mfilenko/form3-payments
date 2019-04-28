package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"

	"github.com/bitly/go-simplejson"
	"github.com/stretchr/testify/assert"
)

func TestHandlers(t *testing.T) {
	s := NewServer()

	s.Router = httprouter.New()
	s.Setup()
	s.OpenDB()

	defer s.Stop()

	var payment string

	t.Run("ReadPayments", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/payments", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		s.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		js, err := simplejson.NewFromReader(rr.Body)

		assert.NotEqual(t, nil, js)
		assert.Equal(t, nil, err)

		_, ok := js.CheckGet("data")
		assert.Equal(t, true, ok)
	})

	t.Run("CreatePayment", func(t *testing.T) {
		sample := `{"type":"Payment","organisation_id":"743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb","attributes":{"amount":"100.21","beneficiary_party":{"account_name":"W Owens","account_number":"31926819","account_number_code":"BBAN","account_type":0,"address":"1 The Beneficiary Localtown SE2","bank_id":"403000","bank_id_code":"GBDSC","name":"Wilfred Jeremiah Owens"},"charges_information":{"bearer_code":"SHAR","sender_charges":[{"amount":"5.00","currency":"GBP"},{"amount":"10.00","currency":"USD"}],"receiver_charges_amount":"1.00","receiver_charges_currency":"USD"},"currency":"GBP","debtor_party":{"account_name":"EJ Brown Black","account_number":"GB29XABC10161234567801","account_number_code":"IBAN","address":"10 Debtor Crescent Sourcetown NE1","bank_id":"203301","bank_id_code":"GBDSC","name":"Emelia Jane Brown"},"end_to_end_reference":"Wil piano Jan","fx":{"contract_reference":"FX123","exchange_rate":"2.00000","original_amount":"200.42","original_currency":"USD"},"numeric_reference":"1002001","payment_id":"123456789012345678","payment_purpose":"Paying for goods/services","payment_scheme":"FPS","payment_type":"Credit","processing_date":"2017-01-18","reference":"Payment for Em's piano lessons","scheme_payment_sub_type":"InternetBanking","scheme_payment_type":"ImmediatePayment","sponsor_party":{"account_number":"56781234","bank_id":"123123","bank_id_code":"GBDSC"}}}`

		tt := []struct {
			testcase string
			payload  []byte
			status   int
		}{
			{"400", []byte(`foo=bar`), http.StatusBadRequest},
			{"500", []byte(`{}`), http.StatusInternalServerError},
			{"201", []byte(sample), http.StatusCreated},
		}

		for _, tc := range tt {
			t.Run(tc.testcase, func(t *testing.T) {
				var payload io.Reader
				if tc.payload != nil {
					payload = bytes.NewBuffer(tc.payload)
				}

				req, err := http.NewRequest("POST", "/payments", payload)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				s.Router.ServeHTTP(rr, req)

				assert.Equal(t, tc.status, rr.Code)

				switch rr.Code {
				case 201:
					js, _ := simplejson.NewFromReader(rr.Body)
					payment, _ = js.Get("data").Get("id").String()
				default:
					js, _ := simplejson.NewFromReader(rr.Body)
					_, ok := js.CheckGet("message")
					assert.Equal(t, true, ok)
				}
			})
		}
	})

	t.Run("ReadPayment", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/payments"+"/"+payment, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		s.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		js, err := simplejson.NewFromReader(rr.Body)

		assert.NotEqual(t, nil, js)
		assert.Equal(t, nil, err)

		_, ok := js.CheckGet("data")
		assert.Equal(t, true, ok)
	})

	t.Run("UpdatePayment", func(t *testing.T) {
		sample := `{"attributes":{"amount":"1.21","end_to_end_reference":"Wil guitar Jan","fx":{"exchange_rate":"2.00000001"}}}`

		tt := []struct {
			testcase string
			payment  string
			payload  []byte
			status   int
		}{
			{"400", payment, []byte(`foo=bar`), http.StatusBadRequest},
			{"404", "00000000-0000-4000-8000-000000000000", []byte(sample), http.StatusNotFound},
			{"500", payment, []byte(`{"type":"0xff"}`), http.StatusInternalServerError},
			{"200", payment, []byte(sample), http.StatusOK},
		}

		for _, tc := range tt {
			t.Run(tc.testcase, func(t *testing.T) {
				var payload io.Reader
				if tc.payload != nil {
					payload = bytes.NewBuffer(tc.payload)
				}

				req, err := http.NewRequest("PATCH", "/payments"+"/"+tc.payment, payload)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				s.Router.ServeHTTP(rr, req)

				assert.Equal(t, tc.status, rr.Code)

				switch rr.Code {
				case 200:
					js, _ := simplejson.NewFromReader(rr.Body)
					ver, _ := js.Get("data").Get("version").Int()
					amount, _ := js.Get("data").Get("attributes").Get("amount").String()
					ref, _ := js.Get("data").Get("attributes").Get("end_to_end_reference").String()
					fxRate, _ := js.Get("data").Get("attributes").Get("fx").Get("exchange_rate").String()
					assert.Equal(t, 1, ver)
					assert.Equal(t, "1.21", amount)
					assert.Equal(t, "Wil guitar Jan", ref)
					assert.Equal(t, "2.00000001", fxRate)
				default:
					js, _ := simplejson.NewFromReader(rr.Body)
					_, ok := js.CheckGet("message")
					assert.Equal(t, true, ok)
				}
			})
		}
	})

	t.Run("UpgradePayment", func(t *testing.T) {
		sample := `{"type":"Payment","organisation_id":"ffffffff-ffff-4fff-bfff-ffffffffffff","attributes":{"debtor_party":{"account_name":"Mr Sending Test"}}}`

		tt := []struct {
			testcase string
			payment  string
			payload  []byte
			status   int
		}{
			{"400", payment, []byte(`foo=bar`), http.StatusBadRequest},
			{"404", "00000000-0000-4000-8000-000000000000", []byte(sample), http.StatusNotFound},
			{"500", payment, []byte(`{"type":"0xff"}`), http.StatusInternalServerError},
			{"200", payment, []byte(sample), http.StatusOK},
		}

		for _, tc := range tt {
			t.Run(tc.testcase, func(t *testing.T) {
				var payload io.Reader
				if tc.payload != nil {
					payload = bytes.NewBuffer(tc.payload)
				}

				req, err := http.NewRequest("PUT", "/payments"+"/"+tc.payment, payload)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()

				s.Router.ServeHTTP(rr, req)

				assert.Equal(t, tc.status, rr.Code)

				switch rr.Code {
				case 200:
					js, _ := simplejson.NewFromReader(rr.Body)
					ver, _ := js.Get("data").Get("version").Int()
					org, _ := js.Get("data").Get("organisation_id").String()
					debtorAcct, _ := js.Get("data").Get("attributes").Get("debtor_party").Get("account_name").String()
					assert.Equal(t, 2, ver)
					assert.Equal(t, "ffffffff-ffff-4fff-bfff-ffffffffffff", org)
					assert.Equal(t, "Mr Sending Test", debtorAcct)
					// Stuff that we've changed on the previous step should be wiped.
					amount, _ := js.Get("data").Get("attributes").Get("amount").String()
					ref, _ := js.Get("data").Get("attributes").Get("end_to_end_reference").String()
					fxRate, _ := js.Get("data").Get("attributes").Get("fx").Get("exchange_rate").String()
					assert.Equal(t, "", amount)
					assert.Equal(t, "", ref)
					assert.Equal(t, "0", fxRate)
				default:
					js, _ := simplejson.NewFromReader(rr.Body)
					_, ok := js.CheckGet("message")
					assert.Equal(t, true, ok)
				}
			})
		}
	})

	t.Run("DeletePayment", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/payments"+"/"+payment, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		s.Router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})
}
