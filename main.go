package main

import (
	"fmt"
	"os"

	"github.com/guysports/go-betfair-api/pkg/cmd"
	"github.com/guysports/go-betfair-api/pkg/types"

	"github.com/alecthomas/kong"
)

var cli struct {
	Test cmd.Test `cmd:"" help:"Login to Betfair and perform an operation"`
}

func main() {
	appkey := os.Getenv("BETFAIR_APP_KEY")
	if appkey == "" {
		fmt.Printf("BETFAIR_APP_KEY is required to be set in environment")
		os.Exit(1)
	}

	ctx := kong.Parse(&cli)
	err := ctx.Run(&types.Globals{
		AppKey: appkey,
	})
	ctx.FatalIfErrorf(err)

}
