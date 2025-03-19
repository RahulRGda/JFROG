// nolint:golint,errcheck
package logging

type ExtraFieldsMap map[string]interface{}

type StandardFields struct {
	Environment    string `json:"environment"`
	Application_id string `json:"application_id"`
	Product_id     string `json:"product_id"`
	Account_id     int    `json:"account_id"`
	Trace_id       string `json:"trace_id"`
	Trace_context  string `json:"trace_context"`
	Channel_id     string `json:"channel_id"`
	User_id        int    `json:"user_id"`
	Activity       string `json:"activity"`
	Line           int    `json:"line"`
	Event          string `json:"event"`
	Action         string `json:"action"`
	Account_name   string `json:"account_name"`
	User_email     string `json:"user_email"`
}
