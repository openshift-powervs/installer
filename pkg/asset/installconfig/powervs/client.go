package powervs

import (
	gohttp "net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/account/accountv2"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	bxsession "github.com/IBM-Cloud/bluemix-go/session"
	"github.com/form3tech-oss/jwt-go"
	_ "github.com/pkg/errors"
)

//Client is struct which provides bluemix session details
type Client struct {
	*bxsession.Session
	User         *User
	AccountAPIV2 accountv2.Accounts
}

//User is struct with user details
type User struct {
	ID      string
	Email   string
	Account string
}

func authenticateAPIKey(sess *bxsession.Session) error {
	config := sess.Config
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

func fetchUserDetails(sess *bxsession.Session) (*User, error) {
	config := sess.Config
	user := User{}
	var bluemixToken string

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

	claims := token.Claims.(jwt.MapClaims)
	if email, ok := claims["email"]; ok {
		user.Email = email.(string)
	}
	user.ID = claims["id"].(string)
	user.Account = claims["account"].(map[string]interface{})["bss"].(string)

	return &user, nil
}

//NewClient func returns new bluemix client
func NewClient() (*Client, error) {
	//var apikey string
	c := &Client{}

	var pisv PISessionVars
	// Grab variables from the installer written authFilePath
	logrus.Debug("Gathering variables from AuthFile")
	err := getPISessionVarsFromAuthFile(&pisv)
	if err != nil {
		return nil, err
	}

	// Frab variables from the users environment
	logrus.Debug("Gathering variables from user environment")
	err = getPISessionVarsFromEnv(&pisv)
	if err != nil {
		return nil, err
	}

	// Prompt the user for the remaining variables
	logrus.Debug("Gathering variables from user")
	err = getPISessionVarsFromUser(&pisv)
	if err != nil {
		return nil, err
	}

	os.Setenv("IC_API_KEY", pisv.APIKey)

	/*if apikey = os.Getenv("API_KEY"); len(apikey) == 0 {
		return nil, errors.New("empty API_KEY variable")
	}*/

	bxSess, err := bxsession.New(&bluemix.Config{
		BluemixAPIKey: pisv.APIKey,
	})
	if err != nil {
		return nil, err
	}

	c.Session = bxSess

	err = authenticateAPIKey(bxSess)
	if err != nil {
		return nil, err
	}

	c.User, err = fetchUserDetails(bxSess)
	if err != nil {
		return nil, err
	}

	accClient, err := accountv2.New(bxSess)
	if err != nil {
		return nil, err
	}

	c.AccountAPIV2 = accClient.Accounts()
	return c, nil
}

//GetAccountType func return the type of account TRAIL/PAID
func (c *Client) GetAccountType() (string, error) {
	myAccount, err := c.AccountAPIV2.Get((*c.User).Account)
	if err != nil {
		return "", err
	}

	return myAccount.Type, nil
}
