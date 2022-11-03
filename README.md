[![domain-availability-go license](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![domain-availability-go made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://pkg.go.dev/github.com/whois-api-llc/domain-availability-go)
[![domain-availability-go test](https://github.com/whois-api-llc/domain-availability-go/workflows/Test/badge.svg)](https://github.com/whois-api-llc/domain-availability-go/actions/)

# Overview

The client library for
[Domain Availability API](https://domain-availability.whoisxmlapi.com/)
in Go language.

The minimum go version is 1.17.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/domain-availability-go
```

# Examples

Full API documentation available [here](https://domain-availability.whoisxmlapi.com/api/documentation/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := domainavailability.NewBasicClient(apiKey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := domainavailability.NewClient(apiKey, domainavailability.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

Domain Availability API lets you get the domain registration state.

```go

// Make request to get the Domain Availability API response as a model instance.
domainAvailabilityResp, _, err := client.Get(ctx, "whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

if domainAvailabilityResp.IsAvailable != nil {
    if *domainAvailabilityResp.IsAvailable {
        log.Println(domainAvailabilityResp.DomainName, "is available") 
    } else {
        log.Println(domainAvailabilityResp.DomainName, "is unavailable")
    }
}

// Make request to get raw data in XML.
resp, err := client.GetRaw(context.Background(), "whoisxmlapi.com",
    domainavailability.OptionOutputFormat("XML"))
if err != nil {
    log.Fatal(err)
}

log.Println(string(resp.Body))

```