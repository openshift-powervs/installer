package powervs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	gohttp "net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/pkg/errors"

	survey "github.com/AlecAivazis/survey/v2"
	bx "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/account/accountv2"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	bxep "github.com/IBM-Cloud/bluemix-go/endpoints"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	bxsession "github.com/IBM-Cloud/bluemix-go/session"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/form3tech-oss/jwt-go"

	"github.com/sirupsen/logrus"
)

var (
	defSessionTimeout    time.Duration = 9000000000000000000.0
	defaultAuthDir                     = filepath.Join(os.Getenv("HOME"), ".powervs")
	defaultAuthFilePath                = filepath.Join(defaultAuthDir, "config.json")
	defaultPVSEPFilePath               = filepath.Join(defaultAuthDir, "power_vs_endpoints.json")
	defaultICEPFilePath                = filepath.Join(defaultAuthDir, "ic_endpoints.json")
)

type APIEndpointType string

const (
	PowerVSEP  APIEndpointType = "powervs"
	IBMCloudEP APIEndpointType = "ibmcloud"
)

//BxClient is struct which provides bluemix session details
type BxClient struct {
	bxConfig     *bx.Config
	PISession    *ibmpisession.IBMPISession
	User         *User // replace with config's IBMID
	AccountAPIV2 accountv2.Accounts
}

//User is struct with user details
type User struct {
	ID      string
	Email   string
	Account string
}

// PISessionVars is an object that holds the variables required to create an ibmpisession object.
// TODO: This should be renamed as it's being used to create the bluemix session
type PISessionVars struct {
	ID        string `json:"id,omitempty"`
	ICAPIKey  string `json:"icapikey,omitempty"`
	PVSAPIKey string `json:"pvsapikey,omitempty"`
	ICEP      string `json:"icep,omitempty"`
	PVSEP     string `json:"pvsep,omitempty"`
	Region    string `json:"region,omitempty"`
	Zone      string `json:"zone,omitempty"`
}

func authenticateAPIKey(config *bx.Config) error {
	tokenRefresher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if err != nil {
		return err
	}
	return tokenRefresher.AuthenticateAPIKey(config.BluemixAPIKey)
}

func fetchUserDetails(config *bx.Config) (*User, error) {
	user := User{}
	var bluemixToken string

	//logrus.Debugf("config.IAMAccessToken: %v", config.IAMAccessToken)

	if strings.HasPrefix(config.IAMAccessToken, "Bearer") {
		bluemixToken = config.IAMAccessToken[7:len(config.IAMAccessToken)]
	} else {
		bluemixToken = config.IAMAccessToken
	}

	token, err := jwt.Parse(bluemixToken, func(token *jwt.Token) (interface{}, error) {
		return "", nil
	})
	if err != nil && !strings.Contains(err.Error(), "key is of invalid type") {
		return &user, err
	}

	logrus.Debug("")
	claims := token.Claims.(jwt.MapClaims)
	if email, ok := claims["email"]; ok {
		logrus.Debug("user email")
		user.Email = email.(string)
	}
	user.ID = claims["id"].(string)
	logrus.Debugf("user id: %s", user.ID)
	user.Account = claims["account"].(map[string]interface{})["bss"].(string)
	logrus.Debugf("user account: %s", user.Account)
	return &user, nil
}

// GetIAMEndpointURL retrievs the IAM Endpoint URL specifid
// by the user	 via the endpoints file, if it exists locally.
func (c *BxClient) GetIAMEndpointURL() string {
	logrus.Debug("GetIAMEndpointURL")
	if c.bxConfig.EndpointsFile == "" {
		return ""
	}
	// find IBMCLOUD_IAM_API_ENDPOINT value. We've already
	// verified that this file exists in NewBxClient
	epLocator := bxep.NewEndpointLocator(c.bxConfig.Region, "public", c.bxConfig.EndpointsFile)
	if ep, err := epLocator.IAMEndpoint(); err == nil {
		return ep
	} else {
		logrus.Debug(err.Error())
	}
	logrus.Debug("GetIAMEndpointURL: no ep")
	return ""
}

//NewBxClient func returns bluemix client
func NewBxClient(t APIEndpointType) (*BxClient, error) {

	logrus.Debug("NewBxClient entering")

	var (
		pisv      PISessionVars
		bxSession *bxsession.Session
		bxClient  = &BxClient{bxConfig: &bx.Config{}}
	)

	// Grab variables from the users environment
	logrus.Debug("Gathering variables from user environment")
	if err := getPISessionVarsFromEnv(&pisv); err != nil {
		return nil, err
	}

	// Grab variables from the installer written authFilePath
	logrus.Debug("Gathering variables from AuthFile")
	if err := getPISessionVarsFromAuthFile(&pisv); err != nil {
		return nil, err
	}

	switch t {
	case PowerVSEP:
		bxClient.bxConfig.BluemixAPIKey = pisv.PVSAPIKey
		bxClient.bxConfig.EndpointsFile = defaultPVSEPFilePath
		logrus.Debugf("NewBxClient, API Key is Power Key: %s", bxClient.bxConfig.BluemixAPIKey)
	case IBMCloudEP:
		bxClient.bxConfig.BluemixAPIKey = pisv.ICAPIKey
		bxClient.bxConfig.EndpointsFile = defaultICEPFilePath
		logrus.Debugf("NewBxClient, API Key is IC Key: %s", bxClient.bxConfig.BluemixAPIKey)
	}

	if _, err := os.Stat(bxClient.bxConfig.EndpointsFile); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, errors.Wrap(err, "error accessing endpoints file:")
		} else {
			logrus.Debugf("endpoints file not present: %s", bxClient.bxConfig.EndpointsFile)
			bxClient.bxConfig.EndpointsFile = ""
		}
	} else {
		logrus.Debugf("using endpoints file at %s", bxClient.bxConfig.EndpointsFile)
	}

	var err error

	bxSession, err = bxsession.New(bxClient.bxConfig)
	if err != nil {
		return nil, err
	}

	endpoint, err := bxClient.bxConfig.EndpointLocator.IAMEndpoint()
	if err != nil {
		logrus.Debug("UH OH. DEV DEBUG ERROR COND")
		return nil, err
	}
	logrus.Debugf("iamendpoint: %v", endpoint)
	if err = authenticateAPIKey(bxClient.bxConfig); err != nil {
		return nil, err
	}

	if bxClient.User, err = fetchUserDetails(bxClient.bxConfig); err != nil {
		logrus.Debug("fetching user details")
		return nil, err
	}

	// is this right? or should it be the e-mail address?
	//pisv.ID = "clnperez@us.ibm.com"
	pisv.ID = bxClient.User.ID

	bxClient.bxConfig.Region = powervs.Regions[pisv.Region].VPCRegion

	// Prompt the user for the remaining variables.
	if err := getPISessionVarsFromUser(&pisv); err != nil {
		return nil, err
	}

	// Save variables to disk.
	if err := savePISessionVars(&pisv); err != nil {
		return nil, err
	}

	if accClient, err := accountv2.New(bxSession); err != nil {
		return nil, errors.Wrap(err, "error retreiving accountv2 client:")
	} else {
		bxClient.AccountAPIV2 = accClient.Accounts()
	}

	logrus.Debugf("succesfully created client")
	return bxClient, nil
}

//GetAccountType func return the type of account TRAIL/PAID
func (c *BxClient) GetAccountType() (string, error) {
	myAccount, err := c.AccountAPIV2.Get((*c.User).Account)
	if err != nil {
		return "", err
	}

	return myAccount.Type, nil
}

//ValidateAccountPermissions Checks permission for provisioning Power VS resources
func (c *BxClient) ValidateAccountPermissions() error {
	accType, err := c.GetAccountType()
	if err != nil {
		return err
	}
	if accType == "TRIAL" {
		return fmt.Errorf("account type must be of Pay-As-You-Go/Subscription type for provision Power VS resources")
	}
	return nil
}

//ValidateDhcpService checks for existing Dhcp service for the provided PowerVS cloud instance
func (c *BxClient) ValidateDhcpService(ctx context.Context, svcInsID string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	//create Power VS DHCP Client
	dhcpClient := instance.NewIBMPIDhcpClient(ctx, c.PISession, svcInsID)
	//Get all DHCP Services
	dhcpServices, err := dhcpClient.GetAll()
	if err != nil {
		return errors.Wrap(err, "failed to get DHCP service details")
	}
	if len(dhcpServices) > 0 {
		return fmt.Errorf("DHCP service already exists for provided cloud instance")
	}
	return nil
}

//ValidateCloudConnectionInPowerVSRegion counts cloud connection in PowerVS Region
func (c *BxClient) ValidateCloudConnectionInPowerVSRegion(ctx context.Context, svcInsID string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	var cloudConnectionsIDs []string
	cloudConnectionClient := instance.NewIBMPICloudConnectionClient(ctx, c.PISession, svcInsID)

	//check number of cloudconnections
	getAllResp, err := cloudConnectionClient.GetAll()
	if err != nil {
		return errors.Wrap(err, "failed to get existing Cloud connection details")
	}

	if len(getAllResp.CloudConnections) >= 2 {
		return fmt.Errorf("cannot create 	new Cloud connection in Power VS. Only two Cloud connections are allowed per zone")
	}

	for _, cc := range getAllResp.CloudConnections {
		cloudConnectionsIDs = append(cloudConnectionsIDs, *cc.CloudConnectionID)
	}

	//check for Cloud connection attached to DHCP Service
	for _, cc := range cloudConnectionsIDs {
		cloudConn, err := cloudConnectionClient.Get(cc)
		if err != nil {
			return errors.Wrap(err, "failed to get Cloud connection details")
		}
		if cloudConn != nil {
			for _, nw := range cloudConn.Networks {
				if nw.DhcpManaged {
					return fmt.Errorf("only one Cloud connection can be attached to any DHCP network per account per zone")
				}
			}
		}
	}
	return nil
}

// NewPISession updates pisession details, return error on fail
func (c *BxClient) NewPISession() error {
	var (
		pisv  PISessionVars
		epURL string = c.GetIAMEndpointURL()
	)
	// Grab variables from the installer written authFilePath
	logrus.Debug("Gathering variables from AuthFile")
	err := getPISessionVarsFromAuthFile(&pisv)
	if err != nil {
		return err
	}

	var authenticator core.Authenticator = &core.IamAuthenticator{
		ApiKey: c.bxConfig.BluemixAPIKey,
		URL:    epURL,
	}

	// Create the session
	options := &ibmpisession.IBMPIOptions{
		Authenticator: authenticator,
		UserAccount:   c.User.Account,
		Region:        pisv.Region,
		// @TODO: don't hardcode
		URL:   "https://dal.power-iaas.test.cloud.ibm.com",
		Zone:  pisv.Zone,
		Debug: os.Getenv("IBM_POWERVS_DEV") == "TRUE",
	}

	c.PISession, err = ibmpisession.NewIBMPISession(options)
	if err != nil {
		return err
	}
	return nil
}

// GetBxClientAPIKey returns the API key used by the Blue Mix Client.
func (c *BxClient) GetBxClientAPIKey() string {
	return c.bxConfig.BluemixAPIKey
}

func getPISessionVarsFromAuthFile(pisv *PISessionVars) error {

	if pisv == nil {
		return errors.New("nil var: PISessionVars")
	}

	authFilePath := defaultAuthFilePath
	if f := os.Getenv("POWERVS_AUTH_FILEPATH"); len(f) > 0 {
		authFilePath = f
	}

	if _, err := os.Stat(authFilePath); os.IsNotExist(err) {
		return nil
	}

	if content, err := ioutil.ReadFile(authFilePath); err == nil {
		if err := json.Unmarshal(content, pisv); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func getPISessionVarsFromEnv(pisv *PISessionVars) error {

	if pisv == nil {
		return errors.New("getPISessionVarsFromEnv nil var: PiSessionVars")
	}

	if len(pisv.ID) == 0 {
		pisv.ID = os.Getenv("IBMID")
	}

	if len(pisv.ICAPIKey) == 0 {
		// APIKeyEnvVars is a list of environment variable names containing an IBM Cloud API key.
		var APIKeyEnvVars = []string{"IC_API_KEY", "IBMCLOUD_API_KEY", "BM_API_KEY", "BLUEMIX_API_KEY"}
		pisv.ICAPIKey = getEnv(APIKeyEnvVars)
	}

	if len(pisv.PVSAPIKey) == 0 {
		pisv.PVSAPIKey = os.Getenv("POWER_VS_API_KEY")
	}

	if len(pisv.Region) == 0 {
		var regionEnvVars = []string{"IBMCLOUD_RqEGION", "IC_REGION"}
		pisv.Region = getEnv(regionEnvVars)
	}

	if len(pisv.Zone) == 0 {
		var zoneEnvVars = []string{"IBMCLOUD_ZONE"}
		pisv.Zone = getEnv(zoneEnvVars)
	}

	logrus.Debugf("getPISessionVarsFromEnv: pisv: %v", pisv)
	return nil
}

func getPISessionVarsFromUser(pisv *PISessionVars) error {
	var err error

	if pisv == nil {
		return errors.New("nil var: PiSessionVars")
	}

	if len(pisv.ID) == 0 {
		err = survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Input{
					Message: "IBM Cloud User ID",
					Help:    "The login for \nhttps://cloud.ibm.com/",
				},
			},
		}, &pisv.ID)
		if err != nil {
			return errors.Wrap(err, "error saving the IBM Cloud User ID")
		}

	}

	if len(pisv.ICAPIKey) == 0 {
		err = survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Password{
					Message: "IBM Cloud API Key",
					Help:    "The API key installation.\nhttps://cloud.ibm.com/iam/apikeys",
				},
			},
		}, &pisv.ICAPIKey)
		if err != nil {
			return errors.Wrap(err, "error saving the API Key")
		}

	}

	if len(pisv.Region) == 0 {
		pisv.Region, err = GetRegion()
		if err != nil {
			return err
		}

	}

	if len(pisv.Zone) == 0 {
		pisv.Zone, err = GetZone(pisv.Region)
		if err != nil {
			return err
		}
	}

	return nil
}

func savePISessionVars(pisv *PISessionVars) error {

	authFilePath := defaultAuthFilePath
	if f := os.Getenv("POWERVS_AUTH_FILEPATH"); len(f) > 0 {
		authFilePath = f
	}

	jsonVars, err := json.Marshal(*pisv)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(authFilePath), 0700)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(authFilePath, jsonVars, 0600)
}

func getEnv(envs []string) string {
	for _, k := range envs {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}
