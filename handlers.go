package maimai

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"euphoria.io/heim/proto"
	"github.com/cpalone/gobot"
	"golang.org/x/net/html"
)

var linkMatcher = regexp.MustCompile("(https?://)?[\\S]+\\.[\\S][\\S]+[\\S^\\.]")

var commandMatcher = regexp.MustCompile("![\\S]+ @[\\S]+")

func getCommandAndUser(text string) (string, string, error) {
	matches := commandMatcher.FindAllString(text, -1)
	if matches == nil {
		return "", "", fmt.Errorf("getCommandAndUser: no matches found")
	}
	// Take first match only
	s := matches[0]
	splits := strings.Split(s, " ")
	if len(splits) != 2 {
		return "", "", fmt.Errorf("getCommandAndUser: invalid command")
	}
	return splits[0], splits[1][1:], nil
}

func isCommand(text string) bool {
	if commandMatcher.FindAllString(text, -1) == nil {
		return false
	}
	return true
}

// LinkTitleHandler searches each SendEvent for urls and posts the link title
// (if available) for the first valid link it finds.
//
// It is an empty struct because it does not need to maintain state.
type LinkTitleHandler struct{}

func extractTitleFromTree(z *html.Tokenizer) string {
	depth := 0
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return ""
		case html.TextToken:
			if depth > 0 {
				title := strings.TrimSpace(string(z.Text()))
				lower := strings.ToLower(title)
				if strings.HasPrefix(lower, "imgur") {
					return ""
				}
				return title
			}
		case html.StartTagToken:
			tn, _ := z.TagName()
			if string(tn) == "title" {
				depth++
			}
		}
	}
}

func getLinkTitle(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Bad response code: %v", resp.StatusCode)
	}
	z := html.NewTokenizer(resp.Body)
	return extractTitleFromTree(z), nil
}

// HandleIncoming checks incoming SendEvents for URLs and posts the link title
// for the first URL returning a valid title.
func (l *LinkTitleHandler) HandleIncoming(r *gobot.Room, p *proto.Packet) (*proto.Packet, error) {
	if p.Type != proto.SendEventType {
		return nil, nil
	}
	r.Logger.Debugf("Handler received SendEvent")
	payload, err := p.Payload()
	if err != nil {
		return nil, err
	}
	msg, ok := payload.(*proto.SendEvent)
	if !ok {
		return nil, fmt.Errorf("Could not assert SendEvent as such.")
	}
	if msg.Sender.Name == "euphoriabot" {
		return nil, nil
	}
	r.Logger.Debugf("Received message with content: %s", msg.Content)
	urls := linkMatcher.FindAllString(msg.Content, -1)
	for _, url := range urls {
		r.Logger.Debugf("Trying URL %s", url)
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}
		title, err := getLinkTitle(url)
		if err == nil && title != "" {
			r.SendText(&msg.ID, "Link title: "+title)
			break
		}
	}
	return nil, nil
}

// Run is a no-op.
func (l *LinkTitleHandler) Run(r *gobot.Room) {
	return
}

// Stop is a no-op.
func (l *LinkTitleHandler) Stop(r *gobot.Room) {
	return
}
