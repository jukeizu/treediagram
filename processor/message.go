package processor

type Message struct {
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
