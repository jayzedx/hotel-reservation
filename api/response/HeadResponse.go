package response

type HeadResponse struct {
	Status     string                  `json:"status"`
	StatusCode int                     `json:"statusCode"`
	Msg        string                  `json:"msg"`
	Data       *map[string]interface{} `json:"data"`
}
