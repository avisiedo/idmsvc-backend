package presenter

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/config"
	"github.com/hmsidm/internal/domain/model"
	"github.com/lib/pq"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuardsRegisterIpa(t *testing.T) {
	var (
		err error
	)
	p := &domainPresenter{cfg: &cfg}
	assert.Panics(t, func() {
		err = p.sharedDomainFillRhelIdm(nil, nil)
	})

	domain := &model.Domain{}
	domain.Type = pointy.Uint(999)
	err = p.sharedDomainFillRhelIdm(domain, nil)
	assert.EqualError(t, err, fmt.Sprintf("'domain.Type' is not '%s'", model.DomainTypeIpaString))

	*domain.Type = model.DomainTypeIpa
	err = p.sharedDomainFillRhelIdm(domain, nil)
	assert.EqualError(t, err, "'domain.IpaDomain' is nil")

	domain.IpaDomain = &model.Ipa{}
	assert.Panics(t, func() {
		err = p.sharedDomainFillRhelIdm(domain, nil)
	})

	output := &public.RegisterDomainResponse{}
	err = p.sharedDomainFillRhelIdm(domain, output)
	assert.NoError(t, err)

	domain.IpaDomain.CaCerts = []model.IpaCert{}
	assert.NotPanics(t, func() {
		err = p.sharedDomainFillRhelIdm(domain, output)
	})

	output.Type = public.DomainTypeRhelIdm
	output.RhelIdm = &public.DomainIpa{}
	err = p.sharedDomainFillRhelIdm(domain, output)
	assert.NoError(t, err)
}

func TestRegisterRhelIdm(t *testing.T) {
	tokenExpiration := &time.Time{}
	*tokenExpiration = time.Now()
	type TestCaseExpected struct {
		Domain *public.RegisterDomainResponse
		Err    error
	}
	type TestCase struct {
		Name     string
		Given    *model.Domain
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "Success minimal rhel-idm content",
			Given: &model.Domain{
				Type: pointy.Uint(model.DomainTypeIpa),
				IpaDomain: &model.Ipa{
					RealmName:       pointy.String(""),
					CaCerts:         []model.IpaCert{},
					Servers:         []model.IpaServer{},
					RealmDomains:    pq.StringArray{},
					Token:           pointy.String("71ad4978-c768-11ed-ad69-482ae3863d30"),
					TokenExpiration: tokenExpiration,
				},
			},
			Expected: TestCaseExpected{
				Domain: &public.RegisterDomainResponse{
					Type: public.DomainTypeRhelIdm,
					RhelIdm: &public.DomainIpa{
						RealmName:    "",
						CaCerts:      []public.DomainIpaCert{},
						Servers:      []public.DomainIpaServer{},
						RealmDomains: []string{},
					},
				},
				Err: nil,
			},
		},
		{
			Name: "Success full rhel-idm content",
			Given: &model.Domain{
				Type: pointy.Uint(model.DomainTypeIpa),
				IpaDomain: &model.Ipa{
					RealmName:       pointy.String("MYDOMAIN.EXAMPLE"),
					Token:           pointy.String("71ad4978-c768-11ed-ad69-482ae3863d30"),
					TokenExpiration: tokenExpiration,
					RealmDomains:    pq.StringArray{"mydomain.example"},
					CaCerts: []model.IpaCert{
						{
							Nickname:     "MYDOMAIN.EXAMPLE IPA CA",
							Issuer:       "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
							Subject:      "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
							SerialNumber: "1",
							Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
						},
					},
					Servers: []model.IpaServer{
						{
							FQDN:                "server1.mydomain.example",
							RHSMId:              "c4a80438-c768-11ed-a60e-482ae3863d30",
							PKInitServer:        true,
							CaServer:            true,
							HCCEnrollmentServer: true,
							HCCUpdateServer:     true,
						},
					},
				},
			},
			Expected: TestCaseExpected{
				Domain: &public.RegisterDomainResponse{
					RhelIdm: &public.DomainIpa{
						RealmName:    "MYDOMAIN.EXAMPLE",
						RealmDomains: pq.StringArray{"mydomain.example"},
						CaCerts: []public.DomainIpaCert{
							{
								Nickname:     "MYDOMAIN.EXAMPLE IPA CA",
								Issuer:       "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
								Subject:      "CN=Certificate Authority,O=MYDOMAIN.EXAMPLE",
								SerialNumber: "1",
								Pem:          "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
							},
						},
						Servers: []public.DomainIpaServer{
							{
								Fqdn:                  "server1.mydomain.example",
								SubscriptionManagerId: "c4a80438-c768-11ed-a60e-482ae3863d30",
								PkinitServer:          true,
								CaServer:              true,
								HccEnrollmentServer:   true,
								HccUpdateServer:       true,
							},
						},
					},
				},
				Err: nil,
			},
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		p := &domainPresenter{cfg: &cfg}
		ipa, err := p.Register(testCase.Given)
		if testCase.Expected.Err != nil {
			assert.EqualError(t, err, testCase.Expected.Err.Error())
			assert.Nil(t, ipa)
		} else {
			assert.NoError(t, err)
			require.NotNil(t, ipa)
			assert.Equal(t, testCase.Expected.Domain.RhelIdm.RealmName, ipa.RhelIdm.RealmName)
			require.Equal(t, len(testCase.Expected.Domain.RhelIdm.RealmDomains), len(ipa.RhelIdm.RealmDomains))
			for i := range ipa.RhelIdm.RealmDomains {
				assert.Equal(t, testCase.Expected.Domain.RhelIdm.RealmDomains[i], ipa.RhelIdm.RealmDomains[i])
			}
			require.Equal(t, len(testCase.Expected.Domain.RhelIdm.CaCerts), len(ipa.RhelIdm.CaCerts))
			for i := range ipa.RhelIdm.CaCerts {
				assert.Equal(t, testCase.Expected.Domain.RhelIdm.CaCerts[i], ipa.RhelIdm.CaCerts[i])
			}
			require.Equal(t, len(testCase.Expected.Domain.RhelIdm.Servers), len(ipa.RhelIdm.Servers))
			for i := range ipa.RhelIdm.Servers {
				assert.Equal(t, testCase.Expected.Domain.RhelIdm.Servers[i], ipa.RhelIdm.Servers[i])
			}
		}
	}
}

func TestSharedDomainFill(t *testing.T) {
	p := &domainPresenter{cfg: &cfg}
	assert.Panics(t, func() {
		p.sharedDomainFill(nil, nil)
	})

	domain := &model.Domain{}
	assert.Panics(t, func() {
		p.sharedDomainFill(domain, nil)
	})

	output := public.RegisterDomainResponse{}
	domainUUID := "6d9575f2-de94-11ed-af6e-482ae3863d30"
	domain.DomainUuid = uuid.MustParse(domainUUID)
	domain.AutoEnrollmentEnabled = pointy.Bool(true)
	domain.DomainName = pointy.String("mydomain.example")
	domain.Title = pointy.String("My Domain Example")
	domain.Description = pointy.String("My Domain Example Description")
	p.sharedDomainFill(domain, &output)
	assert.Equal(t, domainUUID, output.DomainUuid)
	assert.Equal(t, true, output.AutoEnrollmentEnabled)
	assert.Equal(t, "mydomain.example", output.DomainName)
	assert.Equal(t, "My Domain Example", output.Title)
	assert.Equal(t, "My Domain Example Description", output.Description)
}

func TestFillRhelIdmCerts(t *testing.T) {
	p := &domainPresenter{cfg: &cfg}

	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(nil, nil)
	})

	output := public.Domain{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(&output, nil)
	})

	domain := model.Domain{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(&output, &domain)
	})

	output.RhelIdm = &public.DomainIpa{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmCerts(&output, &domain)
	})
}

func TestGuardSharedDomain(t *testing.T) {
	p := &domainPresenter{cfg: &cfg}

	err := p.guardSharedDomain(nil)
	assert.EqualError(t, err, "'domain' is nil")

	domain := model.Domain{}
	err = p.guardSharedDomain(&domain)
	assert.EqualError(t, err, "'domain.Type' is nil")

	domain.Type = pointy.Uint(model.DomainTypeUndefined)
	err = p.guardSharedDomain(&domain)
	assert.EqualError(t, err, "'domain.Type' is invalid")

	*domain.Type = model.DomainTypeIpa
	err = p.guardSharedDomain(&domain)
	assert.NoError(t, err)
}

func TestSharedDomain(t *testing.T) {
	var (
		err    error
		output *public.RegisterDomainResponse
	)
	p := &domainPresenter{cfg: &cfg}

	// Fail some guard check
	output, err = p.sharedDomain(nil)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'domain' is nil")

	// Fail Type not filled
	domain := &model.Domain{}
	output, err = p.sharedDomain(domain)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'domain.Type' is nil")

	// Fail nil IpaDomain
	domain.Type = pointy.Uint(model.DomainTypeIpa)
	output, err = p.sharedDomain(domain)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'domain.IpaDomain' is nil")

	// Not valid Type
	*domain.Type = 999
	output, err = p.sharedDomain(domain)
	assert.Nil(t, output)
	assert.EqualError(t, err, "'domain.Type=999' is invalid")

	// Success minimal values
	*domain.Type = model.DomainTypeIpa
	domain.IpaDomain = &model.Ipa{}
	output, err = p.sharedDomain(domain)
	expected := public.Domain{
		AutoEnrollmentEnabled: false,
		Title:                 "",
		Description:           "",
		DomainName:            "",
		DomainUuid:            model.NilUUID.String(),
		Type:                  public.DomainTypeRhelIdm,
		RhelIdm: &public.DomainIpa{
			RealmName:    "",
			CaCerts:      []public.DomainIpaCert{},
			Servers:      []public.DomainIpaServer{},
			RealmDomains: []string{},
		},
	}
	assert.NoError(t, err)
	require.NotNil(t, output)
	equalPresenterDomain(t, &expected, output)

	// Success with full information
	*domain.Type = model.DomainTypeIpa
	domain.Title = pointy.String("Test Title")
	domain.Description = pointy.String("Test Description")
	domain.DomainName = pointy.String("mydomain.example")
	testUUID := uuid.New()
	testOrgID := "12345"
	domain.DomainUuid = testUUID
	domain.OrgId = testOrgID
	domain.AutoEnrollmentEnabled = pointy.Bool(true)
	testNotValidBefore := time.Now()
	testNotValidAfter := testNotValidBefore.Add(24 * time.Hour)
	domain.IpaDomain.RealmDomains = pq.StringArray{"mydomain.example"}
	domain.IpaDomain.RealmName = pointy.String("MYDOMAIN.EXAMPLE")
	testToken := uuid.New()
	domain.IpaDomain.Token = pointy.String(testToken.String())
	testTokenExpiration := time.Now().Add(15 * time.Minute)
	domain.IpaDomain.TokenExpiration = &testTokenExpiration
	domain.IpaDomain.CaCerts = []model.IpaCert{
		{
			Issuer:         "Ca Cert Issuer test",
			Nickname:       "Ca Cert Nickname test",
			NotValidBefore: testNotValidBefore,
			NotValidAfter:  testNotValidAfter,
			SerialNumber:   "1",
			Subject:        "Ca Cert Subject",
			Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
		},
	}
	testSubscriptionManagerId := "93a46bde-e760-11ed-9a5a-482ae3863d30"
	domain.IpaDomain.Servers = []model.IpaServer{
		{
			FQDN:                "server1.mydomain.example",
			RHSMId:              testSubscriptionManagerId,
			CaServer:            true,
			PKInitServer:        true,
			HCCEnrollmentServer: true,
			HCCUpdateServer:     true,
		},
	}
	expected = public.Domain{
		AutoEnrollmentEnabled: true,
		Title:                 "Test Title",
		Description:           "Test Description",
		DomainName:            "mydomain.example",
		DomainUuid:            testUUID.String(),
		Type:                  public.DomainTypeRhelIdm,
		RhelIdm: &public.DomainIpa{
			RealmName:    "MYDOMAIN.EXAMPLE",
			RealmDomains: []string{"mydomain.example"},
			CaCerts: []public.DomainIpaCert{
				{
					Issuer:         "Ca Cert Issuer test",
					Nickname:       "Ca Cert Nickname test",
					NotValidBefore: testNotValidBefore,
					NotValidAfter:  testNotValidAfter,
					SerialNumber:   "1",
					Subject:        "Ca Cert Subject",
					Pem:            "-----BEGIN CERTIFICATE-----\nMII...\n-----END CERTIFICATE-----\n",
				},
			},
			Servers: []public.DomainIpaServer{
				{
					Fqdn:                  "server1.mydomain.example",
					SubscriptionManagerId: testSubscriptionManagerId,
					CaServer:              true,
					PkinitServer:          true,
					HccEnrollmentServer:   true,
					HccUpdateServer:       true,
				},
			},
		},
	}
	output, err = p.sharedDomain(domain)
	require.NotNil(t, output)
	assert.NoError(t, err)
	equalPresenterDomain(t, output, &expected)
}

func equalPrensenterDomainRhelIdmCaCerts(t *testing.T, expected []public.DomainIpaCert, actual []public.DomainIpaCert) {
	if expected == nil && actual == nil {
		return
	}
	if expected == nil {
		t.Errorf("'expected' is nil")
		return
	}
	if actual == nil {
		t.Errorf("'actual' is nil")
		return
	}

	require.Equal(t, len(expected), len(actual))
	for i := range actual {
		assert.Equal(t, expected[i].Nickname, actual[i].Nickname)
		assert.Equal(t, expected[i].Issuer, actual[i].Issuer)
		assert.Equal(t, expected[i].SerialNumber, actual[i].SerialNumber)
		assert.Equal(t, expected[i].Subject, actual[i].Subject)
		assert.Equal(t, expected[i].NotValidAfter, actual[i].NotValidAfter)
		assert.Equal(t, expected[i].NotValidBefore, actual[i].NotValidBefore)
		assert.Equal(t, expected[i].Pem, actual[i].Pem)
	}
}

func equalPrensenterDomainRhelIdmServers(t *testing.T, expected []public.DomainIpaServer, actual []public.DomainIpaServer) {
	if expected == nil && actual == nil {
		return
	}
	require.NotNil(t, expected)
	require.NotNil(t, actual)

	require.Equal(t, len(expected), len(actual))
	for i := range actual {
		assert.Equal(t, expected[i].Fqdn, actual[i].Fqdn)
		assert.Equal(t, expected[i].SubscriptionManagerId, actual[i].SubscriptionManagerId)
		assert.Equal(t, expected[i].CaServer, actual[i].CaServer)
		assert.Equal(t, expected[i].PkinitServer, actual[i].PkinitServer)
		assert.Equal(t, expected[i].HccEnrollmentServer, actual[i].HccEnrollmentServer)
		assert.Equal(t, expected[i].HccUpdateServer, actual[i].HccUpdateServer)
	}
}

// equalPresenterDomainRhelIdm compare expected public.DomainIpa with actual model.Ipa
func equalPresenterDomainRhelIdm(t *testing.T, expected *public.DomainIpa, actual *public.DomainIpa) {
	if expected == nil && actual == nil {
		return
	}
	require.NotNil(t, expected)
	require.NotNil(t, actual)
	assert.Equal(t, expected.RealmName, actual.RealmName)
	assert.Equal(t, expected.RealmDomains, actual.RealmDomains)
	equalPrensenterDomainRhelIdmCaCerts(t, expected.CaCerts, actual.CaCerts)
	equalPrensenterDomainRhelIdmServers(t, expected.Servers, actual.Servers)
}

// equalPresenterDomain compare expected public.Domain with actual model.Domain
func equalPresenterDomain(t *testing.T, expected *public.Domain, actual *public.Domain) {
	if expected == nil && actual == nil {
		return
	}
	require.NotNil(t, expected)
	require.NotNil(t, actual)
	assert.Equal(t, expected.AutoEnrollmentEnabled, actual.AutoEnrollmentEnabled)
	assert.Equal(t, expected.Title, actual.Title)
	assert.Equal(t, expected.Description, actual.Description)
	assert.Equal(t, expected.DomainName, actual.DomainName)
	assert.Equal(t, expected.DomainUuid, actual.DomainUuid)
	assert.Equal(t, expected.Type, actual.Type)
	switch expected.Type {
	case public.DomainTypeRhelIdm:
		equalPresenterDomainRhelIdm(t, expected.RhelIdm, actual.RhelIdm)
	case "":
	default:
		t.Errorf("asserting agains an invalid Type='%s'", expected.Type)
	}
}

func TestFillRhelIdmServersError(t *testing.T) {
	p := &domainPresenter{cfg: &cfg}

	assert.NotPanics(t, func() {
		p.fillRhelIdmServers(nil, nil)
	})

	output := public.Domain{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmServers(&output, nil)
	})

	domain := model.Domain{}
	assert.NotPanics(t, func() {
		p.fillRhelIdmServers(&output, &domain)
	})
}

func TestFillRhelIdmServers(t *testing.T) {
	type TestCaseGiven struct {
		To   *public.Domain
		From *model.Domain
	}
	type TestCaseExpected struct {
		Err error
		To  *public.Domain
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "Full success copy",
			Given: TestCaseGiven{
				To: &public.Domain{
					RhelIdm: &public.DomainIpa{},
				},
				From: &model.Domain{
					Type: pointy.Uint(model.DomainTypeIpa),
					IpaDomain: &model.Ipa{
						Servers: []model.IpaServer{
							{
								FQDN:                "server1.mydomain.example",
								RHSMId:              "547ce70c-9eb5-4783-a619-086aa26f88e5",
								CaServer:            true,
								HCCEnrollmentServer: true,
								HCCUpdateServer:     true,
								PKInitServer:        true,
							},
						},
					},
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				To: &public.Domain{
					Type: public.DomainTypeRhelIdm,
					RhelIdm: &public.DomainIpa{
						Servers: []public.DomainIpaServer{
							{
								Fqdn:                  "server1.mydomain.example",
								SubscriptionManagerId: "547ce70c-9eb5-4783-a619-086aa26f88e5",
								CaServer:              true,
								HccEnrollmentServer:   true,
								HccUpdateServer:       true,
								PkinitServer:          true,
							},
						},
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		// I instantiate directly because the public methods
		// are not part of the interface
		p := &domainPresenter{cfg: &cfg}
		if testCase.Expected.Err != nil {
			// assert.EqualError(t, err, testCase.Expected.Err.Error())
			assert.Panics(t, func() {
				p.fillRhelIdmServers(testCase.Given.To, testCase.Given.From)
			})
		} else {
			// assert.NoError(t, err)
			assert.NotPanics(t, func() {
				p.fillRhelIdmServers(testCase.Given.To, testCase.Given.From)
			})
			require.NotNil(t, testCase.Expected.To)
			require.NotNil(t, testCase.Expected.To.RhelIdm)
			require.NotNil(t, testCase.Expected.To.RhelIdm.Servers)
			require.NotNil(t, testCase.Given.To)
			require.NotNil(t, testCase.Given.To.RhelIdm)
			require.NotNil(t, testCase.Given.To.RhelIdm.Servers)
			assert.Equal(t, testCase.Expected.To.RhelIdm.Servers, testCase.Given.To.RhelIdm.Servers)
		}
	}
}

func TestBuildPaginationLink(t *testing.T) {
	const prefix = "/api/hmsidm/v1"
	cfg := config.Config{
		Application: config.Application{
			PaginationDefaultLimit: 10,
			PaginationMaxLimit:     100,
		},
	}
	p := &domainPresenter{cfg: &cfg}

	offset := 0
	limit := 10
	output := p.buildPaginationLink(prefix, offset, limit)
	assert.Equal(t, prefix+"/domains?limit=10&offset=0", output)

	offset = -1
	limit = 10
	output = p.buildPaginationLink(prefix, offset, limit)
	assert.Equal(t, prefix+"/domains?limit=10&offset=0", output)

	offset = 0
	limit = 0
	output = p.buildPaginationLink(prefix, offset, limit)
	assert.Equal(t, prefix+"/domains?limit=10&offset=0", output)

	offset = 0
	limit = p.cfg.Application.PaginationMaxLimit + 1
	output = p.buildPaginationLink(prefix, offset, limit)
	assert.Equal(t, fmt.Sprintf(prefix+"/domains?limit=%d&offset=0", p.cfg.Application.PaginationMaxLimit), output)
}

func TestListFillLinks(t *testing.T) {

	p := &domainPresenter{cfg: &cfg}

	// output nil
	assert.Panics(t, func() {
		p.listFillLinks(nil, "/prefix", 10, 0, 1)
	})

	// links at page 1
	output := public.ListDomainsResponse{}
	p.listFillLinks(&output, "/prefix", 10, 0, 1)
	require.NotNil(t, output.Links.First)
	assert.Equal(t, "/prefix/domains?limit=1&offset=0", *output.Links.First)
	require.NotNil(t, output.Links.Previous)
	assert.Equal(t, "/prefix/domains?limit=1&offset=0", *output.Links.Previous)
	require.NotNil(t, output.Links.Next)
	assert.Equal(t, "/prefix/domains?limit=1&offset=1", *output.Links.Next)
	require.NotNil(t, output.Links.Last)
	assert.Equal(t, "/prefix/domains?limit=1&offset=9", *output.Links.Last)

	// links at page 2
	output = public.ListDomainsResponse{}
	p.listFillLinks(&output, "/prefix", 10, 1, 1)
	require.NotNil(t, output.Links.First)
	assert.Equal(t, "/prefix/domains?limit=1&offset=0", *output.Links.First)
	require.NotNil(t, output.Links.Previous)
	assert.Equal(t, "/prefix/domains?limit=1&offset=0", *output.Links.Previous)
	require.NotNil(t, output.Links.Next)
	assert.Equal(t, "/prefix/domains?limit=1&offset=2", *output.Links.Next)
	require.NotNil(t, output.Links.Last)
	assert.Equal(t, "/prefix/domains?limit=1&offset=9", *output.Links.Last)

	// links at before last page
	output = public.ListDomainsResponse{}
	p.listFillLinks(&output, "/prefix", 10, 8, 1)
	require.NotNil(t, output.Links.First)
	assert.Equal(t, "/prefix/domains?limit=1&offset=0", *output.Links.First)
	require.NotNil(t, output.Links.Previous)
	assert.Equal(t, "/prefix/domains?limit=1&offset=7", *output.Links.Previous)
	require.NotNil(t, output.Links.Next)
	assert.Equal(t, "/prefix/domains?limit=1&offset=9", *output.Links.Next)
	require.NotNil(t, output.Links.Last)
	assert.Equal(t, "/prefix/domains?limit=1&offset=9", *output.Links.Last)

	// links at last page
	output = public.ListDomainsResponse{}
	p.listFillLinks(&output, "/prefix", 10, 9, 1)
	require.NotNil(t, output.Links.First)
	assert.Equal(t, "/prefix/domains?limit=1&offset=0", *output.Links.First)
	require.NotNil(t, output.Links.Previous)
	assert.Equal(t, "/prefix/domains?limit=1&offset=8", *output.Links.Previous)
	require.NotNil(t, output.Links.Next)
	assert.Equal(t, "/prefix/domains?limit=1&offset=9", *output.Links.Next)
	require.NotNil(t, output.Links.Last)
	assert.Equal(t, "/prefix/domains?limit=1&offset=9", *output.Links.Last)
}

func TestListFillMeta(t *testing.T) {
	p := &domainPresenter{cfg: &cfg}

	assert.Panics(t, func() {
		p.listFillMeta(nil, 10, 0, 1)
	})

	output := public.ListDomainsResponse{}
	p.listFillMeta(&output, 10, 0, 1)
	assert.Equal(t, int64(10), output.Meta.Count)
	assert.Equal(t, 0, output.Meta.Offset)
	assert.Equal(t, 1, output.Meta.Limit)
}

func TestListFillItem(t *testing.T) {
	p := &domainPresenter{cfg: &cfg}

	assert.Panics(t, func() {
		p.listFillItem(nil, nil)
	})

	output := public.ListDomainsData{}
	assert.Panics(t, func() {
		p.listFillItem(&output, nil)
	})

	// path with all the data
	domain := model.Domain{
		OrgId:                 "12345",
		DomainUuid:            uuid.MustParse("d89b6b9a-ecf4-11ed-9e6c-482ae3863d30"),
		DomainName:            nil,
		AutoEnrollmentEnabled: nil,
		Type:                  pointy.Uint(model.DomainTypeIpa),
	}
	p.listFillItem(&output, &domain)
	assert.Nil(t, domain.AutoEnrollmentEnabled)
	assert.Nil(t, domain.DomainName)
	require.NotNil(t, output.DomainType)
	assert.Equal(t, string(public.RhelIdm), string(output.DomainType))

	// path with all the data
	domain = model.Domain{
		OrgId:                 "12345",
		DomainUuid:            uuid.MustParse("d89b6b9a-ecf4-11ed-9e6c-482ae3863d30"),
		DomainName:            pointy.String("mydomain.example"),
		AutoEnrollmentEnabled: pointy.Bool(true),
		Type:                  pointy.Uint(model.DomainTypeIpa),
	}
	p.listFillItem(&output, &domain)
	require.NotNil(t, domain.AutoEnrollmentEnabled)
	assert.Equal(t, true, *domain.AutoEnrollmentEnabled)
	require.NotNil(t, domain.DomainName)
	assert.Equal(t, "mydomain.example", *domain.DomainName)
	require.NotNil(t, output.DomainType)
	assert.Equal(t, string(public.RhelIdm), string(output.DomainType))
}
