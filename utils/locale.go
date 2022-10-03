package utils

import (
	"github.com/crunchy-labs/crunchyroll-go/v3"
)

// AllLocales is an array of all available locales.
var AllLocales = []crunchyroll.LOCALE{
	crunchyroll.JP,
	crunchyroll.US,
	crunchyroll.LA,
	crunchyroll.LA2,
	crunchyroll.ES,
	crunchyroll.FR,
	crunchyroll.PT,
	crunchyroll.BR,
	crunchyroll.IT,
	crunchyroll.DE,
	crunchyroll.RU,
	crunchyroll.AR,
	crunchyroll.ME,
	crunchyroll.CN,
}

// ValidateLocale validates if the given locale actually exist.
func ValidateLocale(locale crunchyroll.LOCALE) bool {
	for _, l := range AllLocales {
		if l == locale {
			return true
		}
	}
	return false
}

// LocaleLanguage returns the country by its locale.
func LocaleLanguage(locale crunchyroll.LOCALE) string {
	switch locale {
	case crunchyroll.JP:
		return "Japanese"
	case crunchyroll.US:
		return "English (US)"
	case crunchyroll.LA, crunchyroll.LA2:
		return "Spanish (Latin America)"
	case crunchyroll.ES:
		return "Spanish (Spain)"
	case crunchyroll.FR:
		return "French"
	case crunchyroll.PT:
		return "Portuguese (Europe)"
	case crunchyroll.BR:
		return "Portuguese (Brazil)"
	case crunchyroll.IT:
		return "Italian"
	case crunchyroll.DE:
		return "German"
	case crunchyroll.RU:
		return "Russian"
	case crunchyroll.AR, crunchyroll.ME:
		return "Arabic"
	case crunchyroll.CN:
		return "Chinese (China)"
	default:
		return ""
	}
}
