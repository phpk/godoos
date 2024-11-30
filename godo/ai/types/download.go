package types

type OllamaModelsList struct {
	Models []OllamaModelsInfo `json:"models"`
}

type OllamaDetails struct {
	ParameterSize     string `json:"parameter_size"`
	QuantizationLevel string `json:"quantization_level"`
}
type OllamaModelsInfo struct {
	Model   string        `json:"model"`
	Details OllamaDetails `json:"details"`
	Size    int64         `json:"size"`
}
type OllamaModelDetail struct {
	Parameters string                 `json:"parameters"`
	Template   string                 `json:"template"`
	Details    map[string]interface{} `json:"details"`
	ModelInfo  map[string]interface{} `json:"model_info"`
}

// type ResModelInfo struct {
// 	Parameters      string `json:"parameters"`
// 	Template        string `json:"template"`
// 	ContextLength   int64  `json:"context_length"`
// 	EmbeddingLength int64  `json:"embedding_length"`
// 	Size            string `json:"size"`
// 	Quant           string `json:"quant"`
// 	Desk            string `json:"desk"`
// 	Cpu             string `json:"cpu"`
// 	Gpu             string `json:"gpu"`
// }

type Layer struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
	From      string `json:"from,omitempty"`
	status    string
}
type ManifestV2 struct {
	SchemaVersion int      `json:"schemaVersion"`
	MediaType     string   `json:"mediaType"`
	Config        *Layer   `json:"config"`
	Layers        []*Layer `json:"layers"`
}
type OmodelPath struct {
	Space   string
	LibPath string
	Name    string
	Tag     string
}
