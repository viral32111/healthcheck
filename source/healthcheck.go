package main

/* Requirements:
 Needs to be ran with just the URL argument.
 Needs to return only 0 for success, and 1 for failure.
 Needs to output only a single status message.
 Needs to only attempt connections once.
 Needs to have infinite timeout because Docker will handle that for us.
*/

/* Usage:
 hc http://192.168.10.20:5000/metrics
 hc http://abcxyzhiddenservice.onion 127.0.0.1:9150
*/

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"io"
	"strings"
	"golang.org/x/net/proxy"
)

const (
	AUTHOR_NAME = "viral32111"
	AUTHOR_WEBSITE = "https://viral32111.com"

	VERSION = "1.1.0"

	PROXY_ADDRESS = "127.0.0.1:9050"
)


func main() {
	argumentCount := len( os.Args[ 1: ] )

	if argumentCount < 1 {
		fmt.Printf( "Container Health Check, v%s, by %s (%s)\n", VERSION, AUTHOR_NAME, AUTHOR_WEBSITE )
		fmt.Fprintf( os.Stderr, "Usage: %s <url> [socks5 proxy]", os.Args[ 0 ] )
		os.Exit( 1 )
	}

	targetUrl, urlParseError := url.Parse( os.Args[ 1 ] )
	if urlParseError != nil {
		fmt.Fprintf( os.Stderr, "Failed to parse URL %s: %s", os.Args[ 1 ], urlParseError.Error() )
		os.Exit( 1 )
	}

	httpClient := &http.Client{ Timeout: 0, CheckRedirect: func( req *http.Request, via []*http.Request ) error {
		return http.ErrUseLastResponse
	} }

	httpRequest, httpRequestError := http.NewRequest( http.MethodGet, targetUrl.String(), nil )
	if httpRequestError != nil {
		fmt.Fprintf( os.Stderr, "Failed to create HTTP request: %s", httpRequestError.Error() )
		os.Exit( 1 )
	}

	if strings.HasSuffix( targetUrl.Host, ".onion" ) {
		proxyAddress := PROXY_ADDRESS

		if argumentCount == 2 {
			proxyAddress = os.Args[ 2 ]
		}

		proxyDialer, proxyConnectError := proxy.SOCKS5( "tcp", proxyAddress, nil, proxy.Direct )
		if proxyConnectError != nil {
			fmt.Fprintf( os.Stderr, "Failed to connect to proxy %s: %s", proxyAddress, proxyConnectError.Error() )
			os.Exit( 1 )
		}

		httpClient.Transport = &http.Transport{ Dial: proxyDialer.Dial }
	}

	httpResponse, httpExecuteError := httpClient.Do( httpRequest )
	if httpExecuteError != nil {
		fmt.Fprintf( os.Stderr, "Failed to execute HTTP request: %s", httpExecuteError.Error() )
		os.Exit( 1 )
	}

	httpResponseBody, readBodyError := io.ReadAll( httpResponse.Body )
	defer httpResponse.Body.Close()
	if readBodyError != nil {
		fmt.Fprintf( os.Stderr, "Failed to read response body: %s", readBodyError.Error() )
		os.Exit( 1 )
	}

	fmt.Printf( "%s: %s, %d bytes", targetUrl.Host, httpResponse.Status, len( httpResponseBody ) )

	if httpResponse.StatusCode < 200 || httpResponse.StatusCode > 299 {
		os.Exit( 1 )
	}

	os.Exit( 0 )
}
