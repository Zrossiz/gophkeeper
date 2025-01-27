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

func New() *Storage {
	return &Storage{
		User:     *NewUserStorage(),
		Card:     *NewCardStorage(),
		LogoPass: *NewLogoPassStorage(),
		Binary:   *NewBinaryStorage(),
	}
}

func Connect(DBURI string) (*sql.DB, error) {
	db, err := sql.Open("postgres", DBURI)
	if err != nil {
		return nil, err
	}

	return db, nil
}
