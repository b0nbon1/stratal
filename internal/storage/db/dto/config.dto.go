package dto

type TaskConfig struct {
	DependsOn  []string          `json:"depends_on,omitempty" yaml:"depends_on,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Secrets    map[string]string `json:"secrets,omitempty" yaml:"secrets,omitempty"` // secret_name -> env_var_name
	Script     *ScriptConfig     `json:"script,omitempty" yaml:"script,omitempty"`
}

type ScriptConfig struct {
	Language string `json:"language" yaml:"language"`
	Code     string `json:"code" yaml:"code"`
}