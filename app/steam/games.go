package steam

import (
	"context"
	"fmt"
	"net/http"
)

// GamesService handles communication with the Game related
// methods of the Steam API.
type GamesService service

type AppDetailsResponse struct {
	Success bool `json:"success,omitempty"`
	Data    Game `json:"data,omitempty"`
}

type Game struct {
	Type                string       `json:"type,omitempty"`
	Name                string       `json:"name,omitempty"`
	SteamAppID          int64        `json:"steam_appid,omitempty"`
	RequiredAge         int          `json:"required_age,omitempty"`
	IsFree              bool         `json:"is_free,omitempty"`
	DetailedDescription string       `json:"detailed_description,omitempty"`
	AboutTheGame        string       `json:"about_the_game,omitempty"`
	ShortDescription    string       `json:"short_description,omitempty"`
	SupportedLanguages  string       `json:"supported_languages,omitempty"`
	HeaderImage         string       `json:"header_image,omitempty"`
	CapsuleImage        string       `json:"capsule_image,omitempty"`
	CapsuleImagev5      string       `json:"capsule_imagev5,omitempty"`
	Website             string       `json:"website,omitempty"`
	PCRequirements      Requirements `json:"pc_requirements,omitempty"`
	MacRequirements     Requirements `json:"mac_requirements,omitempty"`

	// LinuxRequirements is currently not active.
	// The reason is the game "" (see https://store.steampowered.com/api/appdetails?appids=3900)
	// The JSON states `"linux_requirements":[],` which causes an error when decoding into the struct.
	// We expect a Requirements struct with Minimum and Recommended fields.
	// The API provides an empty array/list instead.
	// I don't know if this is a one-off issue or if this happens more often.
	// However, for our use case right now, we can live without this field.
	//
	// LinuxRequirements   Requirements  `json:"linux_requirements,omitempty"`

	Developers    []string      `json:"developers,omitempty"`
	Publishers    []string      `json:"publishers,omitempty"`
	PriceOverview PriceOverview `json:"price_overview,omitempty"`
	Platforms     Platforms     `json:"platforms,omitempty"`
	Categories    []Category    `json:"categories,omitempty"`
	Genres        []Genre       `json:"genres,omitempty"`
	Screenshots   []Screenshort `json:"screenshots,omitempty"`
	Movies        []Movie       `json:"movies,omitempty"`
	ReleaseDate   ReleaseDate   `json:"release_date,omitempty"`
	Background    string        `json:"background,omitempty"`
	BackgroundRaw string        `json:"background_raw,omitempty"`

	// Missing fields
	// - packages
	// - package_groups
	// - recommendations
	// - achievements
	// - support_info
	// - content_descriptors
	// - ratings
}

type Requirements struct {
	Minimum     string `json:"minimum,omitempty"`
	Recommended string `json:"recommended,omitempty"`
}

type PriceOverview struct {
	Currency         string `json:"currency,omitempty"`
	Initial          int    `json:"initial,omitempty"`
	Final            int    `json:"final,omitempty"`
	DiscountPercent  int    `json:"discount_percent,omitempty"`
	InitialFormatted string `json:"initial_formatted,omitempty"`
	FinalFormatted   string `json:"final_formatted,omitempty"`
}

type Platforms struct {
	Windows bool `json:"windows,omitempty"`
	Mac     bool `json:"mac,omitempty"`
	Linux   bool `json:"linux,omitempty"`
}

type Category struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

type Genre struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

type Screenshort struct {
	ID            int    `json:"id,omitempty"`
	PathThumbnail string `json:"path_thumbnail,omitempty"`
	PathFull      string `json:"path_full,omitempty"`
}

type Movie struct {
	ID        int         `json:"id,omitempty"`
	Name      string      `json:"name,omitempty"`
	Thumbnail string      `json:"thumbnail,omitempty"`
	Webm      MovieFormat `json:"webm,omitempty"`
	Mp4       MovieFormat `json:"mp4,omitempty"`
	Highlight bool        `json:"highlight,omitempty"`
}

type MovieFormat struct {
	Size480 string `json:"480,omitempty"`
	SizeMax string `json:"max,omitempty"`
}

type ReleaseDate struct {
	ComingSoon bool   `json:"coming_soon,omitempty"`
	Date       string `json:"date,omitempty"`
}

// Get a single Game by Steam ID.
func (s *GamesService) GetAppDetails(ctx context.Context, steamID int64, language string) (*Game, *http.Response, error) {
	u := fmt.Sprintf("appdetails?appids=%d", steamID)
	req, err := s.client.NewRequest("GET", u, language, nil)
	if err != nil {
		return nil, nil, err
	}

	apiResponse := make(map[string]AppDetailsResponse)
	resp, err := s.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return nil, resp, err
	}

	var game *Game
	if v, ok := apiResponse[fmt.Sprintf("%d", steamID)]; ok {
		game = &v.Data
	} else {
		return nil, resp, fmt.Errorf("no data found for game with id %d", steamID)
	}

	return game, resp, nil
}
