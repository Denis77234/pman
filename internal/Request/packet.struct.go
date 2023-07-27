package reqstruct

import (
	"regexp"
	"strings"
)

type Packet struct {
	Name string `json:"name"`
	Ver  string `json:"ver"`
}

func (p Packet) VerOperator() string {
	operator := strings.Replace(p.Ver, p.VerNum(), "", -1)
	return operator
}

func (p Packet) VerNum() string {
	oper, _ := regexp.Compile(`[<>=]*`)
	num := oper.ReplaceAllString(p.Ver, "")

	return num
}
