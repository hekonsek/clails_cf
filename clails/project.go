package clails

type Project struct {
	Name     string
	Services []*Service
}

type Service struct {
	Type         string
	Distribution string
}
