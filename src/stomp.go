// stomp datasource

package main

import (
  "github.com/go-stomp/stomp"
  "log"
  "time"
)

// Main stomp config
type STOMP struct {
  // The server, e.g. "192.168.1.1:61613"
  Server          string            `yaml:"server"`
  // User credentials
  Username        string            `yaml:"username"`
  Password        string            `yaml:"password"`
  // Unique ID used for durable subscriptions
  ClientId        string            `yaml:"clientId"`
  // Heart Beat, defaults to 5s
  HeartBeat struct {
    SendTimeout     time.Duration   `yaml:"sendTimeout"`
    ReceiveTimeout  time.Duration   `yaml:"receiveTimeout"`
  } `yaml:"heartBeat"`
  // Host header
  Host            string            `yaml:"host"`
  // Subscriptions
  Subscription    []SUBSCRIPTION    `yaml:"subscription"`
  // ==== Internal values
  enabled         bool              `yaml:"-"`
  connection      *stomp.Conn       `yaml:"-"`
}

// Subscription config
type SUBSCRIPTION struct {
  // Topic to subscribe to
  Topic           string            `yaml:"topic"`
  // Unique name for this subscription
  Name            string            `yaml:"name"`
  // The routing key when sending messages
  routingKey      string            `yaml:"routingKey"`
}

func stompInit() {
  settings.Stomp.enabled = settings.Stomp.Server != "" && settings.Stomp.Username != "" && settings.Stomp.Password != "" && settings.Stomp.ClientId != ""

  if( settings.Stomp.enabled ) {
    if( settings.Stomp.HeartBeat.SendTimeout <=0 ) {
      settings.Stomp.HeartBeat.SendTimeout = time.Duration(5) * time.Second;
    }

    if( settings.Stomp.HeartBeat.ReceiveTimeout <=0 ) {
      settings.Stomp.HeartBeat.ReceiveTimeout = time.Duration(5) * time.Second;
    }

  }

}

func stompConnect() {
  log.Println( "Connecting to", settings.Stomp.Server )

  con, err := stomp.Dial(
    "tcp",
    settings.Stomp.Server,
    stomp.ConnOpt.Login( settings.Stomp.Username, settings.Stomp.Password ),
    stomp.ConnOpt.Host( settings.Stomp.Host ),
    stomp.ConnOpt.AcceptVersion( stomp.V10 ),
    stomp.ConnOpt.AcceptVersion( stomp.V11 ),
    stomp.ConnOpt.AcceptVersion( stomp.V12 ),
    stomp.ConnOpt.HeartBeat( settings.Stomp.HeartBeat.SendTimeout, settings.Stomp.HeartBeat.ReceiveTimeout ),
    stomp.ConnOpt.Header( "client-id", settings.Stomp.ClientId ) )
  fatalOnError( err )
  settings.Stomp.connection = con

  log.Println( "Connected to", settings.Stomp.Server )
}

func stompRun() {
  log.Println( "Stomp" )


}