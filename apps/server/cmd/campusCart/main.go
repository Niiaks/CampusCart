package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Niiaks/campusCart/internal/config"
	"github.com/Niiaks/campusCart/internal/database"
	"github.com/Niiaks/campusCart/internal/handler"
	"github.com/Niiaks/campusCart/internal/lib/file"
	"github.com/Niiaks/campusCart/internal/logger"
	"github.com/Niiaks/campusCart/internal/middleware"
	"github.com/Niiaks/campusCart/internal/repository"
	"github.com/Niiaks/campusCart/internal/router"
	"github.com/Niiaks/campusCart/internal/server"
	"github.com/Niiaks/campusCart/internal/service"
	"github.com/cloudinary/cloudinary-go/v2"
)

const DefaultContextTimeout = 30

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize New Relic logger service
	loggerService := logger.NewLoggerService(cfg.Observability)
	defer loggerService.Shutdown()

	log := logger.NewLoggerWithService(cfg.Observability, loggerService)

	if cfg.Primary.Env != "development" {
		if err := database.Migrate(context.Background(), &log, cfg); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate database")
		}
	}
	// Initialize server
	srv, err := server.New(cfg, &log, loggerService)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize server")
	}

	// Initialize cloudinary
	cld, err := cloudinary.NewFromParams(cfg.Cloudinary.CloudName, cfg.Cloudinary.ApiKey, cfg.Cloudinary.ApiSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize cloudinary")
	}

	fileClient := file.NewClient(cld, &log, cfg.Primary.Env)

	// Initialize repositories
	userRepo := repository.NewUserRepository(srv.DB.Pool)
	sessionRepo := repository.NewSessionRepository(srv.DB.Pool)
	categoryRepo := repository.NewCategoryRepository(srv.DB.Pool)
	listingRepo := repository.NewListingRepository(srv.DB.Pool)
	brandRepo := repository.NewBrandRepository(srv.DB.Pool)

	// Initialize services
	authService := service.NewAuthService(userRepo, sessionRepo, srv.Job)
	categoryService := service.NewCategoryService(categoryRepo, fileClient)
	listingService := service.NewListingService(listingRepo, fileClient)
	brandService := service.NewBrandService(brandRepo, fileClient)

	h := handler.NewHandlers(srv, authService, categoryService, listingService, brandService)
	mw := middleware.NewMiddlewares(srv, sessionRepo, brandRepo)
	r := router.NewRouter(h, mw)

	srv.SetupHTTPServer(r)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	// Start server
	go func() {
		if err = srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	<-ctx.Done()
	stop()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server exited properly")
}
