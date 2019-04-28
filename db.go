package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var db *sqlx.DB

func (server *Server) OpenDB() {
	var err error

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		server.Config.DBHost, server.Config.DBPort, server.Config.DBUser, server.Config.DBPassword, server.Config.DBName, server.Config.DBSslMode)

	// This opens _and_ pings the database.
	server.DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) CloseDB() {
	server.DB.Close()
}

// Transform JSON to a database driver compatible type.
func (attributes PaymentAttributes) Value() (driver.Value, error) {
	obj, err := json.Marshal(attributes)
	return obj, err
}

// Transform `JSONB ([]byte)` that comes from a database into our type.
func (attributes *PaymentAttributes) Scan(data interface{}) error {
	source, ok := data.([]byte)
	if !ok {
		return errors.New("Database returned unexpected data.")
	}

	return json.Unmarshal(source, &attributes)
}

// GET /payments
func GetPayments(db *sqlx.DB) ([]Payment, error) {
	payments := []Payment{}

	query := `
		SELECT
			payments.id, payments.version, payments.organisation_id, payments.attributes,
			payment_types.type
		FROM
			payments
		INNER JOIN
			payment_types
		ON
			payments.type_id = payment_types.id`

	if err := db.Select(&payments, query); err != nil {
		return nil, err
	}
	return payments, nil
}

// GET /payments/:id
func (payment *Payment) Get(db *sqlx.DB) error {
	query := `
		SELECT
			payments.id, payments.version, payments.organisation_id, payments.attributes,
			payment_types.type
		FROM
			payments
		INNER JOIN
			payment_types
		ON
			payments.type_id = payment_types.id
		WHERE
			payments.id = $1`

	return db.Get(payment, query, payment.ID)
}

// POST /payments
func (payment *Payment) Save(db *sqlx.DB) error {
	query := `
		INSERT INTO
			payments (type_id, organisation_id, attributes)
		VALUES
			((SELECT id FROM payment_types WHERE type = $1), $2, $3)
		RETURNING
			id`

	return db.QueryRowx(query, payment.Type, payment.OrganisationID, payment.Attributes).Scan(&payment.ID)
}

// PUT, PATCH /payments/:id
func (payment *Payment) Overwrite(db *sqlx.DB) error {
	query := `
		UPDATE
			payments
		SET
			version = :version,
			type_id = (SELECT id FROM payment_types WHERE type = :type),
			organisation_id = :organisation_id,
			attributes = :attributes
		WHERE
			id = :id AND version = :version - 1`

	_, err := db.NamedExec(query, &payment)

	return err
}

// DELETE /payments/:id
func (payment *Payment) Remove(db *sqlx.DB) error {
	query := `
		DELETE FROM
			payments
		WHERE
			id = $1`

	_, err := db.Exec(query, payment.ID)

	return err
}
