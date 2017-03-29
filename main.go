package main

import (
	"context"
	"image"
	"image/color"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480

	Speed = 5
)

func genEvents(ctx context.Context, ed screen.EventDeque) {
	tick := time.NewTicker(time.Second / 60)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case ts := <-tick.C:
			ed.Send(ts)
		}
	}
}

func Move(r *image.Rectangle, x, y int) Node {
	return NodeFunc(func() Result {
		r.Min.X += x
		r.Max.X += x

		if r.Min.X < 0 {
			r.Max.X -= r.Min.X
			r.Min.X = 0
			return Success
		}
		if r.Max.X > ScreenWidth {
			r.Min.X += ScreenWidth - r.Max.X
			r.Max.X = ScreenWidth
			return Success
		}

		r.Min.Y += y
		r.Max.Y += y

		if r.Min.Y < 0 {
			r.Max.Y -= r.Min.Y
			r.Min.Y = 0
			return Success
		}
		if r.Max.Y > ScreenHeight {
			r.Min.Y += ScreenHeight - r.Max.Y
			r.Max.Y = ScreenHeight
			return Success
		}

		return NotDone
	})
}

func main() {
	driver.Main(func(s screen.Screen) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		win, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  ScreenWidth,
			Height: ScreenHeight,
		})
		if err != nil {
			panic(err)
		}
		defer win.Release()
		go genEvents(ctx, win)

		r := image.Rect(10, 10, 110, 60)
		tree := Sequence(
			Move(&r, Speed, 0),
			Move(&r, 0, Speed),
			Move(&r, -Speed, 0),
			Move(&r, 0, -Speed),
		)

		for {
			switch ev := win.NextEvent().(type) {
			case lifecycle.Event:
				if ev.To == lifecycle.StageDead {
					return
				}

			case time.Time:
				RunBT(tree)

				win.Fill(image.Rect(0, 0, ScreenWidth, ScreenHeight), color.Black, screen.Src)
				win.Fill(r, &color.NRGBA{255, 0, 255, 255}, screen.Over)
				win.Publish()
			}
		}
	})
}
