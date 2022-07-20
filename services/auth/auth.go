package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
)

type AuthInterface interface {
	CreateAuth(string, *TokenDetails) error
	FetchAuth(string) (string, error)
	DeleteTokens(*AccessDetails) error
	DeleteRefresh(string) error
}

var _ AuthInterface = &service{}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	TokenUuid string
	UserId    string
}

type service struct {
	client *redis.Client
}

func NewAuth(client *redis.Client) *service {
	return &service{
		client: client,
	}
}

var ctx = context.Background()

func (tk *service) CreateAuth(userId string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := tk.client.Set(ctx, td.TokenUuid, userId, at.Sub(now)).Result()
	if err != nil {
		return err
	}

	rtCreated, err := tk.client.Set(ctx, td.RefreshUuid, userId, rt.Sub(now)).Result()
	if err != nil {
		return err
	}

	if atCreated == "0" && rtCreated == "0" {
		return errors.New("No record inserted")
	}

	return nil
}

func (tk *service) FetchAuth(tokenUuid string) (string, error) {
	userId, err := tk.client.Get(ctx, tokenUuid).Result()
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (tk *service) DeleteTokens(authD *AccessDetails) error {
	refreshUuid := fmt.Sprintf("%s++%s", authD.TokenUuid, authD.UserId)

	deletedAt, err := tk.client.Del(ctx, authD.TokenUuid).Result()
	if err != nil {
		return err
	}

	deletedRt, err := tk.client.Del(ctx, refreshUuid).Result()
	if err != nil {
		return err
	}

	if deletedAt != 1 && deletedRt != 1 {
		return errors.New("Something went wrong")
	}

	return nil
}

func (tk *service) DeleteRefresh(refreshUuid string) error {
	deleted, err := tk.client.Del(ctx, refreshUuid).Result()
	if err != nil {
		return err
	}

	if deleted == 0 {
		return errors.New("Something went wrong")
	}

	return nil
}
