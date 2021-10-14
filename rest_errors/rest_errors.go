package rest_errors

import "net/http"

type RestErr struct {
	Message    string        `json:"message"`
	StatusCode int           `json:"statuscode"`
	Causes     []interface{} `json:"causes"`
}

func NewError(msg string, code int) *RestErr {
	return &RestErr{
		Message:    msg,
		StatusCode: code,
	}
}

func NewBadRequestError(msg string) *RestErr {
	return &RestErr{
		Message:    msg,
		StatusCode: http.StatusBadRequest,
	}
}

func NewNotFoundError(msg string) *RestErr {
	return &RestErr{
		Message:    msg,
		StatusCode: http.StatusNotFound,
	}
}

func NewInternalServerError(msg string, err error) *RestErr {
	result := &RestErr{
		Message:    msg,
		StatusCode: http.StatusInternalServerError,
	}
	if err != nil {
		result.Causes = append(result.Causes, err.Error())
	}
	return result
}
