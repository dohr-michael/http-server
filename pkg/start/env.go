package start

import (
	"github.com/dohr-michael/http-server/config"
	"os"
	"strings"
	"time"
)

func NewEnvContext(overridden []string) *envContext {
	result := &envContext{
		publicEnv: make(map[string]string),
		NoCacheFiles: []string{"/index.html", "/@/health", "/@/info"},
		AppEnv: "unknown",
		AppName: "http-server",
		AppVersion: config.Config.BuildVersion(),
		AppCommit: config.Config.BuildRevision(),
		AppBuildAt: config.Config.BuildTime(),
		AppStartedAt: time.Now().Format(time.RFC3339),
	}
	envItems := append(os.Environ(), overridden...)

	for _, i := range envItems {
		sep := strings.Index(i, "=")
		key := i[0:sep]
		value := i[sep+1:]
		if strings.HasPrefix(key, "HTTP_SERVER_CONFIG_") {
			result.publicEnv[strings.TrimPrefix(key, "HTTP_SERVER_CONFIG_")] = value
		} else {
			if len(strings.TrimSpace(value)) == 0 {
				continue
			}
			switch key {
			case "NO_CACHE_FILES":
				items := strings.Split(value, ",")
				for _, item := range items {
					if s := strings.TrimSpace(item); len(s) > 0 {
						result.NoCacheFiles = append(result.NoCacheFiles, s)
					}
				}
			case "APP_ENV", "ENV", "PLAY_ID":
				result.AppEnv = value
			case "APP_NAME":
				result.AppName = value
			case "APP_COMMIT":
				result.AppCommit = value
			case "APP_BUILD_AT":
				result.AppBuildAt = value
			case "APP_VERSION":
				result.AppVersion = value
			}
		}
	}

	return result
}

type envContext struct {
	publicEnv    map[string]string
	NoCacheFiles []string
	AppEnv       string
	AppName      string
	AppBuildAt   string
	AppStartedAt string
	AppVersion   string
	AppCommit    string
}

func (c *envContext) PublicEnv() map[string]string {
	return c.publicEnv
}
