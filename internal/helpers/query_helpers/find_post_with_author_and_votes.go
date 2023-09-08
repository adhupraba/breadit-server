package queryhelpers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	rawmessageparser "github.com/adhupraba/breadit-server/internal/helpers/raw_message_parser"
	"github.com/adhupraba/breadit-server/internal/types"
	"github.com/adhupraba/breadit-server/lib"
)

func FindPostWithAuthorAndVotes(ctx context.Context, postId int32) (postData types.PostWithAuthorAndVotes, err error, errCode int) {
	post, err := lib.DB.FindPostWithAuthorAndVotes(ctx, postId)

	if err != nil && strings.Contains(err.Error(), "no rows") {
		return types.PostWithAuthorAndVotes{}, errors.New("Post not found"), http.StatusNotFound
	}

	if err != nil {
		return types.PostWithAuthorAndVotes{}, errors.New("Error when fetching post details"), http.StatusInternalServerError
	}

	votes, err := rawmessageparser.ParseJsonVotes(post.Votes)

	if err != nil {
		return types.PostWithAuthorAndVotes{}, err, http.StatusInternalServerError
	}

	transformedPost := types.PostWithAuthorAndVotes{
		Post:   post.Post,
		Author: post.User,
		Votes:  votes,
	}

	return transformedPost, nil, 0
}
