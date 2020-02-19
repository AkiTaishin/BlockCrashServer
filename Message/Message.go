package Message

import (
	"bytes"
	"log"
	"net/http"
)

// DataBody グローバル格納用
var DataBody string = ""

// SendMessage メッセージを送りたい
func SendMessage(w http.ResponseWriter, r *http.Request) {

	// 送りたいメッセージ
	data := "{\"test\",SuccessSendMessage}"
	// 書き込むためにバイトに変換
	responses := []byte(data)

	// ヘッダーを付与してメッセージを書き込む
	// ヘッダの情報をもとにブラウザとかがどんな動きをするか考えてくれる
	w.Header().Set("Content-Type", "application/json")
	w.Write(responses)

	log.Print("[送信完了 : ", w, "]")
}

// SaveMessage ユーザーからのリクエストデータを保存する
func SaveMessage(w http.ResponseWriter, r *http.Request) {

	// rの中身を格納する
	saveNewData := new(bytes.Buffer)
	saveNewData.ReadFrom(r.Body)

	// 受け取った内容をストリングに変換して格納
	data := saveNewData.String()
	// グローバルに格納
	DataBody = data

	log.Print("DataBody_:", DataBody)
}
