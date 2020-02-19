// Common 汎用的なパッケージ変数
package Common

// UserInfo ユーザー情報
type UserInfo struct {
	Name      string
	ID        string
	HighScore int
}

// UserInfoArray ↑の配列
type UserInfoArray []*UserInfo

var GetUserInfoArray UserInfoArray
