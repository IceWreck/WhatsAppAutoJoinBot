package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/Rhymen/go-whatsapp"
)

// whitelist of topics which should be in title to join
var whitelist = []string{
	"quiz",
	"assignment",
	"solution",
	"answer",
	"worksheet",
	"study",
	"sessional",
	"exam",
	"midsem",
	"endsem",
	"semester",
}

// whitelist of topics which should be in title and should not be a subword
var whitelist_word = []string{
	"mst",
	"est",
	"ests",
	"msts",
}

// regex for whatsapp invite link
var inviteRegex = regexp.MustCompile(`(https?:\/\/)?chat\.whatsapp\.com\/(?:invite\/)?([a-zA-Z0-9_-]{22})`)

// GroupMetaData is needed to parse metadata json (memebers, subject, creator, etc)
// type GroupMetaData struct {
// 	Subject string `json:"subject"`
// }

type waHandler struct {
	wac       *whatsapp.Conn
	startTime uint64
}

func (wh *waHandler) HandleError(err error) {
	// log to stdout but not to the main logs cause these are expected errors
	log.Println("unhandled message:", err)
}

// HandleTextMessage receives whatsapp text messages
func (wh *waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	if message.Info.FromMe || message.Info.Timestamp < wh.startTime {
		return
	}

	// find if message has invite link
	indices := inviteRegex.FindStringIndex(message.Text)
	if indices == nil {
		log.Println("ignoring message", message.Info.RemoteJid)
		return
	}
	// seperate out invite code fromm message
	inviteLink := message.Text[indices[0]:indices[1]]
	inviteList := strings.Split(inviteLink, "/")
	inviteCode := inviteList[len(inviteList)-1]

	// we cant view metadata before joining from api but whatsapp website can do that
	// so scrape subject/title from there
	webPage, err := http.Get(inviteLink)
	if err != nil {
		return
	}
	webBody, err := ioutil.ReadAll(webPage.Body)
	if err != nil {
		return
	}
	groupSubject := getStringInBetween(string(webBody), `<meta property="og:title" content="`, `" />`)
	newLog("group subject is ", groupSubject)

	shouldJoin := false
	for _, item := range whitelist {
		if strings.Contains(strings.ToLower(groupSubject), item) {
			shouldJoin = true
			break
		}
	}

	for _, word := range strings.Fields(groupSubject) {
		for _, item := range whitelist_word {
			if strings.ToLower(word) == item {
				shouldJoin = true
				break
			}
		}
	}

	// join group
	if shouldJoin {

		newLog("invite code is", inviteCode)
		jid, err := wh.wac.GroupAcceptInviteCode(inviteCode)
		if err != nil {
			newLog("cant join group ", err, jid)
			return
		}
	}
	// View Metadata
	// mdata, err := wh.wac.GetGroupMetaData(jid)
	// if err != nil {
	// 	NewLog("error while getting group metadata", err)
	// 	return
	// }
	// fmt.Println(<-mdata)

}

func getStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return ""
	}
	return str[s : s+e]
}
