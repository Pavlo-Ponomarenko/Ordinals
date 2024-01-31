package service

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/kit/kv"
	"log"
	"net/http"
	"ordinals/internal/config"
	"ordinals/internal/data/pg"
	"ordinals/internal/service/handlers"
)

func (s *service) router() chi.Router {
	client, err := getBtcdClient()
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()
	conf := config.New(kv.MustFromEnv())
	q := pg.NewInscriptionQ(conf.DB())
	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxBtcdClient(client),
			handlers.CtxInscriptionQ(&q),
		),
	)
	err = client.NotifyNewTransactions(true)
	if err != nil {
		log.Fatal(err)
	}
	r.Route("/ordinals", func(r chi.Router) {
		// configure endpoints here
		r.Get("/menu", handlers.GetMenuPage)
		r.Get("/", handlers.GetMainPage)
		r.Post("/create_inscription", handlers.CreateInscription)
	})
	r.Handle("/ordinals/js/*", http.StripPrefix("/ordinals/js/", http.FileServer(http.Dir("templates/static/js/"))))
	return r
}
