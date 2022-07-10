package crunchyroll

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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

type Review struct {
	crunchy *Crunchyroll

	SeriesID string

	Review struct {
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

func (r *Review) IsOwner() bool {
	return r.crunchy.Config.AccountID == r.Author.ID
}

func (r *Review) Edit(title, content string, spoiler bool) error {
	if !r.IsOwner() {
		return fmt.Errorf("cannot edit, current user is not the review author")
	}
	endpoint := fmt.Sprintf("https://beta.crunchyroll.com/content-reviews/v2/en-US/user/%s/review/series/%s", r.crunchy.Config.AccountID, r.SeriesID)
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
	resp, err := r.crunchy.requestFull(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(r)

	return nil
}

func (r *Review) Delete() error {
	if !r.IsOwner() {
		return fmt.Errorf("cannot delete, current user is not the review author")
	}
	endpoint := fmt.Sprintf("https://beta.crunchyroll.com/content-reviews/v2/en-US/user/%s/review/series/%s", r.crunchy.Config.AccountID, r.SeriesID)
	_, err := r.crunchy.request(endpoint, http.MethodDelete)
	return err
}

func (r *Review) Helpful() error {
	return r.rate(true)
}

func (r *Review) NotHelpful() error {
	return r.rate(false)
}

func (r *Review) rate(positive bool) error {
	if r.Ratings.Rating != "" {
		var humanReadable string
		switch r.Ratings.Rating {
		case "yes":
			humanReadable = "helpful"
		case "no":
			humanReadable = "not helpful"
		}
		return fmt.Errorf("review is already marked as %s", humanReadable)
	}

	endpoint := fmt.Sprintf("https://beta.crunchyroll.com/content-reviews/v2/user/%s/rating/review/%s", r.crunchy.Config.AccountID, r.Review.ID)
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
	_, err = r.crunchy.requestFull(req)
	return err
}

func (r *Review) Report() error {
	if r.Ratings.Reported {
		return fmt.Errorf("review is already reported")
	}
	endpoint := fmt.Sprintf("https://beta.crunchyroll.com/content-reviews/v2/user/%s/report/review/%s", r.crunchy.Config.AccountID, r.Review.ID)
	_, err := r.crunchy.request(endpoint, http.MethodPut)
	return err
}

func (r *Review) RemoveReport() error {
	if !r.Ratings.Reported {
		return fmt.Errorf("review is not reported")
	}
	endpoint := fmt.Sprintf("https://beta.crunchyroll.com/content-reviews/v2/user/%s/report/review/%s", r.crunchy.Config.AccountID, r.Review.ID)
	_, err := r.crunchy.request(endpoint, http.MethodDelete)
	return err
}
