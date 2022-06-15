package dal

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/HeRaNO/cdoj-vjudge/config"
	"github.com/HeRaNO/cdoj-vjudge/model"
	"github.com/go-redis/redis/v8"
)

func MakeUserPasswordKey(username *string) string {
	return fmt.Sprintf("%s:%s", model.UserNamePasswordKey, *username)
}

func MakeResultKey(submissionID *string) string {
	return fmt.Sprintf("%s:%d", model.SubmissionResultKey, submissionID)
}

func IsAuthValid(ctx context.Context, username *string, password *string) (bool, error) {
	key := MakeUserPasswordKey(username)
	ret, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("[INFO] IsAuthValid(): key is nil, username: %s\n", *username)
		return false, nil
	} else if err != nil {
		log.Printf("[ERROR] IsAuthValid(): redis query error, err: %s", err.Error())
		return false, err
	}
	if ret != *password {
		return false, nil
	}
	return true, nil
}

func SetSubmissionResult(ctx context.Context, submissionID *string, result *string) error {
	key := MakeResultKey(submissionID)
	err := config.RedisClient.Set(ctx, key, *result, 0).Err()
	if err != nil {
		log.Printf("[ERROR] SetSubmissionResult(): redis set error, err: %s", err.Error())
	}
	return err
}

func GetSubmissionResult(ctx context.Context, submissionID *string) (string, error) {
	key := MakeResultKey(submissionID)
	ret, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("[INFO] GetSubmissionResult(): key is nil, submissionID: %d\n", submissionID)
		return "", errors.New("no such submission")
	} else if err != nil {
		log.Printf("[ERROR] GetSubmissionResult(): redis query error, err: %s", err.Error())
		return "", err
	}

	return ret, nil
}
