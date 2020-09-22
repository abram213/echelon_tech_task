package storage

import (
	"context"
	"github.com/volatiletech/authboss/v3"
	"log"
	"task_ws_et/models"
)

type Storage struct {
	Users  map[string]models.User
}

func New() *Storage {
	return &Storage{
		Users: map[string]models.User{
			"admin@admin.com": {
				Email:              "admin@admin.com",
				Password:           "$2a$10$XtW/BrS5HeYIuOCXYe8DFuInetDMdaarMUJEOg/VA/JAIDgw3l4aG", // pass = 1234
				Role: 				"admin",
			},
			"user@user.com": {
				Email:              "user@user.com",
				Password:           "$2a$10$XtW/BrS5HeYIuOCXYe8DFuInetDMdaarMUJEOg/VA/JAIDgw3l4aG", // pass = 1234
				Role: 				"user",
			},
		},
	}
}

// Save the user. Needed to implement authboss ServerStorer interface
func (s Storage) Save(ctx context.Context, user authboss.User) error {
	u := user.(*models.User)
	s.Users[u.Email] = *u

	return nil
}

// Load the user. Needed to implement authboss ServerStorer interface
func (s Storage) Load(ctx context.Context, key string) (user authboss.User, err error) {
	u, ok := s.Users[key]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	return &u, nil
}

// New user. Needed to implement authboss CreatingServerStorer interface
func (s Storage) New(ctx context.Context) authboss.User {
	return &models.User{}
}

// Create the user. Needed to implement authboss CreatingServerStorer interface
func (s Storage) Create(ctx context.Context, user authboss.User) error {
	u := user.(*models.User)

	if _, ok := s.Users[u.Email]; ok {
		return authboss.ErrUserFound
	}

	log.Printf("Created new user: %s\n", u.Email)
	s.Users[u.Email] = *u
	return nil
}