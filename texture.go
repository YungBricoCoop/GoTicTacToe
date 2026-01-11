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

func LoadTextures() (map[uint8][]*ebiten.Image, error) {
	entries, err := texturesFS.ReadDir(TextureFolder)
	if err != nil {
		return nil, fmt.Errorf("read texture folder %q: %w", TextureFolder, err)
	}

	out := make(map[uint8][]*ebiten.Image, len(entries))

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".png") {
			continue
		}

		key, errP := parseNumericFilenameUint8(name)
		if errP != nil {
			continue
		}

		fullPath := filepath.ToSlash(filepath.Join(TextureFolder, name))
		f, errO := texturesFS.Open(fullPath)
		if errO != nil {
			return nil, fmt.Errorf("open %q: %w", fullPath, errO)
		}

		img, _, errI := ebitenutil.NewImageFromReader(f)
		_ = f.Close()
		if errI != nil {
			return nil, fmt.Errorf("decode %q: %w", fullPath, errI)
		}

		strips, errS := sliceIntoVerticalStrips(img)
		if errS != nil {
			return nil, fmt.Errorf("slice %q: %w", fullPath, errS)
		}

		out[key] = strips
	}

	return out, nil
}

func LoadHUDImage(name string) (*ebiten.Image, error) {
	fullPath := filepath.ToSlash(filepath.Join(TextureFolder, name))
	f, errO := texturesFS.Open(fullPath)
	if errO != nil {
		return nil, fmt.Errorf("open HUD image %q: %w", fullPath, errO)
	}

	img, _, errI := ebitenutil.NewImageFromReader(f)
	_ = f.Close()
	if errI != nil {
		return nil, fmt.Errorf("decode HUD image %q: %w", fullPath, errI)
	}
	return img, nil
}

// parseNumericFilenameUint8 returns the uint8 value represented by the filename
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

// sliceIntoVerticalStrips returns {TextureSize} images of width 1
func sliceIntoVerticalStrips(src *ebiten.Image) ([]*ebiten.Image, error) {
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	if width != TextureSize {
		return nil, fmt.Errorf("expected texture width %d, got %d", TextureSize, width)
	}
	if height <= 0 {
		return nil, fmt.Errorf("invalid texture height %d", height)
	}

	strips := make([]*ebiten.Image, TextureSize)

	for x := range TextureSize {
		subImageStrip, _ := src.SubImage(image.Rect(x, 0, x+1, height)).(*ebiten.Image)
		strip := ebiten.NewImage(1, height)
		strip.DrawImage(subImageStrip, nil)
		strips[x] = strip
	}

	return strips, nil
}
