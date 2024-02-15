package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/abema/go-mp4"
	"github.com/golden-vcr/remix"
)

var titleFlag = flag.String("title", "", "User-facing title for the clip")
var tapeIdFlag = flag.Int("tape", 0, "ID of the tape this clip comes from")
var localFlag = flag.Bool("local", false, "Sync this clip to the local db")
var prodFlag = flag.Bool("prod", false, "Sync this clip to the production db")
var accessTokenFlag = flag.String("token", "", "Access token issued to the broadcaster")

func main() {
	// Parse flags
	flag.Parse()
	title := *titleFlag
	if title == "" {
		log.Fatalf("-title must be supplied")
	}
	tapeId := *tapeIdFlag
	if tapeId <= 0 {
		log.Fatal("-tape must be supplied")
	}
	isLocal := *localFlag
	isProd := *prodFlag
	if isLocal == isProd {
		log.Fatalf("one of -local or -prod must be supplied")
	}
	accessToken := *accessTokenFlag
	if accessToken == "" {
		log.Fatalf("-token must be supplied")
	}

	// Determine the URL of the remix server we want to target
	remixUrl := "http://localhost:5010"
	if isProd {
		remixUrl = "https://goldenvcr.com/api/remix"
	}

	// Require the clip filename as a single positional argument
	if flag.NArg() != 1 {
		log.Fatalf("clip filename must be supplied after flags")
	}
	clipPath := os.Args[len(os.Args)-1]

	// Verify that the input file exists and is an mp4
	if !strings.HasSuffix(clipPath, ".mp4") {
		log.Fatalf("clip file must be an .mp4")
	}
	stat, err := os.Stat(clipPath)
	if err != nil || stat.IsDir() {
		log.Fatalf("clip file not found: %s", clipPath)
	}

	// Open the clip for read so we can check its duration
	file, err := os.Open(clipPath)
	if err != nil {
		log.Fatalf("failed to open clip file at %s: %v", clipPath, err)
	}
	defer file.Close()

	// Read far enough into the mp4 file to parse the clip's duration from the mvhd box
	var duration = 0
	_, err = mp4.ReadBoxStructure(file, func(h *mp4.ReadHandle) (interface{}, error) {
		if h.BoxInfo.Type == mp4.BoxTypeMoov() {
			return h.Expand()
		}
		if h.BoxInfo.Type == mp4.BoxTypeMvhd() {
			box, _, err := h.ReadPayload()
			if err != nil {
				return nil, err
			}
			mvhd, ok := box.(*mp4.Mvhd)
			if !ok {
				return nil, fmt.Errorf("mvhd payload is not of type *mp4.Mvhd")
			}
			durationInSeconds := int(math.Round(float64(mvhd.DurationV0) / float64(mvhd.Timescale)))
			duration = durationInSeconds
		}
		return nil, nil
	})
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	if duration == 0 {
		log.Fatalf("failed to parse duration")
	}

	// Summarize the sync
	clipFilename := filepath.Base(clipPath)
	clipId := strings.TrimSuffix(clipFilename, ".mp4")
	fmt.Printf("Syncing clip...\n")
	fmt.Printf("--------------------------------------------------------\n")
	fmt.Printf("      ID: %s\n", clipId)
	fmt.Printf("   Title: %s\n", title)
	fmt.Printf("Duration: %d seconds\n", duration)
	fmt.Printf(" Tape ID: %d\n", tapeId)
	fmt.Printf("--------------------------------------------------------\n")

	// Build a request payload and serialize it to JSON
	data, err := json.Marshal(remix.Clip{
		Id:       clipId,
		Title:    title,
		Duration: duration,
		TapeId:   tapeId,
	})
	if err != nil {
		log.Fatalf("failed to serialize request payload: %v", err)
	}

	// Prepare the HTTP request that will sync our clip data
	url := remixUrl + "/admin/clip"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		log.Fatalf("failed to initialize request: %v", err)
	}
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Make the request and ensure that it was successful
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("request failed: %v", err)
	}
	fmt.Printf("%d\n", res.StatusCode)
	if res.StatusCode != http.StatusNoContent {
		suffix := ""
		if message, err := io.ReadAll(res.Body); err == nil {
			suffix = fmt.Sprintf(": %s", message)
		}
		log.Fatalf("got response %d%s", res.StatusCode, suffix)
	}

	fmt.Printf("OK.\n")
}
