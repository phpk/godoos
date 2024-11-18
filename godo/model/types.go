package model

import "time"

// 参数	类型	描述
// llama_model_path	字符串	LLaMA模型的文件路径。
// ngl	整数	使用的GPU层数。
// ctx_len	整数	模型操作的上下文长度。
// embedding	布尔值	是否在模型中使用嵌入。
// n_parallel	整数	并行操作的数量。
// cont_batching	布尔值	是否使用连续批处理。
// user_prompt	字符串	用于用户的提示。
// ai_prompt	字符串	用于AI助手的提示。
// system_prompt	字符串	用于系统规则的提示。
// pre_prompt	字符串	用于内部配置的提示。
// cpu_threads	整数	推理时使用的线程数（仅CPU模式）。
// n_batch	整数	提示评估步骤的批次大小。
// caching_enabled	布尔值	是否启用提示缓存。
// clean_cache_threshold	整数	触发清理缓存操作的聊天数量。
// grp_attn_n	整数	自我扩展中组注意力因子。
// grp_attn_w	整数	自我扩展中组注意力宽度。
// mlock	布尔值	在macOS中防止系统将模型交换到磁盘。
// grammar_file	字符串	通过提供语法文件路径，您可以使用GBNF语法约束采样。
// model_type	字符串	我们想要使用的模型类型：llm 或 embedding，默认值为 llm。
type ModelConfig struct {
	ModelAlias     string `json:"model_alias"`
	PromptTemplate string `json:"prompt_template"`
	LlamaModelPath string `json:"llama_model_path"` // The file path to the LLaMA model.
	Mmproj         string `json:"mmproj"`
	ModelType      string `json:"model_type"`  // Model type we want to use: llm or embedding, default value is llm
	CPUThreads     int    `json:"cpu_threads"` // The number of threads to use for inferencing (CPU MODE ONLY)
	NGL            int    `json:"ngl"`         // The number of GPU layers to use.
	CtxLen         int    `json:"ctx_len"`     // The context length for the model operations.
	Embedding      bool   `json:"embedding"`   // Whether to use embedding in the model.

	UserPrompt   string `json:"user_prompt"`   // The prompt to use for the user.
	AIPrompt     string `json:"ai_prompt"`     // The prompt to use for the AI assistant.
	SystemPrompt string `json:"system_prompt"` // The prompt to use for system rules.
	// PrePrompt           string `json:"pre_prompt"`            // The prompt to use for internal configuration.

	// NParallel           int    `json:"n_parallel"`            // The number of parallel operations.
	// ContBatching        bool   `json:"cont_batching"`         // Whether to use continuous batching.
	// NBatch              int    `json:"n_batch"`               // The batch size for prompt eval step
	// CachingEnabled      bool   `json:"caching_enabled"`       // To enable prompt caching or not
	// CleanCacheThreshold int    `json:"clean_cache_threshold"` // Number of chats that will trigger clean cache action
	GrpAttnN int `json:"grp_attn_n"` // Group attention factor in self-extend
	GrpAttnW int `json:"grp_attn_w"` // Group attention width in self-extend
	// Mlock               bool   `json:"mlock"`                 // Prevent system swapping of the model to disk in macOS
	GrammarFile string `json:"grammar_file"` // You can constrain the sampling using GBNF grammars by providing path to a grammar file

}
type FileProgress struct {
	Progress   float64 `json:"progress"` // 将进度改为浮点数，以百分比表示
	IsFinished bool    `json:"is_finished"`
	Total      int64   `json:"total"`
	Current    int64   `json:"completed"`
	Status     string  `json:"status"`
}
type ModelStruct struct {
	Model string `json:"model"`
}
type ReqBody struct {
	//DownloadUrl string                 `json:"url"`
	//Options  ModelConfig            `json:"options"`
	Model     string                 `json:"model"`
	Url       []string               `json:"url"`
	Engine    string                 `json:"engine"`
	Type      string                 `json:"type"`
	From      string                 `json:"from"`
	Action    []string               `json:"action"`
	Label     string                 `json:"label"`
	Info      map[string]interface{} `json:"info"`
	Status    string                 `json:"status"`
	Paths     []string               `json:"paths"`
	Params    map[string]interface{} `json:"params"`
	FileName  string                 `json:"file_name"`
	CreatedAt time.Time              `json:"created_at"`
}
type DownloadsRequest struct {
	Urls []string `json:"urls"`
	Dir  string   `json:"model_path"`
}

//	type DelBody struct {
//		DownloadUrl string `json:"url"`
//		ModelDir    string `json:"name"`
//	}
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
