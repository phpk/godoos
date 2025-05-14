package payload

import "time"

type User struct {
	StarredURL        string    `json:"starred_url"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Name              string    `json:"name"`
	Remark            string    `json:"remark"`
	FollowingURL      string    `json:"following_url"`
	Bio               string    `json:"bio"`
	CreatedAt         time.Time `json:"created_at"`
	GistsURL          string    `json:"gists_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	Weibo             string    `json:"weibo"`
	ID                int64     `json:"id"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	Email             string    `json:"email"`
	AvatarURL         string    `json:"avatar_url"`
	EventsURL         string    `json:"events_url"`
	Watched           int       `json:"watched"`
	Blog              string    `json:"blog"`
	Stared            int       `json:"stared"`
	UpdatedAt         time.Time `json:"updated_at"`
	ReposURL          string    `json:"repos_url"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	Login             string    `json:"login"`
	URL               string    `json:"url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
}
