package router

import (
	"avitoTech/internal/controller"
	"avitoTech/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(services *service.Services) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	bc := controller.NewBidController(services.Bid)
	tc := controller.NewTenderController(services.Tender)

	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", controller.Ping)

		r.Route("/tenders", func(r chi.Router) {
			r.Get("/", tc.GetTenders)

			r.Post("/new", tc.CreateTender)

			r.Get("/my", tc.GetUserTenders)
			r.Get("/{tenderId}/status", tc.GetTenderStatus)

			r.Put("/{tenderId}/status", tc.UpdateTenderStatus)

			r.Patch("/{tenderId}/edit", tc.EditTender)

			r.Put("/{tenderId}/rollback/{version}", tc.RollbackTender)

		})

		r.Route("/bids", func(r chi.Router) {
			r.Post("/new", bc.CreateBid)

			r.Get("/my", bc.GetUserBids)
			r.Get("/{tenderId}/list", bc.GetBidsForTender)
			r.Get("/{bidId}/status", bc.GetBidStatus)

			r.Put("/{bidId}/status", bc.UpdateBidStatus)

			r.Patch("/{bidId}/edit", bc.EditBid)

			r.Put("/{bidId}/submit_decision", bc.SubmitBidDecision)
			r.Put("/{bidId}/feedback", bc.SubmitBidFeedback)
			r.Put("/{bidId}/rollback/{version}", bc.RollbackBid)

			r.Get("/{tenderId}/reviews", bc.GetBidReviews)
		})
	})

	return r
}
