package example

import (
	"context"
	"errors"
	domainavailability "github.com/whois-api-llc/domain-availability-go"
	"log"
)

func GetData(apikey string) {
	client := domainavailability.NewBasicClient(apikey)

	// Get parsed Domain Availability API response by a domain name as a model instance.
	domainAvailabilityResp, resp, err := client.Get(context.Background(),
		"whoisxmlapi.com",
		// this option is ignored, as the inner parser works with JSON only.
		domainavailability.OptionOutputFormat("XML"),
		// this option causes both DNS and WHOIS checking to be performed.
		domainavailability.OptionMode("DNS_AND_WHOIS"))

	if err != nil {
		// Handle error message returned by server.
		var apiErr *domainavailability.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	// Then print the domain registration state.
	if domainAvailabilityResp.IsAvailable != nil {
		if *domainAvailabilityResp.IsAvailable {
			log.Println(domainAvailabilityResp.DomainName, "is available")
		} else {
			log.Println(domainAvailabilityResp.DomainName, "is unavailable")
		}
	}

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func GetRawData(apikey string) {
	client := domainavailability.NewBasicClient(apikey)

	// Get raw API response.
	resp, err := client.GetRaw(context.Background(),
		"whoisxmlapi.com",
		// this option causes the Domain Availability API credits will be taken into account.
		domainavailability.OptionCredits("DA"),
		domainavailability.OptionOutputFormat("XML"))

	if err != nil {
		// Handle error message returned by server.
		log.Fatal(err)
	}

	log.Println(string(resp.Body))
}
