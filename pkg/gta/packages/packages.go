package packages

import "fmt"

type Reason int

const (
	ReasonModification Reason = iota
	ReasonContainModifiedPackage
)

type Packages map[string]PresenceReason

type PresenceReason struct {
	Reason Reason

	ModifiedPackage *string
}

func (p Packages) String() string {
	str := ""
	for k, v := range p {
		if v.Reason == ReasonContainModifiedPackage {
			str += fmt.Sprintln(k, "import", `"`+*v.ModifiedPackage+`"`)
		} else {
			str += fmt.Sprintln(k)
		}
	}
	return str
}
