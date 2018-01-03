// Handles http polling mode

package main

import (
  "io/ioutil"
  "log"
  "net/http"
)

type HTTP struct {
  enabled     bool              `yaml:"-"`
  Url         string            `yaml:"url"`    // URL to retrieve
  Headers     map[string]string `yaml:"headers"`  // Headers to send
  BasicAuth struct {
    User      string
    Password  string
  } `yaml:"basicAuth"`
}

func httpInit( ) {
  settings.Http.enabled = settings.Http.Url != ""
  log.Println( "enabled ", settings.Http.enabled, " url ", settings.Http.Url, " expr ", settings.Http.Url != "" )
  if( settings.Http.enabled ) {
    debug("Enabling http")
  }
}

func httpRetrieve() {
  req, err := http.NewRequest( "GET", settings.Http.Url, nil )
  if err != nil {
    log.Fatal( err )
  }

  // BasicAuth if defined
  if( settings.Http.BasicAuth.User != "" && settings.Http.BasicAuth.Password != "" ) {
    req.SetBasicAuth( settings.Http.BasicAuth.User, settings.Http.BasicAuth.Password )
  }

  // Add any headers
  for key, value := range settings.Http.Headers {
    req.Header.Add( key, value )
  }

  log.Println("Retrieving ", settings.Http.Url )
  resp, err := http.DefaultClient.Do( req )
  if err != nil {
    log.Fatal( err )
  }

  log.Println("Response ", resp.Status, " length ", resp.ContentLength, "Uncompressed ", resp.Uncompressed )

  b, err := ioutil.ReadAll( resp.Body )
  if err != nil {
    log.Fatal( err )
  }

  amqpPublish( b )

  defer resp.Body.Close()
}

func httpRun() {
  httpRetrieve()
}
