/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package flagloader

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

// SetFlagsFromEnv parses all registered flags in the given flag set,
// and if they are not already set it attempts to set their values from
// environment variables. Environment variables take the name of the flag but
// are UPPERCASE, and any dashes are replaced by underscores. Environment
// variables additionally are prefixed by the given string followed by
// and underscore. For example, if prefix=PREFIX: some-flag => PREFIX_SOME_FLAG
func SetFlagsFromEnv(fs *flag.FlagSet, prefix string) (err error) {
	alreadySet := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		alreadySet[f.Name] = true
	})
	fs.VisitAll(func(f *flag.Flag) {
		if !alreadySet[f.Name] {
			key := prefix + "_" + strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
			val := os.Getenv(key)
			if val != "" {
				if sErr := fs.Set(f.Name, val); sErr != nil {
					err = fmt.Errorf("invalid value %q for %s: %v", val, key, sErr)
				}
			}
		}
	})
	return err
}
