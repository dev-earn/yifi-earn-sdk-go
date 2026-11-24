package yifi

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dev-earn/yifi-earn-sdk-go/models"
)

// API hosts. Override these only if you are pointing at a self-hosted or
// staging deployment; NewClient chooses between them based on UseTestnet.
var (
	BaseAPIMainURL    = "https://core.yifi.io"
	BaseAPITestnetURL = "https://core-dev.yifi.io"
)

// UseTestnet routes all requests to the testnet host instead of production.
// Set it before calling NewClient.
var UseTestnet = false

// getAPIEndpoint return the base endpoint of the Rest API according the UseTestnet flag
func getAPIEndpoint() string {
	if UseTestnet {
		return BaseAPITestnetURL
	}

	return BaseAPIMainURL
}

// DoReq executes req, decodes a JSON models.Response[T] on 2xx, and returns an
// error carrying the status and (truncated) body on any non-2xx response. It is
// exported for callers building requests manually; the SDK methods use it internally.
func DoReq[T any](cl *http.Client, req *http.Request) (*models.Response[T], error) {

	resp, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		s := strings.TrimSpace(string(respBody))

		if len(s) > 500 {
			s = s[:500] + "..."
		}

		return nil, fmt.Errorf(
			"http %d %s: %s",
			resp.StatusCode,
			http.StatusText(resp.StatusCode),
			s,
		)
	}

	var out models.Response[T]

	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	if !out.Success {
		return nil, fmt.Errorf("api error %s: %s", out.Code, out.Message)
	}

	return &out, nil
}

func normalizeBasePath(p string) string {
	p = strings.TrimSpace(p)

	if p == "" {
		return "/api/v1"
	}

	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}

	return strings.TrimRight(p, "/")
}

func joinPaths(parts ...string) string {
	clean := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)

		if p == "" {
			continue
		}

		clean = append(clean, strings.Trim(p, "/"))
	}

	return "/" + strings.Join(clean, "/")
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func signValidateRequest(secret string, req models.ValidateRequest) string {

	message := req.APIKey +
		strconv.FormatInt(req.Nonce, 10) +
		req.Method +
		req.Path +
		req.BodyHash

	h := hmac.New(sha256.New, []byte(secret))
	_, _ = h.Write([]byte(message))

	return hex.EncodeToString(h.Sum(nil))
}

func (c *Client) partnerHeaders(method, path string, body []byte) http.Header {
	if c == nil || c.partnerAuth == nil {
		return nil
	}

	if c.partnerAuth.PartnerID == "" ||
		c.partnerAuth.APIKey == "" ||
		c.partnerAuth.Secret == "" {
		return nil
	}

	h := make(http.Header)

	signedSecretBytes, err := json.Marshal(c.partnerAuth.Secret)
	if err != nil {
		panic(err)
	}

	nonce := time.Now().UTC().UnixMilli()

	h.Set("X-YIFI-PARTNER-ID", c.partnerAuth.PartnerID)
	h.Set("X-YIFI-API-KEY", c.partnerAuth.APIKey)
	h.Set("X-YIFI-NONCE", strconv.FormatInt(nonce, 10))
	h.Set("X-YIFI-SIGN", signValidateRequest(
		string(signedSecretBytes),
		models.ValidateRequest{
			APIKey:   c.partnerAuth.APIKey,
			Nonce:    nonce,
			Method:   method,
			Path:     path,
			BodyHash: sha256Hex(body),
		},
	))

	return h
}

func (c *Client) buildURL(p string, q url.Values) string {
	u := *c.baseURL

	u.Path = joinPaths(c.basePath, p)

	if q != nil {
		u.RawQuery = q.Encode()
	}

	return u.String()
}

func (c *Client) do(ctx context.Context, method string, p string, q url.Values, body any) (*http.Request, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
	}

	reqPath := joinPaths(c.basePath, p)

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.buildURL(p, q),
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	headers := c.partnerHeaders(method, reqPath, bodyBytes)

	for k, vals := range headers {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	return req, nil
}
