package dns

import (
	"errors"
	"fmt"
	"sync"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/config"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/apex/log"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	provider struct {
		*schema.Provider
	}
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
				Optional: true,
				Type:     schema.TypeSet,
				Elem:     config.Options("dns"),
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

func getConfigDNSV2Service(d tools.ResourceDataFetcher) (*edgegrid.Config, error) {
	var DNSv2Config edgegrid.Config
	var err error
	dns, err := tools.GetSetValue("dns", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		dnsConfig := dns.List()
		if len(dnsConfig) == 0 {
			return nil, fmt.Errorf("'dns' property in provider must have at least one entry")
		}
		configMap, ok := dnsConfig[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("dns config entry is of invalid type; should be 'map[string]interface{}'")
		}
		host, ok := configMap["host"].(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "host", "string")
		}
		accessToken, ok := configMap["access_token"].(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "access_token", "string")
		}
		clientToken, ok := configMap["client_token"].(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "client_token", "string")
		}
		clientSecret, ok := configMap["client_secret"].(string)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "client_secret", "string")
		}
		maxBody, ok := configMap["max_body"].(int)
		if !ok {
			return nil, fmt.Errorf("%w: %s, %q", tools.ErrInvalidType, "max_body", "int")
		}
		DNSv2Config = edgegrid.Config{
			Host:         host,
			AccessToken:  accessToken,
			ClientToken:  clientToken,
			ClientSecret: clientSecret,
			MaxBody:      maxBody,
		}
		dnsv2.Init(DNSv2Config)
		return &DNSv2Config, nil
	}

	edgerc, err := tools.GetStringValue("edgerc", d)
	if err != nil && !errors.Is(err, tools.ErrNotFound) {
		return nil, err
	}

	var section string

	for _, s := range tools.FindStringValues(d, "dns_section", "config_section") {
		if s != "default" {
			section = s
			break
		}
	}

	DNSv2Config, err = edgegrid.Init(edgerc, section)
	if err != nil {
		return nil, err
	}

	dnsv2.Init(DNSv2Config)
	edgegrid.SetupLogging()

	return &DNSv2Config, nil
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

	_, err := getConfigDNSV2Service(d)
	if err != nil {
		return nil
	}
	return nil
}
