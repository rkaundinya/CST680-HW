package schema

import "time"

type voterPoll struct {
	PollID   uint
	VoteDate time.Time
}

type Voter struct {
	VoterID     uint
	FirstName   string
	LastName    string
	VoteHistory []voterPoll
}

type pollOption struct {
	PollOptionID   uint
	PollOptionText string
}

type Poll struct {
	PollID       uint         `json:"id"`
	PollTitle    string       `json:"title"`
	PollQuestion string       `json:"question"`
	PollOptions  []pollOption `json:"options"`
}
