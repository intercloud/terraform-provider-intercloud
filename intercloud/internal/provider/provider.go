package provider

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/spf13/viper"

	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/api"
	"github.com/intercloud/terraform-provider-intercloud/intercloud/internal/config"
	homedir "github.com/mitchellh/go-homedir"
)

// This is a global MutexKV for use within this plugin.
// var intercloudMutexKV = mutexkv.NewMutexKV()

// intercloudProvider returns a terraform.ResourceProvider.
func intercloudProvider() terraform.ResourceProvider {

	// avoids this showing up in test output
	var debugLog = func(f string, v ...interface{}) {
		if os.Getenv("TF_LOG") == "" {
			return
		}

		if os.Getenv("TF_ACC") != "" {
			return
		}

		log.Printf(f, v...)
	}

	dataSources := make(map[string]*schema.Resource)
	resources := make(map[string]*schema.Resource)
	for _, service := range SupportedServices() {
		debugLog("[DEBUG] Registering Data Sources for %q..", service.Name())
		for k, v := range service.SupportedDataSources() {
			if existing := dataSources[k]; existing != nil {
				panic(fmt.Sprintf("An existing Data Source exists for %q", k))
			}

			dataSources[k] = v
		}

		debugLog("[DEBUG] Registering Resources for %q..", service.Name())
		for k, v := range service.SupportedResources() {
			if existing := resources[k]; existing != nil {
				panic(fmt.Sprintf("An existing Resource exists for %q", k))
			}

			resources[k] = v
		}
	}

	icProvider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("INTERCLOUD_ACCESS_TOKEN", nil),
				Description:  descriptions["access_token"],
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Sensitive:    true,
			},
			"shared_credentials_file": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  descriptions["shared_credentials_file"],
				ValidateFunc: validateSharedCredentialsFile,
			},
			"api_endpoint": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				DefaultFunc:  schema.EnvDefaultFunc("INTERCLOUD_API_ENDPOINT", "https://api-console.intercloud.io"),
				Description:  descriptions["api_endpoint"],
				ValidateFunc: validation.IsURLWithHTTPS,
			},
			"organization_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  descriptions["organization_id"],
				ValidateFunc: validation.IsUUID,
			},
		},
		DataSourcesMap: dataSources,
		ResourcesMap:   resources,
	}

	icProvider.ConfigureFunc = providerConfigure(icProvider)

	return icProvider
}

func validateSharedCredentialsFile(v interface{}, k string) (warnings []string, errors []error) {
	// no check on empty value
	if v == nil || v.(string) == "" {
		return
	}
	// file existence
	credsFile := v.(string)
	credsPath, err := homedir.Expand(credsFile)
	if err != nil {
		errors = append(errors,
			fmt.Errorf("Error expanding credentials file path %q", credsFile))
		return
	}
	if _, err := os.Stat(credsPath); err != nil {
		errors = append(errors,
			fmt.Errorf("Error on credentials file %q: %s", credsFile, err))
		return
	}
	return
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	terraformVersion := p.TerraformVersion
	if terraformVersion == "" {
		// Terraform 0.12 introduced this field to the protocol
		// We can therefore assume that if it's missing it's 0.10 or 0.11
		terraformVersion = "0.11+compatible"
	}

	return func(d *schema.ResourceData) (r interface{}, err error) {
		// retrieve access token
		accessToken, err := getProviderAccessToken(d)
		if err != nil {
			return nil, err
		}

		// configure provider
		cfg := config.ProviderConfigInput{
			ApiEndpoint:        d.Get("api_endpoint").(string),
			PrivateAccessToken: accessToken,
			TerraformVersion:   terraformVersion,
			OrganizationID:     uuid.MustParse(d.Get("organization_id").(string)),
		}

		r = config.NewProviderConfig(cfg)

		log.Printf("[DEBUG] configuration: %+v", r)

		// check access token validity
		if _, err := isAccessTokenValid(r.(*config.ProviderConfig)); err != nil {
			return nil, err
		}

		return r, nil
	}
}

// getProviderAccessToken returns the access token respecting the following
// order :
// 1. Static credentials
// 2. Environment variable
// 3. Shared credentials file
func getProviderAccessToken(d *schema.ResourceData) (access_token string, err error) {
	// static credentials / environment variables
	if d.Get("access_token") != "" {
		log.Print("[DEBUG] access token retrieved from static configuration or environment variable")
		return d.Get("access_token").(string), nil
	}
	log.Printf("[DEBUG] trying to get access_token from shared credentials file")
	// shared credentials file
	credsFile := "~/.intercloud/credentials.json"
	if userFile := d.Get("shared_credentials_file"); userFile != nil && d.Get("shared_credentials_file").(string) != "" {
		credsFile = userFile.(string)
	}
	if err := configureViper(credsFile); err != nil {
		return "", errors.New("intercloud: could not read shared credentials file. See the provider configuration documentation more information.")
	}
	if access_token := viper.GetString("access_token"); access_token == "" {
		return "", errors.New("intercloud: could not find default credentials. See the provider configuration documentation more information.")
	}
	log.Printf("[DEBUG] access_token retrieved from shared credentials file: %s", credsFile)
	return viper.GetString("access_token"), nil
}

// configureViper configures viper to read config from specific credentials file
func configureViper(credsFile string) error {
	credsPath, err := homedir.Expand(credsFile)
	if err != nil {
		return err
	}
	viper.SetConfigFile(credsPath)
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// isAccessTokenValid indicates if the provided access token is valid on a configuration
func isAccessTokenValid(config *config.ProviderConfig) (bool, error) {
	_, err := config.ApiClient().ReadAccountInformations()
	if err != nil {
		uerr := errors.Unwrap(err)
		switch uerr {
		case api.ErrUnauthorized, api.ErrForbidden:
			return false, errors.New("intercloud: provided access token is not valid. Please provide a valid access token.")
		}
		return false, errors.New("intercloud: access token cannot be verified.")
	}
	return true, nil
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_endpoint":            "@TODO Desc. api_endpoint",
		"access_token":            "@TODO Desc. access_token",
		"shared_credentials_file": "@TODO Desc. shared_credentials_file",
		"organization_id":         "@TODO Desc. organization_id",
	}
}

// IntercloudProvider returns Intercloud resource Provider
func IntercloudProvider() terraform.ResourceProvider {
	return intercloudProvider()
}
