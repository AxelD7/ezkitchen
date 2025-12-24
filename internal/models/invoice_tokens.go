package models

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"
)

type InvoiceToken struct {
	InvoiceTokenID int
	EstimateID     int
	TokenHash      string
	ExpiresAt      time.Time
	UsedAt         sql.NullTime
	CreatedAt      time.Time
}

type InvoiceTokenModel struct {
	DB *sql.DB
}

func (m *InvoiceTokenModel) Insert(estimateID int, expiresAt time.Time) (string, error) {
	rawToken, tokenHash, err := generateInvoiceToken()
	if err != nil {
		return "", err
	}

	stmt := `
		INSERT INTO invoice_access_tokens (estimate_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`

	_, err = m.DB.Exec(stmt, estimateID, tokenHash, expiresAt)
	if err != nil {
		return "", err
	}

	return rawToken, nil
}

func (m *InvoiceTokenModel) GetByRawToken(rawToken string) (*InvoiceToken, error) {
	sum := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(sum[:])

	stmt := `
		SELECT invoice_token_id, estimate_id, token_hash, expires_at, used_at, created_at
		FROM invoice_access_tokens
		WHERE token_hash = $1
	`

	var it InvoiceToken

	err := m.DB.QueryRow(stmt, tokenHash).Scan(
		&it.InvoiceTokenID,
		&it.EstimateID,
		&it.TokenHash,
		&it.ExpiresAt,
		&it.UsedAt,
		&it.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return &it, nil
}

func (m *InvoiceTokenModel) MarkUsed(id int) error {
	return m.markUsed(m.DB, id)
}

func (m *InvoiceTokenModel) MarkUsedTx(tx *sql.Tx, id int) error {
	return m.markUsed(tx, id)
}

func (m *InvoiceTokenModel) markUsed(exec executor, id int) error {
	stmt := `
		UPDATE invoice_access_tokens
		SET used_at = $2
		WHERE invoice_token_id = $1
	`

	result, err := exec.Exec(stmt, id, time.Now())
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoRecord
	}

	return nil
}

func generateInvoiceToken() (raw string, hash string, err error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}

	raw = base64.RawURLEncoding.EncodeToString(b)

	sum := sha256.Sum256([]byte(raw))
	hash = hex.EncodeToString(sum[:])

	return raw, hash, nil

}
