package shape

import (
	"fmt"
	"strings"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

func init() {
	launcher.NewBrowser().MustGet()
}

func NewBrowser(proxy string) *rod.Browser {
	var browser *rod.Browser

	if proxy != "" {
		//incomplete code, proxy won't work yet

		l := launcher.New()
		l = l.Set(flags.ProxyServer, proxy)

		controlURL, _ := l.Launch()
		browser := rod.New().ControlURL(controlURL).MustConnect()

		go browser.MustHandleAuth("user", "password")()

		browser.MustIgnoreCertErrors(true)
	}

	browser = rod.New().MustConnect()

	return browser
}

func NewPage(browser *rod.Browser) *rod.Page {
	page := stealth.MustPage(browser)

	return page
}

func LoadSite(page *rod.Page, address string) {
	page.MustNavigate(address)
}

type ShapeHeaders struct {
	XDy6zuyp5a string `json:"X-Dy6zuyp5-a"`
	XDy6zuyp5b string `json:"X-Dy6zuyp5-b"`
	XDy6zuyp5c string `json:"X-Dy6zuyp5-c"`
	XDy6zuyp5d string `json:"X-Dy6zuyp5-d"`
	XDy6zuyp5f string `json:"X-Dy6zuyp5-f"`
	XDy6zuyp5z string `json:"X-Dy6zuyp5-z"`
}

type ShapeHarvester struct {
	Proxy          string
	Url            string
	ShapeUrl       string
	Identifier     string
	Method         string
	Body           string
	Headers        ShapeHeaders
	Page           *rod.Page
	Browser        *rod.Browser
	BlockResources bool
}

func (harvester *ShapeHarvester) HarvestHeaders() {
	harvester.Page.MustEval(fmt.Sprintf(`function shape() {
		try {
			fetch("%s", {
				"method" : "%s",
				"referrerPolicy": "no-referrer-when-downgrade",
				"credentials": "include",
				"body": "%s",
				"headers": {
					"accept": "application/json",
					"accept-language": "en-US,en;q=0.9",
					"content-type": "application/json",
					"sec-ch-ua": "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"",
					"sec-ch-ua-mobile": "?0",
					"sec-ch-ua-platform": "\"Windows\"",
					"sec-fetch-dest": "empty",
					"sec-fetch-mode": "cors",
					"sec-fetch-site": "same-site",
					"x-application-name": "web",
				},
			})
		} catch {}
	  }`, harvester.ShapeUrl, harvester.Method, harvester.Body))
}

func (harvester *ShapeHarvester) InitializeHijacking() {
	router := harvester.Page.HijackRequests()

	router.MustAdd("*", func(ctx *rod.Hijack) {
		if harvester.BlockResources {
			if ctx.Request.Method() == "GET" {
				if ctx.Request.Type() == proto.NetworkResourceTypeImage || ctx.Request.Type() == proto.NetworkResourceTypeStylesheet {
					ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
				}

				ctx.ContinueRequest(&proto.FetchContinueRequest{})
			}
		}

		if strings.Contains(ctx.Request.URL().Path, harvester.Identifier) {
			if ctx.Request.Method() == "OPTIONS" {
				ctx.ContinueRequest(&proto.FetchContinueRequest{})
			} else if ctx.Request.Method() == "POST" {
				ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
				for key, value := range ctx.Request.Headers() {
					switch key {
					case "X-Dy6zuyp5-a":
						harvester.Headers.XDy6zuyp5a = value.String()

					case "X-Dy6zuyp5-b":
						harvester.Headers.XDy6zuyp5b = value.String()

					case "X-Dy6zuyp5-c":
						harvester.Headers.XDy6zuyp5c = value.String()

					case "X-Dy6zuyp5-d":
						harvester.Headers.XDy6zuyp5d = value.String()

					case "X-Dy6zuyp5-f":
						harvester.Headers.XDy6zuyp5f = value.String()

					case "X-Dy6zuyp5-z":
						harvester.Headers.XDy6zuyp5z = value.String()
					}
				}
			}
		}

	})

	go router.Run()
}

func (harvester *ShapeHarvester) InitializeHarvester() {
	harvester.Browser = NewBrowser("")
	harvester.Page = NewPage(harvester.Browser)

	harvester.Page.MustNavigate(harvester.Url).MustWaitLoad()

	harvester.InitializeHijacking()
	harvester.HarvestHeaders()
}
