# RabbitMQ Bridge

A simple Stomp to RabbitMQ Bridge

This is a simple go application that can bridge a remote service and submit messages to a local RabbitMQ instance.

Note: This is still alpha quality so bugs will probably exist until I have fully tested it.

## Introduction

This application supports multiple data sources using multiple protocols ranging from simple http polling to connecting to remote Message Brokers using Stomp.

## Configuration

The application requires a single configuration file which is in yaml format. This file consists of several sections listed below. A template is available in the github repository.

### amqp

This section is mandatory and defines the connection to RabbitMQ.

    # Required: Details on logging in to RabbitMQ
    amqp:
      # Required: RabbitMQ connection string: amqp://user:pass@hostname
      url: amqp://user:password@rabbitHostname
      # Required: The routing key to use for all messages
      routingKey: test
      # Optional: The exchange to use, defaults to amq.topic
      #exchange: amq.topic

### debug

An optional setting that turns on additional debugging. For production set this to false or leave out.

    # Turn on debugging
    debug: true


### HTTP polling

This mode is intended for remotely accessing a remote webservice. Any response received is then subbitted to the local RabbitMQ instance.

    http:
      # Required: The url to retrieve
      url: https://api.example.com/
      # Optional: Basic authentication
      #basicAuth:
      #  user: MyUser
      #  password: MyPassword
      # Optional: Headers to send in the request. This is simply a map of key/value pairs sent for each request
      #headers:
      #  Authorization: SomeAuthorizationKey
      #  key: value

## Building

To build run the following:

    docker build -t myimage .

This will build a new image called myimage with everything compiled.

## Running

Ensure you have config.yaml setup then run the following replacing path/to/config.yaml to the absolute path to your config.yaml file:

    docker run -d --name myimage -v path/to/config.yaml:/config.yaml:ro myimage
