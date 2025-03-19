package access_key

import "time"

type CreateAccessKeyRequest struct {
	Name      string `json:"name" binding:"required"`
	ExpiresAt int    `json:"expires_at" binding:"required,gte=3600,lte=31536000"`
}

func (r CreateAccessKeyRequest) ExpireTime() *time.Time {
	t := time.Now().Add(time.Duration(r.ExpiresAt) * time.Second)
	return &t
}
