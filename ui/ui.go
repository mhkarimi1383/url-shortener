package ui

import (
	"embed"
	"os"

	"github.com/labstack/echo/v4"
)

func init() {
	os.Setenv("BASE_URL", "/ui")
}

//go:generate npm ci --prefer-offline
//go:generate npm run build

var (
	//go:embed dist
	FS embed.FS

	MainFS   = echo.MustSubFS(FS, "dist")
	AssetsFS = echo.MustSubFS(FS, "dist/assets")
)
