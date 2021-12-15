package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"usernet/internal/api/handler"
	"usernet/internal/api/server"
	"usernet/internal/app/repos/user"
	"usernet/internal/app/starter"
	"usernet/internal/db/mem/usermemstore"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	ust := usermemstore.NewUsers()
	a := starter.NewApp(ust)
	us := user.NewUsers(ust)
	h := handler.NewRouter(us)
	srv := server.NewServer(":8000", h)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go a.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
