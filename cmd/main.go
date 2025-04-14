package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/magneless/pvz/internal/app"
	"github.com/magneless/pvz/internal/delivery"
	"github.com/magneless/pvz/internal/middleware/auth"
)

func main() {
	app := app.New()
	app.Logger.Info("starting pvz", slog.String("env", app.Config.Env))

	router := chi.NewRouter()

	PVZHandler := delivery.NewPVZHandler(app.PVZUsecase, app.Logger)
	ReceptionHandler := delivery.NewReceptionHandler(app.ReceptionUsecase, app.Logger)
	DumyyLoginHandler := delivery.NewDumyyLoginHandler(app.DummyLoginUsecase, app.Logger)
	ProductHandler := delivery.NewProductHandler(app.ProductUsecase, app.Logger)
	
	router.Post("/dummyLogin", DumyyLoginHandler.CreateToken())

	router.Group(func(r chi.Router) {
		r.Use(auth.New(app.Logger))
	
		r.Post("/pvz", PVZHandler.CreatePVZ())
		r.Get("/pvz", PVZHandler.GetPVZFullInfo())
		r.Post("/receptions", ReceptionHandler.CreateReception())
		r.Post("/pvz/{pvzId}/close_last_reception", ReceptionHandler.CloseReception())
		r.Post("/products", ProductHandler.CreateProduct())
		r.Post("/pvz/{pvzId}/delete_last_product", ProductHandler.DeleteProduct())
	})


	server := &http.Server{
		Addr:         app.Config.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  app.Config.HTTPServer.Timeout,
		IdleTimeout:  app.Config.HTTPServer.IdleTimeout,
		WriteTimeout: app.Config.HTTPServer.Timeout,
	}

	app.Logger.Info("server starting", slog.String("addr", app.Config.HTTPServer.Address))

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		app.Logger.Error("server error", slog.String("err", err.Error()))
	}
}

