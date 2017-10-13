package azure

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
)

// DoRetryForStatusCodes is the Azure specific implementation for retring requests
// It also tries to register the resource provider in case it is unregistered
func DoRetryForStatusCodes(azc Client, codes ...int) autorest.SendDecorator {
	return func(s autorest.Sender) autorest.Sender {
		return autorest.SenderFunc(func(r *http.Request) (resp *http.Response, err error) {
			rr := autorest.NewRetriableRequest(r)
			for attempt := 0; attempt <= azc.RetryAttempts; attempt++ {
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
						azc.Sender = s
						err = azc.register(r, re)
						if err != nil {
							return resp, fmt.Errorf("failed auto registering Resource Provider: %s", err)
						}
					}
				}
				if !autorest.ResponseHasStatusCode(resp, codes...) {
					return resp, err
				}
				delayed := autorest.DelayWithRetryAfter(resp, r.Cancel)
				if !delayed {
					autorest.DelayForBackoff(azc.RetryDuration, attempt, r.Cancel)
				}
			}
			return resp, err
		})
	}
}

func getProvider(re RequestError) (string, error) {
	if re.ServiceError != nil {
		if re.ServiceError.Details != nil && len(*re.ServiceError.Details) > 0 {
			detail := (*re.ServiceError.Details)[0].(map[string]interface{})
			return detail["target"].(string), nil
		}
	}
	return "", errors.New("provider was not found in the response")
}

func (azc Client) register(originalReq *http.Request, re RequestError) error {
	subID := getSubscription(originalReq.URL.Path)
	if subID == "" {
		return errors.New("missing parameter subscriptionID to register resource provider")
	}
	providerName, err := getProvider(re)
	if err != nil {
		return fmt.Errorf("missing parameter provider to register resource provider: %s", err)
	}
	newURL := url.URL{
		Scheme: originalReq.URL.Scheme,
		Host:   originalReq.URL.Host,
	}

	// taken from the resources SDK
	// https://github.com/Azure/azure-sdk-for-go/blob/9f366792afa3e0ddaecdc860e793ba9d75e76c27/arm/resources/resources/providers.go#L252
	pathParameters := map[string]interface{}{
		"resourceProviderNamespace": autorest.Encode("path", providerName),
		"subscriptionId":            autorest.Encode("path", subID),
	}

	const APIVersion = "2016-09-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		azc.azurePreparation(),
		autorest.AsPost(),
		autorest.WithBaseURL(newURL.String()),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/{resourceProviderNamespace}/register", pathParameters),
		autorest.WithQueryParameters(queryParameters),
	)
	req, err := preparer.Prepare(&http.Request{})
	if err != nil {
		return err
	}
	req.Cancel = originalReq.Cancel

	resp, err := autorest.SendWithSender(azc.Sender, req)
	if err != nil {
		return err
	}

	type Provider struct {
		RegistrationState *string `json:"registrationState,omitempty"`
	}
	var provider Provider

	err = autorest.Respond(
		resp,
		WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&provider),
		autorest.ByClosing(),
	)
	if err != nil {
		return err
	}

	// poll for registered provisioning state
	now := time.Now()
	for err == nil && time.Since(now) < azc.PollingDuration {
		// taken from the resources SDK
		// https://github.com/Azure/azure-sdk-for-go/blob/9f366792afa3e0ddaecdc860e793ba9d75e76c27/arm/resources/resources/providers.go#L45
		preparer := autorest.CreatePreparer(
			azc.azurePreparation(),
			autorest.AsGet(),
			autorest.WithBaseURL(newURL.String()),
			autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/{resourceProviderNamespace}", pathParameters),
			autorest.WithQueryParameters(queryParameters),
		)
		req, err = preparer.Prepare(&http.Request{})
		if err != nil {
			return err
		}
		req.Cancel = originalReq.Cancel

		resp, err := autorest.SendWithSender(azc.Sender, req)
		if err != nil {
			return err
		}

		err = autorest.Respond(
			resp,
			WithErrorUnlessStatusCode(http.StatusOK),
			autorest.ByUnmarshallingJSON(&provider),
			autorest.ByClosing(),
		)
		if err != nil {
			return err
		}

		if provider.RegistrationState != nil &&
			*provider.RegistrationState == "Registered" {
			break
		}

		delayed := autorest.DelayWithRetryAfter(resp, originalReq.Cancel)
		if !delayed {
			autorest.DelayForBackoff(azc.PollingDelay, 0, originalReq.Cancel)
		}
	}
	if !(time.Since(now) < azc.PollingDuration) {
		return errors.New("polling for resource provider registration has exceeded the polling duration")
	}
	return err
}

func getSubscription(path string) string {
	parts := strings.Split(path, "/")
	for i, v := range parts {
		if v == "subscriptions" && (i+1) < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}
