package config

import (
	"flag"
	"os"
)

type logConfig struct {
	HTTPPort string
}

func GetLogConfig() (*logConfig, error) {
	//Default configurations values
	lcfg := &logConfig{
		HTTPPort: ":6000",
	}

	fs := flag.NewFlagSet("whatflix-logservice", flag.ExitOnError)
	fs.StringVar(&lcfg.HTTPPort, "logservice", lcfg.HTTPPort, "Address of the log service provider")
	_ = fs.Parse(os.Args[1:])
	return lcfg, nil
}
