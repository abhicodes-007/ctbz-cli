// Comando ctbz: CLI para a plataforma web da Contabilizei.
package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/edusouza/ctbz-cli/internal/cli"
)

// version é definida no build: go build -ldflags "-X main.version=1.2.3".
var version string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	code := cli.Execute(ctx, version, os.Args[1:])
	stop()
	os.Exit(code)
}
