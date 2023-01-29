package main

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res string

	if err := chromedp.Run(ctx,
		chromedp.Navigate(`https://lncn.org/`),
		chromedp.OuterHTML(`body`, &res, chromedp.NodeVisible, chromedp.ByQuery),
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

}
