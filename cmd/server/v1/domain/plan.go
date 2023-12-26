package serverDomain

type PlanExecReq struct {
	PlanId int    `json:"planId"`
	Uuid   string `json:"uuid"`
	Title  string `json:"title"`
	Vus    int    `json:"vus"`

	NsqServerAddress string `json:"nsqServerAddress,omitempty"`
	NsqLookupAddress string `json:"nsqLookupAddress,omitempty"`
}
