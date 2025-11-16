package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/marrria_mme/pr-reviewer-service/config"
	"github.com/marrria_mme/pr-reviewer-service/internal/repository"
	prRepo "github.com/marrria_mme/pr-reviewer-service/internal/repository/pr"
	teamRepo "github.com/marrria_mme/pr-reviewer-service/internal/repository/team"
	userRepo "github.com/marrria_mme/pr-reviewer-service/internal/repository/user"
	"github.com/marrria_mme/pr-reviewer-service/internal/transport/middleware"
	prTransport "github.com/marrria_mme/pr-reviewer-service/internal/transport/pr"
	teamTransport "github.com/marrria_mme/pr-reviewer-service/internal/transport/team"
	userTransport "github.com/marrria_mme/pr-reviewer-service/internal/transport/user"
	prUs "github.com/marrria_mme/pr-reviewer-service/internal/usecase/pr"
	teamUs "github.com/marrria_mme/pr-reviewer-service/internal/usecase/team"
	userUs "github.com/marrria_mme/pr-reviewer-service/internal/usecase/user"
)

type App struct {
	conf   *config.Config
	logger *logrus.Logger
	db     *sql.DB
	router *mux.Router
}

func NewApp(conf *config.Config) (*App, error) {
	logger := logrus.New()

	str, err := repository.GetConnectionString(conf.DBConfig)
	if err != nil {
		return nil, fmt.Errorf("connection string error: %w", err)
	}
	db, err := sql.Open("postgres", str)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	config.ConfigureDB(db, conf.DBConfig)

	teamRepository := teamRepo.NewTeamRepository(db)
	userRepository := userRepo.NewUserRepository(db)
	prRepository := prRepo.NewPRRepository(db)

	teamUsecase := teamUs.NewTeamUsecase(teamRepository)
	userUsecase := userUs.NewUserUsecase(userRepository)
	prUsecase := prUs.NewPRUsecase(prRepository, userRepository)

	teamHandler := teamTransport.NewTeamHandler(teamUsecase)
	userHandler := userTransport.NewUserHandler(userUsecase)
	prHandler := prTransport.NewPRHandler(prUsecase)

	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)

	router.HandleFunc("/team/add", teamHandler.CreateTeam).Methods("POST")
	router.HandleFunc("/team/get", teamHandler.GetTeam).Methods("GET")

	router.HandleFunc("/users/setIsActive", userHandler.SetUserActivity).Methods("POST")
	router.HandleFunc("/users/getReview", userHandler.GetUserReviewPRs).Methods("GET")

	router.HandleFunc("/pullRequest/create", prHandler.CreatePR).Methods("POST")
	router.HandleFunc("/pullRequest/merge", prHandler.MergePR).Methods("POST")
	router.HandleFunc("/pullRequest/reassign", prHandler.ReassignReviewer).Methods("POST")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy"}`))
	}).Methods("GET")

	app := &App{
		conf:   conf,
		logger: logger,
		db:     db,
		router: router,
	}

	return app, nil
}

func (a *App) Run() {

	server := &http.Server{
		Handler:      a.router,
		Addr:         fmt.Sprintf(":%s", a.conf.ServerConfig.Port),
		WriteTimeout: a.conf.ServerConfig.WriteTimeout,
		ReadTimeout:  a.conf.ServerConfig.ReadTimeout,
		IdleTimeout:  a.conf.ServerConfig.IdleTimeout,
	}

	a.logger.Infof("starting server on port %s", a.conf.ServerConfig.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.Fatalf("server failed: %v", err)
	}
}
