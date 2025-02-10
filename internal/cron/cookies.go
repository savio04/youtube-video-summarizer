package cron

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func Execute() {
	for {
		fmt.Println("Capturando cookie")
		err := captureCookies()
		if err != nil {
			log.Println("⚠️ Erro:", err)
		}

		time.Sleep(10 * time.Minute)
	}
}

func captureCookies() error {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), chromedp.DefaultExecAllocatorOptions[:]...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var cookies []*network.Cookie
	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate("https://www.youtube.com"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}),
	)
	if err != nil {
		log.Fatal("Erro ao capturar cookies:", err)
	}

	err = saveCookiesAsNetscape("cookies.txt", cookies)
	if err != nil {
		log.Fatal("Erro ao salvar cookies:", err)
	}

	fmt.Println("Cookies salvos com sucesso!")

	return nil
}

func saveCookiesAsNetscape(filename string, cookies []*network.Cookie) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("# Netscape HTTP Cookie File\n")

	for _, c := range cookies {
		secure := "FALSE"
		if c.Secure {
			secure = "TRUE"
		}

		line := fmt.Sprintf("%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
			c.Domain, "TRUE", c.Path, secure, int64(c.Expires), c.Name, c.Value)
		file.WriteString(line)
	}

	return nil
}
