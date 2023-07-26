package reqstruct

import "encoding/json"

type Target struct {
	Path    string `json:"path"`
	Exclude string `json:"exclude,omitempty"`
}

func (t *Target) UnmarshalJSON(data []byte) error {

	var Obj struct {
		Path    string `json:"path"`
		Exclude string `json:"exclude,omitempty"`
	}

	err := json.Unmarshal(data, &Obj)
	if err == nil {
		t.Path = Obj.Path
		t.Exclude = Obj.Exclude
		return nil
	}

	var Path string

	err1 := json.Unmarshal(data, &Path)
	if err1 != nil {
		return err1
	}
	t.Path = Path
	t.Exclude = ""
	return nil
}
