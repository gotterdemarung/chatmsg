package chatmsg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	assert := assert.New(t)

	res, err := ParseSimple("Good morning @chris! (megusta) (coffee)")
	assert.NotNil(res)
	assert.NoError(err)

	assert.Len(res.Mentions, 1)
	assert.Equal("chris", res.Mentions[0])

	assert.Len(res.Emoticons, 2)
	assert.Equal("megusta", res.Emoticons[0])
	assert.Equal("coffee", res.Emoticons[1])

	res, err = ParseSimple("Olympics are starting soon; http://www.nbcolympics.com")
	assert.NotNil(res)
	assert.NoError(err)

	assert.Len(res.Mentions, 0)
	assert.Len(res.Emoticons, 0)
	assert.Len(res.Links, 1)
	assert.Equal("http://www.nbcolympics.com", res.Links[0].URL)
	assert.Equal("2016 Rio Olympic Games | NBC Olympics", res.Links[0].Title)

	res, err = ParseSimple("@bob @john (success) such a cool feature;\nhttps://twitter.com/jdorfman/status/430511497475670016")
	assert.NotNil(res)
	assert.NoError(err)
	assert.Len(res.Mentions, 2)
	assert.Equal("bob", res.Mentions[0])
	assert.Equal("john", res.Mentions[1])
	assert.Len(res.Emoticons, 1)
	assert.Equal("success", res.Emoticons[0])
	assert.Len(res.Links, 1)
	assert.Equal("https://twitter.com/jdorfman/status/430511497475670016", res.Links[0].URL)
	assert.Equal("Justin Dorfman on Twitter: \"nice @littlebigdetail from @HipChat (shows hex colors when pasted in chat). http://t.co/7cI6Gjy5pq\"", res.Links[0].Title)
}
