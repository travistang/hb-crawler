package hiking_buddies

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/storage"
	"github.com/chromedp/chromedp"
	"github.com/corpix/uarand"
)

type Credential struct {
	Email, Password string
}

const (
	usernameInput = "//*/input[@name='username']"
	passwordInput = "//*/input[@name='password']"
	loginButton   = "//*/button[@type='submit']"

	logoutButton = `//*/a[@href='/routes/logout_user/']`

	loginCookieName = "sessionid"

	screenshotQuality = 0o644
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

func Login(credential *Credential) (error, *CookieCredential) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var cookie CookieCredential
	err := chromedp.Run(ctx, performLoginSteps(credential, &cookie))
	if err != nil {
		return err, nil
	}

	return nil, &cookie
}
