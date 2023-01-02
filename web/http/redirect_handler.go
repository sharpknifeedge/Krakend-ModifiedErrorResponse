package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.boomerangapp.ir/back/pg/internal/gateway/atipay"
	"gitlab.boomerangapp.ir/back/pg/types/consts"
	"gitlab.boomerangapp.ir/back/utils/pkg/sentrylog"
)

func handleRedirect(w http.ResponseWriter, r *http.Request) {

	//get token from url
	vars := mux.Vars(r)
	token := vars["token"]

	//detect form via related psp
	var form []byte
	switch vars["psp"] {
	// case consts.PSPSep:
	// 	form = sep.GetForm(token)
	case consts.PSPAtipay:
		form = atipay.GetForm(token)
	}

	//write form to user
	_, err := w.Write(form)
	if err != nil {
		sentrylog.Info(err)
	}
}
