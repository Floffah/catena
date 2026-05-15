package main

import (
	"github.com/alecthomas/kong"
)

var CLI struct {
	Login struct {
		Instance string `arg:"" default:"http://localhost:3000" help:"URL of the Catena instance to authenticate with"`
	} `cmd:"" help:"Authenticate Catena CLI & Git with a Catena instance"`
}

func main() {
	ctx := kong.Parse(&CLI)

	// TODO: make cli login work
	// Probably need to reimplement using an oauth flow for clerk instead of the ticketing strategy as it might be deprecated (but documentation is sparse so unclear)

	switch ctx.Command() {
	case "login":
		println("Not implemented yet")
	default:
		panic(ctx.Command())
	}
}
