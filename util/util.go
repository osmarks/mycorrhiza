package util

import (
	"crypto/rand"
	"encoding/hex"
	"hash/fnv"
	"log/slog"
	"net/http"
	"strings"
	"time"
	insecureRand "math/rand"

	"github.com/hashicorp/golang-lru/v2"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/files"

	"git.sr.ht/~bouncepaw/mycomarkup/v5/util"
)

// PrepareRq strips the trailing / in rq.URL.Path. In the future it might do more stuff for making all request structs uniform.
func PrepareRq(rq *http.Request) {
	rq.URL.Path = strings.TrimSuffix(rq.URL.Path, "/")
}

// ShorterPath is used by handlerList to display shorter path to the files. It
// simply strips the hyphae directory name.
func ShorterPath(path string) string {
	if strings.HasPrefix(path, files.HyphaeDir()) {
		tmp := strings.TrimPrefix(path, files.HyphaeDir())
		if tmp == "" {
			return ""
		}
		return tmp[1:]
	}
	return path
}

// HTTP404Page writes a 404 error in the status, needed when no content is found on the page.
// TODO: demolish
func HTTP404Page(w http.ResponseWriter, page string) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte(page))
}

// HTTP200Page wraps some frequently used things for successful 200 responses.
// TODO: demolish
func HTTP200Page(w http.ResponseWriter, page string) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(page))
}

// RandomString generates a random string of the given length. It is cryptographically secure to some extent.
func RandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// BeautifulName makes the ugly name beautiful by replacing _ with spaces and using title case.
func BeautifulName(uglyName string) string {
	// Why not reuse
	result := util.BeautifulName(uglyName)
	for i, replace := range cfg.ReplaceFrom {
		result = strings.ReplaceAll(result, replace, cfg.ReplaceTo[i])
	}
	return result
}

// CanonicalName makes sure the `name` is canonical. A name is canonical if it is lowercase and all spaces are replaced with underscores.
func CanonicalName(name string) string {
	return util.CanonicalName(name)
}

// IsProfileName if the given hypha name is a profile name. It takes configuration into consideration.
//
// With default configuration, u/ is the prefix such names have. For example, u/wikimind matches. Note that u/wikimind/sub does not.
func IsProfileName(hyphaName string) bool {
	return strings.HasPrefix(hyphaName, cfg.UserHypha+"/") && strings.Count(hyphaName, "/") == 1
}

// HyphaNameFromRq extracts hypha name from http request. You have to also pass the action which is embedded in the url or several actions. For url /hypha/hypha, the action would be "hypha".
func HyphaNameFromRq(rq *http.Request, actions ...string) string {
	p := rq.URL.Path
	for _, action := range actions {
		if strings.HasPrefix(p, "/"+action+"/") {
			return CanonicalName(strings.TrimPrefix(p, "/"+action+"/"))
		}
	}
	slog.Info("HyphaNameFromRq: this request is invalid, fall back to home hypha")
	return cfg.HomeHypha
}

// FormData is a convenient struct for passing user input and errors to HTML
// forms and showing to the user.
type FormData struct {
	err    error
	fields map[string]string
}

// NewFormData constructs empty form data instance.
func NewFormData() FormData {
	return FormData{
		err:    nil,
		fields: map[string]string{},
	}
}

// FormDataFromRequest extracts a form data from request, using a set of keys.
func FormDataFromRequest(r *http.Request, keys []string) FormData {
	formData := NewFormData()
	for _, key := range keys {
		formData.Put(key, r.FormValue(key))
	}
	return formData
}

// HasError is true if there is indeed an error.
func (f FormData) HasError() bool {
	return f.err != nil
}

// Error returns an error text or empty string, if there are no errors in form data.
func (f FormData) Error() string {
	if f.err == nil {
		return ""
	}
	return f.err.Error()
}

// WithError puts an error into form data and returns itself.
func (f FormData) WithError(err error) FormData {
	f.err = err
	return f
}

// Get accesses form data with a key
func (f FormData) Get(key string) string {
	return f.fields[key]
}

// Put writes a form value for provided key
func (f FormData) Put(key, value string) {
	f.fields[key] = value
}

// IsRevHash checks if the revision hash is valid.
func IsRevHash(revHash string) bool {
	if len(revHash) < 7 {
		return false
	}
	paddedRevHash := revHash
	if len(paddedRevHash)%2 != 0 {
		paddedRevHash = paddedRevHash[:len(paddedRevHash)-1]
	}
	if _, err := hex.DecodeString(paddedRevHash); err != nil {
		return false
	}
	return true
}

func GetMotd() string {
	now := time.Now().UTC().Unix()
	dayIndex := now / 86400
	return cfg.Motds[dayIndex%int64(len(cfg.Motds))]
}

func RequestHeaderFingerprint(rq *http.Request) uint64 {
	fprintHeaders := []string{"accept", "accept-encoding", "accept-language", "dnt", "host", "user-agent", "x-tls-fp"}
	hasher := fnv.New64()

	for _, hdr := range fprintHeaders {
		if value := rq.Header.Get(hdr); value != "" {
			hasher.Write([]byte(hdr))
			hasher.Write([]byte(value))
		}
	}

	return hasher.Sum64()
}

var ipLookup *lru.Cache[string, uint64]
var fprintLookup *lru.Cache[uint64, uint64]

func EstimateTrackingIdentifier(rq *http.Request) uint64 {
	if ipLookup == nil || fprintLookup == nil {
		var err error
		ipLookup, err = lru.New[string, uint64](1<<20)
		if err != nil {
			slog.Error("cache create failed?", "error", err)
		}
		// golang...
		fprintLookup, err = lru.New[uint64, uint64](1<<20)
		if err != nil {
			slog.Error("cache create failed?", "error", err)
		}
	}

	ip := strings.Split(rq.RemoteAddr, ":")[0]
	if val := rq.Header.Get("x-forwarded-for"); val != "" {
		ip = val
	}

	fp := RequestHeaderFingerprint(rq)

	// Try to look up by IP address. If that fails, look up by fingerprint, and update the record for that IP address.
	// If that also fails, assign a random identifier to both.
	id := insecureRand.Uint64()
	if cid, ok := ipLookup.Get(ip); ok {
		id = cid
	} else {
		if cid, ok := fprintLookup.Get(fp); ok {
			id = cid
			ipLookup.Add(ip, cid)
		} else {
			ipLookup.Add(ip, id)
			fprintLookup.Add(fp, id)
		}
	}

	return id
}

var visitHistory *lru.Cache[uint64, []string]
const maxTraceLen int = 64

func EnsureVisitHistoryExists() {
	if visitHistory == nil {
		var err error
		visitHistory, err = lru.New[uint64, []string](1<<18)
		if err != nil {
			slog.Error("cache create failed?", "error", err)
		}
	}
}

func WriteTrace(rq *http.Request, hyphaName string) {
	EnsureVisitHistoryExists()

	trk := EstimateTrackingIdentifier(rq)
	trace := []string{}
	if htrace, ok := visitHistory.Get(trk); ok {
		trace = htrace
	}
	trace = append(trace, hyphaName)
	if len(trace) > maxTraceLen {
		trace = trace[1:]
	}
	visitHistory.Add(trk, trace)
}

func ReadTrace(rq *http.Request) []string {
	EnsureVisitHistoryExists()

	trk := EstimateTrackingIdentifier(rq)
	trace := []string{}
	if htrace, ok := visitHistory.Get(trk); ok {
		trace = htrace
	}

	return trace
}
