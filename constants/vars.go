package constants

import (
	"time"

	"github.com/adhupraba/breadit-server/lib"
)

const (
	AccessTokenTTL                  = time.Hour * 24 * 6 // 3 days
	RefreshTokenTTL                 = time.Hour * 24 * 7 // 7 days -> equivalent to SESSION_MAX_AGE in ui
	CacheAfterUpvotes               = 10
	InfiniteScrollPaginationResults = 4
)

var UseSecureCookies = lib.EnvConfig.Env == "staging" || lib.EnvConfig.Env == "production"
