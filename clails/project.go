package clails

type Project struct {
	Name         string
	Environments []string
	Services     []*Service
}

type Service struct {
	Type         string
	Distribution string
}
