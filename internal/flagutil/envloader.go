package flagutil

import (
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
)

// SetFlagsFromEnv parses all registered flags in the given flag set,
// and if they are not already set it attempts to set their values from
// environment variables. Environment variables take the name of the flag but
// are UPPERCASE, and any dashes are replaced by underscores. Environment
// variables additionally are prefixed by the given string followed by
// and underscore. For example, if prefix=PREFIX: some-flag => PREFIX_SOME_FLAG
// returns zap.Field for logging
func SetFlagsFromEnv(fs *flag.FlagSet, prefix string) zap.Field {
	alreadySet := make(map[string]bool)
	var invalid zap.Field
	fs.Visit(func(f *flag.Flag) {
		alreadySet[f.Name] = true
	})
	fs.VisitAll(func(f *flag.Flag) {
		key := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
		if prefix != "" {
			key = prefix + "_" + key
		}
		f.Usage += " [env: " + key + "]"
		if !alreadySet[f.Name] {
			val := os.Getenv(key)
			if val != "" {
				if sErr := fs.Set(f.Name, val); sErr != nil {
					invalid = zap.String(key, val)
				}
			}
		}
	})
	return invalid
}
