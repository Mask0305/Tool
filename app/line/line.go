package line

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-resty/resty/v2"
	"github.com/kataras/iris/v12"
)

const (
	LineHost string = "https://api.line.me/oauth2/v2.1"
)

type LineTokenReq struct {
	GrantType    string `url:"grant_type" form:"grant_type"`
	Code         string `url:"code" form:"code"`
	RedirectUri  string `url:"redirect_uri" form:"redirect_uri"`
	ClientID     string `url:"client_id" form:"client_id"`
	ClientSecret string `url:"client_secret" form:"client_secret"`
}

type LineTokenRes struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

func LineAPI() {

	app := iris.New()
	party := app.Party("/")
	party.Get("login", func(c iris.Context) {
		host := "https://access.line.me/oauth2/v2.1/authorize?"
		argstr := ""
		args := map[string]string{
			"response_type": "code",
			// Line Developer Console找
			"client_id":    "99999999",
			"redirect_uri": url.QueryEscape("https://a695-61-216-167-153.ngrok.io/auth"),
			"state":        time.Now().Format("200601021504050000"),
			"scope":        "profile%20openid",
		}

		for i, v := range args {
			argstr += fmt.Sprintf("%s=%s&", i, v)
		}

		argstr = strings.TrimRight(argstr, "&")

		host += argstr

		c.JSON(map[string]string{
			"url": host,
		}, iris.JSON{UnescapeHTML: true})
	})

	party.Get("/auth", func(c iris.Context) {
		fmt.Println(c.Request().URL.RawQuery)
		body, _ := c.GetBody()
		fmt.Println(string(body))
	})

	party.Post("/take", func(c iris.Context) {
		req := new(LineTokenReq)
		_ = c.ReadForm(req)

		spew.Dump(req)

		bodyForm := map[string]string{}
		bodyForm["grant_type"] = req.GrantType
		bodyForm["code"] = req.Code
		bodyForm["redirect_uri"] = req.RedirectUri
		bodyForm["client_id"] = req.ClientID
		bodyForm["client_secret"] = req.ClientSecret

		httpClient := resty.New()
		request := httpClient.R()
		request.SetHeader("Content-Type", "application/x-www-form-urlencoded")

		request.SetFormData(bodyForm)

		// set response
		res := new(LineTokenRes)
		request.SetResult(res)

		// post
		response, err := request.Post(LineHost + "/token")
		if err != nil {
			spew.Dump(err)
			return
		}

		_ = json.Unmarshal(response.Body(), res)

		spew.Dump(res)
	})

	app.Listen(":8089")
}
