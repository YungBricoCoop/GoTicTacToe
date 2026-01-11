package main

import (
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) initUI() {
	// Root : vertical (top + HUD)
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(WindowSizeX, WindowSizeY)),
	)

	// Top area (game view)
	topArea := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(
			WindowSizeX,
			WindowSizeY-HUDHeight,
		)),
	)

	// HUD bar (bottom, 3 blocks)
	hudBar := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(WindowSizeX, HUDHeight)),
	)

	// LEFT block (player icon drawn by Ebiten in draw.go, so this stays EMPTY)
	left := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(HUDLeftW, HUDHeight)),
	)

	// CENTER block (score drawn by Ebiten in draw.go, so empty container is OK)
	center := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(HUDCenterW, HUDHeight)),
	)

	// RIGHT block (WASD image)
	right := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(HUDRightW, HUDHeight)),
	)

	wasd := widget.NewGraphic(
		widget.GraphicOpts.Image(g.assets.WASDHUDImage),
	)
	right.AddChild(wasd)

	// Assemble HUD
	hudBar.AddChild(left)
	hudBar.AddChild(center)
	hudBar.AddChild(right)

	root.AddChild(topArea)
	root.AddChild(hudBar)

	g.ui = &ebitenui.UI{Container: root}
}

func (g *Game) uiUpdate() {
	if g.ui != nil {
		g.ui.Update()
	}
}

func (g *Game) uiDraw(screen *ebiten.Image) {
	if g.ui != nil {
		g.ui.Draw(screen)
	}
}
