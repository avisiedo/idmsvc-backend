package impl

import (
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/config"
	"github.com/hmsidm/internal/handler"
	"github.com/hmsidm/internal/interface/client"
	metrics "github.com/hmsidm/internal/metrics"
	usecase_interactor "github.com/hmsidm/internal/usecase/interactor"
	usecase_presenter "github.com/hmsidm/internal/usecase/presenter"
	usecase_repository "github.com/hmsidm/internal/usecase/repository"
	"gorm.io/gorm"
)

type application struct {
	public.ServerInterface
	config  *config.Config
	db      *gorm.DB
	metrics *metrics.Metrics
}

// NewHandler create a new application to handle
// public, private and metrics API.
// config is the configuration to use.
// db is the database conector.
// m hold the metrics available for the application.
// inventory represent the client to communicate to
// the host inventory service.
// Return a handler.Application instance, or panic if something
// goes wrong.
func NewHandler(config *config.Config, db *gorm.DB, m *metrics.Metrics, inventory client.HostInventory) handler.Application {
	if config == nil {
		panic("config is nil")
	}
	if db == nil {
		panic("db is nil")
	}
	i := usecase_interactor.NewDomainInteractor()
	r := usecase_repository.NewDomainRepository()
	p := usecase_presenter.NewDomainPresenter(config)

	// Instantiate application
	return &application{
		config:          config,
		db:              db,
		metrics:         m,
		ServerInterface: NewDomain(db, i, r, p, inventory),
	}
}
