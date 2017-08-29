package azure

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
)

func RegisterResourceProvider(attempts int, backoff time.Duration, codes ...int) autorest.SendDecorator {
	return func(s autorest.Sender) autorest.Sender {
		return autorest.SenderFunc(func(r *http.Request) (resp *http.Response, err error) {
			fmt.Println("RegisterResourceProvider")
			rr := autorest.NewRetriableRequest(r)
			err = rr.Prepare()
			if err != nil {
				fmt.Println("err = rr.Prepare() failed")
				return resp, err
			}

			fmt.Println(rr.Request())
			resp, err = s.Do(rr.Request())
			if err != nil {
				return resp, err
			}

			if resp.StatusCode == http.StatusConflict {
				fmt.Println("omgomg status conflict")
				fmt.Println(resp)

				var re RequestError
				err = autorest.Respond(
					resp,
					autorest.ByUnmarshallingJSON(&re),
					autorest.ByClosing(),
				)
				if err != nil {
					return resp, err
				}

				if re.ServiceError != nil && re.ServiceError.Code == "MissingSubscriptionRegistration" {
					fmt.Println("Dun dun dun, you need to register")
					err = register(s, r)
					if err != nil {
						fmt.Println("Errore")
						fmt.Println(err)
						return resp, err
					}
					fmt.Println("Registered!")
				}
			}
			fmt.Println("Now lets retry...")
			resp, err = autorest.SendWithSender(s, rr.Request(),
				autorest.DoRetryForStatusCodes(attempts, backoff, codes...))

			fmt.Println("retries done!")
			return resp, err
		})
	}
}

func register(sender autorest.Sender, originalReq *http.Request) error {
	parts := strings.Split(originalReq.URL.Path, "/")
	subID := findParameter(parts, "subscriptions")
	provider := findParameter(parts, "providers")
	if subID != "" && provider != "" {
		newURL := url.URL{
			Scheme: originalReq.URL.Scheme,
			Host:   originalReq.URL.Host,
		}
		// taken from the resources SDK
		pathParameters := map[string]interface{}{
			"resourceProviderNamespace": autorest.Encode("path", provider),
			"subscriptionId":            autorest.Encode("path", subID),
		}

		const APIVersion = "2016-09-01"
		queryParameters := map[string]interface{}{
			"api-version": APIVersion,
		}

		preparer := autorest.CreatePreparer(
			autorest.AsPost(),
			autorest.WithBaseURL(newURL.String()),
			autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/{resourceProviderNamespace}/register", pathParameters),
			autorest.WithQueryParameters(queryParameters),
		)
		req, err := preparer.Prepare(&http.Request{})
		if err != nil {
			return err
		}
		//
		fmt.Println(req)
		resp, err := sender.Do(req)
		fmt.Println(resp)
		if err != nil {
			return err
		}
	}
	return nil
}

func findParameter(pathParts []string, id string) string {
	for i, v := range pathParts {
		if v == id && (i+1) < len(pathParts) {
			return pathParts[i+1]
		}
	}
	return ""
}
