package config

import (
	"flag"
	"os"

	"github.com/pkg/errors"
	"github.com/whatflix/pkg/envutils"
)

type Config struct {
	Environment string
	//HTTPHost            string
	HTTPPort            string
	DBConString         string
	DBName              string
	JWTAccessSecretKey  string
	JWTRefreshSecretKey string
	CacheServiceURL     string
	LoadBalancerURL     string
	LogServiceURL       string
}

func GetConfig() (*Config, error) {
	//Default configurations values
	cfg := &Config{
		Environment: envutils.Development,
		//HTTPHost:            "https://127.0.0.13",
		HTTPPort:            ":3000",
		DBConString:         "mongodb://localhost:27017/",
		DBName:              "whatflix",
		JWTAccessSecretKey:  "jdnfksdmfksd",
		JWTRefreshSecretKey: "mcmvmkmsdnfsdmfdsjf",
		CacheServiceURL:     "http://127.0.0.1:5000",
		LoadBalancerURL:     "http://127.0.0.1:2001",
		LogServiceURL:       "http://127.0.0.1:6000",
	}

	fs := flag.NewFlagSet("whatflix-webserver", flag.ExitOnError)
	fs.StringVar(&cfg.Environment, "environment", cfg.Environment, "Environment")
	//fs.StringVar(&cfg.HTTPHost, "http-host", cfg.HTTPHost, "Whatflix-Webserver-Host")
	fs.StringVar(&cfg.HTTPPort, "http-port", cfg.HTTPPort, "Whatflix-Webserver-Port")
	fs.StringVar(&cfg.DBConString, "Conn-String", cfg.DBConString, "MongoDB Con String")
	fs.StringVar(&cfg.DBName, "DB Name", cfg.DBName, "Whatflix DB Name")
	fs.StringVar(&cfg.JWTAccessSecretKey, "Access Secret Key", cfg.JWTAccessSecretKey, "JWT Access Secret Key")
	fs.StringVar(&cfg.JWTRefreshSecretKey, "Refresh Secret Key", cfg.JWTRefreshSecretKey, "JWT Refresh Secret Key")
	fs.StringVar(&cfg.CacheServiceURL, "Cache Service URL", cfg.CacheServiceURL, "Cache Service URL")
	fs.StringVar(&cfg.LoadBalancerURL, "Load Balancer URL", cfg.LoadBalancerURL, "Load Balancer URL")
	fs.StringVar(&cfg.LogServiceURL, "Log Service URL", cfg.LogServiceURL, "Log Service URL")
	_ = fs.Parse(os.Args[1:])
	err := checkConfigEnv(cfg)
	if err != nil {
		return nil, errors.WithMessage(err, "check configs")
	}
	return cfg, nil
}

func checkConfigEnv(cfg *Config) error {
	err := envutils.Check(cfg.Environment)
	if err != nil {
		return errors.WithMessage(err, "environment")
	}
	return nil
}
