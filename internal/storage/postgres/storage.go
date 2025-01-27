package postgres

import (
	"database/sql"
)

type Storage struct {
	Binary   BinaryStorage
	Card     CardStorage
	LogoPass LogoPassStorage
	User     UserStorage
}

func New(conn *sql.DB) *Storage {
	return &Storage{
		User:     *NewUserStorage(conn),
		Card:     *NewCardStorage(conn),
		LogoPass: *NewLogoPassStorage(conn),
		Binary:   *NewBinaryStorage(conn),
	}
}

func Connect(DBURI string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DBURI)
	if err != nil {
		return nil, err
	}

	return db, nil
}
