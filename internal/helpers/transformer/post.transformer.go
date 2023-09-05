package transformer

import (
	"github.com/adhupraba/breadit-server/internal/database"
	rawmessageparser "github.com/adhupraba/breadit-server/internal/helpers/raw_message_parser"
)

type PostWithAuthorAndVotes struct {
	database.Post
	Author database.User   `json:"author"`
	Votes  []database.Vote `json:"votes"`
}

func TransformPostWithAuthorAndVotes(post database.FindPostWithAuthorAndVotesRow) (PostWithAuthorAndVotes, error) {
	votes, err := rawmessageparser.ParseJsonVotes(post.Votes)

	if err != nil {
		return PostWithAuthorAndVotes{}, err
	}

	postFmt := PostWithAuthorAndVotes{
		Post:   post.Post,
		Author: post.User,
		Votes:  votes,
	}

	return postFmt, nil
}
