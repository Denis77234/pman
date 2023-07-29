package reqstruct

import (
	"encoding/json"
)

type UpdatePack struct {
	Name    string `json:"name"`
	Version string `json:"ver"`
}

type Update struct {
	Updates []UpdatePack `json:"packages"`
}

func (u *UpdatePack) UnmarshalJSON(data []byte) error {
	var Obj struct {
		Name    string `json:"name"`
		Version string `json:"ver"`
	}

	err := json.Unmarshal(data, &Obj)
	if err == nil {
		u.Name = Obj.Name
		u.Version = Obj.Version
		return nil
	}

	var Name string

	err1 := json.Unmarshal(data, &Name)
	if err1 != nil {
		return err1
	}
	u.Name = Name
	u.Version = ""
	return nil
}
