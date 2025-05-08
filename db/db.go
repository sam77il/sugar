package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Account struct {
	Email    string
	Password string
}

func GetAllAccounts(db *pgxpool.Pool, ctx context.Context) []Account {
	var accounts []Account
	rows, err := db.Query(ctx, "SELECT email, password FROM accounts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var account Account

		rows.Scan(&account.Email, &account.Password)
		accounts = append(accounts, account)
	}
	return accounts
}