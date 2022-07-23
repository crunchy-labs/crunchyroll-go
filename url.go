package crunchyroll

import (
	"fmt"
	"net/http"
)

// ExtractEpisodesFromUrl extracts all episodes from an url.
// If audio is not empty, the episodes gets filtered after the given locale.
func (c *Crunchyroll) ExtractEpisodesFromUrl(url string, audio ...LOCALE) ([]*Episode, error) {
	series, episodes, err := c.ParseUrl(url)
	if err != nil {
		return nil, err
	}

	var eps []*Episode

	if series != nil {
		seasons, err := series.Seasons()
		if err != nil {
			return nil, err
		}
		for _, season := range seasons {
			if audio != nil {
				locale, err := season.AudioLocale()
				if err != nil {
					return nil, err
				}

				var found bool
				for _, l := range audio {
					if locale == l {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			e, err := season.Episodes()
			if err != nil {
				return nil, err
			}
			eps = append(eps, e...)
		}
	} else if episodes != nil {
		if audio == nil {
			return episodes, nil
		}

		for _, episode := range episodes {
			locale, err := episode.AudioLocale()
			if err != nil {
				return nil, err
			}
			if audio != nil {
				var found bool
				for _, l := range audio {
					if locale == l {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			eps = append(eps, episode)
		}
	}

	if len(eps) == 0 {
		return nil, fmt.Errorf("could not find any matching episode")
	}

	return eps, nil
}

// ParseUrl parses the given url into a series or episode.
// The returning episode is a slice because non-beta urls have the same episode with different languages.
func (c *Crunchyroll) ParseUrl(url string) (*Series, []*Episode, error) {
	return parseUrl(c, url, false)
}

func parseUrl(crunchy *Crunchyroll, url string, recursive bool) (*Series, []*Episode, error) {
	if seriesId, ok := ParseBetaSeriesURL(url); ok {
		series, err := SeriesFromID(crunchy, seriesId)
		if err != nil {
			return nil, nil, err
		}
		return series, nil, nil
	} else if episodeId, ok := ParseBetaEpisodeURL(url); ok {
		episode, err := EpisodeFromID(crunchy, episodeId)
		if err != nil {
			return nil, nil, err
		}
		return nil, []*Episode{episode}, nil
	} else if seriesName, ok := ParseVideoURL(url); ok {
		if recursive {
			return nil, nil, fmt.Errorf("unexpected recursion for url %s", url)
		}

		oldRedirect := crunchy.Client.CheckRedirect
		crunchy.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		defer func() {
			crunchy.Client.CheckRedirect = oldRedirect
		}()
		resp, err := crunchy.request(url)
		if err != nil {
			return nil, nil, err
		}
		if redirectUrl := resp.Header.Get("Location"); redirectUrl != "" {
			return parseUrl(crunchy, redirectUrl, true)
		} else {
			video, err := crunchy.FindVideoByName(seriesName)
			if err != nil {
				return nil, nil, err
			}
			return video.(*Series), nil, nil
		}
	} else if seriesName, title, _, _, ok := ParseEpisodeURL(url); ok {
		if recursive {
			return nil, nil, fmt.Errorf("unexpected recursion for url %s", url)
		}

		oldRedirect := crunchy.Client.CheckRedirect
		crunchy.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		defer func() {
			crunchy.Client.CheckRedirect = oldRedirect
		}()
		resp, err := crunchy.request(url)
		if err != nil {
			return nil, nil, err
		}
		if redirectUrl := resp.Header.Get("Location"); redirectUrl != "" {
			return parseUrl(crunchy, redirectUrl, true)
		} else {
			episodes, err := crunchy.FindEpisodeByName(seriesName, title)
			if err != nil {
				return nil, nil, err
			}
			return nil, episodes, nil
		}
	} else {
		return nil, nil, fmt.Errorf("invalid url %s", url)
	}
}
