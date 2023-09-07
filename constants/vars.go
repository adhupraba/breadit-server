package constants

import (
	"time"

	"github.com/adhupraba/breadit-server/lib"
)

const (
	AccessTokenTTL                  = time.Second * 30
	RefreshTokenTTL                 = time.Second * 180
	CacheAfterUpvotes               = 1
	InfiniteScrollPaginationResults = 2
)

var UseSecureCookies = lib.EnvConfig.Env == "staging" || lib.EnvConfig.Env == "production"
