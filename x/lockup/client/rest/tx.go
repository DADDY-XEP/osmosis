package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/c-osmosis/osmosis/x/lockup/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/lockup/locktokens", newLockTokensHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc("/lockup/unlock", newUnlockTokensHandlerFn(clientCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/lockup/unlock/{%s}", LockID), newUnlockByIDHandlerFn(clientCtx)).Methods("POST")
}

func newLockTokensHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LockTokensReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		duration, err := time.ParseDuration(req.Duration)
		if err != nil {
			return
		}

		msg := &types.MsgLockTokens{
			Owner:    req.Owner,
			Duration: duration,
			Coins:    req.Coins,
		}
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func newUnlockTokensHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UnlockTokensReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// create the message
		msg := &types.MsgUnlockTokens{
			Owner: req.Owner,
		}

		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func newUnlockByIDHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strLockID := vars[LockID]

		if len(strLockID) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "proposalId required but not specified")
			return
		}

		id, ok := rest.ParseUint64OrReturnBadRequest(w, strLockID)
		if !ok {
			return
		}

		var req UnlockTokensByIDReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// create the message
		msg := &types.MsgUnlockPeriodLock{
			Owner: req.Owner,
			ID:    id,
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
