package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	v1 "microservices-template-2024/api/v1"
	"microservices-template-2024/internal/biz"
	cache "microservices-template-2024/pkg/cache"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/* public-facing user functions
 */

type grpcCloseConn func() error

func grpcConn() (v1.UsersClient, grpcCloseConn, error) {
	conn, err := grpc.Dial(os.Getenv("CORE_SERVICE_ADDRESS"), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	resp, err := client.GetUser(context.Background(), &v1.GetUserRequest{Id: &uid})
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

func Create(u *v1.User, raiseErrOnFail bool) (string, error) {
	client, closeConn, err := grpcConn()
	if err != nil {
		return "", err
	}
	defer closeConn()
	resp, err := client.CreateUser(context.Background(), &v1.CreateUserRequest{User: u})
	if err != nil {
		fmt.Printf("failed to create user: %v", err)
		if raiseErrOnFail {
			return "", err
		}
	}
	if raiseErrOnFail && resp.Id == "" {
		return "", errors.New("Failed to create user")
	}

	return resp.Id, nil
}

func Update(u *v1.User, raiseErrOnFail bool) (*v1.User, error) {
	client, closeConn, err := grpcConn()
	if err != nil {
		return nil, err
	}
	defer closeConn()
	resp, err := client.UpdateUser(context.Background(), &v1.UpdateUserRequest{User: u})
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
	resp, err := client.DeleteUser(context.Background(), &v1.DeleteUserRequest{Id: uid})
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
	resp, err := client.ListUser(context.Background(), &v1.ListUserRequest{})
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
	resp, err := client.SignUp(context.Background(), req)
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
	resp, err := client.SignIn(context.Background(), &v1.SignInRequest{Email: email, Password: pass})
	if !resp.Ok || err != nil {
		fmt.Printf("failed to sign in: %v", err)
		return nil, err
	}

	return resp.User, nil
}
