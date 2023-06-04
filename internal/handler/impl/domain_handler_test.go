package impl

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/test"
	"github.com/hmsidm/internal/test/mock/infrastructure/middleware"
	mock_client "github.com/hmsidm/internal/test/mock/interface/client"
	mock_interactor "github.com/hmsidm/internal/test/mock/interface/interactor"
	mock_presenter "github.com/hmsidm/internal/test/mock/interface/presenter"
	mock_repository "github.com/hmsidm/internal/test/mock/interface/repository"
	"github.com/openlyinc/pointy"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func helperDomainHandlerMocks(t *testing.T, db *gorm.DB) (
	*domainComponent,
	*mock_interactor.DomainInteractor,
	*mock_repository.DomainRepository,
	*mock_presenter.DomainPresenter,
	*mock_client.HostInventory,
) {
	interactorMock := mock_interactor.NewDomainInteractor(t)
	repositoryMock := mock_repository.NewDomainRepository(t)
	presenterMock := mock_presenter.NewDomainPresenter(t)
	inventoryMock := mock_client.NewHostInventory(t)

	// It returns the specific implementation to let to
	// write unit tests for private methods and methods that
	// does not belong the implemented interface.
	handler := &domainComponent{
		db:         db,
		interactor: interactorMock,
		repository: repositoryMock,
		presenter:  presenterMock,
		inventory:  inventoryMock,
	}

	return handler,
		interactorMock,
		repositoryMock,
		presenterMock,
		inventoryMock
}

func TestNewDomain(t *testing.T) {
	sqlMock, db, err := test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, sqlMock)
	require.NotNil(t, db)

	_, interactor, repository, presenter, inventory :=
		helperDomainHandlerMocks(t, db)

	service := NewDomain(db, interactor, repository, presenter, inventory)
	require.NotNil(t, service)
}

func TestListDomains(t *testing.T) {
	sqlMock, db, err := test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, sqlMock)
	require.NotNil(t, db)

	// handler, interactor, repository, presenter, inventory :=
	// 	helperDomainHandlerMocks(t, db)
	handler, interactor, repository, presenter, _ :=
		helperDomainHandlerMocks(t, db)

	// Panic on type assertion for the context
	assert.Panics(t, func() {
		err = handler.ListDomains(nil, public.ListDomainsParams{})
	}, "")

	// Error
	sqlMock, db, err = test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, sqlMock)
	require.NotNil(t, db)
	handler, interactor, _, _, _ =
		helperDomainHandlerMocks(t, db)
	ctx := middleware.NewDomainContextInterface(t)
	params := public.ListDomainsParams{}
	testXRHID := identity.XRHID{}
	ctx.On("XRHID").Return(&testXRHID)
	interactor.On("List", &testXRHID, &params).
		Return("", 0, 0, fmt.Errorf("FIXME Some error 1"))
	err = handler.ListDomains(ctx, params)
	require.EqualError(t, err, "FIXME Some error 1")

	// Begin transaction error
	testOrg := "12345"
	testOffset := 0
	testLimit := 10
	sqlMock, db, err = test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, sqlMock)
	require.NotNil(t, db)
	handler, interactor, _, _, _ =
		helperDomainHandlerMocks(t, db)
	ctx.On("XRHID").Return(&testXRHID)
	interactor.On("List", &testXRHID, &params).
		Return(testOrg, testOffset, testLimit, nil)
	sqlMock.ExpectBegin().
		WillReturnError(fmt.Errorf("FIXME Some error 2"))
	err = handler.ListDomains(ctx, params)
	require.EqualError(t, err, "FIXME Some error 2")

	// Repository List method error
	sqlMock, db, err = test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, sqlMock)
	require.NotNil(t, db)
	handler, interactor, repository, _, _ =
		helperDomainHandlerMocks(t, db)
	ctx.On("XRHID").Return(&testXRHID)
	interactor.On("List", &testXRHID, &params).
		Return(testOrg, testOffset, testLimit, nil)
	sqlMock.ExpectBegin().WillDelayFor(20 * time.Millisecond)
	repository.On("List", mock.Anything, testOrg, testOffset, testLimit).
		Return(nil, int64(0), fmt.Errorf("FIXME Some error 3"))
	err = handler.ListDomains(ctx, params)
	require.EqualError(t, err, "FIXME Some error 3")

	// Commit error
	testDomainUUID := uuid.MustParse("0ea29236-02ee-11ee-aff2-482ae3863d30")
	testDomainName := pointy.String("mydomain.example")
	testTitle := pointy.String("My Test Domain Title")
	testDescription := pointy.String("My Test Domain Description")
	testAutoEnrollmentEnabled := pointy.Bool(true)
	testType := pointy.Uint(model.DomainTypeIpa)
	testRealmName := pointy.String("MYDOMAIN.EXAMPLE")
	testDomains := []model.Domain{
		{
			OrgId:                 testOrg,
			DomainUuid:            testDomainUUID,
			DomainName:            testDomainName,
			Title:                 testTitle,
			Description:           testDescription,
			AutoEnrollmentEnabled: testAutoEnrollmentEnabled,
			Type:                  testType,
			IpaDomain: &model.Ipa{
				CaCerts:   []model.IpaCert{},
				Servers:   []model.IpaServer{},
				RealmName: testRealmName,
			},
		},
	}
	sqlMock, db, err = test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, sqlMock)
	require.NotNil(t, db)
	handler, interactor, repository, _, _ =
		helperDomainHandlerMocks(t, db)
	ctx.On("XRHID").Return(&testXRHID)
	interactor.On("List", &testXRHID, &params).
		Return(testOrg, testOffset, testLimit, nil)
	sqlMock.ExpectBegin().WillDelayFor(20 * time.Millisecond)
	repository.On("List", mock.Anything, testOrg, testOffset, testLimit).
		Return(testDomains, int64(len(testDomains)), nil)
	sqlMock.ExpectCommit().
		WillReturnError(fmt.Errorf("FIXME Some error 4"))
	err = handler.ListDomains(ctx, params)
	require.EqualError(t, err, "FIXME Some error 4")

	// Presenter List error
	testPrefix := "/api/hmsidm/v1"
	sqlMock, db, err = test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, sqlMock)
	require.NotNil(t, db)
	handler, interactor, repository, presenter, _ =
		helperDomainHandlerMocks(t, db)
	ctx.On("XRHID").Return(&testXRHID)
	interactor.On("List", &testXRHID, &params).
		Return(testOrg, testOffset, testLimit, nil)
	sqlMock.ExpectBegin().WillDelayFor(20 * time.Millisecond)
	repository.On("List", mock.Anything, testOrg, testOffset, testLimit).
		Return(testDomains, int64(len(testDomains)), nil)
	sqlMock.ExpectCommit()
	presenter.On("List", testPrefix, int64(len(testDomains)), testOffset, testLimit, testDomains).
		Return(nil, fmt.Errorf("FIXME Some error 5"))
	err = handler.ListDomains(ctx, params)
	require.EqualError(t, err, "FIXME Some error 5")

	// Handler ListDomains Successful
	presenterOutput := &public.ListDomainsResponse{
		Meta: public.PaginationMeta{
			Count:  1,
			Limit:  10,
			Offset: 0,
		},
		Links: public.PaginationLinks{
			First:    pointy.String(testPrefix + "/domains?offset=0&limit=10"),
			Previous: pointy.String(testPrefix + "/domains?offset=0&limit=10"),
			Next:     pointy.String(testPrefix + "/domains?offset=0&limit=10"),
			Last:     pointy.String(testPrefix + "/domains?offset=0&limit=10"),
		},
		Data: []public.ListDomainsData{
			{
				DomainId:              testDomainUUID,
				Title:                 *testTitle,
				Description:           *testDescription,
				DomainName:            *testDomainName,
				AutoEnrollmentEnabled: *testAutoEnrollmentEnabled,
				DomainType:            public.RhelIdm,
			},
		},
	}
	sqlMock, db, err = test.NewSqlMock(&gorm.Session{})
	require.NoError(t, err)
	require.NotNil(t, sqlMock)
	require.NotNil(t, db)
	handler, interactor, repository, presenter, _ =
		helperDomainHandlerMocks(t, db)
	ctx.On("XRHID").Return(&testXRHID)
	interactor.On("List", &testXRHID, &params).
		Return(testOrg, testOffset, testLimit, nil)
	sqlMock.ExpectBegin().WillDelayFor(20 * time.Millisecond)
	repository.On("List", mock.Anything, testOrg, testOffset, testLimit).
		Return(testDomains, int64(len(testDomains)), nil)
	sqlMock.ExpectCommit()
	presenter.On("List", testPrefix, int64(len(testDomains)), testOffset, testLimit, testDomains).
		Return(presenterOutput, nil)
	ctx.On("JSON", http.StatusOK, *presenterOutput).
		Return(nil)
	err = handler.ListDomains(ctx, params)
	require.NoError(t, err)
}
