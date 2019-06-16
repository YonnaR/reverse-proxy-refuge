package main

import (
	"flag"
	"log"
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/acme/autocert"
)

var (
	dns     string
	appPort string
)

func init() {
	flag.StringVar(&dns, "dns", "", "white hosted dns available for this app")
	flag.StringVar(&appPort, "port", "", "port of running app")

	flag.Parse()
	if dns == "" {
		log.Fatal("you need to set a dns we will use for the validation of https protocol")
	}
	if appPort == "" {
		log.Fatal("Port of the reunning app to proxy")
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())

	url1, err := url.Parse("http://127.0.0.1:" + appPort)
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
	e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(dns)
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
