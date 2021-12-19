package wallpapers

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// A catagory a wallpaper can have
type WallpaperCatagory string

type WallpaperSubsection struct {
	Catagory WallpaperCatagory `json:"Catagory"`
	URLs     []string          `json:"URLs"`
}

type WallpaperRepo struct {
	// A URL that contains all of the wallpapers.
	BaseWallpaperURL string `json:"BaseWallpaperURL"`
	//
	Subsections []WallpaperSubsection `json:"Subsections"`
	// Strings that can be appended to the base URL for access to uncatagorized sections.
	UncatagorizedPaths []string `json:"UncatagorizedPaths"`
	// The type of wallpaper repo
	Type WallpaperRepoType `json:"type"`
}

type WallpaperRepoType string

const (
	Eyy_Indexer WallpaperRepoType = "eyy-indexer"

	// There is no catagory
	No_Catagory WallpaperCatagory = "NO CATAGORY"
	// May also include subcatagory "Roads" - could be moved to seperate
	Cityscapes  WallpaperCatagory = "Cityscapes"
	Creepy      WallpaperCatagory = "Creepy"
	Ancient     WallpaperCatagory = "Ancient"
	AsiaRussia  WallpaperCatagory = "Asia/Russia"
	Blurry      WallpaperCatagory = "Blurry"
	Calm        WallpaperCatagory = "Calm"
	Cyberpunk   WallpaperCatagory = "Cyberpunk"
	Dark        WallpaperCatagory = "Dark"
	Dreamy      WallpaperCatagory = "Dreamy"
	Dystopian   WallpaperCatagory = "Dystopian"
	Fantasy     WallpaperCatagory = "Fantasy"
	Grainy      WallpaperCatagory = "Grainy"
	Nature      WallpaperCatagory = "Nature"
	Purple      WallpaperCatagory = "Purple"
	Perspective WallpaperCatagory = "Perspective"
	Snow        WallpaperCatagory = "Snow"
	Space       WallpaperCatagory = "Space"
	Synthwave   WallpaperCatagory = "Synthwave"
	Technology  WallpaperCatagory = "Technology"
	Triangular  WallpaperCatagory = "Triangular"
	VHS_Box_Art WallpaperCatagory = "VHS Box"
	Games       WallpaperCatagory = "Games"
)

var (
	ALL_CATAGORIES []WallpaperCatagory = []WallpaperCatagory{
		Ancient, AsiaRussia, Blurry, Calm, Cityscapes, Creepy,
		Cyberpunk, Dark, Dreamy, Dystopian, Fantasy, Games, Grainy,
		Nature, No_Catagory, Perspective, Purple, Snow, Space,
		Synthwave, Technology, Triangular, VHS_Box_Art,
	}
	OnlineWallpapers     []WallpaperRepo
	WallpaperDataIsValid = make(chan struct{})
)

// Returns a semi-random URL to a new wallpaper. Error if none can be found.
//
// If Include and Exclude have any overlap, the overlapping catagories are not included.
// By default all catagories are included, however if `len(Include) > 0` then only the specified catagories are used.
func NewRandomURL(Exclude []WallpaperCatagory, Include []WallpaperCatagory) (string, WallpaperCatagory, error) {
	// Caclulate accepted catagories
	accept := ALL_CATAGORIES
	if len(Include) > 0 {
		accept = []WallpaperCatagory{}
	}
	accept = append(accept, Include...)
	for _, wc := range Exclude {
		for i, catagory := range accept {
			if catagory == wc {
				accept = remove(accept, i)
			}
		}
	}
	// Caclulate # of things to stop after.
	all_subsections := []WallpaperSubsection{}

	<-WallpaperDataIsValid
	for _, wr := range OnlineWallpapers {
		all_subsections = append(all_subsections, wr.Subsections...)
	}
	stop_after := 0
	if len(all_subsections) > 0 {
		stop_after = rand.Intn(len(all_subsections) - 1)
	}
	// Pick a random catagory & pick a random corresponding URL
	base := ""
	append := ""
	var repoType WallpaperRepoType
	catagory := No_Catagory
	for _, wp := range OnlineWallpapers {
		repoType = wp.Type
		base = wp.BaseWallpaperURL
		for _, section := range wp.Subsections {
			if len(section.URLs) > 0 {
				append = section.URLs[rand.Intn(len(section.URLs))]
			}
			catagory = section.Catagory
			stop_after--
			if stop_after <= 0 {
				break
			}
		}
		if stop_after <= 0 {
			break
		}
	}
	directory_url := base + append
	// List files in directory
	urls, err := GetFiles(repoType, directory_url)
	if err != nil {
		return "", No_Catagory, err
	}
	if len(urls) <= 0 {
		return "", No_Catagory, errors.New("could not fetch any url")
	}
	return urls[rand.Intn(len(urls))], catagory, nil
}

func SetRandom() error {
	url, _, err := NewRandomURL(nil, nil)
	if err != nil {
		return err
	}
	return SetFromURL(url)
}

func SetRandomUntagged() error {
	maxTries := 5
	for i := 0; i < maxTries; i++ {
		url, _, err := NewRandomURL(nil, nil)
		if err != nil {
			return err
		}

		path, err := CalculatePath(url)
		if err != nil {
			return err
		}

		playlists, err := GetPlaylistsByPath(path)
		if err != nil {
			return fmt.Errorf("get playlists by path: %w", err)
		}
		if len(playlists) == 0 {
			return SetFromURL(url)
		}
		time.Sleep(time.Second / 2)
	}
	return fmt.Errorf("unable to get random untagged wallpaper after %d tries", maxTries)
}

// Doesnt preserve order, but is very fast
// https://stackoverflow.com/questions/37334119
func remove(s []WallpaperCatagory, i int) []WallpaperCatagory {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
