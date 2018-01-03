// Handle the connection to the remote AMQP server to receive messages

package main

import "log"
import "github.com/streadway/amqp"

type Rabbit struct {
  config      *Config           // Pointer to global config
  connection  *amqp.Connection  // amqp connection
  channel     *amqp.Channel     // amqp channel
}

var rabbit Rabbit

// called by main() ensure mandatory config is present
func amqpInit( config *Config ) {
  if( config.amqp == "" ) {
    log.Fatal( "--amqp is mandatory" )
  }

  rabbit.config = config
}

// Connect to amqp server as necessary
func amqpConnect( ) {
  if( rabbit.connection != nil && rabbit.channel != nil ) {
    log.Println( "Already connected" )
    return
  }

  log.Println( "Connecting to amqp" )

  // Connect using the amqp url
  connection, err := amqp.Dial( config.amqp )
  if( err != nil ) {
    log.Fatal( "Failed to connect to AMQP: ", err )
  }
  rabbit.connection = connection

  // To cleanly shutdown by flushing kernel buffers, make sure to close and
  // wait for the response.
  //defer rabbit.connection.Close()

  // Most operations happen on a channel.  If any error is returned on a
  // channel, the channel will no longer be valid, throw it away and try with
  // a different channel.  If you use many channels, it's useful for the
  // server to
  channel, err := rabbit.connection.Channel()
  if( err != nil ) {
    log.Fatal( "Failed to connect to AMQP: ", err )
  }
  rabbit.channel = channel

  log.Println( "AMQP Connected" )

  if err := rabbit.channel.ExchangeDeclare( rabbit.config.exchange, "topic", true, false, false, false, nil); err != nil {
    log.Fatalf("exchange.declare destination: %s", err)
  }

}

// Publish a message
func amqpPublish( msg []byte ) {
  log.Println( "Publishing to ", rabbit.config.exchange, rabbit.config.routingKey )

  err := rabbit.channel.Publish(
    rabbit.config.exchange,
    rabbit.config.routingKey,
    false,
    false,
    amqp.Publishing{
      Body: msg,
    })
  if err != nil {
    log.Fatal( "Failed to publish message: ", err )
  }
}
