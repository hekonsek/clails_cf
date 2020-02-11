package clails

type Driver interface {
	Validate(project *Project) error
	GenerateModel(project *Project) (map[string]interface{}, error)
	Generate(project *Project) (map[string]string, error)
}
