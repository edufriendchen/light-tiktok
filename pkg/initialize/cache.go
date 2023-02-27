package initialize

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/edufriendchen/light-tiktok/pkg/cache"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"go.opencensus.io/resource"
)

func InitCacheManager() {
	cacheManager, err := cache.NewClient(resource.Resource{})
	if err != nil {
		klog.Fatalf("Init cacheManager error:", err)
	}
	global.CacheManager = cacheManager
}
