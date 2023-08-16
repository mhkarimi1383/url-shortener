/*
Copyright © 2023 Muhammed Hussein Karimi info@karimi.dev

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"github.com/mhkarimi1383/url-shortener/internal/flagutil"
	"github.com/mhkarimi1383/url-shortener/internal/log"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "url-shortener",
	Short: "Simple and minimalism URL Shortener",
	Long:  ``,
	Run:   start,
}

func Execute() {
	if invalid := flagutil.SetFlagsFromEnv(rootCmd.PersistentFlags(), "USH"); invalid.String != "" {
		log.Logger.Panic("Invalid environemt values provided", invalid)
	}

	err := rootCmd.Execute()
	if err != nil {
		log.Logger.Panic(err.Error())
	}
}

func init() {
	var cfg configuration.Config
	rootCmd.PersistentFlags().StringVarP(&cfg.ListenAddress, "listen-address", "l", "127.0.0.1:8080", "Host:Port to listen")
}

func start(_ *cobra.Command, _ []string) {
	println("Hello World.!")
}
