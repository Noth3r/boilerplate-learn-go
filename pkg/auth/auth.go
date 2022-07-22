package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
)

type AuthInterface interface {
	CreateAuth(string, *TokenDetails) error
	FetchAuth(string) (string, error)
	DeleteTokens(*AccessDetails) error
	DeleteRefresh(string) error
	CheckRevoked(string) (bool, error)
	RevokeAll(string) error
	RevokeRefresh(string, string) error
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

	rtData, err := json.Marshal(map[string]interface{}{"userId": userId, "revoked": false})

	fmt.Println("token: " + td.TokenUuid + " " + "refresh: " + td.RefreshUuid)

	atCreated, err := tk.client.Set(ctx, td.TokenUuid, userId, at.Sub(now)).Result()
	if err != nil {
		return err
	}

	rtCreated, err := tk.client.Set(ctx, td.RefreshUuid, rtData, rt.Sub(now)).Result()
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
	at := strings.Split(refreshUuid, "++")[0]
	_, err := tk.client.Del(ctx, at).Result()
	deleted, err := tk.client.Del(ctx, refreshUuid).Result()
	if err != nil {
		return err
	}

	if deleted == 0 {
		return errors.New("Something went wrong")
	}

	return nil
}

type data struct {
	Revoked bool   `json:"revoked"`
	UserId  string `json:"userId"`
}

func (tk *service) CheckRevoked(refreshUuid string) (bool, error) {
	res, err := tk.client.Get(ctx, refreshUuid).Result()

	if err != nil {
		return false, err
	}

	data := data{}

	errors := json.Unmarshal([]byte(res), &data)
	if errors != nil {
		return false, errors
	}

	if data.Revoked {
		return true, nil
	}

	return false, nil
}

func (tk *service) RevokeAll(userId string) error {
	data, err := tk.client.Keys(context.Background(), "*"+userId).Result()
	if err != nil {
		fmt.Println(err)
		return err
	}

	if len(data) > 100 {
		last := data[len(data)-1]
		tk.DeleteRefresh(last)
	}

	for _, v := range data {
		tk.RevokeRefresh(v, userId)
	}

	return nil
}

func (tk *service) RevokeRefresh(refreshUuid string, userId string) error {
	rtData, err := json.Marshal(map[string]interface{}{"userId": userId, "revoked": true})
	if err != nil {
		return err
	}

	errSet := tk.client.Set(ctx, refreshUuid, rtData, tk.client.TTL(ctx, refreshUuid).Val())
	if errSet != nil {
		return err
	}
	return nil
}
