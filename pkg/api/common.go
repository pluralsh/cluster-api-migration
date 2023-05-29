package api

// ManagedMachinePoolScaling specifies scaling options.
type ManagedMachinePoolScaling struct {
	MinSize int32 `json:"minSize,omitempty"`
	MaxSize int32 `json:"maxSize,omitempty"`
}
