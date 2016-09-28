package chatmsg

import (
	"errors"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

var mentionsRegex = regexp.MustCompile(`(?i)@(\w+)`)
var emoticonsRegex = regexp.MustCompile(`(?i)\((\w+)\)`)

// urlRegex is RegExp, used to parse URLs
// it is very far from real world examples, cyrillic/emoji domains, etc
var urlRegex = regexp.MustCompile(`(?i)((https?)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?)`)
var titleRegex = regexp.MustCompile(`(?i)<title>(.*?)<\/title>`)

// Parse func parses incoming message using regular expressions and return Result
// or error, if any
// This is proof-of-concept example, real application must optimize Regexes or
// even switch to char-by-char parsing to optimize performance
func Parse(msg string, urlToTitle func(string) (string, error)) (*Result, error) {
	res := new(Result)
	if len(msg) == 0 {
		return res, nil
	}
	if urlToTitle == nil {
		return nil, errors.New("No URL to title reader provided")
	}

	for _, matches := range mentionsRegex.FindAllStringSubmatch(msg, -1) {
		res.Mentions = append(res.Mentions, matches[1])
	}
	for _, matches := range emoticonsRegex.FindAllStringSubmatch(msg, -1) {
		res.Emoticons = append(res.Emoticons, matches[1])
	}

	wg := sync.WaitGroup{}
	m := sync.Mutex{}
	var err error

	for _, matches := range urlRegex.FindAllStringSubmatch(msg, -1) {
		wg.Add(1)
		go func(url string) {
			title, e := urlToTitle(url)
			if e == nil {
				m.Lock()
				res.Links = append(res.Links, ResultLink{
					URL:   url,
					Title: title,
				})
				m.Unlock()
			} else {
				err = e
			}
			wg.Done()
		}(matches[1])
	}

	wg.Wait()
	return res, err
}

// ParseSimple performs simple parsing using simple HTTP client
// It is not recommended to use it in production mode, leave it
// for testing purpose only
func ParseSimple(msg string) (*Result, error) {
	return Parse(msg, httpRead)
}

// httpRead performs simple HTTP GET request and grabs TITLE using regex
// this is simple proof-of-concept, real reader must provide user agent
// and other headers and parse OpenGraph and other meta data
// BTW real reader must implement caching and metrics/monitoring
func httpRead(url string) (string, error) {
	cl := http.Client{}
	resp, err := cl.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if matches := titleRegex.FindStringSubmatch(string(bts)); len(matches) > 0 {
		return html.UnescapeString(matches[1]), nil
	}

	return "", nil
}
