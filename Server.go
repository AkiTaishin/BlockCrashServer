package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"./Common"
	"./CreateId"
	"./Message"
)

func main() {

	log.Print("Running...")

	// 自分が設定したURIがたたかれたら
	http.HandleFunc("/saveUserData", SaveUserData)
	http.HandleFunc("/sendUserData", SendUserData)
	http.HandleFunc("/sendAllUserData", SendAllUserData)

	http.HandleFunc("/pleaseResponse", Message.SendMessage)
	http.HandleFunc("/pleaseSaveData", Message.SaveMessage)

	// 応答待ち
	http.ListenAndServe(":8080", nil)
}

// SaveUserData ユーザーのデータをセーブする
// 起動時と終了時のデータ保存に使用
// 新規作成か上書きかのチェックで分岐
func SaveUserData(w http.ResponseWriter, r *http.Request) {

	// リクエスト読み取り用の一時変数作成、保管
	data := new(bytes.Buffer)
	data.ReadFrom(r.Body)

	// データの文字列をバイトの配列に変換
	datastring := data.String()
	log.Print("dataString_:", datastring)
	dataArray := []byte(datastring)

	// バイトの配列をUserInfoに変換
	var user Common.UserInfo
	err := json.Unmarshal(dataArray, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 新規作成か上書き保存か確認
	// user.IDが空白の時新規作成 | user.IDが空白の時以外は上書き
	response, err := json.Marshal(CreateorOverWriteUserData(user))

	// Json変換に失敗がなければ送信データとして書き込み
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	}
}

// SendUserData リクエストと一致したIDのユーザーデータを送信
// 2回め以降に呼ばれる
func SendUserData(w http.ResponseWriter, r *http.Request) {

	// リクエストを受け取るためのbytesを作成
	request := new(bytes.Buffer)
	// プレイヤーIDをリクエストから取得
	request.ReadFrom(r.Body)

	// GetUserDataに渡すための変数を作成
	// 中身はユーザーデータ
	requestStr := request.String()
	// 作成したユーザーデータ
	log.Print("[requestStr : ", requestStr, "]")

	// 保存しているデータの中にIDが一致するものがあるか探査
	// return:UserInfo
	data, err := GetUserData(requestStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Print("(^o^)ﾉ ＜ そんなデータないぜー")
		return
	}

	// 取得したUserInfoをJsonに変換
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 取得したUserInfoをリクエストユーザーに送信
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	// 送りつけた内容ログ
	log.Print("[Send Data : ", string(response), "]")
}

// CreateorOverWriteUserData IDが登録されていない場合新規ユーザー作成
// それ以外は既存上書き
// タイトルと終了時
func CreateorOverWriteUserData(data Common.UserInfo) *Common.UserInfo {

	work := new(Common.UserInfo)

	if data.ID == "" {

		// 新規作成
		work = CreateNewUser(data.Name)
	} else {
		// 上書き保存
		work = OverWriteUser(data)
	}

	return work
}

// CreateNewUser 登録されていないIDが呼び出された場合に新規ユーザー作成
// 最初にのみ呼び出される
func CreateNewUser(name string) *Common.UserInfo {

	// 新規ID作成
	newID := CreateId.CreateID(5)
	//newID = "aaaaa"

	// 被っていないかチェック
	for i := 0; i < len(Common.GetUserInfoArray); i++ {

		// 登録データを全探査し、ほかのユーザーとIDが被っていたらID変更
		if Common.GetUserInfoArray[i].ID == newID {

			log.Print("被ってる")
			newID = CreateId.CreateID(5)
			log.Print("ChangeID_:", newID)
			continue

		} else {

			break
		}

	}

	// 新規データ作成、配列の最後尾に追加
	work := new(Common.UserInfo)
	work.ID = newID
	work.Name = name
	work.HighScore = 0
	Common.GetUserInfoArray = append(Common.GetUserInfoArray, work)

	return work
}

// OverWriteUser 既存IDが呼び出されたら上書き
// ゲームの終了(スコア保存の時)に呼ばれる
func OverWriteUser(data Common.UserInfo) *Common.UserInfo {

	work := new(Common.UserInfo)
	work.ID = "error"
	work.Name = ""
	work.HighScore = 0

	for i := 0; i < len(Common.GetUserInfoArray); i++ {

		if Common.GetUserInfoArray[i].ID != data.ID {

			continue

		} else {

			Common.GetUserInfoArray[i].HighScore = data.HighScore
			Common.GetUserInfoArray[i].Name = data.Name
			work = Common.GetUserInfoArray[i]
		}
	}

	return work
}

// GetUserData 指定されたIDと同じIDのデータをそのまま渡す
// タイトルで呼ばれる
func GetUserData(SameID string) (*Common.UserInfo, error) {

	user := new(Common.UserInfo)

	for i := 0; i < len(Common.GetUserInfoArray); i++ {

		if Common.GetUserInfoArray[i].ID != SameID {

			continue

		} else {

			user = Common.GetUserInfoArray[i]
			return user, nil
		}

	}

	// エラー
	// ゲーム側でエラー番号を登録してその番号によってエラーメッセージを表示したらやさしいね
	err := errors.New("そのユーザーIDは登録されていません")
	log.Print("[エラー処理_:", err, "]")
	return nil, err

}

// SendAllUserData ランキング用
// 全てのデータを送り付ける
func SendAllUserData(w http.ResponseWriter, r *http.Request) {

	data := Common.GetUserInfoArray
	//配列をJSONに変換
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//全ユーザー情報を送信
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	log.Print("[send all user data : ", string(response), "]")
}
