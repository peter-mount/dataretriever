package main

import "flag"
import "log"
import "os"

type Config struct {
  help    bool        // true to show help and exit
  amqp    string      // amqp url
}

var config Config

func main() {
  log.Println( "RabbitMQ Bridge v0.1" )

  flag.StringVar( &config.amqp, "amqp", "", "The AMQP server to send messages to" )
  flag.BoolVar( &config.help, "h", false, "Show help" )

  flag.Parse()

  if( config.help ) {
    flag.PrintDefaults()
    os.Exit(0)
  }

  amqpInit( &config )
}
