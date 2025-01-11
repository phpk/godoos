package types

type ChatRequest struct {
	Model       string                 `json:"model"`
	Engine      string                 `json:"engine"`
	Stream      bool                   `json:"stream"`
	WebSearch   bool                   `json:"webSearch"`
	FileContent string                 `json:"fileContent"`
	FileName    string                 `json:"fileName"`
	Options     map[string]interface{} `json:"options"`
	Messages    []Message              `json:"messages"`
	KnowledgeId uint                   `json:"knowledgeId"`
}

type AskDocResponse struct {
	Content  string   `json:"content"`
	Score    float32  `json:"score"`
	FilePath string   `json:"file_path"`
	FileName string   `json:"file_name"`
	Position Position `json:"position,omitempty"`
}

type Position struct {
	X0     float64 `json:"x0"`
	X1     float64 `json:"x1"`
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
	Page   int     `json:"page"`
	Line   int     `json:"line"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
}
type Image struct {
	Text   string `json:"text"`
	CosURL string `json:"cos_url"`
}
type OpenAIResponse struct {
	ID        string            `json:"id"`
	Created   int64             `json:"created"`
	RequestID string            `json:"request_id"`
	Model     string            `json:"model"`
	Choices   []Choice          `json:"choices"`
	Usage     Usage             `json:"usage"`
	WebSearch []WebSearchResult `json:"web_search,omitempty"`
	Documents []AskDocResponse  `json:"documents"`
	Problems  []string          `json:"problems"`
	Images    []string          `json:"images"`
}
type Choice struct {
	Index        int        `json:"index"`
	FinishReason string     `json:"finish_reason"`
	Message      Message    `json:"message"`
	ToolCalls    []ToolCall `json:"tool_calls,omitempty"`
}
type Message struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

type ToolCall struct {
	Function FunctionRes `json:"function"`
	ID       string      `json:"id"`
	Type     string      `json:"type"`
}
type FunctionRes struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type WebSearchResult struct {
	Icon    string `json:"icon"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	Media   string `json:"media"`
	Content string `json:"content"`
}
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
