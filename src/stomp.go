// stomp datasource

package main

import (
  "github.com/go-stomp/stomp"
  "github.com/go-stomp/stomp/frame"
  "github.com/peter-mount/golib/statistics"
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
  // Connect delay to stop reconnecting too quickly
  ConnectDelay    time.Duration     `yaml:"connectDelay"`
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
  // The routing key when sending messages
  RoutingKey      string            `yaml:"routingKey"`
  // Headers to send on subscription
  Headers         map[string]string `yaml:"headers"`
  // Label used in recording statistics
  Label           string
  // ==== Internal
  sub             *stomp.Subscription `yaml:"-"`
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
  go func() {
    if settings.Stomp.ConnectDelay >= time.Second {
      log.Println( "Stomp waiting", settings.Stomp.ConnectDelay.String() )
      time.Sleep( settings.Stomp.ConnectDelay )
    }
    stompConnectImpl()
  }()
}

func stompConnectImpl() {

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

  for index,element := range settings.Stomp.Subscription {
    log.Println( index, "Subscribing to", element.Topic )

    var opts []func(*frame.Frame) error

    for key, value := range element.Headers {
      opts = append( opts, stomp.SubscribeOpt.Header( key, value ) )
    }

    sub, err := con.Subscribe(
      element.Topic,
      stomp.AckClient,
      func(f *frame.Frame) error {
        for _, h := range opts {
          err := h(f)
          if err != nil {
            return err
          }
        }
        return nil
      } )
    fatalOnError( err )
    settings.Stomp.Subscription[ index ].sub = sub

    if settings.Stomp.Subscription[ index ].Label == "" {
      settings.Stomp.Subscription[ index ].Label = settings.Stomp.Subscription[ index ].RoutingKey
    }

    go processQueue( &(settings.Stomp.Subscription[ index ]) )
  }
}

func processQueue( subscription *SUBSCRIPTION ) {
  debug( "Listening for", subscription.Topic )
  for {
    debug( "Tick", subscription.Topic )
    msg := <-subscription.sub.C
    fatalOnError( msg.Err )

    settings.Amqp.Publish( subscription.RoutingKey, msg.Body )

    statistics.Incr( subscription.Label )

    err := settings.Stomp.connection.Ack( msg )
    fatalOnError( err )
  }
}
