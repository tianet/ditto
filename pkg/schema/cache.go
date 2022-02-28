package schema

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var C *cache.Cache

func NewCache() {
	C = cache.New(5*time.Minute, 10*time.Minute)
}
