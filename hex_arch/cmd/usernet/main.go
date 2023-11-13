package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"usernet/internal/api/handler"
	"usernet/internal/api/server"
	"usernet/internal/app/repos/aroundment"
	"usernet/internal/app/repos/user"
	"usernet/internal/app/starter"
	"usernet/internal/db/mem/aroundmentmemstore"
	"usernet/internal/db/mem/usermemstore"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	ust := usermemstore.NewUsers()
	ast := aroundmentmemstore.NewAroundments()

	a := starter.NewApp(ust, ast)

	us := user.NewUsers(ust)
	as := aroundment.NewAroundments(ast)

	h := handler.NewRouter(us, as)
	srv := server.NewServer(":8000", h)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
