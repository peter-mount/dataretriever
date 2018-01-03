package main

import (
  "log"
  "io/ioutil"
  "path/filepath"
  "gopkg.in/yaml.v2"
)

type Config struct {
  Debug   bool    // Debug logging
  Amqp    AMQP
  Http    HTTP
}

var settings Config

func loadConfig( configFile *string ) {
  filename, _ := filepath.Abs( *configFile )
  log.Println( "Loading config: ", filename )

  yml, err := ioutil.ReadFile( filename )
  if err != nil {
    log.Fatal( err )
  }

  //settings := Config{}
  err = yaml.Unmarshal( yml, &settings )
  if err != nil {
    log.Fatal( err )
  }

  log.Printf( "Config: %+v\n", settings )

  // Call each supported init method so they can play with the config
  amqpInit()
  httpInit()
}

// log.Println() only if debug is enabled
func debug( v ...interface{} ) {
  if( settings.Debug ) {
    log.Println( v... )
  }
}

// helper, log fatal if err is not nil
func fatalOnError( err interface{} ) {
  if( err != nil ) {
    log.Fatal( err )
  }
}
