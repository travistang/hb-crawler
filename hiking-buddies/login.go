package hiking_buddies

import (
	"context"
	"fmt"
	"hb-crawler/rating-gain/database"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/storage"
	"github.com/chromedp/chromedp"
	"github.com/corpix/uarand"
	"github.com/sirupsen/logrus"
)

type Credential struct {
	Email, Password string
}

const (
	usernameInput = "//*/input[@name='username']"
	passwordInput = "//*/input[@name='password']"
	loginButton   = "//*/button[@type='submit']"

	loginCookieName = "sessionid"
	csrfCookieName  = "csrftoken"
)

func createHeaders() network.Headers {
	return network.Headers{
		"User-Agent":     uarand.GetRandom(),
		"Sec-Fetch-Mode": "cors",
	}
}

func sleep(seconds int) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			time.Sleep(time.Duration(seconds) * time.Second)
			return nil
		}),
	}
}

func retrieveLoginCookie(cookie *CookieCredential) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := storage.GetCookies().Do(ctx)
			if err != nil {
				return err
			}

			for _, cookieField := range cookies {
				value := cookieField.Value
				switch cookieField.Name {
				case "sessionid":
					cookie.SessionId = value
					break
				case "csrftoken":
					cookie.CSRFToken = value
					break
				}

			}
			return nil
		}),
	}
}

func performLoginSteps(credential *Credential, cookie *CookieCredential) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		network.SetExtraHTTPHeaders(createHeaders()),
		chromedp.Navigate(string(LoginEndpoint)),

		chromedp.WaitVisible(usernameInput),
		chromedp.SendKeys(usernameInput, credential.Email),

		chromedp.WaitVisible(passwordInput),
		chromedp.SendKeys(passwordInput, credential.Password),

		chromedp.WaitVisible(loginButton),
		chromedp.Click(loginButton),

		sleep(2),

		retrieveLoginCookie(cookie),
	}
}

func doLogin(credential *Credential) (*CookieCredential, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var cookie CookieCredential
	err := chromedp.Run(ctx, performLoginSteps(credential, &cookie))
	if err != nil {
		return nil, err
	}

	return &cookie, nil
}

func getCachedCredential(repo *database.LoginCredentialRepository, credential *Credential) (*CookieCredential, error) {
	savedCredential := database.LoginCredentialRecord{}
	now := time.Now()
	ageThreshold := now.Add(time.Duration(-3) * time.Hour).Unix()

	logrus.Debugf("Retrieving credentials for user %s newer than time %d", credential.Email, ageThreshold)
	err := repo.GetCredential(&savedCredential, credential.Email, ageThreshold)
	if err != nil {
		return nil, err
	}

	if savedCredential.Username != credential.Email {
		return nil, fmt.Errorf("no eligible credential for user %s", credential.Email)
	}

	return &CookieCredential{
		SessionId: savedCredential.SessionId,
		CSRFToken: "",
	}, nil
}

func cacheCredential(
	repo *database.LoginCredentialRepository,
	credential *Credential,
	cookie *CookieCredential,
) error {
	return repo.SaveCredential(&database.LoginCredentialRecord{
		SessionId: cookie.SessionId,
		Username:  credential.Email,
	})
}

func Login(repo *database.LoginCredentialRepository, credential *Credential) (*CookieCredential, error) {
	cached, _ := getCachedCredential(repo, credential)
	if cached != nil {
		return cached, nil
	}

	cookie, err := doLogin(credential)
	if err != nil {
		return nil, err
	}

	if err := cacheCredential(repo, credential, cookie); err != nil {
		return nil, err
	}

	return cookie, nil
}
