package initialize

import (
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/edufriendchen/light-tiktok/pkg/jwt"
)

// InitJWT to init JWT
func InitJWT() {
	global.Jwt = jwt.NewJWT([]byte(consts.JWTSecretKey))
}
