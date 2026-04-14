package store

import (
	"database/sql"
	"time"

	"github.com/jingxinwangdev/go-prject/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	InsertToken(token *tokens.Token) error
	CreateToken(userID int64, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int64, scope string) error
}

func (pg *PostgresTokenStore) CreateToken(userID int64, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	pg.InsertToken(token)
	return token, err
}

func (pg *PostgresTokenStore) InsertToken(token *tokens.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
	`
	_, err := pg.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (pg *PostgresTokenStore) DeleteAllTokensForUser(userID int64, scope string) error {
	query := `
		DELETE FROM tokens 
		WHERE user_id = $1 AND scope = $2
	`
	_, err := pg.db.Exec(query, userID, scope)
	return err
}
