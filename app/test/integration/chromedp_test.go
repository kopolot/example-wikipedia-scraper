package integration

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

func TestHowContextWork(t *testing.T) {
	ctx := context.Background()

	opts := make([]chromedp.ExecAllocatorOption, 0)
	opts = append(chromedp.DefaultExecAllocatorOptions[:], opts...)

	chromectx, chromecancel := chromedp.NewExecAllocator(ctx, opts...)
	newctx, newcancel := chromedp.NewContext(chromectx)

	err := chromedp.Run(newctx,
		chromedp.Navigate("https://httpbin.org/get"),
	)

	// chcoiaz wystarczy tylko chromecancel() i newcancel() nie jest potrzebny, bo chromecancel() zamyka cały kontekst i child contexty
	newcancel()
	chromecancel()
	time.Sleep(5 * time.Millisecond)

	if err != nil {
		t.Fatalf("Failed to navigate: %v", err)
	}
}

func TestHowContextWork2(t *testing.T) {
	ctx := context.Background()

	opts := make([]chromedp.ExecAllocatorOption, 0)
	opts = append(chromedp.DefaultExecAllocatorOptions[:], opts...)

	newctx, newcancel := chromedp.NewContext(ctx)

	err := chromedp.Run(newctx,
		chromedp.Navigate("https://httpbin.org/get"),
	)

	newcancel()
	time.Sleep(5 * time.Millisecond)

	if err != nil {
		t.Fatalf("Failed to navigate: %v", err)
	}
}
