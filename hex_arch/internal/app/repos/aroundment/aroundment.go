package aroundment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Aroundment struct {
	ID   uuid.UUID
	Name string
	Type AroundmentType
}

type AroundmentType string

const (
	AroundmentTypeProject      AroundmentType = "Project"
	AroundmentTypeOrganization AroundmentType = "Organization"
	AroundmentTypeGroup        AroundmentType = "Corp group"
	AroundmentTypeCommunity    AroundmentType = "Community"
)

type AroundmentStore interface {
	Create(ctx context.Context, a Aroundment) (*uuid.UUID, error)
	Read(ctx context.Context, aid uuid.UUID) (*Aroundment, error)
	Delete(ctx context.Context, aid uuid.UUID) error
	SearchAroundments(ctx context.Context, s string) ([]*Aroundment, error)
}

type Aroundments struct {
	astore AroundmentStore
}

func NewAroundments(astore AroundmentStore) *Aroundments {
	return &Aroundments{
		astore: astore,
	}
}

func (as *Aroundments) Create(ctx context.Context, a Aroundment) (*Aroundment, error) {
	a.ID = uuid.New()
	id, err := as.astore.Create(ctx, a)
	if err != nil {
		return nil, fmt.Errorf("create aroundment error: %w", err)
	}
	a.ID = *id
	return &a, nil
}

func (as *Aroundments) Read(ctx context.Context, aid uuid.UUID) (*Aroundment, error) {
	u, err := as.astore.Read(ctx, aid)
	if err != nil {
		return nil, fmt.Errorf("read aroundment error: %w", err)
	}
	return u, nil
}

func (as *Aroundments) Delete(ctx context.Context, aid uuid.UUID) error {
	a, err := as.astore.Read(ctx, aid)
	if err != nil {
		return fmt.Errorf("search aroundment error: %w", err)
	}
	return as.astore.Delete(ctx, a.ID)
}

func (as *Aroundments) SearchAroundments(ctx context.Context, s string) ([]*Aroundment, error) {

	return as.astore.SearchAroundments(ctx, s)
}
