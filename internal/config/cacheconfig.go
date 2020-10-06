package config

import (
	"flag"
	"os"
)

type cacheConfig struct {
	HTTPPort string
}

func GetCacheConfig() (*cacheConfig, error) {
	//Default configurations values
	ccfg := &cacheConfig{
		HTTPPort: ":5000",
	}

	fs := flag.NewFlagSet("whatflix-cacheservice", flag.ExitOnError)
	fs.StringVar(&ccfg.HTTPPort, "cacheservice", ccfg.HTTPPort, "Address of the caching service provider")
	_ = fs.Parse(os.Args[1:])
	return ccfg, nil
}
