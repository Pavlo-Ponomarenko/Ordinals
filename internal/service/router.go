package service

import (
	"gitlab.com/distributed_lab/kit/kv"
	"ordinals/internal/service/handlers"
	"ordinals/internal/wallet"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()
	config := kv.MustFromEnv()
	accounts, err := config.GetStringMap("accounts")
	var mainAddress string
	if err == nil {
		mainAddress = accounts["main"].(string)
	}

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxAccountInfo(&wallet.AccountInfo{mainAddress}),
		),
	)
	r.Route("/integrations/Ordinals", func(r chi.Router) {
		// configure endpoints here
	})

	return r
}
