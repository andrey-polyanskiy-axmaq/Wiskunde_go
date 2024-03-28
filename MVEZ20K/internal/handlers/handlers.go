package handlers

import (
	"MVEZ20K/internal/config"
	"MVEZ20K/internal/driver"
	"MVEZ20K/internal/models"
	"MVEZ20K/internal/render"
	"MVEZ20K/internal/repository"
	"MVEZ20K/internal/repository/dbrepo"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates the repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page Handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.html", &models.TemplateData{})
}

func (m *Repository) UsInfo(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "info.html", &models.TemplateData{})
}

func (m *Repository) Contacts(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contacts.html", &models.TemplateData{})
}
