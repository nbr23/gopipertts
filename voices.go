package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// These describe the voices.json file

type Voice struct {
	Key          string                 `json:"key"`
	Name         string                 `json:"name"`
	Language     Language               `json:"language"`
	Quality      string                 `json:"quality"`
	NumSpeakers  int                    `json:"num_speakers"`
	SpeakerIDMap map[string]interface{} `json:"speaker_id_map"`
	Files        map[string]File        `json:"files"`
	Aliases      []string               `json:"aliases"`
}

type Language struct {
	Code           string `json:"code"`
	Family         string `json:"family"`
	Region         string `json:"region"`
	NameNative     string `json:"name_native"`
	NameEnglish    string `json:"name_english"`
	CountryEnglish string `json:"country_english"`
}

type File struct {
	SizeBytes int64  `json:"size_bytes"`
	MD5Digest string `json:"md5_digest"`
}

type Voices map[string]Voice

// This describes the individual voice json files

type VoiceDetails struct {
	Audio        VoiceDetailsAudio `json:"audio"`
	SpeakerIdMap map[string]int    `json:"speaker_id_map"`
}

type VoiceDetailsAudio struct {
	SampleRate int `json:"sample_rate"`
}

var DOWNLOADED_VOICES = make(map[string]VoiceDetails)

func loadVoicesDetails() {
	if _, err := os.Stat(VOICES_PATH); os.IsNotExist(err) {
		os.Mkdir(VOICES_PATH, 0755)
	}

	files, err := os.ReadDir(VOICES_PATH)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		voiceName := file.Name()
		if filepath.Ext(voiceName) != ".json" {
			continue
		}
		voiceName = voiceName[:len(voiceName)-5]
		if filepath.Ext(voiceName) != ".onnx" {
			continue
		}

		if _, err := os.Stat(fmt.Sprintf("%s/%s", VOICES_PATH, voiceName)); os.IsNotExist(err) {
			continue
		}

		voiceName = voiceName[:len(voiceName)-5]

		voice, err := parseVoiceDetails(fmt.Sprintf("%s/%s", VOICES_PATH, file.Name()))
		if err != nil {
			log.Println("Failed to get voice details: ", err)
			continue
		}

		DOWNLOADED_VOICES[voiceName] = voice
	}
}

func getVoiceDetails(voices *Voices, voiceName string) (VoiceDetails, error) {
	voice, ok := DOWNLOADED_VOICES[voiceName]
	if !ok {
		err := downloadVoiceFiles(voices, voiceName)
		if err != nil {
			return VoiceDetails{}, err
		}
		voice = DOWNLOADED_VOICES[voiceName]
	}
	return voice, nil
}

func parseVoiceDetails(filePath string) (VoiceDetails, error) {
	voice := VoiceDetails{}
	file, err := os.Open(filePath)
	if err != nil {
		return voice, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&voice)
	if err != nil {
		return voice, err
	}
	return voice, nil
}

func getAvailableVoices(voicesPath string) Voices {
	voices := make(Voices)

	file, err := os.Open(voicesPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&voices)
	if err != nil {
		log.Fatal(err)
	}
	return voices

}

func downloadVoiceFiles(voices *Voices, voiceName string) error {
	if _, ok := DOWNLOADED_VOICES[voiceName]; ok {
		return nil
	}

	voice, ok := (*voices)[voiceName]
	if !ok {
		return fmt.Errorf("Voice not found: %s", voiceName)
	}

	for fileName, _ := range voice.Files {
		baseFilename := filepath.Base(fileName)
		if baseFilename == "MODEL_CARD" {
			continue
		}
		filePath := fmt.Sprintf("%s/%s", VOICES_PATH, baseFilename)
		if _, err := os.Stat(filePath); err == nil {
			log.Println("File already exists, skipping download", filePath)
			continue
		}

		url := fmt.Sprintf("%s/%s", VOICES_REPO_BASE_URL, fileName)
		log.Println("Downloading", url, "to", filePath)

		out, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer out.Close()

		resp, err := http.Get(url)
		if err != nil {
			log.Println("Failed to download", url, ":", err)
			return err
		}
		defer resp.Body.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
		log.Println("Downloaded ", url, " to ", filePath)
	}
	voiceDetails, err := parseVoiceDetails(fmt.Sprintf("%s/%s.onnx.json", VOICES_PATH, voiceName))
	if err != nil {
		return err
	}
	DOWNLOADED_VOICES[voiceName] = voiceDetails
	return nil
}
