package config

import (
	"flag"
	"os"
)

type loadBalConfig struct {
	HTTPPort      string
	IntHTTPPort   string
	LogServiceURL string
}

func GetLoadBalConfig() (*loadBalConfig, error) {
	//Default configurations values
	lbcfg := &loadBalConfig{
		HTTPPort:      ":2000",
		IntHTTPPort:   ":2001",
		LogServiceURL: "http://127.0.0.10:6000",
	}

	fs := flag.NewFlagSet("whatflix-loadbalancerservice", flag.ExitOnError)
	fs.StringVar(&lbcfg.HTTPPort, "loadbalancerservice", lbcfg.HTTPPort, "Address of the load balancer service provider")
	fs.StringVar(&lbcfg.IntHTTPPort, "loadbalancerservice-Internal Endpoint Port", lbcfg.IntHTTPPort, "Address of the load balancer service provider - Internal endpoint port")
	fs.StringVar(&lbcfg.LogServiceURL, "logserviceURL", lbcfg.LogServiceURL, "Log service URL")
	_ = fs.Parse(os.Args[1:])
	return lbcfg, nil
}
