package main

import (
	"arch/ikeppu/github.com/configs"
	"arch/ikeppu/github.com/internal/auth"
	"arch/ikeppu/github.com/internal/link"
	"arch/ikeppu/github.com/internal/stat"
	"arch/ikeppu/github.com/internal/user"
	"arch/ikeppu/github.com/pkg/db"
	"arch/ikeppu/github.com/pkg/event"
	"arch/ikeppu/github.com/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()

	database := db.NewDb(conf)

	router := http.NewServeMux()

	eventBus := event.NewEventBus()

	// Repositories
	linkRepo := link.NewLinkRepository(database)
	userRepo := user.NewUserRepository(database)
	statRepo := stat.NewStatRepository(database)

	// Services
	authService := auth.NewAuthRepository(userRepo)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepo,
	})

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		AuthService: authService,
		Config:      conf,
	})

	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepo,
		EventBus:       eventBus,
		Config:         conf,
	})

	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepo,
		Config:         conf,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.Cors,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	go statService.AddClick()

	fmt.Println("Server is listening on port 8081")

	err := server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}

}
