package wallpapers

import "fmt"

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
}

const (
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
func NewURL(Exclude []WallpaperCatagory, Include []WallpaperCatagory) (string, error) {
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
	fmt.Println(accept)
	return "", nil
}

// Doesnt preserve order, but is very fast
// https://stackoverflow.com/questions/37334119
func remove(s []WallpaperCatagory, i int) []WallpaperCatagory {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
