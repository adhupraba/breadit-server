package queryhelpers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adhupraba/breadit-server/constants"
	"github.com/adhupraba/breadit-server/internal/database"
	rawmessageparser "github.com/adhupraba/breadit-server/internal/helpers/raw_message_parser"
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

	params := database.FindPostsOfSubredditParams{
		Offset:        0,
		Limit:         constants.InfiniteScrollPaginationResults,
		SubredditName: subredditName,
	}
	// fmt.Printf("GetSubredditWithPosts db params => %#v\n", params)

	postsWithData, err, errCode := GetPostsOfSubreddit(ctx, params)

	if err != nil {
		return SubredditWithPosts{}, err, errCode
	}

	subredditWithPosts := SubredditWithPosts{
		Subreddit:   subreddit,
		MemberCount: subscriptionCount,
		Posts:       postsWithData,
	}

	return subredditWithPosts, nil, 0
}

func GetPostsOfSubreddit(ctx context.Context, params database.FindPostsOfSubredditParams) (data []PostWithData, err error, errCode int) {
	posts, err := lib.DB.FindPostsOfSubreddit(ctx, params)

	if err != nil {
		fmt.Println("error getting posts of subreddit =>", err)
		return []PostWithData{}, fmt.Errorf("Error while fetching posts of the subreddit."), http.StatusInternalServerError
	}

	postsWithData := make([]PostWithData, len(posts))

	for idx, post := range posts {
		votes, err := rawmessageparser.ParseJsonVotes(post.Votes)

		if err != nil {
			return data, err, http.StatusInternalServerError
		}

		comments, err := rawmessageparser.ParseJsonComments(post.Comments)

		if err != nil {
			return data, err, http.StatusInternalServerError
		}

		postWithData := PostWithData{
			Post:      post.Post,
			Author:    post.User,
			Votes:     votes,
			Comments:  comments,
			Subreddit: post.Subreddit,
		}

		postsWithData[idx] = postWithData
	}

	return postsWithData, nil, 0
}
