package azure

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
)

// RegisterResourceProvider tries to register the Azure resource provider
// in case it is not registered yet.
func RegisterResourceProvider(attempts int, backoff time.Duration) autorest.SendDecorator {
	return func(s autorest.Sender) autorest.Sender {
		return autorest.SenderFunc(func(r *http.Request) (resp *http.Response, err error) {
			rr := autorest.NewRetriableRequest(r)
			err = rr.Prepare()
			if err != nil {
				return resp, err
			}

			resp, err = s.Do(rr.Request())
			if err != nil {
				return resp, err
			}

			if resp.StatusCode == http.StatusConflict {
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
					err = register(s, r)
					if err != nil {
						return resp, fmt.Errorf("failed auto registering Resource Provider: %s", err)
					}
				}
			}
			err = rr.Prepare()
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
		resp, err := sender.Do(req)
		if err != nil {
			return err
		}
		return autorest.Respond(
			resp,
			WithErrorUnlessStatusCode(http.StatusOK),
			autorest.ByClosing(),
		)

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
