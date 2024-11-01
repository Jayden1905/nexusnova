package user

import (
	"database/sql"
	"fmt"

	"golang.org/x/net/context"

	"github.com/jayden1905/nexusnova/cmd/pkg/database"
	"github.com/jayden1905/nexusnova/types"
)

type Store struct {
	db *database.Queries
}

// NewStore initializes the Store with the database queries
func NewStore(db *database.Queries) *Store {
	return &Store{db: db}
}

// GetUserByEmail fetches a user by email from the database
func (s *Store) GetUserByEmail(email string) (*database.User, error) {
	user, err := s.db.GetUserByEmail(context.Background(), email) // Use the SQLC-generated method
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByID fetches a user by ID from the database
func (s *Store) GetUserByID(id int32) (*database.User, error) {
	user, err := s.db.GetUserByID(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user in the database
func (s *Store) CreateUser(ctx context.Context, user *types.User) error {
	err := s.db.CreateUser(ctx, database.CreateUserParams{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	})
	if err != nil {
		return err
	}

	return nil
}
