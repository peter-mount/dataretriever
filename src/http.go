// Handles http polling mode

package main

import (
  "github.com/peter-mount/golib/statistics"
  "gopkg.in/robfig/cron.v2"
  "io/ioutil"
  "log"
  "net/http"
  "time"
)

type HTTP struct {
  // URL to retrieve
  Url             string            `yaml:"url"`
  // One of Schedule or Duration is required
  // Crontab schedule
  Schedule        string            `yaml:"schedule"`
  // Duration between requests
  Duration        time.Duration     `yaml:"duration"`
  // If set then we don't retrieve on startup
  NotOnStart      bool              `yaml:"notOnStart"`
  // Publish on error
  PublishOnError  bool              `yaml:"publishOnError"`
  // Headers to send
  Headers         map[string]string `yaml:"headers"`
  // Optional Basic authentication
  BasicAuth struct {
    User      string
    Password  string
  } `yaml:"basicAuth"`
  // The routing key to use
  RoutingKey  string `yaml:"routingKey"`
  // ==== Internal values
  enabled         bool              `yaml:"-"`
}

func httpInit( ) {
  // Only enable if we have a url set and a duration >= 1 second
  settings.Http.enabled = settings.Http.Url != ""

  // Check if duration is in use else use the schedule
  if( settings.Http.enabled ) {
    if( settings.Http.Duration < time.Second ) {
      settings.Http.enabled = settings.Http.Schedule != ""
    }
  }

  if( settings.Http.enabled ) {
    if( settings.Http.RoutingKey == "" ) {
      log.Fatal( "http.routingKey is mandatory" )
    }

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

  // Publish only on 200 unless PublishOnError is set
  if resp.StatusCode >= 200 && resp.StatusCode < 300 {
    settings.Amqp.Publish( settings.Http.RoutingKey, b )
    statistics.Incr( settings.Http.RoutingKey + ".ok" )
  } else {
    if settings.Http.PublishOnError {
      settings.Amqp.Publish( settings.Http.RoutingKey, b )
    }
    statistics.Incr( settings.Http.RoutingKey + ".error" )
  }

  defer resp.Body.Close()
}

// Uses a simple ticker to run every duration
func httpRunDuration() {
  debug( "Running with duration", settings.Http.Duration.String() )

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

// Run's on a cron schedule
func httpRunCron() {
  debug( "Running with schedule", settings.Http.Schedule )

  c := cron.New()
  c.AddFunc( settings.Http.Schedule, func() {
    httpRetrieve()
  })
  c.Start()
}

func httpRun() {
  if( settings.Http.Duration >= time.Second ) {
    httpRunDuration()
  } else {
    httpRunCron()
  }
}
