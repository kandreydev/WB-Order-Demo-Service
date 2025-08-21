package main

import (
	"log/slog"
	"os"

	"github.com/GkadyrG/L0/backend/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("app run", slog.Any("err", err))
		os.Exit(1)
	}
}
