package clails

type Driver interface {
	Validate(project *Project) error
	GenerateModel(project *Project) (monitoring map[string]interface{}, environments map[string]interface{}, err error)
	Generate(project *Project) (monitoring string, environments map[string]string, err error)
}
