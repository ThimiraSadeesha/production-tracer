package interceptors

type Result struct {
	Response  int         `json:"response"`
	Path      string      `json:"path"`
	Timestamp string      `json:"timestamp"`
	Error     *int        `json:"error,omitempty"`
	Message   interface{} `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type LegacyResult struct {
	Status      int         `json:"status,omitempty"`
	Environment string      `json:"environment,omitempty"`
	RequestID   interface{} `json:"request_id,omitempty"`
	Path        interface{} `json:"path,omitempty"`
	Timestamp   string      `json:"timestamp,omitempty"`
	Message     string      `json:"message,omitempty"`
	Data        interface{} `json:"data,omitempty"`
	Errors      interface{} `json:"errors,omitempty"`
	Pagination  interface{} `json:"pagination,omitempty"`
}
