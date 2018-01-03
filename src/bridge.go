package main

import "flag"
import "log"
import "os"

type Config struct {
  help        bool        // true to show help and exit
  // RabbitMQ config
  amqp        string      // amqp url
  exchange    string      // destination exchange
  routingKey  string      // routing key for outbound messages
}

var config Config

func main() {
  log.Println( "RabbitMQ Bridge v0.1" )

  flag.BoolVar( &config.help, "h", false, "Show help" )

  flag.StringVar( &config.amqp, "amqp", "", "The AMQP server to send messages to" )
  flag.StringVar( &config.exchange, "exchange", "amq.topic", "Destination exchange" )
  flag.StringVar( &config.routingKey, "routingKey", "", "Routing key for outbound messages" )

  flag.Parse()

  if( config.help ) {
    flag.PrintDefaults()
    os.Exit(0)
  }

  amqpInit( &config )

  amqpConnect()
  amqpPublish( []byte("Test message") )
}
