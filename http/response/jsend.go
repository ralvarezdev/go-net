package response

type (
	// JSendBody struct
	JSendBody struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}

	// JSendSuccessBody struct
	JSendSuccessBody struct {
		JSendBody
	}

	// JSendFailBody struct
	JSendFailBody struct {
		JSendBody
		Code *string `json:"code,omitempty"`
	}

	// JSendErrorBody struct
	JSendErrorBody struct {
		JSendBody
		Message *string `json:"message,omitempty"`
		Code    *string `json:"code,omitempty"`
	}
)

// NewJSendSuccessBody creates a new JSend success response body
func NewJSendSuccessBody(
	data interface{},
) *JSendSuccessBody {
	return &JSendSuccessBody{
		JSendBody: JSendBody{
			Status: "success",
			Data:   data,
		},
	}
}

// NewJSendFailBody creates a new JSend fail response body
func NewJSendFailBody(
	data interface{},
	code *string,
) *JSendFailBody {
	return &JSendFailBody{
		JSendBody: JSendBody{
			Status: "fail",
			Data:   data,
		},
		Code: code,
	}
}

// NewJSendErrorBody creates a new JSend error response body
func NewJSendErrorBody(
	message string,
	data interface{},
	code *string,
) *JSendErrorBody {
	return &JSendErrorBody{
		JSendBody: JSendBody{
			Status: "error",
			Data:   data,
		},
		Message: &message,
		Code:    code,
	}
}
