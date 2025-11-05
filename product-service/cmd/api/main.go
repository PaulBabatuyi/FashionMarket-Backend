package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PaulBabatuyi/FashionMarket-Backend/product-service/internal/cache"
	"github.com/PaulBabatuyi/FashionMarket-Backend/product-service/internal/data"
	"github.com/PaulBabatuyi/FashionMarket-Backend/product-service/internal/jsonlog"
	"github.com/PaulBabatuyi/FashionMarket-Backend/product-service/internal/jwt"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	jwt struct {
		publicKeyPath    string
		expectedIssuer   string
		expectedAudience string
	}
	userService struct {
		url string
	}
	cache struct {
		userTTL time.Duration
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config       config
	logger       *jsonlog.Logger
	models       data.Models
	wg           sync.WaitGroup
	jwtValidator *jwt.JWTValidator
	httpClient   *http.Client
	userCache    *cache.UserCache
}

func main() {
	var cfg config

	// Server config
	flag.IntVar(&cfg.port, "port", 5000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment")

	// Database config
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("FASHION_PRODUCT_SERVICE"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")

	// Rate limiter config
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// JWT config
	flag.StringVar(&cfg.jwt.publicKeyPath, "jwt-public-key", os.Getenv("JWT_PUBLIC_KEY"), "Path to JWT public key")
	flag.StringVar(&cfg.jwt.expectedIssuer, "jwt-issuer", "github.com/PaulBabatuyi/FashionMarket-Backend/user-service", "Expected JWT issuer")
	flag.StringVar(&cfg.jwt.expectedAudience, "jwt-audience", "github.com/PaulBabatuyi/FashionMarket-Backend/*", "Expected JWT audience")

	// User service config
	flag.StringVar(&cfg.userService.url, "user-service-url", os.Getenv("USER_SERVICE_URL"), "User service URL")

	// Cache config
	flag.DurationVar(&cfg.cache.userTTL, "cache-user-ttl", 5*time.Minute, "User cache TTL")

	// Use the flag.Func() function to process the -cors-trusted-origins command line
	// flag.
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	if cfg.userService.url == "" {
		logger.PrintFatal(errors.New("USER_SERVICE_URL is required"), nil)
	}

	// Database connection
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	// Initialize JWT validator
	jwtValidator, err := jwt.NewJWTValidator(jwt.Config{
		PublicKeyPath:    cfg.jwt.publicKeyPath,
		ExpectedIssuer:   cfg.jwt.expectedIssuer,
		ExpectedAudience: cfg.jwt.expectedAudience,
	})
	if err != nil {
		logger.PrintFatal(err, map[string]string{
			"component": "jwt_validator",
		})
	}
	logger.PrintInfo("JWT validator initialized", nil)

	// HTTP client for service-to-service calls
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Initialize user cache (THIS WAS MISSING!)
	userCache := cache.NewUserCache(cfg.cache.userTTL)
	logger.PrintInfo("user cache initialized", map[string]string{
		"ttl": cfg.cache.userTTL.String(),
	})

	app := &application{
		config:       cfg,
		logger:       logger,
		models:       data.NewModels(db),
		jwtValidator: jwtValidator,
		httpClient:   httpClient,
		userCache:    userCache,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
