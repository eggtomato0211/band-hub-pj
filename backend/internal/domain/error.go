package domain

import "errors"

var (
	// ユーザやイベントなどが見つからない場合のエラー
	ErrNotFound = errors.New("見つかりませんでした")

	// メールアドレスが重複している場合のエラー
	ErrDuplicateEmail = errors.New("メールアドレスが既に存在しています")

	// ログイン失敗時のエラー
	ErrInvalidCredentials = errors.New("メールアドレスまたはパスワードが正しくありません")

	// 認証されていない場合のエラー
	ErrUnauthorized = errors.New("権限がありません")

	// 権限がない場合のエラー
	ErrForbidden = errors.New("アクセスが禁止されています")
)