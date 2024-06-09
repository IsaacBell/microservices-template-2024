package users

import (
	"context"
	v1 "core/api/v1"
	"core/internal/auth"
	"core/internal/biz"
	"core/internal/util"
	cache "core/pkg/cache"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc/metadata"
)

/* users client
 */

type grpcCloseConn func() error

func grpcConn() (v1.UsersClient, grpcCloseConn, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(os.Getenv("CORE_SERVICE_ADDRESS")),
		grpc.WithMiddleware(
			circuitbreaker.Client(),
		),
	)
	if err != nil {
		return nil, nil, err
	}
	client := v1.NewUsersClient(conn)
	return client, func() error {
		return conn.Close()
	}, nil
}

func UserCacheKey(id string) string {
	return "user:" + id
}

func CacheExpiry() time.Duration {
	return time.Hour * 6
}

func UserFromCache(uid string) *biz.User {
	var u *biz.User
	cacheKey := UserCacheKey(uid)
	if cached, err := cache.Cache(context.Background()).Get(cacheKey).Result(); err == nil {
		if err := json.Unmarshal([]byte(cached), &u); err != nil {
			fmt.Println("cache miss: ", err)
			return nil
		}
	}
	return u
}

func Get(uid string, raiseErrIfNotFound bool) (*v1.User, error) {
	client, closeConn, err := grpcConn()
	if err != nil {
		return nil, err
	}
	defer closeConn()
	ctx := context.WithValue(context.Background(), util.ContextKeyMethod, "GET")

	resp, err := client.GetUser(ctx, &v1.GetUserRequest{Id: &uid})
	if err != nil {
		fmt.Printf("failed to get user: %v", err)
		if raiseErrIfNotFound {
			return nil, err
		}
	}
	if raiseErrIfNotFound && resp.User == nil {
		return nil, errors.New("User not found: " + uid)
	}

	return resp.User, nil
}

func AuthorizeUser(u *v1.User) (context.Context, string, error) {
	dataModel := biz.ProtoToUserData(u)
	ctx, token, err := auth.Encode(context.Background(), dataModel)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate JWT token: %v", err)
	}

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer "+token))
	ctx = context.WithValue(ctx, util.ContextKeyMethod, "POST")
	return ctx, token, nil
}

func Create(u *v1.User, raiseErrOnFail bool) (string, string, error) {
	client, closeConn, err := grpcConn()
	if err != nil {
		return "", "", err
	}
	defer closeConn()

	ctx, token, err := AuthorizeUser(u)

	resp, err := client.CreateUser(ctx, &v1.CreateUserRequest{User: u})
	if err != nil {
		fmt.Printf("failed to create user: %v", err)
		if raiseErrOnFail {
			return "", "", err
		}
	}
	if raiseErrOnFail && resp.Id == "" {
		return "", "", errors.New("Failed to create user")
	}

	return resp.Id, token, nil
}

func Update(u *v1.User, raiseErrOnFail bool) (*v1.User, error) {
	client, closeConn, err := grpcConn()
	if err != nil {
		return nil, err
	}
	defer closeConn()
	ctx := context.WithValue(context.Background(), util.ContextKeyMethod, "PUT")
	resp, err := client.UpdateUser(ctx, &v1.UpdateUserRequest{User: u})
	if err != nil {
		fmt.Printf("failed to update user: %v", err)
		if raiseErrOnFail {
			return nil, err
		}
	}
	if raiseErrOnFail && !resp.Ok {
		return nil, errors.New("Failed to update user")
	}

	return u, nil
}

func Delete(uid string, raiseErrOnFail bool) (bool, error) {
	client, closeConn, err := grpcConn()
	if err != nil {
		return false, err
	}
	defer closeConn()
	ctx := context.WithValue(context.Background(), util.ContextKeyMethod, "DELETE")
	resp, err := client.DeleteUser(ctx, &v1.DeleteUserRequest{Id: uid})
	if err != nil {
		fmt.Printf("failed to delete user: %v", err)
		if raiseErrOnFail {
			return false, err
		}
	}

	return resp.Ok, nil
}

func List() ([]*v1.User, error) {
	client, closeConn, err := grpcConn()
	if err != nil {
		return nil, err
	}
	defer closeConn()
	ctx := context.WithValue(context.Background(), util.ContextKeyMethod, "GET")
	resp, err := client.ListUser(ctx, &v1.ListUserRequest{})
	if err != nil {
		fmt.Printf("failed to list users: %v", err)
		return nil, err
	}

	return resp.Users, nil
}

func SignUp(req *v1.SignUpRequest) (string, error) {
	client, closeConn, err := grpcConn()
	if err != nil {
		return "", err
	}
	defer closeConn()
	ctx := context.WithValue(context.Background(), util.ContextKeyMethod, "GET")
	resp, err := client.SignUp(ctx, req)
	if err != nil {
		fmt.Printf("failed to sign up: %v", err)
		return "", err
	}

	return resp.Id, nil
}

func SignIn(email, pass string) (*v1.User, error) {
	client, closeConn, err := grpcConn()
	if err != nil {
		return nil, err
	}
	defer closeConn()
	ctx := context.WithValue(context.Background(), util.ContextKeyMethod, "GET")
	resp, err := client.SignIn(ctx, &v1.SignInRequest{Email: email, Password: pass})
	if !resp.Ok || err != nil {
		fmt.Printf("failed to sign in: %v", err)
		return nil, err
	}

	return resp.User, nil
}
