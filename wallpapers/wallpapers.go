package wallpapers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
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
	// VHS Box art
	VHS_Box WallpaperCatagory = "VHS Box"
	Games   WallpaperCatagory = "Games"
)

var (
	ALL_CATAGORIES []WallpaperCatagory = []WallpaperCatagory{
		Ancient, AsiaRussia, Blurry, Calm, Cityscapes, Creepy,
		Cyberpunk, Dark, Dreamy, Dystopian, Fantasy, Games, Grainy,
		Nature, No_Catagory, Perspective, Purple, Snow, Space,
		Synthwave, Technology, Triangular, VHS_Box,
	}
	OnlineWallpapers []WallpaperRepo
)

// Returns a semi-random URL to a new wallpaper.
//
// If Include and Exclude have any overlap, the overlapping catagories are not included.
// By default all catagories are included, however if `len(Include) > 0` then only the specified catagories are used.
func NewURL(Exclude []WallpaperCatagory, Include []WallpaperCatagory) (string, WallpaperCatagory, error) {
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
	fmt.Printf("catagory: %v\n", catagory)
	fmt.Printf("url: %v\n", directory_url)
	// List files in directory
	urls, err := GetFiles(repoType, directory_url)
	must(err)
	fmt.Printf("urls: %v\n", urls)
	// Pick one FILE (not subdirectory)
	// Return
	return "", catagory, nil
}

// Doesnt preserve order, but is very fast
// https://stackoverflow.com/questions/37334119
func remove(s []WallpaperCatagory, i int) []WallpaperCatagory {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	// Should be pretty random
	rand.Seed(int64(time.Now().Local().Nanosecond() * time.Now().Day() * time.Now().Hour() * time.Now().Second()))
	// do some random shit?
	for range make([]struct{}, 230) {
		rand.ExpFloat64()
		rand.Int63()
	}
	// Load OnlineWallpapers with data from ./wallpapers.json
	f, err := os.Open("./wallpapers/wallpapers.json")
	must(err)
	bytes, err := io.ReadAll(f)
	must(err)
	must(json.Unmarshal(bytes, &OnlineWallpapers))
}
