package main

import (
	"compass-backend/pkg/api/bootstrap"
	"go.uber.org/fx"
)

func main() {
	fx.New(bootstrap.CreateApp()).Run()
}
