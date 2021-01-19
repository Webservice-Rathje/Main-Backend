package generalModels

type ErrorResponseModel struct {
	Error          string `json:"error"`
	CausedBy       string `json:"caused_by"`
	CouldBeFixedBy string `json:"could_be_fixed_by"`
	Alert          string `json:"alert"`
}
