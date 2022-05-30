package backend

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/buger/jsonparser"
)

type UseCase struct{}

func (uc UseCase) Influencer(userName string) (Model, error) {
	userName = cleanUserName(userName)
	url := fmt.Sprintf(InstagramURL, userName)

	body, err := doRequest(url)
	if err != nil {
		return Model{}, err
	}

	return parseBody(body)
}

func cleanUserName(userName string) string {
	return strings.Trim(userName, "@")
}

func doRequest(url string) ([]byte, error) {
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		errWrap := fmt.Errorf("can't create a new request: %v", err)
		log.Print(errWrap)
		return nil, errWrap
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		errWrap := fmt.Errorf("can't get data from request: %v", err)
		log.Print(errWrap)
		return nil, errWrap
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Printf("couldn't close the body: %v", err)
		}
	}()

	if response.StatusCode != 200 {
		errWrap := fmt.Errorf("we get an unexpected status code: %d", response.StatusCode)
		return nil, errWrap
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		errWrap := fmt.Errorf("can't read reader interface to parse to []bytes: %v", err)
		log.Print(errWrap)

		return nil, errWrap
	}

	data := buf.Bytes()

	return data, nil
}

func parseBody(data []byte) (Model, error) {
	if string(data[:14]) == "<!DOCTYPE html" {
		return Model{}, fmt.Errorf("oh, looks like we are blocked ☹️")
	}

	title, _ := jsonparser.GetString(data, "title")
	if strings.EqualFold(title, "Restricted profile") {
		return Model{}, fmt.Errorf("we couldn't get this user's info. Sorry")
	}

	userName, err := jsonparser.GetString(data, "graphql", "user", "username")
	if err != nil {
		log.Printf("can't parse user name: %v", err)
	}

	followers, err := jsonparser.GetInt(data, "graphql", "user", "edge_followed_by", "count")
	if err != nil {
		log.Printf("can't parse followers: %v", err)
	}

	following, err := jsonparser.GetInt(data, "graphql", "user", "edge_follow", "count")
	if err != nil {
		log.Printf("can't parse following: %v", err)
	}

	biography, err := jsonparser.GetString(data, "graphql", "user", "biography")
	if err != nil {
		log.Printf("can't parse biography: %v", err)
	}

	photo, err := jsonparser.GetString(data, "graphql", "user", "profile_pic_url_hd")
	if err != nil {
		log.Printf("can't parse photo: %v", err)
	}

	amountPosts, err := jsonparser.GetInt(data, "graphql", "user", "edge_owner_to_timeline_media", "count")
	if err != nil {
		log.Printf("can't parse amount posts: %v", err)
	}

	isPrivate, err := jsonparser.GetBoolean(data, "graphql", "user", "is_private")
	if err != nil {
		log.Printf("can't parse is_private field: %v", err)
	}

	var posts Posts
	var avgLikes int64
	var avgComments int64
	var engagementRate float64

	// We can only read this information if the profile is public
	if !isPrivate {
		_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, arrErr error) {
			postPhoto, internalErr := jsonparser.GetString(value, "node", "display_url")
			if internalErr != nil {
				log.Printf("can't parse string post photo: %v", internalErr)
			}

			likes, internalErr := jsonparser.GetInt(value, "node", "edge_liked_by", "count")
			if internalErr != nil {
				log.Printf("can't parse int post likes: %v", internalErr)
			}

			comments, internalErr := jsonparser.GetInt(value, "node", "edge_media_to_comment", "count")
			if internalErr != nil {
				log.Printf("can't parse int post comments: %v", internalErr)
			}

			post := Post{
				Photo:    postPhoto,
				Likes:    likes,
				Comments: comments,
			}

			posts = append(posts, post)
		}, "graphql", "user", "edge_owner_to_timeline_media", "edges")
		if err != nil {
			log.Printf("can't read posts data: %v", err)
		}

		avgLikes = posts.AvgLikes()
		avgComments = posts.AvgComments()
		engagementRate = float64(avgLikes+avgComments) / float64(followers) * 100
	}

	model := Model{
		Username:       userName,
		Photo:          photo,
		Biography:      biography,
		Followers:      followers,
		Following:      following,
		AmountPosts:    amountPosts,
		IsPrivate:      isPrivate,
		EngagementRate: engagementRate,
		AvgLikes:       avgLikes,
		AvgComments:    avgComments,
		Posts:          posts,
	}

	return model, nil
}
