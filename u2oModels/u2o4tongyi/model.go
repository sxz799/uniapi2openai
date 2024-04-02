package u2o4tongyi

type TongyiBody struct {
	Model      string     `json:"model"`
	Input      Input      `json:"input"`
	Parameters Parameters `json:"parameters"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Input struct {
	Messages []Message `json:"messages"`
}
type Parameters struct {
	ResultFormat      string `json:"result_format"`
	IncrementalOutput bool   `json:"incremental_output"`
}

type TongYiResponseBody struct {
	Output    Output `json:"output"`
	Usage     Usage  `json:"usage"`
	RequestID string `json:"request_id"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choices struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}
type Output struct {
	Choices []Choices `json:"choices"`
}

func NewTongyiBody(modelName string, messages []Message) *TongyiBody {
	return &TongyiBody{
		Model: modelName,
		Input: Input{
			Messages: messages,
		},
		Parameters: Parameters{
			ResultFormat:      "message",
			IncrementalOutput: true,
		},
	}
}
