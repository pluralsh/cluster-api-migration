package api

// ManagedMachinePoolScaling specifies scaling options.
type ManagedMachinePoolScaling struct {
	MinSize int32 `json:"minSize"`
	MaxSize int32 `json:"maxSize"`
}

type Labels map[string]string
