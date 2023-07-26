package reqstruct

import "fmt"

type Packet struct {
	Name string `json:"name"`
	Ver  string `json:"ver"`
}

type Request struct {
	Name    string   `json:"name"`
	Ver     string   `json:"ver"`
	Targets []Target `json:"targets"`
	Packets []Packet `json:"packets"`
}

func (r Request) ArchiveName(extension string) string {
	name := fmt.Sprintf("%v.%v", r.Name, extension)
	return name
}
