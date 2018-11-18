package processor

import "github.com/jukeizu/treediagram/api/protobuf-spec/processing"

type Request struct {
	Id        string `json:"id"`
	Source    string `json:"source"`
	Bot       User   `json:"bot"`
	Author    User   `json:"author"`
	ChannelId string `json:"channelId"`
	ServerId  string `json:"serverId"`
	Mentions  []User `json:"mentions"`
	Content   string `json:"content"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func toRequest(pr *processing.TreediagramRequest) Request {
	r := Request{
		Source:    pr.Source,
		Bot:       toUser(pr.Bot),
		Author:    toUser(pr.Author),
		ChannelId: pr.ChannelId,
		ServerId:  pr.ServerId,
		Content:   pr.Content,
	}

	for _, m := range pr.Mentions {
		r.Mentions = append(r.Mentions, toUser(m))
	}

	return r
}

func toUser(pu *processing.User) User {
	u := User{
		Id:   pu.Id,
		Name: pu.Name,
	}

	return u
}
