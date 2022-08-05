package main

/* Requirements:
Needs to be ran with just the URL argument.
Needs to return only 0 for success, and 1 for failure.
Needs to output only a single status message.
Needs to only attempt connections once.
Needs to have infinite timeout because Docker will handle that for us.
*/

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/proxy"
)

const (
	PROJECT_NAME = "Container Health Check"
	PROJECT_VERSION = "2.0.0"

	AUTHOR_NAME = "viral32111"
	AUTHOR_WEBSITE = "https://viral32111.com"
)

func main() {

	// Holds the values of the command-line flags, and provides the defaults
	flagExpectStatusCode := 200
	flagUseMethod := "GET"
	flagProxyServer := ""

	// Setup the command-line flags
	flag.IntVar( &flagExpectStatusCode, "expect", flagExpectStatusCode, "The HTTP status code to consider a successful response, e.g. 204." )
	flag.StringVar( &flagUseMethod, "method", flagUseMethod, "The HTTP method to use in the request, e.g. GET." )
	flag.StringVar( &flagProxyServer, "proxy", flagProxyServer, "The IPv4 address & port number of a SOCKS5 proxy server, useful for Tor, e.g. 127.0.0.1:9050." )

	// Set a custom help message
	flag.Usage = func() {
		fmt.Printf( "%s, v%s, by %s (%s).\n", PROJECT_NAME, PROJECT_VERSION, AUTHOR_NAME, AUTHOR_WEBSITE )
		fmt.Printf( "\nUsage: %s [-expect <status>] [-method <method>] [-proxy <IP:PORT>] <URL>\n", os.Args[ 0 ] )

		flag.PrintDefaults()

		os.Exit( 1 ) // By default it exits with code 2
	}

	// Parse the command-line flags and arguments
	flag.Parse()
	argumentUrl := strings.Join( flag.Args(), "" )

	// Ensure the HTTP status code is valid
	if ( flagExpectStatusCode <= 0 || flagExpectStatusCode >= 1000 ) {
		fmt.Fprintln( os.Stderr, "Invalid HTTP status code to expect specified, it must be within the range 1 to 999." )
		os.Exit( 1 )
	}

	// Ensure the HTTP method is valid
	if ( flagUseMethod == "" ) {
		fmt.Fprintln( os.Stderr, "Invalid HTTP method to use specified, it must not be empty." )
		os.Exit( 1 )
	}

	// Ensure the target HTTP URL is valid
	targetUrl, parseError := url.Parse( argumentUrl )
	if ( parseError != nil ) {
		fmt.Fprintln( os.Stderr, "Invalid HTTP target URL specified:", parseError.Error() )
		os.Exit( 1 )
	} else if ( argumentUrl == "" ) {
		fmt.Fprintln( os.Stderr, "Invalid HTTP target URL specified, it must not be empty." )
		os.Exit( 1 )
	}

	// Create a HTTP client that does not follow redirects or timeout
	httpClient := &http.Client {
		Timeout: 0,
		CheckRedirect: func( request *http.Request, via []*http.Request ) error {
			return http.ErrUseLastResponse
		},
	}

	// Create a HTTP request to be sent later on
	httpRequest, requestError := http.NewRequest( strings.ToUpper( flagUseMethod ), targetUrl.String(), nil )
	if ( requestError != nil ) {
		fmt.Fprintln( os.Stderr, "Error creating HTTP request:", requestError.Error() )
		os.Exit( 1 )
	}

	// If a proxy server was provided...
	if ( flagProxyServer != "" ) {

		// Separate the address by IP address and port number
		flagProxyServerComponents := strings.SplitN( flagProxyServer, ":", 2 )

		// Ensure there are exactly two components in the array
		if ( len( flagProxyServerComponents ) != 2 ) {
			fmt.Fprintln( os.Stderr, "Invalid IPv4 address and port combination specified for the proxy server, it should be structured as IP:PORT." )
			os.Exit( 1 )
		}

		// Ensure the IP address is a valid IPv4 address
		proxyServerAddress := net.ParseIP( flagProxyServerComponents[ 0 ] )
		if ( proxyServerAddress == nil || proxyServerAddress.To4() == nil ) {
			fmt.Fprintln( os.Stderr, "Invalid IPv4 address specified for the proxy server." )
			os.Exit( 1 )
		}

		// Ensure the port number is valid
		proxyServerPort, parseError := strconv.ParseInt( flagProxyServerComponents[ 1 ], 10, 16 )
		if ( parseError != nil ) {
			fmt.Fprintln( os.Stderr, "Error parsing port number as an integer:", parseError.Error() )
			os.Exit( 1 )
		} else if ( proxyServerPort <= 0 || proxyServerPort >= 65536 ) {
			fmt.Fprintln( os.Stderr, "Invalid port number specified for the proxy server, it must be within the range 1 to 65535." )
			os.Exit( 1 )
		}

		// Create a dialer to the proxy server
		proxyDialer, proxyError := proxy.SOCKS5( "tcp", fmt.Sprintf( "%s:%d", proxyServerAddress, proxyServerPort ), nil, proxy.Direct )
		if ( proxyError != nil ) {
			fmt.Fprintln( os.Stderr, "Error creating dialer to proxy server:", proxyError.Error() )
			os.Exit( 1 )
		}

		// Set the HTTP client to use this proxy dialer
		httpClient.Transport = &http.Transport {
			Dial: proxyDialer.Dial,
		}

	}

	// Execute the HTTP request using the HTTP client
	httpResponse, executeError := httpClient.Do( httpRequest )
	if ( executeError != nil ) {
		fmt.Fprintln( os.Stderr, "Error executing HTTP request:", executeError.Error() )
		os.Exit( 1 )
	}

	// The healthcheck failed if the response status code does not match what was expected
	if ( httpResponse.StatusCode != flagExpectStatusCode ) {
		fmt.Println( "FAILURE,", httpResponse.Status )
		os.Exit( 1 )
	}

	// If we made it this far then the healthcheck is successful
	fmt.Println( "SUCCESS,", httpResponse.Status )
	os.Exit( 0 )

}
