package function

type FunctionInfo struct {
	Package    string              `json:"package"`
	Module     string              `json:"module"`
	Name       string              `json:"name"`
	Parameters []FunctionParameter `json:"parameters"`
}

type FunctionParameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (fi *FunctionInfo) GetParameters() ([]string, []string) {
	names := make([]string, len(fi.Parameters))
	types := make([]string, len(fi.Parameters))
	for i, param := range fi.Parameters {
		names[i] = param.Name
		types[i] = param.Type
	}

	return names, types
}
