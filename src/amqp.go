// Handle the connection to the remote AMQP server to receive messages

package main

import (
  "log"
  "github.com/streadway/amqp"
)

type AMQP struct {
  Url         string `yaml:"url"`
  Exchange    string `yaml:"exchange"`
  RoutingKey  string `yaml:"routingKey"`
  connection  *amqp.Connection  `yaml:"-"`  // amqp connection
  channel     *amqp.Channel     `yaml:"-"`  // amqp channel
}

// called by main() ensure mandatory config is present
func amqpInit( ) {
  if( settings.Amqp.Url == "" ) {
    log.Fatal( "amqp.url is mandatory" )
  }

  if( settings.Amqp.Exchange == "" ) {
    settings.Amqp.Exchange = "amq.topic"
  }

  if( settings.Amqp.RoutingKey == "" ) {
    log.Fatal( "amqp.routingKey is mandatory" )
  }
}

// Connect to amqp server as necessary
func amqpConnect( ) {
  if( settings.Amqp.connection != nil && settings.Amqp.channel != nil ) {
    log.Println( "Already connected" )
    return
  }

  log.Println( "Connecting to amqp" )

  // Connect using the amqp url
  connection, err := amqp.Dial( settings.Amqp.Url )
  if( err != nil ) {
    log.Fatal( "Failed to connect to AMQP: ", err )
  }
  settings.Amqp.connection = connection

  // To cleanly shutdown by flushing kernel buffers, make sure to close and
  // wait for the response.
  //defer rabbit.connection.Close()

  // Most operations happen on a channel.  If any error is returned on a
  // channel, the channel will no longer be valid, throw it away and try with
  // a different channel.  If you use many channels, it's useful for the
  // server to
  channel, err := settings.Amqp.connection.Channel()
  if( err != nil ) {
    log.Fatal( "Failed to connect to AMQP: ", err )
  }
  settings.Amqp.channel = channel

  log.Println( "AMQP Connected" )

  if err := settings.Amqp.channel.ExchangeDeclare( settings.Amqp.Exchange, "topic", true, false, false, false, nil); err != nil {
    log.Fatalf("exchange.declare destination: %s", err)
  }

}

// Publish a message
func amqpPublish( msg []byte ) {
  //log.Println( "Publishing to ", settings.amqp.exchange, settings.amqp.routingKey )

  err := settings.Amqp.channel.Publish(
    settings.Amqp.Exchange,
    settings.Amqp.RoutingKey,
    false,
    false,
    amqp.Publishing{
      Body: msg,
    })
  if err != nil {
    log.Fatal( "Failed to publish message: ", err )
  }
}
