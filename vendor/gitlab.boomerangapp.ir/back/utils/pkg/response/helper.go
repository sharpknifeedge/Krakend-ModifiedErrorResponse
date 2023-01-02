package response

//ErrorResponse make error `Response`
//handle Message with defined messages
func ErrorResponse(code uint16) Response {
	var out = Response{}
	out.Code = code

	if code < uint16(len(messages)) {
		out.Message = &messages[code]
	}
	return out
}
