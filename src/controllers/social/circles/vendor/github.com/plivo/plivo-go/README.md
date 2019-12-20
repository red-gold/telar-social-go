# plivo-go

[![Build Status](https://travis-ci.org/plivo/plivo-go.svg?branch=master)](https://travis-ci.org/plivo/plivo-go)
[![GoDoc](https://godoc.org/github.com/plivo/plivo-go?status.svg)](https://godoc.org/github.com/plivo/plivo-go)

The Plivo Go SDK makes it simpler to integrate communications into your Go applications using the Plivo REST API. Using the SDK, you will be able to make voice calls, send SMS and generate Plivo XML to control your call flows.

## Prerequisites

- Go >= 1.7.x

## Installation

You can use the SDK using the `go` command.

    $ go get github.com/plivo/plivo-go

You can also install by cloning this repository into your `GOPATH`.

## Getting started

### Authentication
To make the API requests, you need to create a `Client` and provide it with authentication credentials (which can be found at [https://manage.plivo.com/dashboard/](https://manage.plivo.com/dashboard/)).

We recommend that you store your credentials in the `PLIVO_AUTH_ID` and the `PLIVO_AUTH_TOKEN` environment variables, so as to avoid the possibility of accidentally committing them to source control. If you do this, you can initialise the client with no arguments and it will automatically fetch them from the environment variables:

```go
package main

import "github.com/plivo/plivo-go"

func main()  {
  client, err := plivo.NewClient("", "", &plivo.ClientOptions{})
  if err != nil {
    panic(err)
  }
}
```
Alternatively, you can specifiy the authentication credentials while initializing the `Client`.

```go
package main

import "github.com/plivo/plivo-go"

func main()  {
 client, err := plivo.NewClient("your_auth_id", "your_auth_token", &plivo.ClientOptions{})
 if err != nil {
   panic(err)
 }
}
```

## The Basics
The SDK uses consistent interfaces to create, retrieve, update, delete and list resources. The pattern followed is as follows:

```go
client.Resources.Create(Params{}) // Create
client.Resources.Get(Id) // Get
client.Resources.Update(Id, Params{}) // Update
client.Resources.Delete(Id) // Delete
client.Resources.List() // List all resources, max 20 at a time
```

Using `client.Resources.List()` would list the first 20 resources by default (which is the first page, with `limit` as 20, and `offset` as 0). To get more, you will have to use `limit` and `offset` to get the second page of resources.

## Examples

### Send a message

```go
package main

import "github.com/plivo/plivo-go"

func main()  {
  client, err := plivo.NewClient("", "", &plivo.ClientOptions{})
  if err != nil {
    panic(err)
  }
  client.Messages.Create(plivo.MessageCreateParams{
    Src: "the_source_number",
    Dst: "the_destination_number",
    Text: "Hello, world!",
  })
}
```

### Make a call

```go
package main

import "github.com/plivo/plivo-go"

func main()  {
  client, err := plivo.NewClient("", "", &plivo.ClientOptions{})
  if err != nil {
    panic(err)
  }
  client.Calls.Create(plivo.CallCreateParams{
    From: "the_source_number",
    To: "the_destination_number",
    AnswerURL: "http://answer.url",
  })
}
```

### Generate Plivo XML

```go
package main

import "github.com/plivo/plivo-go/plivo/xml"

func main()  {
  println(xml.ResponseElement{
    Contents: []interface{}{
      new(xml.SpeakElement).SetContents("Hello, world!"),
    },
    }.String())
}
```

This generates the following XML:

```xml
<Response>
  <Speak>Hello, world!</Speak>
</Response>
```

### More examples
Refer to the [Plivo API Reference](https://api-reference.plivo.com/latest/go/introduction/overview) for more examples. Also refer to the [guide to setting up dev environment](https://developers.plivo.com/getting-started/setting-up-dev-environment/) on [Plivo Developers Portal](https://developers.plivo.com) to setup a simple Go server & use it to test out your integration in under 5 minutes.

## Reporting issues
Report any feedback or problems with this version by [opening an issue on Github](https://github.com/plivo/plivo-go/issues).
