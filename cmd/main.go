package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-kit/log"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/markhaur/messapp-backend/pkg"
	"github.com/markhaur/messapp-backend/pkg/auth"
	"github.com/markhaur/messapp-backend/pkg/mysql"
	"github.com/markhaur/messapp-backend/pkg/reservations"
	"github.com/markhaur/messapp-backend/pkg/userlist"
)

func main() {
	logger := log.NewJSONLogger(os.Stderr)
	defer logger.Log("msg", "terminated")

	path, found := os.LookupEnv("MESSAPP_CONFIG_PATH")
	if found {
		if err := godotenv.Load(path); err != nil {
			logger.Log("msg", "could not load .env file", "path", path, "err", err)
		}
	} else {
		logger.Log("msg could not find MESSAPP_CONFIG_PATH env variable")
	}

	var config struct {
		ServerAddress              string        `envconfig:"SERVER_ADDRESS" default:"localhost:8085"`
		ServerWriteTimeout         time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"15s"`
		ServerReadTimeout          time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"15s"`
		ServerIdleTimeout          time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"60s"`
		GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT" default:"30s"`
		DBSource                   string        `envconfig:"DB_SOURCE"`
		DBConnectTimeout           time.Duration `envconfig:"DB_CONNECT_TIMEOUT"`
		OTELExporterJaegerEndpoint string        `envconfig:"OTEL_EXPORTER_JAEGER_ENDPOINT"`
	}
	if err := envconfig.Process("MESSAPP", &config); err != nil {
		logger.Log("msg", "could not load env vars", "err", err)
		os.Exit(1)
	}

	var userRepository pkg.UserRepository
	var reservationRepository pkg.ReservationRepository

	if config.DBSource != "" {
		ctx, cancel := context.WithTimeout(context.Background(), config.DBConnectTimeout)
		defer cancel()
		db, err := mysql.NewDB(ctx, config.DBSource)
		if err != nil {
			logger.Log("msg", "could not connect to mysql", "err", err)
			os.Exit(1)
		}

		if err = mysql.Migrate("file://pkg/mysql/migrations", db); err != nil {
			logger.Log("msg", "could not run mysql schema migrations", "err", err)
			os.Exit(1)
		}

		userRepository = mysql.NewUserRepository(db)
		reservationRepository = mysql.NewReservationRepository(db)

		defer func() {
			if err := db.Close(); err != nil {
				logger.Log("msg", "could not close db connection", "err", err)
			}
		}()
	}

	var userService userlist.Service
	userService = userlist.NewService(userRepository)
	userService = userlist.LoggingMiddleware(logger)(userService)

	var reservationService reservations.Service
	reservationService = reservations.NewService(reservationRepository)
	reservationService = reservations.LoggingMiddleware(logger)(reservationService)

	var authService auth.Service
	authService = auth.NewService(userRepository)
	authService = auth.LoggingMiddleware(logger)(authService)

	mux := http.NewServeMux()
	mux.Handle("/userlist/v1/", userlist.NewServer(userService, logger))
	mux.Handle("/resvlist/v1/", reservations.NewServer(reservationService, logger))
	mux.Handle("/auth/v1/", auth.NewServer(authService, logger))

	handler := enableCors(mux)

	server := &http.Server{
		Addr:         config.ServerAddress,
		WriteTimeout: config.ServerWriteTimeout,
		ReadTimeout:  config.ServerReadTimeout,
		IdleTimeout:  config.ServerIdleTimeout,
		Handler:      handler,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		logger.Log("transport", "http", "address", config.ServerAddress, "msg", "listening")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log("transport", "http", "address", config.ServerAddress, "msg", "failed", "err", err)
			sig <- os.Interrupt
		}
	}()

	logger.Log("received", <-sig, "msg", "terminating")
	if err := server.Shutdown(context.Background()); err != nil {
		logger.Log("msg", "could not shutdown http server", "err", err)
	}

}

func enableCors(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		handler.ServeHTTP(w, req)
	})
}
