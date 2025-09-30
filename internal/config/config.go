package config

type ServerConfiguration struct {
	Port         int    `json:"port,omitempty"`
	DatabaseName string `json:"database_name,omitempty"`
	LogsDir      string `json:"logs_dir,omitempty"`
}
