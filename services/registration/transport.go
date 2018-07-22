package registration

import "gopkg.in/mgo.v2/bson"

type Command struct {
	Id             bson.ObjectId `bson:"_id, omitempty"`
	Server         string
	Name           string
	Regex          string
	RequireMention bool
	Endpoint       string
	Help           string
	Enabled        bool
}

type CommandQuery struct {
	Server   string
	LastId   string
	PageSize int
}

type CommandQueryResult struct {
	Commands []Command
	HasMore  bool
}
