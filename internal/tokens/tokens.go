package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Scope     string    `json:"scope"`
	Expiry    time.Time `json:"expiry"`
}

const (
	ScopeAuthentication = "authentication"
)

func GenerateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Scope:  scope,
		Expiry: time.Now().Add(ttl),
	}
	// 这行代码创建了一个长度为 32 的字节切片，其中每个字节都被初始化为 0。
	emptyBytes := make([]byte, 32)
	//  用随机字节填充
	_, err := rand.Read(emptyBytes)
	if err != nil {
		return nil, err
	}
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}
