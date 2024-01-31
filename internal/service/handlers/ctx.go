package handlers

import (
	"context"
	"github.com/btcsuite/btcd/rpcclient"
	"gitlab.com/distributed_lab/logan/v3"
	"net/http"
	"ordinals/internal/data"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	btcdClientCtxKey
	inscriptionQCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxBtcdClient(client *rpcclient.Client) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, btcdClientCtxKey, client)
	}
}

func BtcdClient(r *http.Request) *rpcclient.Client {
	return r.Context().Value(btcdClientCtxKey).(*rpcclient.Client)
}

func CtxInscriptionQ(q *data.InscriptionQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, inscriptionQCtxKey, *q)
	}
}

func InscriptionQ(r *http.Request) data.InscriptionQ {
	return r.Context().Value(inscriptionQCtxKey).(data.InscriptionQ).New()
}
