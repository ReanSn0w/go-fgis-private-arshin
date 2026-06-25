package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultBaseURL    = "https://fgis.gost.ru/fundmetrology/api/"
	DefaultPublicURL  = "https://fgis.gost.ru/fundmetrology/"
	DefaultOrgID      = "CURRENT_ORG"
	DefaultUserAgent  = "go-fgis-private-arshin/0.1 (+https://github.com/ReanSn0w/go-fgis-private-arshin)"
	DefaultRatePeriod = time.Second
)

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	publicURL  *url.URL
	userAgent  string
	limiter    *RateLimiter
}

type Option func(*Client) error

func NewClient(options ...Option) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	baseURL, err := url.Parse(DefaultBaseURL)
	if err != nil {
		return nil, err
	}

	publicURL, err := url.Parse(DefaultPublicURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		publicURL:  publicURL,
		userAgent:  DefaultUserAgent,
		limiter:    NewRateLimiter(DefaultRatePeriod),
	}

	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		if httpClient == nil {
			return fmt.Errorf("http client is nil")
		}
		c.httpClient = httpClient
		return nil
	}
}

func WithBaseURL(rawURL string) Option {
	return func(c *Client) error {
		parsed, err := url.Parse(rawURL)
		if err != nil {
			return err
		}
		if !strings.HasSuffix(parsed.Path, "/") {
			parsed.Path += "/"
		}
		c.baseURL = parsed
		return nil
	}
}

func WithPublicURL(rawURL string) Option {
	return func(c *Client) error {
		parsed, err := url.Parse(rawURL)
		if err != nil {
			return err
		}
		if !strings.HasSuffix(parsed.Path, "/") {
			parsed.Path += "/"
		}
		c.publicURL = parsed
		return nil
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *Client) error {
		if strings.TrimSpace(userAgent) == "" {
			return fmt.Errorf("user agent is empty")
		}
		c.userAgent = userAgent
		return nil
	}
}

func WithRateLimit(interval time.Duration) Option {
	return func(c *Client) error {
		if interval < 0 {
			return fmt.Errorf("rate limit interval is negative")
		}
		c.limiter = NewRateLimiter(interval)
		return nil
	}
}

func (c *Client) newAPIURL(path string, values url.Values) string {
	u := c.baseURL.ResolveReference(&url.URL{Path: strings.TrimPrefix(path, "/")})
	u.RawQuery = values.Encode()
	return u.String()
}

func (c *Client) registryReferer(registryID string) string {
	u := c.publicURL.ResolveReference(&url.URL{Path: "registry/" + registryID})
	return u.String()
}

func (c *Client) setDefaultHeaders(req *http.Request, registryID string) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9")
	req.Header.Set("User-Agent", c.userAgent)
	if registryID != "" {
		req.Header.Set("Referer", c.registryReferer(registryID))
	}
}

func (c *Client) getJSON(ctx context.Context, path string, values url.Values, registryID string, target any) error {
	if err := c.limiter.Wait(ctx); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.newAPIURL(path, values), nil)
	if err != nil {
		return err
	}
	c.setDefaultHeaders(req, registryID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	return decodeJSONResponse(resp, target)
}

func decodeJSONResponse(resp *http.Response, target any) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		contentType := resp.Header.Get("Content-Type")
		if strings.Contains(strings.ToLower(contentType), "json") {
			var apiErr APIError
			if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Status != 0 {
				return &apiErr
			}
		}
		return NewHTTPError(resp)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" && !strings.Contains(strings.ToLower(contentType), "json") {
		return &UnexpectedContentTypeError{
			StatusCode:  resp.StatusCode,
			ContentType: contentType,
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return err
	}

	return nil
}
