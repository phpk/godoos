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

type InvokeResp struct {
	RequestID      string   `json:"requestId"`
	Content        string   `json:"content"`
	Problems       []string `json:"problems"`
	DocumentSlices []struct {
		Document      Document `json:"document"`
		SliceInfo     []Slice  `json:"slice_info"`
		HidePositions bool     `json:"hide_positions"`
		Images        []Image  `json:"images"`
	} `json:"documents"`
}
type AskDocResponse struct {
	Content  string  `json:"content"`
	Score    float32 `json:"score"`
	FilePath string  `json:"file_path"`
	FileName string  `json:"file_name"`
}
type Document struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Dtype int    `json:"dtype"`
}

type Slice struct {
	DocumentID string    `json:"document_id"`
	Position   *Position `json:"position,omitempty"`
	Line       int       `json:"line,omitempty"`
	SheetName  string    `json:"sheet_name,omitempty"`
	Text       string    `json:"text"`
}
type Position struct {
	X0     float64 `json:"x0"`
	X1     float64 `json:"x1"`
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
	Page   int     `json:"page"`
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
