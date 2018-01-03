// Handles http polling mode

package main

import (
  "io/ioutil"
  "log"
  "net/http"
  "time"
)

type HTTP struct {
  enabled     bool              `yaml:"-"`
  // URL to retrieve
  Url         string            `yaml:"url"`
  // Duration between requests
  Duration    time.Duration     `yaml:"duration"`
  // If set then we don't retrieve on startup
  NotOnStart  bool              `yaml:"notOnStart"`
  // Headers to send
  Headers     map[string]string `yaml:"headers"`
  // Optional Basic authentication
  BasicAuth struct {
    User      string
    Password  string
  } `yaml:"basicAuth"`
}

func httpInit( ) {
  // Only enable if we have a url set and a duration >= 1 second
  settings.Http.enabled = settings.Http.Url != "" && settings.Http.Duration >= time.Second
  if( settings.Http.enabled ) {
    debug("Enabling http for", settings.Http.Url, "every", settings.Http.Duration )
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
  // Initial retrieve
  if( !settings.Http.NotOnStart ) {
    httpRetrieve()
  }

  // Now run every duration
  ticker := time.NewTicker( settings.Http.Duration )
  go func() {
    for {
      debug( "Tick" )
      select {
        case <- ticker.C:
          httpRetrieve()
      }
    }
  }()
}
