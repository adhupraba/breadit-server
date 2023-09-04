package helpers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adhupraba/breadit-server/internal/database"
	"github.com/adhupraba/breadit-server/lib"
)

type PostWithData struct {
	database.Post
	Author    database.User      `json:"author"`
	Votes     []database.Vote    `json:"votes"`
	Comments  []database.Comment `json:"comments"`
	Subreddit database.Subreddit `json:"subreddit"`
}

type SubredditWithPosts struct {
	database.Subreddit
	MemberCount int64          `json:"memberCount"`
	Posts       []PostWithData `json:"posts"`
}

func GetSubredditWithPosts(ctx context.Context, subredditName string) (data SubredditWithPosts, err error, errCode int) {
	subreddit, err := lib.DB.FindSubredditByName(ctx, subredditName)

	if subreddit.ID == 0 {
		return SubredditWithPosts{}, fmt.Errorf("Subreddit not found."), http.StatusNotFound
	}

	if err != nil {
		return SubredditWithPosts{}, fmt.Errorf("Error when finding the subreddit."), http.StatusInternalServerError
	}

	subscriptionCount, err := lib.DB.FindSubscriptionCountOfSubreddit(ctx, subreddit.ID)

	if err != nil {
		return SubredditWithPosts{}, fmt.Errorf("Error when calculating the subscribers count."), http.StatusInternalServerError
	}

	posts, err := lib.DB.FindPostsOfSubredditWithAuthor(ctx, database.FindPostsOfSubredditWithAuthorParams{
		SubredditID: subreddit.ID,
		Offset:      0,
		Limit:       10,
	})

	if err != nil {
		return SubredditWithPosts{}, fmt.Errorf("Error while fetching posts of the subreddit."), http.StatusInternalServerError
	}

	postIds := []int32{}
	for _, post := range posts {
		postIds = append(postIds, post.Post.ID)
	}

	comments, err := lib.DB.FindCommentsOfPosts(ctx, postIds)

	if err != nil {
		return SubredditWithPosts{}, fmt.Errorf("Error while fetching comments of the posts."), http.StatusInternalServerError
	}

	votes, err := lib.DB.FindVotesOfPosts(ctx, postIds)

	if err != nil {
		return SubredditWithPosts{}, fmt.Errorf("Error while fetching votes of the posts."), http.StatusInternalServerError
	}

	commentsOfPosts := make(map[int32][]database.Comment)

	for _, comment := range comments {
		commentsOfPosts[comment.PostID] = append(commentsOfPosts[comment.PostID], comment)
	}

	votesOfPosts := make(map[int32][]database.Vote)

	for _, vote := range votes {
		votesOfPosts[vote.PostID] = append(votesOfPosts[vote.PostID], vote)
	}

	var postsWithData []PostWithData

	for _, post := range posts {
		vp := votesOfPosts[post.Post.ID]
		cp := commentsOfPosts[post.Post.ID]

		votes := make([]database.Vote, len(vp))

		if len(vp) != 0 {
			votes = vp
		}

		comments := make([]database.Comment, len(cp))

		if len(cp) != 0 {
			comments = cp
		}

		postWithData := PostWithData{
			Post:      post.Post,
			Author:    post.User,
			Votes:     votes,
			Comments:  comments,
			Subreddit: subreddit,
		}

		postsWithData = append(postsWithData, postWithData)
	}

	subredditWithPosts := SubredditWithPosts{
		Subreddit:   subreddit,
		MemberCount: subscriptionCount,
		Posts:       postsWithData,
	}

	return subredditWithPosts, nil, 0
}
