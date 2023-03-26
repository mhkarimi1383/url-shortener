/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package config

import (
	"github.com/mhkarimi1383/url-shortener/types"
)

var config *types.Config

func SetConfig(cfg *types.Config) {
	config = cfg
}

func GetConfig() types.Config {
	return *config
}
