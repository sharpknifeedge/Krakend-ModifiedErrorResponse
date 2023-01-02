package http

import (
	"context"
	"net/http"
	"strconv"

	"gitlab.boomerangapp.ir/back/pg/internal/db"
	"gitlab.boomerangapp.ir/back/pg/types/consts"
	utilsConsts "gitlab.boomerangapp.ir/back/utils/consts"
	"gitlab.boomerangapp.ir/back/utils/pkg/http_helper"
	"gitlab.boomerangapp.ir/back/utils/pkg/response"
	"gitlab.boomerangapp.ir/back/utils/pkg/sentrylog"
	"gitlab.boomerangapp.ir/back/utils/pkg/validation"
)

func handleTransactions(w http.ResponseWriter, r *http.Request) {

	//init response handler
	mr := http_helper.MuxResp{
		Writer: w,
		Hub:    sentrylog.GetHubFromCtx(r.Context()),
	}

	//get user|producer id from header
	role := r.Header.Get(utilsConsts.HTTPROleHeader)
	if role != utilsConsts.RoleProducer && role != utilsConsts.RoleUser {
		mr.HandleErr(nil, response.AccessDenied)
		return
	}
	userID, _ := strconv.ParseUint(r.Header.Get(utilsConsts.HTTPUserID), 10, 32)

	//check payment id
	paymentID := r.URL.Query().Get("id")
	if paymentID != "" {
		id, _ := strconv.ParseUint(paymentID, 10, 32)
		handleWithPaymentID(r.Context(), uint(id), uint(userID), mr, w)
	} else {
		handleWithUserID(uint(userID), role, mr, w, r)
	}
}

func handleWithPaymentID(ctx context.Context, paymentID, userID uint, mr http_helper.MuxResp, w http.ResponseWriter) {

	//load payment
	payment, err := db.LoadPayment(ctx, nil, paymentID)
	if err != nil {
		vErr := validation.ValidateByErr(err)
		if vErr != nil {
			mr.HandleValidationErr(vErr)
		} else {
			mr.HandleErr(err, response.InternalError)
		}
		return
	}

	//check userID
	if (payment.UserID.Valid && payment.UserID.Uint == userID) ||
		(payment.ProducerID.Valid && payment.ProducerID.Uint == userID) {
		mr.SendResp(payment, nil, false)
		return
	}

	mr.HandleServiceErr(consts.ErrCodeUserIDConflict, consts.Messages)
}

func handleWithUserID(userID uint, role string, mr http_helper.MuxResp, w http.ResponseWriter, r *http.Request) {

	//get limit, offset
	l := r.URL.Query().Get("limit")
	o := r.URL.Query().Get("offset")
	var limit, offset uint64

	//check default limit, offset
	if l == "" {
		limit = 20
	} else {
		limit, _ = strconv.ParseUint(l, 10, 32)
	}
	if o != "" {
		offset, _ = strconv.ParseUint(o, 10, 32)
	}

	//load payments
	isProducer := false
	if role == utilsConsts.RoleProducer {
		isProducer = true
	}
	payments, total, err := db.LoadPaymentsByUserID(r.Context(), userID, isProducer, int(limit), int(offset))
	if err != nil {
		vErr := validation.ValidateByErr(err)
		if vErr != nil {
			mr.HandleValidationErr(vErr)
		} else {
			mr.HandleErr(err, response.InternalError)
		}
		return
	}

	//send response
	mr.SendResp(payments, &response.Page{
		Total:  int(total),
		Limit:  int(limit),
		Offset: int(offset),
	}, false)
}
