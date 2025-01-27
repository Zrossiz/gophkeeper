package postgres

type UserStorage struct {
}

func NewUserStorage() *UserStorage {
	return &UserStorage{}
}
