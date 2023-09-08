package constants

import (
	"net/http"
	"time"

	"github.com/adhupraba/breadit-server/lib"
)

const (
	AccessTokenTTL                  = time.Hour * 24 * 6 // 3 days
	RefreshTokenTTL                 = time.Hour * 24 * 7 // 7 days -> equivalent to SESSION_MAX_AGE in ui
	CacheAfterUpvotes               = 10
	InfiniteScrollPaginationResults = 4
)

var UseSecureCookies = func() bool {
	return lib.EnvConfig.Env == "staging" || lib.EnvConfig.Env == "production"
}

var UseSameSiteMethod = func() http.SameSite {
	if lib.EnvConfig.Env == "staging" || lib.EnvConfig.Env == "production" {
		return http.SameSiteNoneMode
	} else {
		return http.SameSiteDefaultMode
	}
}
