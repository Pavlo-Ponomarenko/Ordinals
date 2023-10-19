package handlers

import (
	"context"
	"net/http"
	"ordinals/internal/wallet"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	accountInfoCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxAccountInfo(info *wallet.AccountInfo) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, accountInfoCtxKey, info)
	}
}

func AccountInfo(r *http.Request) *wallet.AccountInfo {
	return r.Context().Value(accountInfoCtxKey).(*wallet.AccountInfo)
}
