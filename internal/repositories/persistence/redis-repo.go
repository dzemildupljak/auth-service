package persistence

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepo(ctx context.Context, client *redis.Client) *RedisRepo {
	return &RedisRepo{
		ctx:    ctx,
		client: client,
	}
}

func (redisrepo *RedisRepo) SetMiddlewareUser(usr domain.UserMiddleware) error {
	redisUsr := map[string]string{
		"id":         usr.Id,
		"username":   usr.Username,
		"email":      usr.Email,
		"isverified": string(strconv.FormatBool(usr.Isverified)),
	}

	for k, v := range redisUsr {
		err := redisrepo.client.HSet(redisrepo.ctx, redisUsr["id"], k, v).Err()
		if err != nil {
			fmt.Println("Redis repo Set user failed")
			return err
		}
	}

	err := redisrepo.client.Expire(redisrepo.ctx, redisUsr["id"], 1*time.Minute).Err()
	if err != nil {
		fmt.Println("Redis repo Set user failed")
		return err
	}
	return nil
}

func (redisrepo *RedisRepo) GetMiddlewareUser(usrId string) domain.UserMiddleware {
	userSes := redisrepo.client.HGetAll(redisrepo.ctx, usrId).Val()

	var midUsr domain.UserMiddleware

	midUsr.Id = userSes["id"]
	midUsr.Username = userSes["username"]
	midUsr.Email = userSes["email"]
	isverified, err := strconv.ParseBool(userSes["isverified"])

	if err != nil {
		midUsr.Isverified = false
	} else {
		midUsr.Isverified = isverified
	}

	return midUsr
}

func (redisrepo *RedisRepo) ClearItemByKey(itemKey string) error {
	err := redisrepo.client.Del(redisrepo.ctx, itemKey).Err()
	if err != nil {
		fmt.Println("Redis repo Clear user failed")
		return err
	}

	return err
}
