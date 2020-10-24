package action

import (
	"github.com/chromedp/chromedp"
	"time"
)

func SangForLogin(url, username, password string, loginTime, pageTime int) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.SendKeys(`#user`, username, chromedp.ByID),
		chromedp.SendKeys(`#password`, password, chromedp.ByID),
		chromedp.Sleep(time.Duration(loginTime) * time.Second),
		chromedp.Click(`#button`, chromedp.ByID),
		chromedp.WaitVisible(`#ext-gen159`, chromedp.ByID),
		chromedp.Sleep(time.Duration(pageTime) * time.Second),
	}
}
