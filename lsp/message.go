package lsp

type Request struct {
    RPC     string  `json:"jsonrpc"`
    ID      int     `json:"id"`
    Method  string  `json:"method"`

    // types of params will be specified later on the Request types
}

type Response struct {
    RPC     string  `json:"jsonrpc"`
    ID      *int    `json:"id,omitempty"`

    // results and errors will be specified later on the Response types
}

type Notification struct {
    RPC     string  `json:"jsonrpc"`
    Method  string  `json:"method"`
}
