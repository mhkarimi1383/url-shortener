package ui

import (
	"embed"
	"os"

	"github.com/labstack/echo/v4"
)

func init() {
	os.Setenv("BASE_URL", "/ui")
}

//go:generate pnpm i --frozen-lockfile
//go:generate pnpm build

var (
	//go:embed dist
	FS embed.FS

	MainFS   = echo.MustSubFS(FS, "dist")
	AssetsFS = echo.MustSubFS(FS, "dist/assets")
)
