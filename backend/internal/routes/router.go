package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/handlers"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/collaboration"
	doc "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/handlers"
)

func SetupRouter(userHandler *handlers.UserHandler, documentHandler *doc.DocumentHandler, roleHandler *doc.RoleHandler, dashboardHandler *doc.DashboardHandler, collaborationHandler *collaboration.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Use(cors.Handler(cors.Options{

	// }))
	UserRoutes(r, userHandler)
	DocumentRoutes(r, documentHandler, roleHandler, dashboardHandler)
	CollaborationRoute(r, collaborationHandler)
	return r
}
