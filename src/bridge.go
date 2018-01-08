package main

import (
  "flag"
  "log"
  "os"
  "time"
)

func main() {
  log.Println( "RabbitMQ Bridge v0.1" )

  help := flag.Bool( "h", false, "Show help" )
  yamlFile := flag.String( "f", "/config.yaml", "The config file to use" )

  flag.Parse()

  if( *help ) {
    flag.PrintDefaults()
    os.Exit(0)
  }

  // Load config
  loadConfig( yamlFile )

  if( !settings.Http.enabled && !settings.Stomp.enabled ) {
    log.Fatal( "No message source configured, bailing out" )
  }

  settings.Amqp.Connect()

  if( settings.Http.enabled ) {
    httpRun()
  } else if( settings.Stomp.enabled ) {
    stompConnect()
  }

  // Now keep running forever
  for {
    time.Sleep(time.Minute)
  }
}
