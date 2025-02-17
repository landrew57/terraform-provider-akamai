// Package dns contains implementation for Akamai Terraform sub-provider responsible for managing DNS zones configuration
package dns

import (
	"errors"
	"fmt"
	"sync"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/dns"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v3/pkg/tools"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider

		client dns.DNS
	}

	// Option is a dns provider option
	Option func(p *provider)
)

var (
	once sync.Once

	inst *provider
)

// Subprovider returns a core sub provider
func Subprovider() akamai.Subprovider {
	once.Do(func() {
		inst = &provider{Provider: Provider()}
	})

	return inst
}

// Provider returns the Akamai terraform.Resource provider.
func Provider() *schema.Provider {

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"dns_section": {
				Optional:   true,
				Type:       schema.TypeString,
				Default:    "default",
				Deprecated: akamai.NoticeDeprecatedUseAlias("dns_section"),
			},
			"dns": {
				Optional:   true,
				Type:       schema.TypeSet,
				Elem:       config.Options("dns"),
				MaxItems:   1,
				Deprecated: akamai.NoticeDeprecatedUseAlias("dns"),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akamai_authorities_set": dataSourceAuthoritiesSet(),
			"akamai_dns_record_set":  dataSourceDNSRecordSet(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"akamai_dns_zone":   resourceDNSv2Zone(),
			"akamai_dns_record": resourceDNSv2Record(),
		},
	}
	return provider
}

// WithClient sets the client interface function, used for mocking and testing
func WithClient(c dns.DNS) Option {
	return func(p *provider) {
		p.client = c
	}
}

// Client returns the DNS interface
func (p *provider) Client(meta akamai.OperationMeta) dns.DNS {
	if p.client != nil {
		return p.client
	}
	return dns.Client(meta.Session())
}

func getConfigDNSV2Service(d *schema.ResourceData) error {
	var inlineConfig *schema.Set
	for _, key := range []string{"dns", "config"} {
		opt, err := tools.GetSetValue(key, d)
		if err != nil {
			if !errors.Is(err, tools.ErrNotFound) {
				return err
			}
			continue
		}
		if inlineConfig != nil {
			return fmt.Errorf("only one inline config section can be defined")
		}
		inlineConfig = opt
	}
	if err := d.Set("config", inlineConfig); err != nil {
		return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
	}

	for _, s := range tools.FindStringValues(d, "dns_section", "config_section") {
		if s != "default" && s != "" {
			if err := d.Set("config_section", s); err != nil {
				return fmt.Errorf("%w: %s", tools.ErrValueSet, err.Error())
			}
			break
		}
	}

	return nil
}

func (p *provider) Name() string {
	return "dns"
}

// DNSProviderVersion update version string anytime provider adds new features
const DNSProviderVersion string = "v0.8.3"

func (p *provider) Version() string {
	return DNSProviderVersion
}

func (p *provider) Schema() map[string]*schema.Schema {
	return p.Provider.Schema
}

func (p *provider) Resources() map[string]*schema.Resource {
	return p.Provider.ResourcesMap
}

func (p *provider) DataSources() map[string]*schema.Resource {
	return p.Provider.DataSourcesMap
}

func (p *provider) Configure(log log.Interface, d *schema.ResourceData) diag.Diagnostics {
	log.Debug("START Configure")

	if err := getConfigDNSV2Service(d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
