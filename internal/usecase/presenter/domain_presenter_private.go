package presenter

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/hmsidm/internal/api/public"
	"github.com/hmsidm/internal/domain/model"
	"github.com/openlyinc/pointy"
)

func (p *domainPresenter) fillRhelIdmServers(
	target *public.Domain,
	source *model.Domain,
) {
	if target == nil || source == nil {
		return
	}
	if target.RhelIdm == nil || source.IpaDomain == nil {
		return
	}
	target.RhelIdm.Servers = make(
		[]public.DomainIpaServer,
		len(source.IpaDomain.Servers),
	)
	for i := range source.IpaDomain.Servers {
		target.RhelIdm.Servers[i].Fqdn =
			source.IpaDomain.Servers[i].FQDN
		target.RhelIdm.Servers[i].SubscriptionManagerId =
			source.IpaDomain.Servers[i].RHSMId
		target.RhelIdm.Servers[i].Location =
			source.IpaDomain.Servers[i].Location
		target.RhelIdm.Servers[i].CaServer =
			source.IpaDomain.Servers[i].CaServer
		target.RhelIdm.Servers[i].HccEnrollmentServer =
			source.IpaDomain.Servers[i].HCCEnrollmentServer
		target.RhelIdm.Servers[i].HccUpdateServer =
			source.IpaDomain.Servers[i].HCCUpdateServer
		target.RhelIdm.Servers[i].PkinitServer =
			source.IpaDomain.Servers[i].PKInitServer
	}
}

func (p *domainPresenter) fillRhelIdmCerts(
	output *public.Domain,
	domain *model.Domain,
) {
	if output == nil || domain == nil || output.RhelIdm == nil || domain.IpaDomain == nil {
		return
	}
	output.RhelIdm.CaCerts = make(
		[]public.DomainIpaCert,
		len(domain.IpaDomain.CaCerts),
	)
	for i := range domain.IpaDomain.CaCerts {
		output.RhelIdm.CaCerts[i].Nickname =
			domain.IpaDomain.CaCerts[i].Nickname
		output.RhelIdm.CaCerts[i].Issuer =
			domain.IpaDomain.CaCerts[i].Issuer
		output.RhelIdm.CaCerts[i].NotValidAfter =
			domain.IpaDomain.CaCerts[i].NotValidAfter
		output.RhelIdm.CaCerts[i].NotValidBefore =
			domain.IpaDomain.CaCerts[i].NotValidBefore
		output.RhelIdm.CaCerts[i].SerialNumber =
			domain.IpaDomain.CaCerts[i].SerialNumber
		output.RhelIdm.CaCerts[i].Subject =
			domain.IpaDomain.CaCerts[i].Subject
		output.RhelIdm.CaCerts[i].Pem =
			domain.IpaDomain.CaCerts[i].Pem
	}
}

func (p *domainPresenter) guardSharedDomain(
	domain *model.Domain,
) error {
	if domain == nil {
		return fmt.Errorf("'domain' is nil")
	}
	if domain.Type == nil {
		return fmt.Errorf("'domain.Type' is nil")
	}
	if *domain.Type == model.DomainTypeUndefined {
		return fmt.Errorf("'domain.Type' is invalid")
	}
	return nil
}

func (p *domainPresenter) sharedDomain(
	domain *model.Domain,
) (output *public.Domain, err error) {
	// Expect domain not nil and domain.Type filled
	if err = p.guardSharedDomain(domain); err != nil {
		return nil, err
	}

	// Domain common code
	output = &public.Domain{}
	p.sharedDomainFill(domain, output)

	switch *domain.Type {
	case model.DomainTypeIpa:
		// Specific rhel-idm domain code
		output.Type = model.DomainTypeIpaString
		output.RhelIdm = &public.DomainIpa{}
		err = p.sharedDomainFillRhelIdm(domain, output)
	default:
		err = fmt.Errorf("'domain.Type=%d' is invalid", *domain.Type)
	}
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (p *domainPresenter) sharedDomainFill(
	domain *model.Domain,
	output *public.Domain,
) {
	output.DomainUuid = domain.DomainUuid.String()
	if domain.AutoEnrollmentEnabled != nil {
		output.AutoEnrollmentEnabled = *domain.AutoEnrollmentEnabled
	}
	if domain.DomainName != nil {
		output.DomainName = *domain.DomainName
	}
	if domain.Title != nil {
		output.Title = *domain.Title
	}
	if domain.Description != nil {
		output.Description = *domain.Description
	}
}

func (p *domainPresenter) sharedDomainFillRhelIdm(
	domain *model.Domain,
	output *public.Domain,
) (err error) {
	if domain.Type != nil && *domain.Type != model.DomainTypeIpa {
		return fmt.Errorf(
			"'domain.Type' is not '%s'",
			model.DomainTypeIpaString,
		)
	}
	if domain.IpaDomain == nil {
		return fmt.Errorf("'domain.IpaDomain' is nil")
	}
	if output == nil {
		panic("'output' is nil")
	}
	output.RhelIdm = &public.DomainIpa{}
	if domain.IpaDomain.RealmName != nil {
		output.RhelIdm.RealmName = *domain.IpaDomain.RealmName
	}

	if domain.IpaDomain.RealmDomains != nil {
		output.RhelIdm.RealmDomains = append(
			[]string{},
			domain.IpaDomain.RealmDomains...)
	} else {
		output.RhelIdm.RealmDomains = []string{}
	}

	p.fillRhelIdmCerts(output, domain)

	p.fillRhelIdmServers(output, domain)

	return nil
}

func (p *domainPresenter) buildPaginationLink(prefix string, offset int, limit int) string {
	if limit == 0 {
		limit = p.cfg.Application.PaginationDefaultLimit
	}
	if limit > p.cfg.Application.PaginationMaxLimit {
		limit = p.cfg.Application.PaginationMaxLimit
	}
	if offset < 0 {
		offset = 0
	}

	q := url.Values{}
	q.Add("limit", strconv.FormatInt(int64(limit), 10))
	q.Add("offset", strconv.FormatInt(int64(offset), 10))

	return fmt.Sprintf("%s/domains?%s", prefix, q.Encode())
}

func (p *domainPresenter) listFillLinks(output *public.ListDomainsResponse, prefix string, count int64, offset int, limit int) {
	if output == nil {
		panic("'output' is nil")
	}

	// Calculate the offsets
	currentOffset := ((offset + limit - 1) / limit) * limit
	firstOffset := 0
	prevOffset := currentOffset - limit
	nextOffset := currentOffset + limit
	lastOffset := ((int(count)+limit-1)/limit)*limit - limit
	if firstOffset > prevOffset {
		prevOffset = firstOffset
	}
	if nextOffset > lastOffset {
		nextOffset = lastOffset
	}

	// Build the link
	output.Links.First = pointy.String(p.buildPaginationLink(prefix, firstOffset, limit))
	output.Links.Previous = pointy.String(p.buildPaginationLink(prefix, prevOffset, limit))
	output.Links.Next = pointy.String(p.buildPaginationLink(prefix, nextOffset, limit))
	output.Links.Last = pointy.String(p.buildPaginationLink(prefix, lastOffset, limit))
}

func (p *domainPresenter) listFillMeta(output *public.ListDomainsResponse, count int64, offset int, limit int) {
	if output == nil {
		panic("'output' is nil")
	}
	output.Meta.Count = count
	output.Meta.Offset = offset
	output.Meta.Limit = limit
}

func (p *domainPresenter) listFillItem(output *public.ListDomainsData, domain *model.Domain) {
	if output == nil {
		panic("'output' is nil")
	}
	if domain == nil {
		panic("'domain' is nil")
	}
	if domain.AutoEnrollmentEnabled == nil {
		output.AutoEnrollmentEnabled = false
	} else {
		output.AutoEnrollmentEnabled = *domain.AutoEnrollmentEnabled
	}
	if domain.DomainName != nil {
		output.DomainName = *domain.DomainName
	}
	output.DomainType = public.ListDomainsDataDomainType(model.DomainTypeString(*domain.Type))
	output.DomainUuid = domain.DomainUuid.String()
}
