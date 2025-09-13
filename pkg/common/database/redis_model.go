package database

import (
	"context"
	"time"
)

const (
	getuiToken = "GETUI_TOKEN"
)

func (d *DataBases) SetGetuiToken(token string, expireTime int64) error {
	return d.RDB.Set(context.Background(), getuiToken, token, time.Duration(expireTime)*time.Second).Err()
}
