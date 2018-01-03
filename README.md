# RabbitMQ Bridge

This is a simple go application that can bridge a remote service and submit messages to a local RabbitMQ instance.

Note: This is still alpha quality so bugs will probably exist until I have fully tested it.

## Introduction

This application supports multiple data sources using multiple protocols ranging from simple http polling to connecting to remote Message Brokers using Stomp.

## Configuration

The application requires a single configuration file which is in yaml format. This file consists of several sections listed below. A template is available in the github repository.

### amqp

This section is mandatory and defines the connection to RabbitMQ.

```
# Required: Details on logging in to RabbitMQ
amqp:
  # Required: RabbitMQ connection string: amqp://user:pass@hostname
  url: amqp://user:password@rabbitHostname
  # Required: The routing key to use for all messages
  routingKey: test
  # Optional: The exchange to use, defaults to amq.topic
  #exchange: amq.topic
```

### debug

An optional setting that turns on additional debugging. For production set this to false or leave out.

```
# Turn on debugging
debug: true
```

### HTTP polling

This mode is intended for remotely accessing a remote webservice. Any response received is then subbitted to the local RabbitMQ instance.

```
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
  # The schedule to retrieve
  schedule: "0/10 * * * * *"
```

## Schedule expression format

A cron expression represents a set of times, using 6 space-separated fields.

Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | No         | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
Note: Month and Day-of-week field values are case insensitive. "SUN", "Sun", and "sun" are equally accepted.

### Special Characters
Asterisk ( * )

The asterisk indicates that the cron expression will match for all values of the field; e.g., using an asterisk in the 5th field (month) would indicate every month.

Slash ( / )

Slashes are used to describe increments of ranges. For example 3-59/15 in the 1st field (minutes) would indicate the 3rd minute of the hour and every 15 minutes thereafter. The form "*\/..." is equivalent to the form "first-last/...", that is, an increment over the largest possible range of the field. The form "N/..." is accepted as meaning "N-MAX/...", that is, starting at N, use the increment until the end of that specific range. It does not wrap around.

Comma ( , )

Commas are used to separate items of a list. For example, using "MON,WED,FRI" in the 5th field (day of week) would mean Mondays, Wednesdays and Fridays.

Hyphen ( - )

Hyphens are used to define ranges. For example, 9-17 would indicate every hour between 9am and 5pm inclusive.

Question mark ( ? )

Question mark may be used instead of '*' for leaving either day-of-month or day-of-week blank.

### Predefined schedules
You may use one of several pre-defined schedules in place of a cron expression.

Entry                  | Description                                | Equivalent To
-----                  | -----------                                | -------------
@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
@monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
@weekly                | Run once a week, midnight on Sunday        | 0 0 0 * * 0
@daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
@hourly                | Run once an hour, beginning of hour        | 0 0 * * * *

### Intervals
You may also schedule a job to execute at fixed intervals. This is supported by formatting the cron spec like this:

@every <duration>
where "duration" is a string accepted by time.ParseDuration (http://golang.org/pkg/time/#ParseDuration).

For example, "@every 1h30m10s" would indicate a schedule that activates every 1 hour, 30 minutes, 10 seconds.

Note: The interval does not take the job runtime into account. For example, if a job takes 3 minutes to run, and it is scheduled to run every 5 minutes, it will have only 2 minutes of idle time between each run.

### Time zones
By default, all interpretation and scheduling is done in the machine's local time zone (as provided by the Go time package http://www.golang.org/pkg/time). The time zone may be overridden by providing an additional space-separated field at the beginning of the cron spec, of the form "TZ=Asia/Tokyo"

Be aware that jobs scheduled during daylight-savings leap-ahead transitions will not be run!

## Building

To build run the following:

    docker build -t myimage .

This will build a new image called myimage with everything compiled.

## Running

Ensure you have config.yaml setup then run the following replacing path/to/config.yaml to the absolute path to your config.yaml file:

    docker run -d --name myimage -v path/to/config.yaml:/config.yaml:ro myimage
