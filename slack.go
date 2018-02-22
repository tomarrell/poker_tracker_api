package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

type Slacker struct {
	client       *slack.Client
	channel      string
	testChannel  string
	movioRealmID int
}

func mustNewSlacker(token, channel, testChannel string, realmId int) *Slacker {
	return &Slacker{
		slack.New(token),
		channel,
		testChannel,
		realmId,
	}
}

func (s *Slacker) SendSummary(balances []balanceSummary, pSessions []PlayerSession, pNames map[int]string, realmId int) {
	sStr := ""

	for _, s := range pSessions {
		sStr = fmt.Sprintf("%s%s: %+.2f\n", sStr, pNames[s.PlayerID], float64(s.Walkout.Int64-s.Buyin.Int64)/100)
	}

	bStr := ""
	for _, b := range balances {
		bStr = fmt.Sprintf("%s%s: %+.2f\n", bStr, b.PlayerName, float64(b.Total)/100)
	}

	today := fmt.Sprintf(`
Today:
--------
%s
`, sStr)

	total := fmt.Sprintf(`
Total:
--------
%s
`, bStr)

	msg := fmt.Sprintf("%s%s%s%s",
		"```",
		today,
		strings.TrimRight(total, " "),
		"```",
	)
	log.Info(msg)

	s.send(msg, realmId)
}

func (s *Slacker) send(msg string, realmId int) {
	ch := s.testChannel
	if realmId == s.movioRealmID {
		ch = s.channel
	}

	params := slack.NewPostMessageParameters()
	params.IconURL = "https://profitrobot.me/dist/profitrobot-avatar__0acf03fa61e3b19634b7fe01e25ab495.png"
	params.Username = "pokerbot"
	_, _, err := s.client.PostMessage(
		ch,
		msg,
		params,
	)
	fmt.Println(s.channel)
	if err != nil {
		log.Error(err)
	}
}
