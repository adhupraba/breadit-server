package queryhelpers

import (
	"context"
	"errors"
	"net/http"

	rawmessageparser "github.com/adhupraba/breadit-server/internal/helpers/raw_message_parser"
	"github.com/adhupraba/breadit-server/internal/types"
	"github.com/adhupraba/breadit-server/lib"
)

func FindCommentsAndRepliesOfAPost(ctx context.Context, postId int32) (commentsWithReplies []types.CommentWithReplies, err error, errCode int) {
	comments, err := lib.DB.FindCommentsOfAPost(ctx, postId)

	if err != nil {
		return []types.CommentWithReplies{}, errors.New("Error fetching comments"), http.StatusInternalServerError
	}

	commentIds := []int32{}

	for _, comment := range comments {
		commentIds = append(commentIds, comment.Comment.ID)
	}

	replies, err := lib.DB.FindRepliesForComments(ctx, commentIds)

	if err != nil {
		return []types.CommentWithReplies{}, errors.New("Error fetching comments"), http.StatusInternalServerError
	}

	commentToRepliesMap := make(map[int32][]types.CommentWithAuthorAndVotes)

	for _, reply := range replies {
		parentCommentId := reply.Comment.ReplyToID.Int32

		votes, err := rawmessageparser.ParseJsonCommentVotes(reply.Votes)

		if err != nil {
			return []types.CommentWithReplies{}, errors.New("Error parsing comment votes"), http.StatusInternalServerError
		}

		replyComment := types.CommentWithAuthorAndVotes{
			Comment: reply.Comment,
			Author:  reply.User,
			Votes:   votes,
		}

		commentToRepliesMap[parentCommentId] = append(commentToRepliesMap[parentCommentId], replyComment)
	}

	commentsWithReplies = make([]types.CommentWithReplies, len(comments))

	for idx, comment := range comments {
		votes, err := rawmessageparser.ParseJsonCommentVotes(comment.Votes)

		if err != nil {
			return []types.CommentWithReplies{}, errors.New("Error parsing comment votes"), http.StatusInternalServerError
		}

		replies := commentToRepliesMap[comment.Comment.ID]

		if replies == nil {
			replies = []types.CommentWithAuthorAndVotes{}
		}

		parentComment := types.CommentWithReplies{
			CommentWithAuthorAndVotes: types.CommentWithAuthorAndVotes{
				Comment: comment.Comment,
				Author:  comment.User,
				Votes:   votes,
			},
			Replies: replies,
		}

		commentsWithReplies[idx] = parentComment
	}

	return commentsWithReplies, nil, 0
}
