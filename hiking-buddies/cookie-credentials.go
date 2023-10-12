package hiking_buddies

type CookieCredential struct {
	SessionId, CSRFToken string
}

func (c *CookieCredential) AsCookie() map[string]string {
	return map[string]string{
		"sessionid": c.SessionId,
		"csrftoken": c.CSRFToken,
	}
}
