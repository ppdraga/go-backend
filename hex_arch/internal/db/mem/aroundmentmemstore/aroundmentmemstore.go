package aroundmentmemstore

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"sync"
	"usernet/internal/app/repos/aroundment"
)

var _ aroundment.AroundmentStore = &Aroundments{}

type Aroundments struct {
	sync.Mutex
	m map[uuid.UUID]aroundment.Aroundment
}

func NewUsers() *Aroundments {
	return &Aroundments{
		m: make(map[uuid.UUID]aroundment.Aroundment),
	}
}

func (as *Aroundments) Create(ctx context.Context, a aroundment.Aroundment) (*uuid.UUID, error) {
	as.Lock()
	defer as.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	as.m[a.ID] = a
	return &a.ID, nil
}

func (as *Aroundments) Read(ctx context.Context, uid uuid.UUID) (*aroundment.Aroundment, error) {
	as.Lock()
	defer as.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	a, ok := as.m[uid]
	if ok {
		return &a, nil
	}
	return nil, sql.ErrNoRows
}

// не возвращает ошибку если не нашли
func (as *Aroundments) Delete(ctx context.Context, aid uuid.UUID) error {
	as.Lock()
	defer as.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	delete(as.m, aid)
	return nil
}
