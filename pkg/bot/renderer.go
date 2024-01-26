package bot

import (
	"bytes"
	"fmt"
	"image/color"
	"math"

	"github.com/SkaceKamen/valetudo-telegram-bot/pkg/valetudo"
	"github.com/fogleman/gg"
)

func renderMap(mapData *valetudo.RobotStateMap) []byte {
	w := int(math.Round(float64(mapData.Size.X) / float64(mapData.PixelSize)))
	h := int(math.Round(float64(mapData.Size.Y) / float64(mapData.PixelSize)))

	minX := w
	minY := h
	maxX := 0
	maxY := 0

	for _, layer := range mapData.Layers {
		if layer.Dimensions.X.Min < minX {
			minX = layer.Dimensions.X.Min
		}

		if layer.Dimensions.X.Max > maxX {
			maxX = layer.Dimensions.X.Max
		}

		if layer.Dimensions.Y.Min < minY {
			minY = layer.Dimensions.Y.Min
		}

		if layer.Dimensions.Y.Max > maxY {
			maxY = layer.Dimensions.Y.Max
		}
	}

	minX -= int(float64(w) * 0.01)
	minY -= int(float64(h) * 0.01)
	maxX += int(float64(h) * 0.01)
	maxY += int(float64(h) * 0.01)

	resizedW := maxX - minX
	resizedH := maxY - minY

	ctx := gg.NewContext(resizedW, resizedH)

	for _, layer := range mapData.Layers {
		if layer.Type == "wall" {
			renderLayer(ctx, layer, minX, minY, color.RGBA{0, 0, 0, 255})
		}
		if layer.Type == "floor" {
			renderLayer(ctx, layer, minX, minY, color.RGBA{200, 200, 200, 255})
		}
		if layer.Type == "segment" {
			renderLayer(ctx, layer, minX, minY, color.RGBA{128, 128, 128, 255})
		}

		if layer.Dimensions.X.Min < minX {
			minX = layer.Dimensions.X.Min
		}

		if layer.Dimensions.X.Max > maxX {
			maxX = layer.Dimensions.X.Max
		}

		if layer.Dimensions.Y.Min < minY {
			minY = layer.Dimensions.Y.Min
		}

		if layer.Dimensions.Y.Max > maxY {
			maxY = layer.Dimensions.Y.Max
		}
	}

	for _, entity := range mapData.Entities {
		fmt.Println(entity.Type, "at", float64((*entity.Points)[0]-minX)/float64(mapData.PixelSize), ",", float64((*entity.Points)[1]-minY)/float64(mapData.PixelSize))

		if entity.Type == "charger_location" {
			ctx.SetColor(color.RGBA{0, 255, 0, 255})
			ctx.DrawCircle((float64((*entity.Points)[0])/float64(mapData.PixelSize))-float64(minX), (float64((*entity.Points)[1])/float64(mapData.PixelSize))-float64(minY), 5)
			ctx.Fill()
		}

		if entity.Type == "robot_position" {
			ctx.SetColor(color.RGBA{0, 0, 255, 255})
			ctx.DrawCircle((float64((*entity.Points)[0])/float64(mapData.PixelSize))-float64(minX), (float64((*entity.Points)[1])/float64(mapData.PixelSize))-float64(minY), 10)
			ctx.Fill()
		}
	}

	buffer := bytes.Buffer{}
	ctx.EncodePNG(&buffer)

	return buffer.Bytes()
}

func renderLayer(ctx *gg.Context, layer valetudo.RobotStateMapLayer, xOffset int, yOffset int, color color.Color) {
	ctx.SetColor(color)

	for i := 0; i < len(layer.CompressedPixels); i += 3 {
		xStart := layer.CompressedPixels[i]
		y := layer.CompressedPixels[i+1]
		count := layer.CompressedPixels[i+2]
		for j := 0; j < count; j++ {
			x := xStart + j
			ctx.SetPixel(x-xOffset, y-yOffset)
		}
	}
}
