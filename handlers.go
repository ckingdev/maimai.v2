package maimai

import (
	"fmt"
	"net/http"
	"regexp"
	// "strconv"
	"strings"
	// "time"

	"euphoria.io/heim/proto"
	// "github.com/boltdb/bolt"
	"github.com/cpalone/gobot"
	"golang.org/x/net/html"
)

var linkMatcher = regexp.MustCompile("(https?://)?[\\S]+\\.[\\S][\\S]+[\\S^\\.]")

// type SeenHandler struct {
// 	seenTime map[string]time.Time
// }

// func SetDBKey(db *bolt.DB, bucket, key, val []byte) error {
// 	return db.Update(func(tx *bolt.Tx) error {
// 		b := tx.Bucket(bucket)
// 		return b.Put(key, val)
// 	})
// }

// func IsSeenCommand(msg string) bool {
// 	if strings.HasPrefix(msg, "!seen @") && len(msg) > 7 {
// 		return true
// 	}
// 	return false
// }

// func GetDBKey(db *bolt.DB, bucket, key []byte) ([]byte, error) {
// 	var val []byte
// 	err := db.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket(bucket)
// 		val = b.Get(key)
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return val, nil
// }

// func (sh *SeenHandler) HandleIncoming(r *gobot.Room, p *proto.Packet) (*proto.Packet, error) {
// 	if p.Type != proto.SendEventType {
// 		return nil, nil
// 	}
// 	payload, err := p.Payload()
// 	if err != nil {
// 		return nil, err
// 	}
// 	msg, ok := payload.(*proto.SendEvent)
// 	if !ok {
// 		return nil, fmt.Errorf("Could not assert payload as *proto.SendEvent.")
// 	}
// 	if err := SetDBKey(r.DB,
// 		[]byte("seen"),
// 		[]byte(msg.Sender.Name),
// 		[]byte(strconv.FormatInt(time.Now().Unix(), 10))); err != nil {
// 		return err
// 	}
// 	if !IsSeenCommand(msg.Content) {
// 		return nil, nil
// 	}
// 	user := msg.Content[7:]
// 	t, err := GetDBKey(r.DB, []byte("seen"), []byte(msg.Sender.Name))
// 	if err != nil {
// 		return nil, err
// 	}
// 	var reply string
// 	if t == nil {
// 		reply = "User has not been seen yet."
// 	} else {
// 		lastSeenInt, _ := strconv.Atoi(string(t))
// 		lastSeenTime := time.Unix(int64(lastSeenInt, 0))
// 		since := time.Since(lastSeenTime)
// 		reply = fmt.Sprintf("Seen %v hours ago.", int(since.Hours()))
// 	}
// 	repPayload := proto.SendCommand{Parent: msg.ID, Content: reply}
// 	repPacket, err := gobot.MakePacket(proto.SendType, repPayload)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return repPacket, nil
// }

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
				if title == "Imgur" {
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

func (l *LinkTitleHandler) HandleIncoming(r *gobot.Room, p *proto.Packet) (*proto.Packet, error) {
	if p.Type != proto.SendEventType {
		return nil, nil
	}
	payload, err := p.Payload()
	if err != nil {
		return nil, err
	}
	msg, ok := payload.(*proto.SendEvent)
	if !ok {
		return nil, fmt.Errorf("Could not assert SendEvent as such.")
	}
	urls := linkMatcher.FindAllString(msg.Content, -1)
	for _, url := range urls {
		if !strings.HasPrefix(url, "http") {
			url = "http://" + url
		}
		title, err := getLinkTitle(url)
		if err != nil && title != "" {
			r.SendText(&msg.ID, "Link title: "+title)
			break
		}
	}
	return nil, nil
}
