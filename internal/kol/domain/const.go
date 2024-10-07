package domain

type Sex string

const (
	SexMale   Sex = "m"
	SexFemale Sex = "f"
)

func (s Sex) IsValid() bool {
	switch s {
	case SexMale, SexFemale:
		return true
	default:
		return false
	}
}

func (s Sex) String() string {
	return string(s)
}
