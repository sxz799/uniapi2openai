package u2o4tongyiWeb

type TongYiWebRespBody struct {
	CanFeedback    bool       `json:"canFeedback"`
	CanRegenerate  bool       `json:"canRegenerate"`
	CanShare       bool       `json:"canShare"`
	CanShow        bool       `json:"canShow"`
	ContentFrom    string     `json:"contentFrom"`
	ContentType    string     `json:"contentType"`
	Contents       []Contents `json:"contents"`
	MsgID          string     `json:"msgId"`
	MsgStatus      string     `json:"msgStatus"`
	Params         Params     `json:"params"`
	ParentMsgID    string     `json:"parentMsgId"`
	SessionID      string     `json:"sessionId"`
	SessionOpen    bool       `json:"sessionOpen"`
	SessionShare   bool       `json:"sessionShare"`
	SessionWarnNew bool       `json:"sessionWarnNew"`
	StopReason     string     `json:"stopReason"`
	TraceID        string     `json:"traceId"`
}
type Contents struct {
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
	ID          string `json:"id"`
	Role        string `json:"role"`
	Status      string `json:"status"`
}
type Params struct {
}
