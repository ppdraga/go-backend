package starter

import (
	"context"
	"sync"
	"usernet/internal/app/repos/aroundment"

	"usernet/internal/app/repos/user"
)

type App struct {
	us *user.Users
	as *aroundment.Aroundments
}

func NewApp(ust user.UserStore, ast aroundment.AroundmentStore) *App {
	a := &App{
		us: user.NewUsers(ust),
		as: aroundment.NewAroundments(ast),
	}
	return a
}

type APIServer interface {
	Start(us *user.Users, as *aroundment.Aroundments)
	Stop()
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs APIServer) {
	defer wg.Done()
	hs.Start(a.us, a.as)
	<-ctx.Done()
	hs.Stop()
}
