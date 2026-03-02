package value

import "fmt"

type Part string

// 有効なパートの定数定義
const (
	PartVo  Part = "Vo"
	PartGt  Part = "Gt"
	PartBa  Part = "Ba"
	PartDr  Part = "Dr"
	PartKey Part = "Key"
)

// validateParts はバリデーション用のmap
var validateParts = map[Part]bool{
	PartVo:  true,
	PartGt:  true,
	PartBa:  true,
	PartDr:  true,
	PartKey: true,
}

// バリデーション付きのコンストラクタ関数
// mapを使って有効なパートかどうかを判定する
func NewPart(s string) (Part, error) {
	p := Part(s)
	if validateParts[p] {
		return p, nil
	}
	return "", fmt.Errorf("無効なパートです: %s", s)
}