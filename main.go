package main

import (
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	e := echo.New()
	url1, err := url.Parse("http://127.0.0.1:8080/")
	if err != nil {
		e.Logger.Fatal(err)
	}
	b := middleware.NewRoundRobinBalancer(
		[]*middleware.ProxyTarget{
			{
				URL: url1,
			},
		})
	e.Use(middleware.Proxy(b))
	// dns autorisation
	e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("api.yoannfort.ga")
	// cache file
	e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")

	// Http server
	go func(c *echo.Echo) {
		// https redirection
		e.Pre(middleware.HTTPSRedirect())
		e.Logger.Fatal(e.Start(":80"))
	}(e)

	// Https server
	e.Logger.Fatal(e.StartAutoTLS(":443"))
}
