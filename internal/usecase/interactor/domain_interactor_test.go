package interactor

import (
	"fmt"
	"testing"

	api_public "github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/hmsidm/internal/interface/interactor"
	"github.com/openlyinc/pointy"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTodoInteractor(t *testing.T) {
	var component interactor.DomainInteractor
	assert.NotPanics(t, func() {
		component = NewDomainInteractor()
	})
	assert.NotNil(t, component)
}

func TestCreate(t *testing.T) {
	type TestCaseGiven struct {
		Params *api_public.CreateDomainParams
		In     *api_public.CreateDomain
	}
	type TestCaseExpected struct {
		Err   error
		OrgId string
		Out   *model.Domain
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}
	testCases := []TestCase{
		{
			Name: "nil for params",
			Given: TestCaseGiven{
				Params: nil,
				In:     &api_public.CreateDomain{},
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("'params' cannot be nil"),
				Out: nil,
			},
		},
		{
			Name: "nil for In",
			Given: TestCaseGiven{
				Params: &api_public.CreateDomainParams{},
				In:     nil,
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("'body' cannot be nil"),
				Out: nil,
			},
		},
		{
			Name: "nil for Out",
			Given: TestCaseGiven{
				Params: &api_public.CreateDomainParams{},
				In:     &api_public.CreateDomain{},
			},
			Expected: TestCaseExpected{
				Err: fmt.Errorf("X-Rh-Identity content cannot be an empty string"),
				Out: nil,
			},
		},
		{
			Name: "success case",
			Given: TestCaseGiven{
				Params: &api_public.CreateDomainParams{
					XRhIdentity: EncodeIdentity(
						&identity.Identity{
							OrgID: "12345",
							Internal: identity.Internal{
								OrgID: "12345",
							},
						},
					),
				},
				In: &api_public.CreateDomain{
					AutoEnrollmentEnabled: true,
					DomainName:            "domain.example",
					DomainType:            api_public.CreateDomainDomainTypeIpa,
					Ipa: api_public.CreateDomainIpa{
						CaList:    "",
						RealmName: pointy.String("DOMAIN.EXAMPLE"),
						ServerList: &[]string{
							"server1.domain.example",
							"server2.domain.example",
						},
					},
				},
			},
			Expected: TestCaseExpected{
				Err: nil,
				Out: &model.Domain{
					OrgId:                 "12345",
					DomainName:            pointy.String("domain.example"),
					DomainType:            pointy.Uint(model.DomainTypeIpa),
					AutoEnrollmentEnabled: pointy.Bool(true),
					IpaDomain: &model.Ipa{
						CaList:     pointy.String(""),
						RealmName:  pointy.String("DOMAIN.EXAMPLE"),
						ServerList: pointy.String("server1.domain.example,server2.domain.example"),
					},
				},
			},
		},
	}

	component := NewDomainInteractor()
	for _, testCase := range testCases {
		t.Log(testCase.Name)
		orgId, data, err := component.Create(testCase.Given.Params, testCase.Given.In)
		if testCase.Expected.Err != nil {
			require.Error(t, err)
			require.Equal(t, testCase.Expected.Err.Error(), err.Error())
			assert.Equal(t, "", orgId)
			assert.Nil(t, data)
		} else {
			assert.NoError(t, err)
			require.NotNil(t, testCase.Expected.Out)

			assert.Equal(t, *testCase.Expected.Out.AutoEnrollmentEnabled, *data.AutoEnrollmentEnabled)
			assert.Equal(t, testCase.Expected.Out.OrgId, data.OrgId)
			assert.Equal(t, *testCase.Expected.Out.DomainName, *data.DomainName)
			assert.Equal(t, *testCase.Expected.Out.DomainType, *data.DomainType)
			assert.Equal(t,
				*testCase.Expected.Out.IpaDomain.CaList,
				*data.IpaDomain.CaList)
			assert.Equal(t,
				*testCase.Expected.Out.IpaDomain.RealmName,
				*data.IpaDomain.RealmName)
			assert.Equal(t,
				*testCase.Expected.Out.IpaDomain.ServerList,
				*data.IpaDomain.ServerList)
		}
	}
}