package reqstruct

import (
	"encoding/json"
)

type Packet struct {
	Name string `json:"name"`
	Ver  string `json:"ver"`
}

func (p *Packet) UnmarshalJSON(data []byte) error {
	var Obj struct {
		Name    string `json:"name"`
		Version string `json:"ver"`
	}

	err := json.Unmarshal(data, &Obj)
	if err == nil {
		p.Name = Obj.Name
		p.Ver = Obj.Version
		return nil
	}

	var Name string

	err1 := json.Unmarshal(data, &Name)
	if err1 != nil {
		return err1
	}
	p.Name = Name
	p.Ver = ""
	return nil
}
