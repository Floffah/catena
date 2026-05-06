package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/cli/browser"
	"github.com/floffah/catena/internal/pkg/auth"
	"github.com/gin-gonic/gin"
)

var CLI struct {
	Login struct {
		Instance string `arg:"" default:"http://localhost:3000" help:"URL of the Catena instance to authenticate with"`
	} `cmd:"" help:"Authenticate Catena CLI & Git with a Catena instance"`
}

func main() {
	ctx := kong.Parse(&CLI)

	// TODO: make cli login work
	// Probably need to reimplement using an oauth flow for clerk instead of the ticketing strategy as it might be deprecated (but documentation is sparse so unclear)

	switch ctx.Command() {
	case "login":
		g := gin.Default()

		var srv http.Server

		g.Handle("GET", "/callback", func(c *gin.Context) {
			token := c.Query("token")
			if token == "" {
				c.String(400, "token query parameter is required")
				return
			}

			postBody, err := json.Marshal(map[string]string{
				"strategy": "ticket",
				"ticket":   token,
			})
			println("exchanging token for session with payload: " + string(postBody))
			if err != nil {
				c.String(500, "failed to marshal request body: "+err.Error())
				return
			}
			resp, err := http.Post(auth.ClerkFrontendApiUrl+"/v1/client/sign_ins", "application/json", bytes.NewBuffer(postBody))
			if err != nil {
				c.String(500, "failed to exchange token for session: "+err.Error())
				return
			}

			if resp.StatusCode != 200 {
				c.String(500, "failed to exchange token for session: received non-200 status code "+resp.Status)
				return
			}

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				c.String(500, "failed to read response body: "+err.Error())
				return
			}

			var signInResponse struct {
				Response struct {
					SupportedFirstFactors []struct {
						Strategy string `json:"strategy"`
					} `json:"supported_first_factors"`
					Id string `json:"id"`
				} `json:"response"`
			}
			err = json.Unmarshal(body, &signInResponse)

			// just print all the json
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, body, "", "  ")
			if err != nil {
				c.String(500, "failed to format response body: "+err.Error())
				return
			}

			c.Header("Content-Type", "application/json")
			c.String(200, prettyJSON.String())

			go func() {
				// sleep for a bit to ensure the response is sent before shutting down the server
				time.Sleep(1 * time.Second)

				err := srv.Close()
				if err != nil {
					panic("failed to shutdown callback server: " + err.Error())
				}

				os.Exit(0)
			}()
		})

		srv = http.Server{
			Addr:    ":0",
			Handler: g,
		}
		ln, err := net.Listen("tcp", srv.Addr)
		if err != nil {
			panic("failed to start callback server: " + err.Error())
		}

		_, port, _ := net.SplitHostPort(ln.Addr().String())

		err = browser.OpenURL(CLI.Login.Instance + "/auth/cli?redirect_uri=http://localhost:" + port + "/callback")
		if err != nil {
			panic("failed to open browser: " + err.Error())
		}

		err = http.Serve(ln, g)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic("failed to start callback server: " + err.Error())
		}
	default:
		panic(ctx.Command())
	}
}
