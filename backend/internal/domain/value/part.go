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

var validateParts = map[Part]bool{
	PartVo:  true,
	PartGt:  true,
	PartBa:  true,
	PartDr:  true,
	PartKey: true,
}

func NewPart(s string) (Part, error) {
	p := Part(s)
	if validateParts[p] {
		return p, nil
	}
	return "", fmt.Errorf("無効なパートです: %s", s)
}