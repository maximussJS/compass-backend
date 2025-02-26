package main

import (
	"compass-backend/pkg/sender/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(bootstrap.CreateApp()).Run()
}
