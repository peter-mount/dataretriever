package main

import (
  "flag"
  "log"
  "os"
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

  if( !settings.Http.enabled ) {
    log.Fatal( "No message source configured, bailing out" )
  }

  amqpConnect()
  //amqpPublish( []byte("Test message") )

  if( settings.Http.enabled ) {
    httpRun()
  }

}
