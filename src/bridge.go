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

  loadConfig( yamlFile )

  amqpInit()

  amqpConnect()
  amqpPublish( []byte("Test message") )
}
