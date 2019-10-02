package types

import (
	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	ID        int
	Username  string
	FirstName string
	LastName  string
	Password  string
}
