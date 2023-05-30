package api

type ManagedMachinePoolScaling struct {
	MinSize int32 `json:"minSize"`
	MaxSize int32 `json:"maxSize"`
}

const (
	LinuxOS   = "Linux"
	WindowsOS = "Windows"
)
