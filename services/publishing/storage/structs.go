package storage

type Request struct {
	CorrelationId  string `json:"correlationId"`
	ChannelId      string `json:"channelId"`
	User           User   `json:"user"`
	PrivateMessage bool   `json:"privateMessage"`
	IsRedirect     bool   `json:"isRedirect"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type MessageRequest struct {
	Request
	Id        string `json:"id"`
	MessageId string `json:"messageId"`
}

type Message struct {
	Id      string  `json:"id"`
	Content string  `json:"content"`
	Embed   *Embed  `json:"embed,omitempty"`
	Tts     bool    `json:"tts"`
	Files   []*File `json:"files,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconUrl      string `json:"iconUrl,omitempty"`
	ProxyIconUrl string `json:"proxyIconUrl,omitempty"`
}

type EmbedImage struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxyUrl,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type EmbedThumbnail struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxyUrl,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type EmbedVideo struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxyUrl,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type EmbedProvider struct {
	Url  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

type EmbedAuthor struct {
	Url          string `json:"url,omitempty"`
	Name         string `json:"name,omitempty"`
	IconUrl      string `json:"iconUrl,omitempty"`
	ProxyIconUrl string `json:"proxyIconUrl,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type Embed struct {
	Url         string          `json:"url,omitempty"`
	Type        string          `json:"type,omitempty"`
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	Timestamp   string          `json:"timestamp,omitempty"`
	Color       int             `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []*EmbedField   `json:"fields,omitempty"`
}

type File struct {
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	Bytes       []byte `json:"bytes"`
}
