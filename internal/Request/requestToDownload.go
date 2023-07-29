package reqstruct

type Update struct {
	Updates []Packet `json:"packages"`
}
