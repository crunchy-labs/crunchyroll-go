package crunchyroll

import (
	"regexp"
)

// ParseSeriesURL tries to extract the season id of the given crunchyroll url, pointing to a season.
func ParseSeriesURL(url string) (seasonId string, ok bool) {
	pattern := regexp.MustCompile(`(?m)^https?://((www|beta)\.)?crunchyroll\.com/(\w{2}/)?series/(?P<seasonId>\w+).*`)
	if urlMatch := pattern.FindAllStringSubmatch(url, -1); len(urlMatch) != 0 {
		groups := regexGroups(urlMatch, pattern.SubexpNames()...)
		seasonId = groups["seasonId"]
		ok = true
	}
	return
}

// ParseEpisodeURL tries to extract the episode id of the given crunchyroll url, pointing to an episode.
func ParseEpisodeURL(url string) (episodeId string, ok bool) {
	pattern := regexp.MustCompile(`(?m)^https?://((www|beta)\.)?crunchyroll\.com/(\w{2}/)?watch/(?P<episodeId>\w+).*`)
	if urlMatch := pattern.FindAllStringSubmatch(url, -1); len(urlMatch) != 0 {
		groups := regexGroups(urlMatch, pattern.SubexpNames()...)
		episodeId = groups["episodeId"]
		ok = true
	}
	return
}
