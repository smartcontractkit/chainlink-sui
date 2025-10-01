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
