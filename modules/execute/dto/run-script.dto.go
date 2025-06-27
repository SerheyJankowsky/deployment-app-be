package dto

type RunScriptDto struct {
	ScriptID uint `json:"script_id" validate:"required"`
	ServerID uint `json:"server_id" validate:"required"`
	EnvID    uint `json:"env_id" validate:"omitempty"`
	LoadEnv  bool `json:"load_env" validate:"required"`
}
