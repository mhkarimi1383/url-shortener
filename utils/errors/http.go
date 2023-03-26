/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package errors

import "github.com/labstack/echo/v4"

func ErrorToHTTPErrorAdapter(code int, message string, err error) *echo.HTTPError {
	e := echo.NewHTTPError(code, message)
	e.SetInternal(err)
	return e
}
