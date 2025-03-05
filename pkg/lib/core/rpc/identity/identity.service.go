package identity

import (
	"context"

	gcontext "github.com/LoveCatdd/webctx/pkg/lib/core/context"
	"github.com/LoveCatdd/webctx/pkg/lib/core/goroutine"
	"github.com/golang-jwt/jwt/v4"
)

const (
	// 获取 custonContextKey
	custonContextKey = gcontext.CustonContextKey
)

type IdentityService interface {

	// get uid
	UserId(context.Context) string

	// get username
	UserName(context.Context) string

	// get user identity name
	UserIdentityName(context.Context) string
}

type Impl struct{}

func parse(c context.Context, key string) (any, bool) {
	customContext := c.Value(custonContextKey).(*gcontext.CustomContext)
	if info, ok := customContext.ContextHolder().ContenxtMap().Get(key); ok {
		return info, true
	}
	return nil, false
}

func (Impl) claimMap(c context.Context) (jwt.MapClaims, bool) {
	if info, ok := parse(c, goroutine.JWT_MAP_CLAIM); ok {
		return info.(jwt.MapClaims), true
	}
	return nil, false
}

func (i Impl) UserId(c context.Context) any {
	if info, ok := i.claimMap(c); ok && info != nil {
		return info[IDENTITY_USERTID_KEY]
	}
	return ""
}

func (i Impl) UserName(c context.Context) string {
	if info, ok := i.claimMap(c); ok && info != nil {
		return info[IDENTITY_USERTNAEM_KEY].(string)
	}
	return ""
}

func (i Impl) UserIdentityName(c context.Context) string {
	if info, ok := i.claimMap(c); ok && info != nil {
		return info[IDENTITY_USERIDENTITYNAME_KEY].(string)
	}
	return ""
}
