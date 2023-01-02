package response

//Message Standard Persian Message
type Message struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}

//Page page number and pagination info
type Page struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

//Response Global response
type Response struct {
	Code       uint16      `json:"code"`
	Message    *Message    `json:"message,omitempty"`
	Pagination *Page       `json:"pagination,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func (r Response) Error() string {
	if r.Code == 0 {
		return ""
	}

	if r.Message != nil {
		return r.Message.Header
	}

	return "-"
}
