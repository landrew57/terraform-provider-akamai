package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiRatePolicies_data_basic(t *testing.T) {
	t.Run("match by RatePolicies ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getRatePoliciesResponse := appsec.GetRatePoliciesResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSRatePolicies/RatePolicies.json"), &getRatePoliciesResponse)
		require.NoError(t, err)

		client.On("GetRatePolicies",
			mock.Anything,
			appsec.GetRatePoliciesRequest{ConfigID: 43253, ConfigVersion: 7},
		).Return(&getRatePoliciesResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSRatePolicies/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_rate_policies.test", "id", "43253"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
