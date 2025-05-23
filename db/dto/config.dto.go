package dto

type AutomationConfig struct {
    Name        string        `json:"name" yaml:"name"`
    Description string        `json:"description,omitempty" yaml:"description,omitempty"`
    Schedule    string        `json:"schedule,omitempty" yaml:"schedule,omitempty"` // cron format
    Tasks       []TaskConfig  `json:"tasks" yaml:"tasks"`
}

type TaskConfig struct {
    ID          string            `json:"id" yaml:"id"`                     // unique task ID
    Type        string            `json:"type" yaml:"type"`                 // "builtin" or "custom"
    Name        string            `json:"name" yaml:"name"`                 // "SendEmail", "MyCustomScript"
    DependsOn   []string          `json:"depends_on,omitempty" yaml:"depends_on,omitempty"`
    Parameters  map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"` // task args
    Script      *ScriptConfig     `json:"script,omitempty" yaml:"script,omitempty"`         // only for custom
}

type ScriptConfig struct {
    Language string `json:"language" yaml:"language"`
    Code     string `json:"code" yaml:"code"`
}
