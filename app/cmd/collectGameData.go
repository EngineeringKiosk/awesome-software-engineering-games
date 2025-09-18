package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	libIO "github.com/EngineeringKiosk/awesome-software-engineering-games/io"
	"github.com/EngineeringKiosk/awesome-software-engineering-games/steam"
)

const (
	imageFolder      = "images"
	defaultUserAgent = "EngineeringKiosk-awesome-software-engineering-games"
)

// collectGameDataCmd represents the collectGameData command
var collectGameDataCmd = &cobra.Command{
	Use:   "collectGameData",
	Short: "Collects additional data per game",
	Long: `We only have basic data about each game.
To make the whole project more useful, we aim to collect additional data from the Game API.

This command gathers this additional data per Game and stores them back into
the generated JSON files.`,
	RunE: cmdCollectGameData,
}

func init() {
	rootCmd.AddCommand(collectGameDataCmd)

	collectGameDataCmd.Flags().String("json-directory", "", "Directory on where to store the json files")

	err := collectGameDataCmd.MarkFlagRequired("json-directory")
	if err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
}

func cmdCollectGameData(cmd *cobra.Command, args []string) error {
	jsonDir, err := cmd.Flags().GetString("json-directory")
	if err != nil {
		return err
	}

	log.Printf("Reading files with extension %s from directory %s", libIO.JSONExtension, jsonDir)
	jsonFiles, err := libIO.GetAllFilesFromDirectory(jsonDir, libIO.JSONExtension)
	if err != nil {
		return err
	}
	log.Printf("%d files found with extension %s in directory %s", len(jsonFiles), libIO.JSONExtension, jsonDir)

	steamClient := steam.NewClient(nil)

	for _, f := range jsonFiles {
		absJsonFilePath := filepath.Join(jsonDir, f.Name())
		log.Printf("Processing file %s", absJsonFilePath)
		jsonFileContent, err := os.ReadFile(absJsonFilePath)
		if err != nil {
			return err
		}

		gameInfo := &GameInformation{}
		err = json.Unmarshal(jsonFileContent, gameInfo)
		if err != nil {
			return err
		}

		if gameInfo.SteamID > 0 {
			// Get Game info
			log.Printf("Requesting 'Games.GetAppDetails' data [EN] from Steam API for game %s (ID: %d) ...", gameInfo.Name, gameInfo.SteamID)
			gameEnglish, _, err := steamClient.Games.GetAppDetails(context.Background(), gameInfo.SteamID, "en-US")
			if err != nil {
				return err
			}
			log.Printf("Requesting 'Games.GetAppDetails' data [EN] from Steam API for game %s (ID: %d) ... successful", gameInfo.Name, gameInfo.SteamID)

			log.Printf("Requesting 'Games.GetAppDetails' data [DE] from Steam API for game %s (ID: %d) ...", gameInfo.Name, gameInfo.SteamID)
			gameGerman, _, err := steamClient.Games.GetAppDetails(context.Background(), gameInfo.SteamID, "de-DE")
			if err != nil {
				return err
			}
			log.Printf("Requesting 'Games.GetAppDetails' data [DE] from Steam API for game %s (ID: %d) ... successful", gameInfo.Name, gameInfo.SteamID)

			// Set basic game data
			gameInfo.RequiredAge = gameEnglish.RequiredAge
			gameInfo.IsFree = gameEnglish.IsFree
			gameInfo.Website = gameEnglish.Website
			gameInfo.Platforms = Platforms{
				Windows: gameEnglish.Platforms.Windows,
				Mac:     gameEnglish.Platforms.Mac,
				Linux:   gameEnglish.Platforms.Linux,
			}
			gameInfo.ReleaseDate = ReleaseDate{
				ComingSoon: gameEnglish.ReleaseDate.ComingSoon,
				Date:       gameEnglish.ReleaseDate.Date,
			}
			gameInfo.GermanContent = getLanguageContent(gameGerman)
			gameInfo.EnglishContent = getLanguageContent(gameEnglish)

			// Download cover-image
			imageFileExtension := path.Ext(gameEnglish.HeaderImage)
			// Sometimes we have file extensions like .png?t=1655195362
			// but we only want .png
			if strings.Contains(imageFileExtension, "?") {
				imageFileExtension, _, _ = strings.Cut(imageFileExtension, "?")
			}
			jsonFileExtension := path.Ext(f.Name())
			imageFileName := f.Name()[0:len(f.Name())-len(jsonFileExtension)] + imageFileExtension
			absImageFilePath := filepath.Join(jsonDir, imageFolder, imageFileName)
			gameImageName := imageFileName

			// Sometimes gameEnglish.HeaderImage is empty.
			imageExistsAlready := false
			if len(gameEnglish.HeaderImage) == 0 {
				gameImageName, imageExistsAlready = libIO.DoesImageExistsOnDisk(absImageFilePath, false)
			}
			if len(gameEnglish.HeaderImage) == 0 && imageExistsAlready {
				log.Println("Skipping downloading new version of game image, because there is no image to download")
				log.Println("The pipeline didn't fail, because a previous downloaded image exists")

			} else {
				log.Printf("Downloading %s into %s ...", gameEnglish.HeaderImage, absImageFilePath)
				_, err = downloadFile(gameEnglish.HeaderImage, absImageFilePath)
				if err != nil {
					// Sometimes we get errors like
					// Error: Get "<URL>": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
					log.Printf("Downloading %s into %s ... error: %v", gameEnglish.HeaderImage, absImageFilePath, err)

					// If we get an error, but we have a target image already
					// (like an old one), it is better to use the old one than failing.
					//
					// Having the latest up to date image is not the highest priority here.

					var oldExists bool
					gameImageName, oldExists = libIO.DoesImageExistsOnDisk(absImageFilePath, true)
					if oldExists {
						log.Printf("We were not able to download the new image %s", gameEnglish.HeaderImage)
						log.Printf("The pipeline didn't fail, because the previous version %s exists", absImageFilePath)

					} else {
						return err
					}
				} else {
					log.Printf("Downloading %s into %s ... successful", gameEnglish.HeaderImage, absImageFilePath)
				}
			}
			gameInfo.Image = filepath.Join(imageFolder, gameImageName)
		} else {
			log.Printf("Skipping data retrieval from Steam for %s, because SteamID is %d", absJsonFilePath, gameInfo.SteamID)
		}

		// Write the information back to the JSON file
		// Dump data into JSON file
		log.Printf("Write %s to disk ...", absJsonFilePath)
		err = libIO.WriteJSONFile(absJsonFilePath, gameInfo)
		if err != nil {
			return err
		}
		log.Printf("Write %s to disk ... successful", absJsonFilePath)
	}

	return nil
}

func downloadFile(address, fileName string) (*http.Response, error) {
	client := &http.Client{
		Timeout: 45 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 30 * time.Second,
		},
	}

	req, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", defaultUserAgent)
	response, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != 200 {
		return response, fmt.Errorf("received %d as status code, expected 200", response.StatusCode)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return response, err
	}
	defer func() { _ = file.Close() }()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return response, err
	}

	return response, nil
}

func getLanguageContent(game *steam.Game) LanguageContent {
	content := LanguageContent{
		ShortDescription: game.ShortDescription,
		Categories:       []string{},
		Genres:           []string{},
	}

	for _, c := range game.Categories {
		content.Categories = append(content.Categories, c.Description)
	}
	for _, g := range game.Genres {
		content.Genres = append(content.Genres, g.Description)
	}

	return content
}
