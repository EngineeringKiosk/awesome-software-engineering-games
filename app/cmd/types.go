package cmd

import (
	"strings"
	"time"
)

type GameInformation struct {
	Name    string `yaml:"name" json:"name"`
	SteamID int64  `yaml:"steamID" json:"steamID"`
	Slug    string `json:"slug"`

	Website     string      `json:"website,omitempty"`
	RequiredAge int         `json:"required_age,omitempty"`
	IsFree      bool        `json:"is_free,omitempty"`
	Platforms   Platforms   `json:"platforms,omitempty"`
	ReleaseDate ReleaseDate `json:"release_date,omitempty"`
	Image       string      `json:"image,omitempty"`

	// German
	GermanContent  LanguageContent `json:"german_content,omitempty"`
	EnglishContent LanguageContent `json:"english_content,omitempty"`
}

type LanguageContent struct {
	ShortDescription string   `json:"short_description,omitempty"`
	Categories       []string `json:"categories,omitempty"`
	Genres           []string `json:"genres,omitempty"`
}

type Platforms struct {
	Windows bool `json:"windows"`
	Mac     bool `json:"mac"`
	Linux   bool `json:"linux"`
}

type ReleaseDate struct {
	ComingSoon bool   `json:"coming_soon"`
	Date       string `json:"date,omitempty"`
}

func (g GameInformation) GetReleaseDate() string {
	t, err := time.Parse("2 Jan, 2006", g.ReleaseDate.Date)
	if err != nil {
		return ""
	}
	s := t.Format("Monday, 02 January 2006")

	return s
}

func (g GameInformation) GenresAsList() string {
	s := ""
	if len(g.EnglishContent.Genres) > 0 {
		s = strings.Join(g.EnglishContent.Genres, ", ")
	}

	return s
}
