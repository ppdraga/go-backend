package usermemstore

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"strings"
	"sync"
	"usernet/internal/app/repos/user"
)

var _ user.UserStore = &Users{}

type Users struct {
	sync.Mutex
	m map[uuid.UUID]user.User
}

func NewUsers() *Users {
	return &Users{
		m: make(map[uuid.UUID]user.User),
	}
}

func (us *Users) Create(ctx context.Context, u user.User) (*uuid.UUID, error) {
	us.Lock()
	defer us.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	us.m[u.ID] = u
	return &u.ID, nil
}

func (us *Users) Read(ctx context.Context, uid uuid.UUID) (*user.User, error) {
	us.Lock()
	defer us.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	u, ok := us.m[uid]
	if ok {
		return &u, nil
	}
	return nil, sql.ErrNoRows
}

// не возвращает ошибку если не нашли
func (us *Users) Delete(ctx context.Context, uid uuid.UUID) error {
	us.Lock()
	defer us.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	delete(us.m, uid)
	return nil
}

func (us *Users) SearchUsers(ctx context.Context, s string) ([]*user.User, error) {
	us.Lock()
	defer us.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	ret := make([]*user.User, 0)

	for _, u := range us.m {
		if strings.Contains(u.Name, s) {
			ret = append(ret, &u)
		}
	}

	return ret, nil
}

func (us *Users) AddUserToAroundment(ctx context.Context, uid, aid uuid.UUID) error {
	us.Lock()
	defer us.Unlock()

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	u, ok := us.m[uid]
	if ok {
		u.Aroundments = append(u.Aroundments, aid)
		return nil
	}
	return sql.ErrNoRows

}

func (us *Users) DeleteUserFromAroundment(ctx context.Context, uid, aid uuid.UUID) error {
	us.Lock()
	defer us.Unlock()

	select {
	case <-ctx.Done():
		return nil
	default:
	}
	u, ok := us.m[uid]
	if ok {
		var new []uuid.UUID
		for _, id := range u.Aroundments {
			if id != aid {
				new = append(new, id)
			}
		}
		u.Aroundments = new
		return nil
	}
	return sql.ErrNoRows

}

func (us *Users) ListUsersInAroundment(ctx context.Context, aid uuid.UUID) ([]*user.User, error) {
	//    todo implement
	return nil, nil
}

func (us *Users) SearchAroundmentsByUserMembership(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error) {
	//    todo implement
	return nil, nil
}
