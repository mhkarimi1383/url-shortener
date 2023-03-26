/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package url

import "net/url"

func IsValidUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
