package env

import (
	"os"
	"strings"

	"github.com/rashintha/logger"
)

var CONF = map[string]string{}

func init() {
	logger.Defaultln("Loading environment variables")
	data, err := os.ReadFile(".env")
	if err != nil {
		logger.Warningln("⚠️ No .env file found — loading system environment variables")
		loadSystemEnv()
		return
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, "#") {
			line = strings.Split(line, "#")[0]
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`) // remove quotes if any
		if key != "" {
			CONF[key] = value
			os.Setenv(key, value)
		}
	}
	logger.Defaultf("✅ Loaded %d environment variables", len(CONF))
}

func loadSystemEnv() {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			CONF[pair[0]] = pair[1]
		}
	}
	logger.Defaultf("✅ Loaded %d environment variables from system", len(CONF))
}
