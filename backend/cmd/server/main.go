package main

import (
	"log"
	"net/http"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/handlers"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/models"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/auth/services"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/collaboration"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/config"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/db"
	dashboardhandler "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/handlers"
	dochandler "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/handlers"
	rolehandler "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/handlers"
	docmodel "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/models"
	rolemodel "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/models"
	dashboardservice "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/services"
	docservice "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/services"
	roleservice "github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/documents/services"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/routes"
)

func main() {
	cfg := config.Load()
	log.Println("Port:", cfg.Port)

	// PORT SETUP
	port := cfg.Port
	log.Printf("Starting up on http://localhost:%s", port)

	db.ConnectDB()

	userModel := models.NewUserModel(db.DB)
	documentModel := docmodel.NewDocumentModel(db.DB)
	roleModel := rolemodel.NewRoleModel(db.DB)
	versionModel := docmodel.NewVersionModel(db.DB)

	userService := services.NewUserService(userModel)
	documentService := docservice.NewDocumentService(documentModel, roleModel, versionModel, db.DB)
	roleService := roleservice.NewRoleService(roleModel, userModel, documentModel, db.DB)
	dashboardService := dashboardservice.NewDashboardService(db.DB, documentModel, roleModel, userModel)
	collaborationService := collaboration.NewCollaborationService(documentService)

	userHandler := handlers.NewUserHandler(userService)
	documentHandler := dochandler.NewDocumentHandler(documentService)
	roleHandler := rolehandler.NewRoleHandler(roleService)
	dashboardHandler := dashboardhandler.NewDashboardHandler(dashboardService)

	hub := collaboration.NewHub()
	collaborationHandler := collaboration.NewCollaborationHandler(roleService, collaborationService, hub)

	go hub.Run()

	// ROUTER SETUP
	r := routes.SetupRouter(userHandler, documentHandler, roleHandler, dashboardHandler, collaborationHandler)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
