package appsec

import (
	"encoding/json"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/appsec"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAkamaiAdvancedSettingsAttackPayloadLoggingDataBasic(t *testing.T) {
	t.Run("match by AdvancedSettingsAttackPayloadLogging ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getLoggingResponse := appsec.GetAdvancedSettingsAttackPayloadLoggingResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLogging.json"), &getLoggingResponse)
		require.NoError(t, err)

		client.On("GetAdvancedSettingsAttackPayloadLogging",
			mock.Anything,
			appsec.GetAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: 43253, Version: 7},
		).Return(&getLoggingResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSAdvancedSettingsAttackPayloadLogging/match_by_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_attack_payload_logging.test", "id", "43253:"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}

func TestAkamaiAdvancedSettingsAttackPayloadLoggingDataBasicPolicyId(t *testing.T) {
	t.Run("match by AdvancedSettingsAttackPayloadLoggingPolicy ID", func(t *testing.T) {
		client := &appsec.Mock{}

		config := appsec.GetConfigurationResponse{}
		err := json.Unmarshal(loadFixtureBytes("testdata/TestResConfiguration/LatestConfiguration.json"), &config)
		require.NoError(t, err)

		client.On("GetConfiguration",
			mock.Anything,
			appsec.GetConfigurationRequest{ConfigID: 43253},
		).Return(&config, nil)

		getLoggingResponse := appsec.GetAdvancedSettingsAttackPayloadLoggingResponse{}
		err = json.Unmarshal(loadFixtureBytes("testdata/TestDSAdvancedSettingsAttackPayloadLogging/AdvancedSettingsAttackPayloadLogging.json"), &getLoggingResponse)
		require.NoError(t, err)

		client.On("GetAdvancedSettingsAttackPayloadLogging",
			mock.Anything,
			appsec.GetAdvancedSettingsAttackPayloadLoggingRequest{ConfigID: 43253, Version: 7, PolicyID: "test_policy"},
		).Return(&getLoggingResponse, nil)

		useClient(client, func() {
			resource.Test(t, resource.TestCase{
				IsUnitTest:        true,
				ProviderFactories: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: loadFixtureString("testdata/TestDSAdvancedSettingsAttackPayloadLogging/match_by_policy_id.tf"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.akamai_appsec_advanced_settings_attack_payload_logging.policy", "id", "43253:test_policy"),
						),
					},
				},
			})
		})

		client.AssertExpectations(t)
	})

}
