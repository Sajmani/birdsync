package inat

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func makeTestJWT(t *testing.T, claims any) string {
	t.Helper()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS512"}`))
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("marshaling test JWT claims: %v", err)
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	sig := base64.RawURLEncoding.EncodeToString([]byte("test-signature"))
	return header + "." + payload + "." + sig
}

type failOnRead struct {
	t *testing.T
}

func (f failOnRead) Read([]byte) (int, error) {
	f.t.Error("unexpected read from stdin when credential was already resolved")
	return 0, errors.New("unexpected stdin read")
}

func TestParseAPIToken(t *testing.T) {
	sampleJWT := makeTestJWT(t, map[string]any{"user_id": 42, "exp": 1900000000})

	good := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "bare token",
			raw:  sampleJWT,
			want: sampleJWT,
		},
		{
			name: "bare token with surrounding whitespace",
			raw:  "  " + sampleJWT + "\n",
			want: sampleJWT,
		},
		{
			name: "compact JSON framing",
			raw:  `{"api_token":"` + sampleJWT + `"}`,
			want: sampleJWT,
		},
		{
			name: "spaced JSON framing",
			raw:  `  { "api_token" : "` + sampleJWT + `" }  `,
			want: sampleJWT,
		},
	}
	for _, tc := range good {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseAPIToken(tc.raw)
			if err != nil {
				t.Fatalf("ParseAPIToken() error = %v", err)
			}
			if got != tc.want {
				t.Errorf("ParseAPIToken() = %q, want %q", got, tc.want)
			}
		})
	}

	bad := []struct {
		name string
		raw  string
	}{
		{name: "empty", raw: ""},
		{name: "whitespace only", raw: "   \t\n"},
		{name: "malformed JSON", raw: `{"api_token": "secret-val`},
		{name: "empty api_token in JSON", raw: `{"api_token":""}`},
		{name: "wrong JSON key", raw: `{"token":"abc"}`},
		{name: "bare token with internal space", raw: "token with spaces"},
		{name: "JSON token with internal space", raw: `{"api_token":"token with spaces"}`},
	}
	for _, tc := range bad {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseAPIToken(tc.raw)
			if err == nil {
				t.Fatalf("ParseAPIToken(%q) = %q, want error", tc.raw, got)
			}
		})
	}
}

func TestUserIDFromToken(t *testing.T) {
	cases := []struct {
		name   string
		token  string
		wantID string
		wantOK bool
	}{
		{
			name:   "numeric user_id claim",
			token:  makeTestJWT(t, map[string]any{"user_id": 8543312, "exp": 1900000000}),
			wantID: "8543312",
			wantOK: true,
		},
		{
			name:   "string user_id claim",
			token:  makeTestJWT(t, map[string]any{"user_id": "birder42"}),
			wantID: "birder42",
			wantOK: true,
		},
		{
			name:   "zero user_id rejected",
			token:  makeTestJWT(t, map[string]any{"user_id": 0}),
			wantOK: false,
		},
		{
			name:   "negative user_id rejected",
			token:  makeTestJWT(t, map[string]any{"user_id": -5}),
			wantOK: false,
		},
		{
			name:   "missing user_id claim",
			token:  makeTestJWT(t, map[string]any{"exp": 1900000000}),
			wantOK: false,
		},
		{
			name:   "opaque non-JWT token",
			token:  "not-a-jwt-token",
			wantOK: false,
		},
		{
			name:   "invalid base64 payload",
			token:  "header.!!!invalid!!!.sig",
			wantOK: false,
		},
		{
			name:   "non-JSON payload",
			token:  "header." + base64.RawURLEncoding.EncodeToString([]byte("not json")) + ".sig",
			wantOK: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotID, gotOK := UserIDFromToken(tc.token)
			if gotOK != tc.wantOK || gotID != tc.wantID {
				t.Errorf("UserIDFromToken() = (%q, %v), want (%q, %v)", gotID, gotOK, tc.wantID, tc.wantOK)
			}
		})
	}
}

// TestJWTFromEnvironmentInHeader verifies that both bare and JSON-framed
// INAT_API_TOKEN values are unwrapped and sent verbatim in the Authorization
// header of authenticated requests.
func TestJWTFromEnvironmentInHeader(t *testing.T) {
	wantToken := makeTestJWT(t, map[string]any{"user_id": 12345})

	for _, tc := range []struct {
		name string
		env  string
	}{
		{name: "bare token", env: wantToken},
		{name: "JSON-framed token", env: `{"api_token":"` + wantToken + `"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tok, err := getAPIToken(tc.env, failOnRead{t: t})
			if err != nil {
				t.Fatalf("getAPIToken() error = %v", err)
			}

			var authHeaders []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authHeaders = append(authHeaders, r.Header.Get("Authorization"))
				if r.Method == http.MethodGet {
					json.NewEncoder(w).Encode(Observations{})
					return
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			client := newTestClient(server.URL, tok, "test-agent")
			if _, err := client.DownloadObservations("12345", time.Time{}, time.Time{}); err != nil {
				t.Fatalf("DownloadObservations() error = %v", err)
			}
			if err := client.CreateObservation(Observation{}); err != nil {
				t.Fatalf("CreateObservation() error = %v", err)
			}
			if len(authHeaders) != 2 {
				t.Fatalf("server saw %d requests, want 2", len(authHeaders))
			}
			for i, got := range authHeaders {
				if got != wantToken {
					t.Errorf("request %d Authorization = %q, want %q", i+1, got, wantToken)
				}
			}
		})
	}
}

// TestJWTFromPromptUnwrapped verifies that when INAT_API_TOKEN is unset, both
// JSON-framed (including spaces) and bare tokens entered at the prompt are
// unwrapped and sent in the Authorization header, and invalid lines retry.
func TestJWTFromPromptUnwrapped(t *testing.T) {
	origOut := log.Writer()
	log.SetOutput(&bytes.Buffer{})
	defer log.SetOutput(origOut)

	wantToken := makeTestJWT(t, map[string]any{"user_id": 67890})

	for _, tc := range []struct {
		name  string
		input string
	}{
		{
			name:  "JSON-framed with spaces",
			input: `{"api_token": "` + wantToken + `"}` + "\n",
		},
		{
			name:  "bare token",
			input: wantToken + "\n",
		},
		{
			name:  "retries after invalid line",
			input: "{bad json\n" + `{"api_token":"` + wantToken + `"}` + "\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tok, err := getAPIToken("", strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("getAPIToken() error = %v", err)
			}

			var gotAuth string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotAuth = r.Header.Get("Authorization")
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			client := newTestClient(server.URL, tok, "test-agent")
			if err := client.CreateObservation(Observation{}); err != nil {
				t.Fatalf("CreateObservation() error = %v", err)
			}
			if gotAuth != wantToken {
				t.Errorf("Authorization = %q, want %q", gotAuth, wantToken)
			}
		})
	}
}

// TestNoBuiltInCredentials verifies that when neither INAT_API_TOKEN nor
// interactive input supplies a token, credential resolution fails and no
// built-in fallback token is used.
func TestNoBuiltInCredentials(t *testing.T) {
	origOut := log.Writer()
	log.SetOutput(&bytes.Buffer{})
	defer log.SetOutput(origOut)

	tok, err := getAPIToken("", strings.NewReader(""))
	if err == nil {
		t.Fatalf("getAPIToken() with empty stdin returned (%q, nil), want error", tok)
	}
	if tok != "" {
		t.Errorf("getAPIToken() returned non-empty token %q on error", tok)
	}
}

// TestGetUserID verifies precedence: INAT_USER_ID first, then JWT user_id
// claim without touching stdin, then interactive prompt with EOF error.
func TestGetUserID(t *testing.T) {
	origOut := log.Writer()
	log.SetOutput(&bytes.Buffer{})
	defer log.SetOutput(origOut)

	jwtToken := makeTestJWT(t, map[string]any{"user_id": 8543312})

	t.Run("INAT_USER_ID takes precedence over JWT", func(t *testing.T) {
		got, err := getUserID("birder", jwtToken, failOnRead{t: t})
		if err != nil {
			t.Fatalf("getUserID() error = %v", err)
		}
		if got != "birder" {
			t.Errorf("getUserID() = %q, want %q", got, "birder")
		}
	})

	t.Run("unset INAT_USER_ID extracts user_id from JWT without reading stdin", func(t *testing.T) {
		got, err := getUserID("", jwtToken, failOnRead{t: t})
		if err != nil {
			t.Fatalf("getUserID() error = %v", err)
		}
		if got != "8543312" {
			t.Errorf("getUserID() = %q, want %q", got, "8543312")
		}
	})

	t.Run("unset INAT_USER_ID and non-JWT prompts on stdin", func(t *testing.T) {
		got, err := getUserID("", "opaque-token", strings.NewReader("\nprompteduser\n"))
		if err != nil {
			t.Fatalf("getUserID() error = %v", err)
		}
		if got != "prompteduser" {
			t.Errorf("getUserID() = %q, want %q", got, "prompteduser")
		}
	})

	t.Run("unset INAT_USER_ID and non-JWT returns error on EOF", func(t *testing.T) {
		got, err := getUserID("", "opaque-token", strings.NewReader(""))
		if err == nil {
			t.Fatalf("getUserID() on EOF returned (%q, nil), want error", got)
		}
	})
}

// TestCredentialsNeverLogged checks that the API token never appears in
// captured log output or returned errors across parsing, prompts, and HTTP
// requests including 401 and error responses.
func TestCredentialsNeverLogged(t *testing.T) {
	const sentinel = "SUPER_SECRET_SENTINEL_TOKEN_VALUE"
	jwtToken := makeTestJWT(t, map[string]any{"user_id": 999, "sub": sentinel})

	var logBuf bytes.Buffer
	origOut := log.Writer()
	log.SetOutput(&logBuf)
	defer log.SetOutput(origOut)

	checkNoLeak := func(where string, err error) {
		t.Helper()
		if err != nil {
			msg := err.Error()
			if strings.Contains(msg, sentinel) || strings.Contains(msg, jwtToken) {
				t.Errorf("%s: error leaked credential: %q", where, msg)
			}
		}
		logs := logBuf.String()
		if strings.Contains(logs, sentinel) || strings.Contains(logs, jwtToken) {
			t.Errorf("%s: log output leaked credential", where)
		}
	}

	// Valid and malformed environment tokens.
	_, err := getAPIToken(jwtToken, failOnRead{t: t})
	checkNoLeak("valid env token", err)

	_, err = getAPIToken(`{"api_token":"`+sentinel, failOnRead{t: t})
	if err == nil {
		t.Fatal("malformed env JSON should fail")
	}
	checkNoLeak("malformed env JSON", err)

	_, err = getAPIToken(sentinel+" extra", failOnRead{t: t})
	if err == nil {
		t.Fatal("env token with whitespace should fail")
	}
	checkNoLeak("env token with whitespace", err)

	// Interactive prompt retry on malformed input, then valid JSON token.
	promptIn := `{"api_token":"` + sentinel + "\n" + `{"api_token":"` + jwtToken + `"}` + "\n"
	tok, err := getAPIToken("", strings.NewReader(promptIn))
	if err != nil {
		t.Fatalf("prompt getAPIToken() error = %v", err)
	}
	checkNoLeak("prompt retry and success", err)

	// Authenticated HTTP calls (200, 401, 422).
	status := http.StatusOK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if status != http.StatusOK {
			w.Write([]byte("rejected"))
		}
	}))
	defer server.Close()

	client := newTestClient(server.URL, tok, "test-agent")
	for _, code := range []int{http.StatusOK, http.StatusUnauthorized, http.StatusUnprocessableEntity} {
		status = code
		err = client.CreateObservation(Observation{})
		checkNoLeak(http.StatusText(code), err)
	}
}
