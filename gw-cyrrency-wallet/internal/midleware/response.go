package midleware

type ResponseFetch struct {
	Status      string   `json:"status"`
	StatusCode  int      `json:"status_code"`
	TextError   string   `json:"error,omitempty"`
	Id          []string `json:"id,omitempty"` //возвращается при Insert или Update
	AccessToken string   `json:"access_token,omitempty"`
}

func Bearer(token string) *ResponseFetch {
	return &ResponseFetch{
		Status:      "Ok",
		StatusCode:  200,
		AccessToken: token,
	}
}
