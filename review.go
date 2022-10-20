package crunchyroll

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ReviewRating represents stars for a series rating from one to five.
type ReviewRating string

const (
	OneStar    ReviewRating = "s1"
	TwoStars                = "s2"
	ThreeStars              = "s3"
	FourStars               = "s4"
	FiveStars               = "s5"
)

type ratingStar struct {
	Displayed  string `json:"displayed"`
	Unit       string `json:"unit"`
	Percentage int    `json:"percentage"`
}

// Rating represents the overall rating of a series.
type Rating struct {
	OneStar    ratingStar `json:"1s"`
	TwoStars   ratingStar `json:"2s"`
	ThreeStars ratingStar `json:"3s"`
	FourStars  ratingStar `json:"4s"`
	FiveStars  ratingStar `json:"5s"`
	Average    string     `json:"average"`
	Total      int        `json:"total"`
	Rating     string     `json:"rating"`
}

// Review is the interface which gets implemented by OwnerReview and UserReview.
type Review interface{}

type review struct {
	crunchy *Crunchyroll

	SeriesID string

	ReviewData struct {
		ID              string    `json:"id"`
		Title           string    `json:"title"`
		Body            string    `json:"body"`
		Language        LOCALE    `json:"language"`
		CreatedAt       time.Time `json:"created_at"`
		ModifiedAt      time.Time `json:"modified_at"`
		AuthoredReviews int       `json:"authored_reviews"`
		Spoiler         bool      `json:"spoiler"`
	} `json:"review"`
	AuthorRating ReviewRating `json:"author_rating"`
	Author       struct {
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
		ID       string `json:"ID"`
	} `json:"author"`
	Ratings struct {
		Yes struct {
			Displayed string `json:"displayed"`
			Unit      string `json:"unit"`
		} `json:"yes"`
		No struct {
			Displayed string `json:"displayed"`
			Unit      string `json:"unit"`
		} `json:"no"`
		Total string `json:"total"`
		// yes or no so basically a bool if set
		Rating   string `json:"rating"`
		Reported bool   `json:"reported"`
	} `json:"ratings"`
}

// OwnerReview is a series review which has been written from the current logged-in user.
type OwnerReview struct {
	Review

	*review
}

// Edit edits the review from the logged in account.
func (or *OwnerReview) Edit(title, content string, spoiler bool) error {
	endpoint := fmt.Sprintf("https://beta-api.crunchyroll.com/content-reviews/v2/en-US/user/%s/review/series/%s", or.crunchy.Config.AccountID, or.SeriesID)
	body, _ := json.Marshal(map[string]any{
		"title":   title,
		"body":    content,
		"spoiler": spoiler,
	})
	req, err := http.NewRequest(http.MethodPatch, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := or.crunchy.requestFull(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(or)

	return nil
}

// Delete deletes the review from the logged in account.
func (or *OwnerReview) Delete() error {
	endpoint := fmt.Sprintf("https://beta-api.crunchyroll.com/content-reviews/v2/en-US/user/%s/review/series/%s", or.crunchy.Config.AccountID, or.SeriesID)
	_, err := or.crunchy.request(endpoint, http.MethodDelete)
	return err
}

// UserReview is a series review written from other crunchyroll users.
type UserReview struct {
	Review

	*review
}

// RateHelpful rates the review as helpful. A review can only be rated once
// as helpful (or not helpful) and this cannot be undone, so be careful. Use
// Rated to see if the review was already rated.
func (ur *UserReview) RateHelpful() error {
	return ur.rate(true)
}

// RateNotHelpful rates the review as not helpful. A review can only be rated
// once as helpful (or not helpful) and this cannot be undone, so be careful.
// Use Rated to see if the review was already rated.
func (ur *UserReview) RateNotHelpful() error {
	return ur.rate(false)
}

// Rated returns if the user already rated the review (with RateHelpful or
// RateNotHelpful).
func (ur *UserReview) Rated() bool {
	return ur.Ratings.Rating != ""
}

func (ur *UserReview) rate(positive bool) error {
	if ur.Rated() {
		var humanReadable string
		switch ur.Ratings.Rating {
		case "yes":
			humanReadable = "helpful"
		case "no":
			humanReadable = "not helpful"
		}
		return fmt.Errorf("review is already rated as %s", humanReadable)
	}

	endpoint := fmt.Sprintf("https://beta-api.crunchyroll.com/content-reviews/v2/user/%s/rating/review/%s", ur.crunchy.Config.AccountID, ur.ReviewData.ID)
	var body []byte
	if positive {
		body, _ = json.Marshal(map[string]string{"rate": "yes"})
	} else {
		body, _ = json.Marshal(map[string]string{"rate": "no"})
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := ur.crunchy.requestFull(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&ur.Ratings)

	return nil
}

// Report reports the review. Only works if the review hasn't been reported yet.
// See UserReview.Ratings.Reported if it is already reported.
func (ur *UserReview) Report() error {
	if ur.Ratings.Reported {
		return fmt.Errorf("review is already reported")
	}
	endpoint := fmt.Sprintf("https://beta-api.crunchyroll.com/content-reviews/v2/user/%s/report/review/%s", ur.crunchy.Config.AccountID, ur.ReviewData.ID)
	_, err := ur.crunchy.request(endpoint, http.MethodPut)
	if err != nil {
		return err
	}

	ur.Ratings.Reported = true

	return nil
}

// RemoveReport removes the report request from the review. Only works if the user
// has reported the review. See UserReview.Ratings.Reported if it is already reported.
func (ur *UserReview) RemoveReport() error {
	if !ur.Ratings.Reported {
		return fmt.Errorf("review is not reported")
	}
	endpoint := fmt.Sprintf("https://beta-api.crunchyroll.com/content-reviews/v2/user/%s/report/review/%s", ur.crunchy.Config.AccountID, ur.ReviewData.ID)
	_, err := ur.crunchy.request(endpoint, http.MethodDelete)
	if err != nil {
		return err
	}

	ur.Ratings.Reported = false

	return nil
}
