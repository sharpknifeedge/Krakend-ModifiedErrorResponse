package http_helper

import (
	"encoding/json"
	"net/http"
	"log"

	"github.com/getsentry/sentry-go"
	"gitlab.boomerangapp.ir/back/utils/pkg/response"
)

//MuxResp is a response helper type for logging and writing
type MuxResp struct {
	Writer http.ResponseWriter
	Hub    *sentry.Hub
}

//SendResp to send proper response.
//page is optional so pass nil if don't need.
//if created == true then data should contain ID.
func (mr MuxResp) SendResp(data interface{}, page *response.Page, created bool) {

	if !mr.checkNil() {
		return
	}

	mr.Writer.Header().Set("Content-Type", "application/json")
	if created {
		mr.Writer.WriteHeader(http.StatusCreated)
	} else {
		mr.Writer.WriteHeader(http.StatusOK)
	}

	resp := response.ErrorResponse(response.Success)
	resp.Data = data
	resp.Pagination = page

	err := json.NewEncoder(mr.Writer).Encode(resp)
	if err != nil {
		mr.Hub.CaptureMessage(err.Error())
	}
}

//HandleServiceErr handles service special errors.
func (mr MuxResp) HandleServiceErr(eCode uint16, msgMap map[uint16]response.Message) {

	if !mr.checkNil() {
		return
	}

	mr.Writer.Header().Set("Content-Type", "application/json")
	mr.Writer.WriteHeader(http.StatusBadRequest)

	resp := response.ErrorResponse(eCode)
	msg := msgMap[eCode]
	resp.Message = &msg

	err := json.NewEncoder(mr.Writer).Encode(resp)
	if err != nil {
		mr.Hub.CaptureMessage(err.Error())
	}
}

//HandleValidationErr handles errors which mysql-Validation func returns and other validation errors
func (mr MuxResp) HandleValidationErr(err error) {

	if !mr.checkNil() {
		return
	}

	mr.Writer.Header().Set("Content-Type", "application/json")
	mr.Writer.WriteHeader(http.StatusUnprocessableEntity)

	resp := response.ErrorResponse(response.UnprocessableEntity)
	resp.Message = &response.Message{Header: err.Error()}

	err = json.NewEncoder(mr.Writer).Encode(resp)
	if err != nil {
		mr.Hub.CaptureMessage(err.Error())
	}
}

//HandleErr handles global errors
func (mr MuxResp) HandleErr(err error, eCode uint16) {

	if !mr.checkNil() {
		return
	}

	if err != nil {
		mr.Hub.CaptureException(err)
	}
	mr.Writer.Header().Set("Content-Type", "application/json")
	switch eCode {
	case response.UnprocessableEntity:
		mr.Writer.WriteHeader(422)
	case response.AccessDenied:
		mr.Writer.WriteHeader(403)
	case response.Unauthorized:
		mr.Writer.WriteHeader(401)
	case response.InternalError:
		mr.Writer.WriteHeader(503)
	case response.ConnectionError:
		mr.Writer.WriteHeader(502)
	default:
		mr.Writer.WriteHeader(400)
	}

	err = json.NewEncoder(mr.Writer).Encode(response.ErrorResponse(eCode))
	if err != nil {
		mr.Hub.CaptureMessage(err.Error())
	}
}

func (mr MuxResp) checkNil() bool {

	if mr.Writer == nil || mr.Hub == nil {
		log.Println("unable to handle response")
		return false
	}

	return true
}
