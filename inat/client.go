package inat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	debug = false
	// BaseURL is the standard base URL for the iNaturalist API.
	BaseURL = "https://api.inaturalist.org/v2"
)

type Client struct {
	apiToken  string
	userAgent string
	baseURL   string

	// mu guards the pacing and request-budget state. birdsync is
	// single-threaded today, but a client that silently stopped pacing or
	// counting under concurrency would be a nasty way to earn an IP block.
	mu                 sync.Mutex
	minRequestInterval time.Duration
	lastRequest        time.Time
	maxRequests        int
	requestCount       int
}

func NewClient(baseURL, apiToken, userAgent string) *Client {
	return &Client{
		baseURL:            baseURL,
		apiToken:           apiToken,
		userAgent:          userAgent,
		minRequestInterval: DefaultMinRequestInterval,
		maxRequests:        DefaultMaxRequests,
	}
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

// DefaultMinRequestInterval paces requests to the rate iNaturalist asks for:
//
//	Please keep requests to about 1 per second, and around 10k API requests a
//	day... Requests exceeding this limit might be throttled, and will return an
//	HTTP 429 exception "Too Many Requests"... We may block IPs that
//	consistently exceed these limits.
//
// A first sync of a media-heavy account runs to thousands of requests, and
// nothing else in birdsync limits how fast it issues them. Media uploads
// already take seconds each, so this costs little there; it matters for runs
// that create many observations without media, which would otherwise go as
// fast as the server answers (birdsync.requestRate).
const DefaultMinRequestInterval = time.Second

// DefaultMaxRequests caps the requests one Client sends at 9,000, leaving a
// 10% margin below iNaturalist's "around 10k API requests a day" guideline
// (tech.dailyRequestBudget). Because every HTTP call passes through roundTrip,
// this bounds downloads and --dryrun runs as well as writes.
const DefaultMaxRequests = 9000

// ErrRequestBudgetExceeded is returned when sending another request would
// exceed the client's per-run request cap.
var ErrRequestBudgetExceeded = errors.New("daily iNaturalist API request budget exceeded")

// SetMinRequestInterval overrides the pacing. Tests set it to zero; nothing
// else should raise the rate above the default without a reason to think
// iNaturalist has changed its guidance.
func (c *Client) SetMinRequestInterval(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.minRequestInterval = d
}

// SetMaxRequests overrides the per-client request cap for tests.
func (c *Client) SetMaxRequests(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.maxRequests = n
}

// RequestCount returns the number of HTTP requests this client has sent.
func (c *Client) RequestCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.requestCount
}

// pace checks the per-client request cap and blocks until enough time has
// passed since the previous request.
func (c *Client) pace() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.requestCount >= c.maxRequests {
		return fmt.Errorf("%w (%d requests)", ErrRequestBudgetExceeded, c.maxRequests)
	}
	if !c.lastRequest.IsZero() {
		if wait := c.minRequestInterval - time.Since(c.lastRequest); wait > 0 {
			time.Sleep(wait)
		}
	}
	c.lastRequest = time.Now()
	c.requestCount++
	return nil
}

func (c *Client) roundTrip(req *http.Request) (string, error) {
	// Every request in this package goes through here, which is the only
	// reason one pacing and budget check is enough.
	if err := c.pace(); err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", c.apiToken)

	if debug {
		log.Printf("\nREQUEST: %+v\n", req)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("making HTTP request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("%s: refresh your INAT_API_TOKEN from https://www.inaturalist.org/users/api_token",
			resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		return "", newStatusError(resp)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading HTTP response: %w", err)
	}
	body := string(b)
	if debug {
		log.Println("\nBODY: " + body)
	}
	return body, nil
}

func (c *Client) CreateObservation(obs Observation) error {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(CreateObservation{
		Observation: obs,
	})
	if err != nil {
		return fmt.Errorf("CreateObservation: %w", err)
	}
	req, err := http.NewRequest("POST", c.baseURL+"/observations", buf)
	if err != nil {
		return fmt.Errorf("CreateObservation: %w", err)
	}
	_, err = c.roundTrip(req)
	if err != nil {
		return fmt.Errorf("CreateObservation: %w", err)
	}
	log.Printf("Created %s\n", obs.URLWithSpecies())
	return nil
}

func (c *Client) UpdateObservation(obs Observation) error {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(UpdateObservation{
		IgnorePhotos: true, // don't clobber photos!
		Observation:  obs,
	})
	if err != nil {
		return fmt.Errorf("UpdateObservation: %w", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/observations/%s", c.baseURL, obs.UUID), buf)
	if err != nil {
		return fmt.Errorf("UpdateObservation: %w", err)
	}
	_, err = c.roundTrip(req)
	if err != nil {
		return fmt.Errorf("UpdateObservation: %w", err)
	}
	log.Printf("Updated %s\n", obs.URLWithSpecies())
	return nil
}

func (c *Client) DeleteObservation(id uuid.UUID) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/observations/%s", c.baseURL, id), nil)
	if err != nil {
		return fmt.Errorf("DeleteObservation: %w", err)
	}
	_, err = c.roundTrip(req)
	if err != nil {
		return fmt.Errorf("DeleteObservation: %w", err)
	}
	log.Printf("Deleted http://inaturalist.org/observations/%s\n", id)
	return nil
}

func (c *Client) UploadMedia(filename string, isPhoto bool, mlAssetID string, obsUUID string) error {
	destFilename := "ML" + mlAssetID + path.Ext(filename)
	var fieldName string
	var postURL string
	if isPhoto {
		fieldName = "observation_photo[observation_id]"
		postURL = c.baseURL + "/observation_photos"
		log.Println("Uploading photo as", destFilename)
	} else {
		fieldName = "observation_sound[observation_id]"
		postURL = c.baseURL + "/observation_sounds"
		log.Println("Uploading sound as", destFilename)
	}
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	fileWriter, err := writer.CreateFormFile("file", destFilename)
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(fileWriter, f)
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}
	err = writer.WriteField(fieldName, obsUUID)
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}

	req, err := http.NewRequest("POST", postURL, &requestBody)
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}
	// Set the Content-Type header to the multipart writer's boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_, err = c.roundTrip(req)
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}
	// TODO: log the media URL from the response body
	return nil
}

// maxErrorBody caps how much of a failed response is kept. iNaturalist's
// explanations are short; anything longer is probably an HTML error page.
const maxErrorBody = 512

// StatusError is returned when the API answers with a non-200 status. It
// carries the response body, because iNaturalist explains its refusals there
// and the status line alone cannot distinguish a file that is too large from
// one in an unsupported format (tech.errorsIncludeBody).
type StatusError struct {
	StatusCode int
	Status     string
	Body       string
}

func newStatusError(resp *http.Response) *StatusError {
	e := &StatusError{StatusCode: resp.StatusCode, Status: resp.Status}
	// A failure to read the body must not mask the failure being reported.
	b, err := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
	if err != nil {
		return e
	}
	body := strings.TrimSpace(string(b))
	// An HTML error page comes from a proxy in front of the API and says
	// nothing the status line doesn't, while spilling markup across several
	// log lines. A 413 from nginx arrives this way.
	if body == "" || strings.HasPrefix(body, "<") {
		return e
	}
	// Collapse the message onto one line so a log stays readable.
	e.Body = strings.Join(strings.Fields(body), " ")
	return e
}

func (e *StatusError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("bad HTTP status: %s", e.Status)
	}
	return fmt.Sprintf("bad HTTP status: %s: %s", e.Status, e.Body)
}

// Permanent reports whether retrying is pointless because the service rejected
// the request itself rather than the attempt.
//
// 4xx means the request was unacceptable, with three exceptions that a later
// run can get past: 401 needs a fresh token, 408 and 429 are explicit
// invitations to try again. Everything else, including every 5xx and every
// network error, is treated as transient — retrying costs a request, whereas
// wrongly giving up loses a user's photo (product.permanentUploadFailureRecorded).
func (e *StatusError) Permanent() bool {
	switch e.StatusCode {
	case http.StatusUnauthorized, http.StatusRequestTimeout, http.StatusTooManyRequests:
		return false
	}
	return e.StatusCode >= 400 && e.StatusCode < 500
}
