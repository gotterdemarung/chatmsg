package chatmsg

// ResultLink contains pair of URL with it's title
type ResultLink struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// Result holds general parsing result
// It can be converted into JSON using json.Marshall or
// any other marshaller
type Result struct {
	Mentions  []string     `json:"mentions"`
	Emoticons []string     `json:"emoticons"`
	Links     []ResultLink `json:"links"`
}
