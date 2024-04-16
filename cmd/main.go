package main

import (
	"comixsearch/internal/adapters/pgstorage"
	internal "comixsearch/internal/app"
	"comixsearch/pkg/stem"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func captureSigQuit(ctx context.Context) func() error {
	return func() error {
		sigQuit := make(chan os.Signal, 1)
		signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
		signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v ", s)
		case <-ctx.Done():
			return nil
		}
	}
}

func main() {
	// s := &internal.SearchApp{
	// 	N: stem.NewStem("english", true),
	// 	P: &xkcdparse.XkcdParser{
	// 		LastId: 2890,
	// 	},
	// }
	// s.Run()

	// st, _ := pgstorage.NewStorage(context.Background(), "postgresql://postgres:1234@127.0.0.1:5432/postgres")
	// comixes := []models.Comix{
	// 	{
	// 		Id:         1,
	// 		Title:      "Bla",
	// 		Content:    "",
	// 		AltContent: "",
	// 		Link:       "",
	// 	},
	// 	{
	// 		Id:         2,
	// 		Title:      "Blabla",
	// 		Content:    "",
	// 		AltContent: "",
	// 		Link:       "",
	// 	},
	// }
	// st.Write(context.Background(), comixes)

	// fmt.Println(st.Read(context.Background(), []string{}))
	ctx := context.Background()
	st, _ := pgstorage.NewStorage(ctx, "postgresql://postgres:1234@127.0.0.1:5432/postgres")

	s := internal.NewApp(stem.NewStem("english", true), &internal.XkcdFetcher{
		LastId: 2900,
	}, st, nil, runtime.NumCPU())

	eg, _ := errgroup.WithContext(ctx)
	// eg.Go(captureSigQuit(ctx))
	eg.Go(s.Run(ctx))
	eg.Wait()
}
