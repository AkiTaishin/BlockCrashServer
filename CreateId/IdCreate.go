package CreateId

import "math/rand"

// UserInfo ユーザー情報
type UserInfo struct {
	Name      string
	ID        string
	HighScore int
}

// UserInfoArray ↑の配列
type UserInfoArray []*UserInfo

var userInfoArray UserInfoArray

const (
	rs5Letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rs5LetterIdxBits = 6
	rs5LetterIdxMask = 1<<rs5LetterIdxBits - 1
	rs5LetterIdxMax  = 63 / rs5LetterIdxBits
)

// CreateID ランダムにIDを振り分け
func CreateID(n int) string {

	id := make([]byte, n)
	cache, remain := rand.Int63(), rs5LetterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), rs5LetterIdxMax
		}
		idx := int(cache & rs5LetterIdxMask)
		if idx < len(rs5Letters) {
			id[i] = rs5Letters[idx]
			i--
		}
		cache >>= rs5LetterIdxBits
		remain--
	}

	return string(id)
}
