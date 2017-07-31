package vault

import (
	"encoding/json"
	"errors"
)

func (auth *AuthInfo) WrapData(wrapttl string, raw string) (string, error) {
	client, err := auth.Client()
	if err != nil {
		return "", err
	}

	// unmarshal raw string into a map
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return "", err
	}

	// setup wrapping function
	client.SetWrappingLookupFunc(func(operation, path string) string {
		return wrapttl
	})

	resp, err := client.Logical().Write("/sys/wrapping/wrap", data)
	if err != nil {
		return "", err
	}
	return resp.WrapInfo.Token, nil
}

func (auth *AuthInfo) UnwrapData(wrappingToken string) (map[string]interface{}, error) {
	// unwrap uses the manual client creation in case auth's ID is empty
	client, err := NewVaultClient()
	if err != nil {
		return nil, err
	}
	client.SetToken(auth.ID)

	// if auth is empty, unwrapping is still allowed. It just won't be vault audited
	if auth.ID == "" {
		client.SetToken(wrappingToken)
		wrappingToken = ""
	}

	// unwrap either with auth or without
	resp, err := client.Logical().Unwrap(wrappingToken)
	if err != nil {
		return nil, errors.New("Failed to unwrap provided token " + err.Error())
	}
	return resp.Data, nil
}
