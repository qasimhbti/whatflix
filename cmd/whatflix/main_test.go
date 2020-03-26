package main

import (
	"io/ioutil"
	"log"

	"github.com/davecgh/go-spew/spew"
)

func init() {
	log.SetOutput(ioutil.Discard)
	spew.Config.Indent = "\t"
}
