package request

type CreateAccountRequest struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}
