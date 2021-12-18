package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Name        string
	Data        string
	Permissions int
	Aroundments []uuid.UUID
}

type UserStore interface {
	Create(ctx context.Context, u User) (*uuid.UUID, error)
	Read(ctx context.Context, uid uuid.UUID) (*User, error)
	Delete(ctx context.Context, uid uuid.UUID) error
	SearchUsers(ctx context.Context, s string) ([]*User, error)

	AddUserToAroundment(ctx context.Context, uid, aid uuid.UUID) error
	DeleteUserFromAroundment(ctx context.Context, uid, aid uuid.UUID) error
	ListUsersInAroundment(ctx context.Context, aid uuid.UUID) ([]*User, error)
	SearchAroundmentsByUserMembership(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error)
}

type Users struct {
	ustore UserStore
}

func NewUsers(ustore UserStore) *Users {
	return &Users{
		ustore: ustore,
	}
}

func (us *Users) Create(ctx context.Context, u User) (*User, error) {
	u.ID = uuid.New()
	id, err := us.ustore.Create(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("create user error: %w", err)
	}
	u.ID = *id
	return &u, nil
}

func (us *Users) Read(ctx context.Context, uid uuid.UUID) (*User, error) {
	u, err := us.ustore.Read(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("read user error: %w", err)
	}
	return u, nil
}

func (us *Users) Delete(ctx context.Context, uid uuid.UUID) error {
	u, err := us.ustore.Read(ctx, uid)
	if err != nil {
		return fmt.Errorf("search user error: %w", err)
	}
	return us.ustore.Delete(ctx, u.ID)
}

func (us *Users) SearchUsers(ctx context.Context, s string) ([]*User, error) {
	return us.ustore.SearchUsers(ctx, s)
}

func (us *Users) AddUserToAroundment(ctx context.Context, uid, aid uuid.UUID) error {
	return us.ustore.AddUserToAroundment(ctx, uid, aid)
}

func (us *Users) DeleteUserFromAroundment(ctx context.Context, uid, aid uuid.UUID) error {
	return us.ustore.DeleteUserFromAroundment(ctx, uid, aid)
}

func (us *Users) ListUsersInAroundment(ctx context.Context, aid uuid.UUID) ([]uuid.UUID, error) {
	users, err := us.ustore.ListUsersInAroundment(ctx, aid)
	if err != nil {
		return nil, err
	}
	ret := make([]uuid.UUID, len(users))
	for _, u := range users {
		ret = append(ret, u.ID)
	}
	return ret, nil
}

func (us *Users) SearchAroundmentsByUserMembership(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error) {
	return us.ustore.SearchAroundmentsByUserMembership(ctx, uid)
}
