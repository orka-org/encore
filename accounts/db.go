package accounts

import "encore.dev/storage/sqldb"

var usersDB = sqldb.NewDatabase("accounts", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})
