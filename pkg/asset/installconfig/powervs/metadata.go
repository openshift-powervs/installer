package powervs

import (
	"context"
	"fmt"
	gohttp "net/http"
	"os"
	"sync"

	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	bxsession "github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=./metadata.go -destination=./mock/powervsmetadata_generated.go -package=mock

// MetadataAPI represents functions that eventually call out to the API
type MetadataAPI interface {
	AccountID(ctx context.Context) (string, error)
	APIKey(ctx context.Context, t APIEndpointType) (string, error)
	CISInstanceCRN(ctx context.Context) (string, error)
}

// Metadata holds additional metadata for InstallConfig resources that
// do not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	BaseDomain string

	accountID      string
	endpointURL    string
	apiKey         string
	cisInstanceCRN string
	client         *Client

	mutex sync.Mutex
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(baseDomain string) *Metadata {
	return &Metadata{BaseDomain: baseDomain}
}

// AccountID returns the IBM Cloud account ID associated with the authentication
// credentials.
func (m *Metadata) AccountID(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.client == nil {
		client, err := NewClient(PowerVSEP)
		if err != nil {
			return "", err
		}

		m.client = client
	}

	logrus.Debug("metadata AccountID")
	if m.accountID == "" {
		apiKeyDetails, err := m.client.GetAuthenticatorAPIKeyDetails(ctx)
		if err != nil {
			return "", err
		}

		m.accountID = *apiKeyDetails.AccountID
	}

	logrus.Debug("metadata AccountID exit")

	return m.accountID, nil
}

// APIKey returns the IBM Cloud account API Key associated with the authentication
// credentials.
func (m *Metadata) APIKey(ctx context.Context, t APIEndpointType) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.client == nil {
		client, err := NewClient(t)
		if err != nil {
			return "", err
		}

		m.client = client
	}

	if m.apiKey == "" {
		m.apiKey = m.client.GetAPIKey()
	}

	return m.apiKey, nil
}

// getCISInstanceCRN gets the CRN name for the specified DNS base domain.
func getCISInstanceCRN(APIKey string, BaseDomain string, iamEPURL string) (string, error) {

	logrus.Debugf("getCISInstanceCRN APIKey %s, BaseDomain %s, epURL, %s", APIKey, BaseDomain, iamEPURL)
	var (
		tokenProviderEndpoint         = "https://iam.cloud.ibm.com"
		CISInstanceCRN                string
		bxSession                     *bxsession.Session
		err                           error
		tokenRefresher                *authentication.IAMAuthRepository
		authenticator                 *core.IamAuthenticator
		controllerSvc                 *resourcecontrollerv2.ResourceControllerV2
		listInstanceOptions           *resourcecontrollerv2.ListResourceInstancesOptions
		listResourceInstancesResponse *resourcecontrollerv2.ResourceInstancesList
		instance                      resourcecontrollerv2.ResourceInstance
		zonesService                  *zonesv1.ZonesV1
		listZonesOptions              *zonesv1.ListZonesOptions
		listZonesResponse             *zonesv1.ListZonesResp
	)

	if iamEPURL != "" {
		tokenProviderEndpoint = iamEPURL
	}

	bxSession, err = bxsession.New(&bluemix.Config{
		BluemixAPIKey:         APIKey,
		TokenProviderEndpoint: &tokenProviderEndpoint,
		Debug:                 false,
	})
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: bxsession.New: %v", err)
	}
	tokenRefresher, err = authentication.NewIAMAuthRepository(bxSession.Config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: authentication.NewIAMAuthRepository: %v", err)
	}
	err = tokenRefresher.AuthenticateAPIKey(bxSession.Config.BluemixAPIKey)
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: tokenRefresher.AuthenticateAPIKey: %v", err)
	}
	authenticator = &core.IamAuthenticator{
		ApiKey: APIKey,
		URL:    tokenProviderEndpoint,
	}
	err = authenticator.Validate()
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: authenticator.Validate: %v", err)
	}
	// Instantiate the service with an API key based IAM authenticator
	controllerSvc, err = resourcecontrollerv2.NewResourceControllerV2(&resourcecontrollerv2.ResourceControllerV2Options{
		Authenticator: authenticator,
		ServiceName:   "cloud-object-storage",
		// @TODO: Obvs un-hardcode this. do we even need to set it? it should do the default right thing
		//URL: "https://resource-controller.test.cloud.ibm.com",
	})
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: creating ControllerV2 Service: %v", err)
	}
	listInstanceOptions = controllerSvc.NewListResourceInstancesOptions()
	listInstanceOptions.SetResourceID(cisServiceID)
	listResourceInstancesResponse, _, err = controllerSvc.ListResourceInstances(listInstanceOptions)
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: ListResourceInstances: %v", err)
	}
	for _, instance = range listResourceInstancesResponse.Resources {
		authenticator = &core.IamAuthenticator{
			ApiKey: APIKey,
		}

		err = authenticator.Validate()
		if err != nil {
		}

		zonesService, err = zonesv1.NewZonesV1(&zonesv1.ZonesV1Options{
			Authenticator: authenticator,
			Crn:           instance.CRN,
		})
		if err != nil {
			return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: NewZonesV1: %v", err)
		}
		listZonesOptions = zonesService.NewListZonesOptions()
		listZonesResponse, _, err = zonesService.ListZones(listZonesOptions)
		if listZonesResponse == nil {
			return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: ListZones: %v", err)
		}
		for _, zone := range listZonesResponse.Result {
			if *zone.Status == "active" {
				if *zone.Name == BaseDomain {
					CISInstanceCRN = *instance.CRN
				}
			}
		}
	}

	return CISInstanceCRN, nil
}

// CISInstanceCRN returns the Cloud Internet Services instance CRN that is
// managing the DNS zone for the base domain.
func (m *Metadata) CISInstanceCRN(ctx context.Context) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.cisInstanceCRN != "" {
		return m.cisInstanceCRN, nil
	}

	// b/c ibmcloud and powervs are in different dc's, testing in staging is
	// split and messy. we can't re-use the same endpoint URLs.
	split := os.Getenv("IBM_POWERVS_DEV")
	if m.client != nil && split == "TRUE" && m.client.EPType != IBMCloudEP {
		if tempClient, err := NewClient(IBMCloudEP); err != nil {
			return "", errors.Wrap(err, "CISInstanceCRN: error cretaing new client")
		} else {
			cisInstanceCRN, err := getCISInstanceCRN(tempClient.GetAPIKey(), m.BaseDomain, tempClient.IAMEP)
			if err != nil {
				return "", err
			}
			m.cisInstanceCRN = cisInstanceCRN
			return m.cisInstanceCRN, nil
		}
	}
	if m.client == nil {
		client, err := NewClient(IBMCloudEP)
		if err != nil {
			return "", err
		}
		m.client = client
		logrus.Debugf("Metadata CISInstanceCRN, created new client with IBMCloudEP: %v", client)

	} else {
		logrus.Debugf("Metadata CISInstanceCRN, existing client being used: %v", m.client)
	}

	if m.apiKey == "" {
		m.apiKey = m.client.GetAPIKey()
	}

	if m.endpointURL == "" {
		m.endpointURL = m.client.IAMEP
	}

	cisInstanceCRN, err := getCISInstanceCRN(m.apiKey, m.BaseDomain, m.endpointURL)
	if err != nil {
		return "", err
	}

	m.cisInstanceCRN = cisInstanceCRN

	return m.cisInstanceCRN, nil
}

// SetCISInstanceCRN sets Cloud Internet Services instance CRN to a string value.
func (m *Metadata) SetCISInstanceCRN(crn string) {
	m.cisInstanceCRN = crn
}
