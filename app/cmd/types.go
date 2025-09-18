package cmd

import (
	"strings"
	"time"
)

type GameInformation struct {
	Name         string `yaml:"name" json:"name"`
	SteamID      int64  `yaml:"steamID" json:"steamID"`
	Slug         string `json:"slug"`
	Repository   string `yaml:"repository" json:"repository,omitempty"`
	Programmable bool   `yaml:"programmable" json:"programmable,omitempty"`

	Website     string      `yaml:"website" json:"website,omitempty"`
	RequiredAge int         `yaml:"required_age" json:"required_age,omitempty"`
	IsFree      bool        `yaml:"is_free" json:"is_free,omitempty"`
	Platforms   Platforms   `yaml:"platforms" json:"platforms,omitempty"`
	ReleaseDate ReleaseDate `yaml:"release_date" json:"release_date,omitempty"`
	Image       string      `yaml:"image" json:"image,omitempty"`

	// German
	GermanContent  LanguageContent `yaml:"german_content" json:"german_content,omitempty"`
	EnglishContent LanguageContent `yaml:"english_content" json:"english_content,omitempty"`
}

type LanguageContent struct {
	ShortDescription string   `yaml:"short_description" json:"short_description,omitempty"`
	Categories       []string `yaml:"categories" json:"categories,omitempty"`
	Genres           []string `yaml:"genres" json:"genres,omitempty"`
}

type Platforms struct {
	Windows bool `yaml:"windows" json:"windows"`
	Mac     bool `yaml:"mac" json:"mac"`
	Linux   bool `yaml:"linux" json:"linux"`
}

type ReleaseDate struct {
	ComingSoon bool   `yaml:"coming_soon" json:"coming_soon"`
	Date       string `yaml:"date" json:"date,omitempty"`
}

func (g GameInformation) GetReleaseDate() string {
	layouts := []string{
		"2 Jan, 2006", // e.g. 10 Sep, 2024
		"Jan 2, 2006", // e.g. Aug 23, 2018
	}
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, g.ReleaseDate.Date)
		if err == nil {
			return t.Format("Monday, 02. January 2006")
		}
	}

	return ""
}

func (g GameInformation) GenresAsList() string {
	s := ""
	if len(g.EnglishContent.Genres) > 0 {
		s = strings.Join(g.EnglishContent.Genres, ", ")
	}

	return s
}
