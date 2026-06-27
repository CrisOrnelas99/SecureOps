// Package nvd verifies the NVD API client and response mapping.
package nvd

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestClientLookupCVE verifies request construction and safe DTO mapping.
func TestClientLookupCVE(t *testing.T) {
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.URL.Host != officialNVDHost {
			t.Fatalf("expected official NVD host, got %q", request.URL.Host)
		}
		if request.URL.Query().Get("cveIds") != "CVE-2021-44228" {
			t.Fatalf("expected cveIds query, got %q", request.URL.RawQuery)
		}
		if request.Header.Get("apiKey") != "server-side-key" {
			t.Fatal("expected API key to be sent as a server-side header")
		}

		body := `{
			"totalResults": 1,
			"vulnerabilities": [{
				"cve": {
					"id": "CVE-2021-44228",
					"published": "2021-12-10T10:15:09.067",
					"lastModified": "2024-11-21T12:15:26.783",
					"descriptions": [{"lang": "en", "value": "Apache Log4j remote code execution."}],
					"metrics": {"cvssMetricV31": [{"cvssData": {"baseSeverity": "CRITICAL"}}]}
				}
			}]
		}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewClientWithHTTPClient(
		"https://services.nvd.nist.gov/rest/json/cves/2.0",
		"server-side-key",
		&http.Client{Transport: transport},
		NewRateLimiter(10, time.Second),
	)
	if err != nil {
		t.Fatalf("expected client to build, got %v", err)
	}

	response, err := client.LookupCVE(context.Background(), " cve-2021-44228 ")
	if err != nil {
		t.Fatalf("expected lookup to succeed, got %v", err)
	}
	if response.CVEID != "CVE-2021-44228" {
		t.Fatalf("expected normalized CVE ID, got %q", response.CVEID)
	}
	if response.Severity != "CRITICAL" {
		t.Fatalf("expected severity CRITICAL, got %q", response.Severity)
	}
	if response.NVDURL != "https://nvd.nist.gov/vuln/detail/CVE-2021-44228" {
		t.Fatalf("unexpected NVD URL: %q", response.NVDURL)
	}
}

// TestClientRejectsUnsafeBaseURL verifies outbound host allowlisting.
func TestClientRejectsUnsafeBaseURL(t *testing.T) {
	_, err := NewClientWithHTTPClient("https://example.com/rest/json/cves/2.0", "", nil, nil)
	if !errors.Is(err, ErrInvalidBaseURL) {
		t.Fatalf("expected invalid base URL error, got %v", err)
	}
}

// TestClientRestrictsRedirects verifies the production client does not follow redirects.
func TestClientRestrictsRedirects(t *testing.T) {
	client, err := NewClient("https://services.nvd.nist.gov/rest/json/cves/2.0", "")
	if err != nil {
		t.Fatalf("expected client to build, got %v", err)
	}
	if client.httpClient.CheckRedirect == nil {
		t.Fatal("expected redirect policy to be configured")
	}

	redirectRequest, err := http.NewRequest(http.MethodGet, "https://example.com/redirect", nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	originalRequest, err := http.NewRequest(http.MethodGet, "https://services.nvd.nist.gov/rest/json/cves/2.0", nil)
	if err != nil {
		t.Fatalf("failed to build original request: %v", err)
	}

	err = client.httpClient.CheckRedirect(redirectRequest, []*http.Request{originalRequest})
	if !errors.Is(err, http.ErrUseLastResponse) {
		t.Fatalf("expected redirect to be blocked, got %v", err)
	}
}

// TestClientHandlesNVDNotFound verifies empty NVD results map to not found.
func TestClientHandlesNVDNotFound(t *testing.T) {
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"totalResults":0,"vulnerabilities":[]}`)),
			Header:     make(http.Header),
		}, nil
	})
	client, err := NewClientWithHTTPClient(
		"https://services.nvd.nist.gov/rest/json/cves/2.0",
		"",
		&http.Client{Transport: transport},
		NewRateLimiter(10, time.Second),
	)
	if err != nil {
		t.Fatalf("expected client to build, got %v", err)
	}

	_, err = client.LookupCVE(context.Background(), "CVE-2021-44228")
	if !errors.Is(err, ErrCVEIDNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

// TestClientRejectsMismatchedCVEID verifies external response identity is checked.
func TestClientRejectsMismatchedCVEID(t *testing.T) {
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(strings.NewReader(`{
				"totalResults": 1,
				"vulnerabilities": [{
					"cve": {
						"id": "CVE-1999-0001",
						"descriptions": [{"lang": "en", "value": "wrong record"}]
					}
				}]
			}`)),
			Header: make(http.Header),
		}, nil
	})
	client, err := NewClientWithHTTPClient(
		"https://services.nvd.nist.gov/rest/json/cves/2.0",
		"",
		&http.Client{Transport: transport},
		NewRateLimiter(10, time.Second),
	)
	if err != nil {
		t.Fatalf("expected client to build, got %v", err)
	}

	_, err = client.LookupCVE(context.Background(), "CVE-2021-44228")
	if !errors.Is(err, ErrInvalidNVDResponse) {
		t.Fatalf("expected invalid NVD response, got %v", err)
	}
}

// TestRateLimiter verifies rolling-window rate limiting.
func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(1, time.Minute)
	now := time.Date(2026, 6, 27, 12, 0, 0, 0, time.UTC)

	if !limiter.Allow(now) {
		t.Fatal("expected first request to be allowed")
	}
	if limiter.Allow(now.Add(time.Second)) {
		t.Fatal("expected second request inside the window to be rejected")
	}
	if !limiter.Allow(now.Add(time.Minute + time.Second)) {
		t.Fatal("expected request after the window to be allowed")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
