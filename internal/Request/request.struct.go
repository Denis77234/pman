package reqstruct

import "fmt"

type Request struct {
	Name    string   `json:"name"`
	Ver     string   `json:"ver"`
	Targets []Target `json:"targets"`
	Packets []Packet `json:"packets"`
}

func (r Request) ArchiveName(extension string) string {
	name := fmt.Sprintf("%v.%v", r.Ver, extension)
	return name
}
