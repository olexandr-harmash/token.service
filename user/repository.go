package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pgAdapter "github.com/vgarvardt/go-pg-adapter"
)

type UserStore struct {
	adapter   pgAdapter.Adapter
	tableName string
}

func NewStore(adapter pgAdapter.Adapter) (*UserStore, error) {
	store := &UserStore{
		adapter:   adapter,
		tableName: "oauth2_users",
	}

	err := store.initTable()

	if err != nil {
		return store, fmt.Errorf("cannot create user's store: %s", err.Error())
	}

	return store, nil
}

func (s *UserStore) initTable() error {
	return s.adapter.Exec(context.Background(), fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s (
		login      TEXT   	 NOT NULL,
		created_at TIMESTAMPTZ NOT NULL,
		data       JSONB       NOT NULL,
		CONSTRAINT %[1]s_pkey PRIMARY KEY (login)
		);
	`, s.tableName))
}

func (s *UserStore) Create(ctx context.Context, info *User) error {
	buf, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("cannot insert user: %s", err.Error())
	}

	item := &UserStoreItem{
		Data:      buf,
		CreatedAt: time.Now(),
	}

	return s.adapter.Exec(
		ctx,
		fmt.Sprintf("INSERT INTO %s (created_at, login, data) VALUES ($1, $2, $3)", s.tableName),
		item.CreatedAt,
		item.Login,
		item.Data,
	)
}

func (s *UserStore) GetByLogin(ctx context.Context, login string) (*User, error) {
	if login == "" {
		return nil, nil
	}

	var item UserStoreItem
	if err := s.adapter.SelectOne(ctx, &item, fmt.Sprintf("SELECT * FROM %s WHERE login = $1", s.tableName), login); err != nil {
		return nil, fmt.Errorf("cannot get user by login: %s", err.Error())
	}

	return item.toUserData()
}

type UserStoreItem struct {
	CreatedAt time.Time `db:"created_at"`
	Login     string    `db:"login"`
	Data      []byte    `db:"data"`
}

func (s *UserStoreItem) toUserData() (*User, error) {
	var user *User
	err := json.Unmarshal(s.Data, &user)
	if err != nil {
		return nil, fmt.Errorf("cannot convert data to user: %s", err.Error())
	}
	return user, nil
}
