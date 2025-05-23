package dto

type AutomationConfig struct {
    Name        string        `json:"name" yaml:"name"`
    Description string        `json:"description,omitempty" yaml:"description,omitempty"`
    Tasks       []TaskConfig  `json:"tasks" yaml:"tasks"`
}

type TaskConfig struct {
    ID          string            `json:"id" yaml:"id"`                     
    Type        string            `json:"type" yaml:"type"`                
    Name        string            `json:"name" yaml:"name"`                
    DependsOn   []string          `json:"depends_on,omitempty" yaml:"depends_on,omitempty"`
    Parameters  map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"` 
    Script      *ScriptConfig     `json:"script,omitempty" yaml:"script,omitempty"`         
}

type ScriptConfig struct {
    Language string `json:"language" yaml:"language"`
    Code     string `json:"code" yaml:"code"`
}
