package main

import (
	"embed"
	"fmt"
	"image"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed assets/textures/*.png
var texturesFS embed.FS

// TextureStrips represents an array of vertical strips of a texture image.
type TextureStrips [TextureSize]*ebiten.Image

// Texture represents a texture with its source image and vertical strips.
type Texture struct {
	Source *ebiten.Image
	Strips TextureStrips
}

// TextureMap maps texture IDs to their corresponding Texture.
type TextureMap map[TextureId]Texture

// TextureId represents the ID of a texture.
// 1-127 are reserved for wall textures.
// 128-255 are reserved for sprite textures.
type TextureId uint8

const (
	// walls
	WallBrick       TextureId = 1
	WallBrickHole   TextureId = 2
	WallBrickGopher TextureId = 3

	// sprites
	PlayerXSymbol    TextureId = 128
	PlayerXCharacter TextureId = 129
	PlayerOSymbol    TextureId = 130
	PlayerOCharacter TextureId = 131
	SkeletonSkull    TextureId = 132
	Chains           TextureId = 133
	Light            TextureId = 134
)

// LoadTextures loads all textures defined in imageManifest.
// Each file is decoded once into Source, and Strips are derived for raycasting.
func LoadTextures() (TextureMap, error) {
	out := make(TextureMap, len(imageManifest))

	for id, filename := range imageManifest {
		fullPath := filepath.ToSlash(filepath.Join(TextureFolder, filename))

		f, err := texturesFS.Open(fullPath)
		if err != nil {
			return nil, fmt.Errorf("open %q: %w", fullPath, err)
		}

		img, _, err := ebitenutil.NewImageFromReader(f)
		_ = f.Close()
		if err != nil {
			return nil, fmt.Errorf("decode %q: %w", fullPath, err)
		}

		strips, err := sliceIntoVerticalStrips(img)
		if err != nil {
			return nil, fmt.Errorf("slice %q: %w", fullPath, err)
		}

		out[id] = Texture{
			Source: img,
			Strips: strips,
		}
	}

	return out, nil
}

// parseNumericFilenameUint8 returns the uint8 value represented by the filename.
func parseNumericFilenameUint8(filename string) (uint8, error) {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	n, err := strconv.Atoi(base)
	if err != nil {
		return 0, fmt.Errorf("texture filename %q is not a number: %w", filename, err)
	}
	if n < 0 || n > 255 {
		return 0, fmt.Errorf("texture filename %q out of uint8 range: %d", filename, n)
	}
	return uint8(n), nil
}

// sliceIntoVerticalStrips returns TextureSize images of width 1.
// The source texture must be exactly TextureSize pixels wide.
func sliceIntoVerticalStrips(src *ebiten.Image) (TextureStrips, error) {
	// get source image bounds
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// enforce fixed texture width for raycasting
	if width != TextureSize {
		return TextureStrips{}, fmt.Errorf("expected texture width %d, got %d", TextureSize, width)
	}

	// height must be valid
	if height <= 0 {
		return TextureStrips{}, fmt.Errorf("invalid texture height %d", height)
	}

	// array of 1 pixel wide vertical texture strips
	var strips TextureStrips

	for x := 0; x < TextureSize; x++ {
		// take a 1 pixel wide slice from the source texture
		sub, ok := src.SubImage(image.Rect(x, 0, x+1, height)).(*ebiten.Image)
		if !ok || sub == nil {
			return TextureStrips{}, fmt.Errorf("failed to slice texture column x=%d", x)
		}

		// store the strip directly, subimage shares the same underlying texture
		strips[x] = sub
	}

	return strips, nil
}
