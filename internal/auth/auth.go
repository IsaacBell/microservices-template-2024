package auth

import (
	"context"
	"core/internal/biz"
	"core/internal/constants"
	zap "core/internal/logs"
	"encoding/json"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware"
	jwtMiddleware "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/golang-jwt/jwt"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type options struct {
	signingMethod jwt.SigningMethod
	claims        func() jwt.Claims
	tokenHeader   map[string]interface{}
}

type UserCtxKey struct{}

var (
	currCtx *AuthCtx = NewAuthCtx()
	logger  zap.Logger
)

type AuthCtx struct {
	Ctx context.Context
}

func NewAuthCtx() *AuthCtx {
	return &AuthCtx{Ctx: context.WithValue(context.Background(), "authorized", "NewAuthCtx")}
}

func NewAuthCtxFrom(ctx context.Context) *AuthCtx {
	ctx = context.WithValue(ctx, "authorized", "NewAuthCtxFrom")
	currCtx = &AuthCtx{Ctx: ctx}
	return currCtx
}

type JwtHandler struct {
	ctx context.Context
}

func JwtMiddleware(authCtx *AuthCtx) middleware.Middleware {
	callback := WithContext(authCtx.Ctx).DecodeJwt
	return jwtMiddleware.Server(callback)
}

func WithContext(ctx_ context.Context) *JwtHandler {
	j := &JwtHandler{}
	j.ctx = ctx_
	return j
}

func Encode(ctx context.Context, u *biz.User) (context.Context, string, error) {
	out, err := WithContext(ctx).EncodeJwt(u)
	return ctx, out, err
}
func Decode(ctx context.Context, token *jwtv5.Token) (*context.Context, interface{}, error) {
	out, err := WithContext(ctx).DecodeJwt(token)
	return &ctx, out, err
}

func (j *JwtHandler) WithContext(ctx_ context.Context) *JwtHandler {
	j.ctx = ctx_
	return j
}

func (j *JwtHandler) EncodeJwt(u *biz.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":           u.ID,
		"user_id":       u.ID,
		"user_email":    u.Email,
		"username":      u.Username,
		"first_name":    u.FirstName,
		"last_name":     u.LastName,
		"session_token": u.SessionToken,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key := []byte(constants.JwtKey)
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	// fmt.Printf("...generated user token: %v\n", tokenStr)
	return tokenStr, nil
}

func (j *JwtHandler) DecodeJwt(token *jwtv5.Token) (interface{}, error) {
	claims := jwtv5.MapClaims{}

	// verify signed jwt key
	_, err := jwtv5.ParseWithClaims(token.Raw, &claims, func(token *jwtv5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(constants.JwtKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %v", err)
	}

	user := make(map[string]string)
	fmt.Printf("\n::::jwt claims::::\n")
	for key, val := range claims {
		fmt.Println("-> ", key+": ", val)
		if val == nil {
			user[key] = ""
		} else {
			user[key] = val.(string)
		}
	}
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	// cache.Cache(context.Background()).Set("auth:"+user["sub"], data, time.Hour*2)

	ctx := context.WithValue(currCtx.Ctx, "session", data)
	currCtx.Ctx = ctx

	return []byte(constants.JwtKey), nil
	// return user, nil
}
