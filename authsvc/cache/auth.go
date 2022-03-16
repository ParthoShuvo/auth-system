package cache

import (
	"github.com/parthoshuvo/authsvc/token"
)

func (td *TokenDB) SetRefreshToken(authToken *token.AuthToken) error {
	return td.rdb.Set(td.ctx, authToken.UserID(), authToken.UUID(), authToken.Expires()).Err()
}
