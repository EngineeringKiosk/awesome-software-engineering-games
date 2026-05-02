package cmd

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/EngineeringKiosk/awesome-software-engineering-games/io"
)

// defaultGameImage is the placeholder image used when a non-Steam game's YAML
// file does not specify an 'image' field. The path is expressed relative to the
// 'generated/' folder, matching the convention used by manually authored entries
// such as games/artifacts.yml.
const defaultGameImage = "../assets/dummy-image.png"

// convertYamlToJsonCmd represents the convertYamlToJson command
var convertYamlToJsonCmd = &cobra.Command{
	Use:   "convertYamlToJson",
	Short: "Converts Game YAML files into JSON files",
	Long: `The YAML representation of the basic game info is more for humans.
For machines we have a JSON format with more information about the game available.

This command converts the basic YAML information into JSON format.`,
	RunE: cmdConvertYamlToJson,
}

func init() {
	rootCmd.AddCommand(convertYamlToJsonCmd)

	convertYamlToJsonCmd.Flags().String("yaml-directory", "", "Directory on where to find the yaml files")
	convertYamlToJsonCmd.Flags().String("json-directory", "", "Directory on where to store the json files")

	err := convertYamlToJsonCmd.MarkFlagRequired("yaml-directory")
	if err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
	err = convertYamlToJsonCmd.MarkFlagRequired("json-directory")
	if err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
	convertYamlToJsonCmd.MarkFlagsRequiredTogether("yaml-directory", "json-directory")
}

func cmdConvertYamlToJson(cmd *cobra.Command, args []string) error {
	yamlDir, err := cmd.Flags().GetString("yaml-directory")
	if err != nil {
		return err
	}

	jsonDir, err := cmd.Flags().GetString("json-directory")
	if err != nil {
		return err
	}

	yamlExts := io.GetYAMLExtensions()
	log.Printf("Reading files with extensions %v from directory %s", yamlExts, yamlDir)
	yamlFiles, err := io.GetAllFilesFromDirectoryWithExtensions(yamlDir, yamlExts)
	if err != nil {
		return err
	}
	log.Printf("%d files found with extensions %v in directory %s", len(yamlFiles), yamlExts, yamlDir)

	log.Printf("Reading files with extension %s from directory %s", io.JSONExtension, jsonDir)
	jsonFiles, err := io.GetAllFilesFromDirectory(jsonDir, io.JSONExtension)
	if err != nil {
		return err
	}
	log.Printf("%d files found with extension %s in directory %s", len(jsonFiles), io.JSONExtension, jsonDir)

	// Process every YAML file found and dump it into a JSON
	// file with the same name.
	// If the JSON file already exists, merge it and only update the data
	// that is available in the YAML file.
	for _, f := range yamlFiles {
		absYamlFilePath := filepath.Join(yamlDir, f.Name())
		log.Printf("Processing file %s", absYamlFilePath)
		yamlFileContent, err := os.ReadFile(absYamlFilePath)
		if err != nil {
			return err
		}

		gameInfo := &GameInformation{}
		err = yaml.Unmarshal(yamlFileContent, gameInfo)
		if err != nil {
			return err
		}

		currentFileExtension := path.Ext(f.Name())
		jsonFileName := f.Name()[0:len(f.Name())-len(currentFileExtension)] + io.JSONExtension
		absJsonFilePath := filepath.Join(jsonDir, jsonFileName)

		log.Printf("Converting %s to %s", absYamlFilePath, absJsonFilePath)

		// Check if we have a related json file already
		if _, ok := jsonFiles[jsonFileName]; ok {
			// JSON file exists.
			// Read JSON file into Game Information structure
			// and overwrite yaml information
			jsonFileContent, err := os.ReadFile(absJsonFilePath)
			if err != nil {
				return err
			}

			gameJsonInfo := &GameInformation{}
			err = json.Unmarshal(jsonFileContent, gameJsonInfo)
			if err != nil {
				return err
			}

			gameInfo = mergeGameInformation(gameInfo, gameJsonInfo)
		}

		// Add generated fields
		gameInfo.Slug = slug.Make(gameInfo.Name)

		// Non-Steam games (SteamID == 0) carry their cover artwork via the YAML
		// 'image' field. If that field is missing, fall back to a checked-in
		// placeholder so the pipeline can complete, and surface a warning so the
		// gap is visible in CI output.
		if gameInfo.SteamID == 0 && len(gameInfo.Image) == 0 {
			log.Printf("WARNING: %s has no 'image' field set in YAML; falling back to default %s. Add a real cover at games/images/<slug>.<ext> and set 'image: ../games/images/<slug>.<ext>' in the YAML to replace it.", absYamlFilePath, defaultGameImage)
			gameInfo.Image = defaultGameImage
		}

		// Dump data into JSON file
		log.Printf("Write %s to disk ...", absJsonFilePath)
		err = io.WriteJSONFile(absJsonFilePath, gameInfo)
		if err != nil {
			return err
		}
		log.Printf("Write %s to disk ... successful", absJsonFilePath)
	}

	log.Printf("Converting of YAML to JSON ... successful")
	return nil
}

// mergeGameInformation will overwrite a fixed set of
// fields from source into target.
func mergeGameInformation(source, target *GameInformation) *GameInformation {
	// Those fields are all fields where the yaml file is the source of truth.
	// If the yaml structure will be changesd, this function needs to be updated as well.
	//
	// Maybe there is a smarter implementation of it (reflection?) but for now
	// this is good enough.
	//
	// If you change the code below by adding / removing fields,
	// please update CONTRIBUTING.md as well.
	target.Name = source.Name
	target.SteamID = source.SteamID
	// Only overwrite website if we have one in the source.
	if len(source.Website) > 0 {
		target.Website = source.Website
	}
	target.Repository = source.Repository
	target.Programmable = source.Programmable
	// target.License is intentionally not touched: it is enriched from the
	// GitHub API in collectGameData and YAML is not a source of truth for it.

	if source.SteamID == 0 {
		// If there is no SteamID set, we do not have any platform or release date information.
		// Thus we take it from the source (yaml).
		target.Image = source.Image
		target.Platforms = Platforms{
			Windows: source.Platforms.Windows,
			Mac:     source.Platforms.Mac,
			Linux:   source.Platforms.Linux,
		}
		target.ReleaseDate = ReleaseDate{
			Date: source.ReleaseDate.Date,
		}
		target.EnglishContent = LanguageContent{
			ShortDescription: source.EnglishContent.ShortDescription,
			Genres:           source.EnglishContent.Genres,
		}
		target.GermanContent = LanguageContent{
			ShortDescription: source.GermanContent.ShortDescription,
			Genres:           source.GermanContent.Genres,
		}
	}

	return target
}
