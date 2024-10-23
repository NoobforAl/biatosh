package store

import (
	"biatosh/contract"
	"biatosh/entity"
	"biatosh/store/database"
	"context"
	"database/sql"
	_ "embed"
)

type Store struct {
	db *database.Queries

	log contract.Logger
}

//go:embed sql/schema.sql
var ddl string

func New(log contract.Logger) contract.Store {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}

	return &Store{
		db:  database.New(db),
		log: log,
	}
}

func (s *Store) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	user.Password = genHashPassword(user.Password)

	newUser, err := s.db.CreateUser(ctx, database.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Name:     user.Name,
		Password: user.Password,
	})
	if err != nil {
		return nil, err
	}

	return convertUserToEntityUser(&newUser), nil
}

func (s *Store) GetUser(ctx context.Context, id int) (*entity.User, error) {
	user, err := s.db.GetUser(ctx, int64(id))
	if err != nil {
		return nil, err
	}

	return convertUserToEntityUser(&user), nil
}

func (s *Store) GetUsers(ctx context.Context) ([]*entity.User, error) {
	panic("implement me")
	return nil, nil
}

func (s *Store) LoginUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	user.Password = genHashPassword(user.Password)

	finedUser, err := s.db.GetUserByUsernamePassword(ctx, database.GetUserByUsernamePasswordParams{
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		return nil, err
	}

	return convertUserToEntityUser(&finedUser), nil
}

func (s *Store) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	panic("implement me")
	return nil, nil
}

func (s *Store) DeleteUser(ctx context.Context, id int) error {
	panic("implement me")
	return nil
}
