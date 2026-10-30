package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/handlers"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/collaboration"
	doc "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/handlers"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/middleware"
)

func UserRoutes(r chi.Router, userHandler *handlers.UserHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", userHandler.SignUpUser)
		r.Post("/login", userHandler.LoginUser)
		r.With(middleware.AuthMiddleware).Get("/test", userHandler.Test)
	})
}

func DocumentRoutes(r chi.Router, documentHandler *doc.DocumentHandler, roleHandler *doc.RoleHandler, dashboardHandler *doc.DashboardHandler) {
	r.Route("/dashboard", func(r chi.Router) {
		r.With(middleware.AuthMiddleware).Get("/", dashboardHandler.GetDashboard)
	})

	r.Route("/docs", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Post("/", documentHandler.CreateDoc)
		r.Get("/all", documentHandler.GetMyAllDocs)

		r.Route("/{docID}", func(r chi.Router) {
			r.Get("/", documentHandler.GetDocByID)
			r.Put("/", documentHandler.UpdateDoc)
			r.Delete("/", documentHandler.DeleteDoc)

			r.Route("/versions", func(r chi.Router) {
				r.Get("/", documentHandler.GetVersions)
				r.Get("/{version}", documentHandler.GetVersion)
				r.Post("/save", documentHandler.SaveVersion)
				r.Post("/{version}/restore", documentHandler.RestoreVersion)
			})

			r.Route("/members", func(r chi.Router) {
				r.Get("/{memberID}", roleHandler.GetMemberDetails)
				r.Get("/", roleHandler.GetMembersList)
				r.Put("/{memberID}", roleHandler.UpdateRole)
				r.Delete("/{memberID}", roleHandler.DeleteMember)
			})

			r.Route("/invites", func(r chi.Router) {
				r.Post("/", roleHandler.CreateInvite)
				r.Get("/", roleHandler.GetSentInvites)
				r.Post("/{inviteID}/accept", roleHandler.AcceptInvite)

			})
		})
	})

	r.Route("/user/invites", func(r chi.Router) {
		r.With(middleware.AuthMiddleware).Get("/", roleHandler.GetInvites)
	})
}

func CollaborationRoute(r chi.Router, collaborationHandler *collaboration.Handler) {
	// r.Get("/ws", collaborationHandler.HandleConnection)
	r.With(middleware.AuthMiddleware).Get("/ws/documents/{docID}", collaborationHandler.HandleConnection)
}
