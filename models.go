package main

import (
	"github.com/go-openapi/strfmt"
)

// First, I was hoping I can import and use
// https://github.com/form3tech-oss/go-form3/blob/master/models/payment.go
// as ready-to-use type definitions, but then I've stumbled upon discrepancies
// between those efinitions and sample data
// (http://mockbin.org/bin/41ca3269-d8c4-4063-9fd5-f306814ff03f).
//
// Additionally, there's a need to extend models with tags for `sql` package.

type Payment struct {
	ID             strfmt.UUID       `json:"id" db:"id"`
	Type           string            `json:"type,omitempty" db:"type"`
	Version        uint64            `json:"version" db:"version"`
	OrganisationID strfmt.UUID       `json:"organisation_id" db:"organisation_id"`
	Attributes     PaymentAttributes `json:"attributes" db:"attributes"`
}

type PaymentAttributes struct {
	// `ProcessingDate` is `string` and not `time.Time` because it's stored as a string in a database.
	// It is stored as a string in a database because the whole attributes property is stored as JSONB.
	// I could have destructured it better in the database, but I need some business context for this ;)
	Amount               float64     `json:"amount,omitempty,string"`
	BeneficiaryParty     Beneficiary `json:"beneficiary_party"`
	ChargesInformation   Charges     `json:"charges_information,omitempty"`
	Currency             string      `json:"currency"`
	DebtorParty          Debtor      `json:"debtor_party"`
	EndToEndReference    string      `json:"end_to_end_reference"`
	FX                   Exchange    `json:"fx"`
	SponsorParty         Sponsor     `json:"sponsor_party"`
	NumericReference     string      `json:"numeric_reference"`
	PaymentID            string      `json:"payment_id"`
	PaymentPurpose       string      `json:"payment_purpose"`
	PaymentScheme        string      `json:"payment_scheme"`
	PaymentType          string      `json:"payment_type"`
	ProcessingDate       string      `json:"processing_date"`
	Reference            string      `json:"reference"`
	SchemePaymentSubType string      `json:"scheme_payment_sub_type"`
	SchemePaymentType    string      `json:"scheme_payment_type"`
}

type Beneficiary struct {
	AccountName       string `json:"account_name,omitempty"`
	AccountNumber     string `json:"account_number,omitempty"`
	AccountNumberCode string `json:"account_number_code,omitempty"`
	AccountType       uint64 `json:"account_type,omitempty"`
	Address           string `json:"address,omitempty"`
	Name              string `json:"name,omitempty"`
	BankID            string `json:"bank_id,omitempty"`
	BankIDCode        string `json:"bank_id_code,omitempty"`
}

type Charges struct {
	BearerCode              string   `json:"bearer_code"`
	SenderCharges           []Charge `json:"sender_charges"`
	ReceiverChargesAmount   float64  `json:"receiver_charges_amount,string"`
	ReceiverChargesCurrency string   `json:"receiver_charges_currency"`
}

type Charge struct {
	Amount   float64 `json:"amount,string"`
	Currency string  `json:"currency"`
}

type Debtor struct {
	AccountName       string `json:"account_name,omitempty"`
	AccountNumber     string `json:"account_number,omitempty"`
	AccountNumberCode string `json:"account_number_code,omitempty"`
	Address           string `json:"address,omitempty"`
	Name              string `json:"name,omitempty"`
	BankID            string `json:"bank_id,omitempty"`
	BankIDCode        string `json:"bank_id_code,omitempty"`
}

type Exchange struct {
	ContractReference string  `json:"contract_reference"`
	OriginalAmount    float64 `json:"original_amount,string"`
	OriginalCurrency  string  `json:"original_currency"`
	ExchangeRate      float64 `json:"exchange_rate,string"`
}

type Sponsor struct {
	AccountNumber string `json:"account_number"`
	BankID        string `json:"bank_id"`
	BankIDCode    string `json:"bank_id_code"`
}
