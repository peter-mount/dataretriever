# RabbitMQ Bridge

A simple Stomp to RabbitMQ Bridge

This is a simple go application that can bridge a remote service and submit messages to a local RabbitMQ instance.

# Introduction

This application supports multiple data sources using multiple protocols ranging from simple http polling to connecting to remote Message Brokers using Stomp.

## HTTP polling

This mode is intended for remotely accessing a remote webservice. Any response received is then subbitted to the local RabbitMQ instance.
