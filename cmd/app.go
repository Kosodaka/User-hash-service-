package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	v1 "mainHashService/app/controller/http/v1"
	"mainHashService/app/entity"
	checkerRepo "mainHashService/app/repo/checker/impl"
	"mainHashService/app/repo/postgres/fetchdata"
	"mainHashService/app/repo/postgres/unhasher"
	s3repoImpl "mainHashService/app/repo/s3/impl"
	checkerUc "mainHashService/app/usecase/checker"
	unhasherUc "mainHashService/app/usecase/unhasher"
	"mainHashService/app/usecase/utills/writer"
	"mainHashService/app/usecase/utills/zipper"
	"mainHashService/pkg/httpserver"
	"mainHashService/pkg/logger"
	"mainHashService/pkg/postgres"
	"mainHashService/pkg/s3"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	Logger     logger.Logger
	CheckerUC  checkerUc.CheckerUCImpl
	UnhasherUC unhasherUc.UnhasherUCImpl
	HTTP       entity.HTTP
}

func NewDB(lg *logger.Logger, cfg *entity.Config) *pgxpool.Pool {
	dbCfg := cfg.GetConfigForDB()
	err := postgres.ParseConfig(&dbCfg)
	if err != nil {
		lg.Logger.Fatal().Msgf("could not parse postgres config %v: ", err)
	}

	pgxPool, err := postgres.NewPgxPool(&dbCfg)
	if err != nil {
		lg.Logger.Fatal().Msgf("could not setup pgx %v: ", err)
	}
	return pgxPool
}

func NewS3(lg *logger.Logger, cfg *entity.Config) (svc *s3.S3) {
	s3Cfg := cfg.GetConfigForS3()

	svc, err := s3.New(&s3Cfg)
	if err != nil {
		lg.Logger.Fatal().Msgf("coult not setup s3: %v", err)
	}

	return svc
}

func New(cfg *entity.Config) *App {
	lg := logger.NewConsoleLogger(cfg.Log.Level)
	db := NewDB(lg, cfg)
	minio := NewS3(lg, cfg)
	fileWriter := writer.NewFileWriter()
	zipper := zipper.NewZipper(lg, cfg.ZipPass)
	checkerRepository := checkerRepo.New(lg, cfg.Endpoint, cfg.HMACSecert)

	unhasherRepo := unhasher.New(lg, cfg.Endpoint)
	fetchRepo := fetchdata.New(lg, db)
	s3repo := s3repoImpl.New(lg, minio)

	return &App{
		Logger:     *lg,
		CheckerUC:  *checkerUc.New(lg, checkerRepository),
		UnhasherUC: *unhasherUc.New(lg, unhasherRepo, fetchRepo, s3repo, *fileWriter, *zipper),
		HTTP:       cfg.HTTP,
	}
}

func (a *App) SetupHTTP() (handler *gin.Engine) {
	handler = gin.New()

	v1.NewRouter(handler, &v1.RouterParams{
		Logger:    &a.Logger,
		CheckerUc: &a.CheckerUC,
		UnesherUc: &a.UnhasherUC,
	})
	return handler
}

func (a *App) Run() {
	handler := a.SetupHTTP()
	httpServer := httpserver.New(handler, a.HTTP.Port)

	// Graceful Shutdown
	go func() {
		err := a.gracefullyShutdown(httpServer)
		if err != nil {
			a.Logger.Error(err.Error())
		}

	}()

	// Start HTTP
	httpServer.Serve()
}

func (a *App) gracefullyShutdown(
	httpServer *httpserver.Server,
) (err error) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		a.Logger.Logger.Info().Msgf("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		a.Logger.Logger.Error().Msg("app - Run - httpServer.Notify: " + err.Error())
		return err
	}

	err = httpServer.Shutdown()
	if err != nil {
		a.Logger.Logger.Error().Msg("app - Run - httpServer.Shutdown: " + err.Error())
		return err
	}

	return nil
}
