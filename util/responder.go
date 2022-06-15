package util

func SuccessResponse(data interface{}) map[string]interface{} {
	respData := map[string]interface{}{
		"code": "0",
		"msg":  "success",
	}
	if data == nil {
		respData["data"] = []interface{}{}
	} else {
		respData["data"] = data
	}
	return respData
}

func SuccessResponseWithTotal(data interface{}, total int) map[string]interface{} {
	respData := map[string]interface{}{
		"code": "0",
		"msg":  "success",
	}
	if data == nil {
		respData["data"] = []interface{}{}
	} else {
		respData["data"] = data
		respData["total"] = total
	}
	return respData
}

func ErrorResponse(code string, msg string) map[string]interface{} {
	return map[string]interface{}{
		"code": code,
		"msg":  msg,
	}
}
