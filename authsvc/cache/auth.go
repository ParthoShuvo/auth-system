package cache

import (
	"github.com/parthoshuvo/authsvc/token"
)

func (td *TokenDB) SetRefreshToken(authToken *token.AuthToken) error {
	return td.rdb.Set(td.ctx, authToken.UserID(), authToken.UUID(), authToken.Expires()).Err()
}

func (td *TokenDB) GetRefreshToken(key string) (string, error) {
	return td.rdb.Get(td.ctx, key).Result()
}

func (td *TokenDB) RevokeRefreshToken(key string) error {
	return td.rdb.Del(td.ctx, key).Err()
}
