package token_pool

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

func encodeHexByteArray(bytes []byte) string {
	hexString := hex.EncodeToString(bytes)
	if !strings.HasPrefix(hexString, "0x") {
		hexString = "0x" + hexString
	}
	return hexString
}

func formatSuiObjectString(objectString string) string {
	if !strings.HasPrefix(objectString, "0x") {
		objectString = "0x" + objectString
	}
	return objectString
}

func DecodeBase64ParamsArray(base64Params []string) ([]string, error) {
	decodedParams := []string{}
	for _, param := range base64Params {
		decodedParam, err := base64.StdEncoding.DecodeString(param)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 param: %w", err)
		}
		decodedParams = append(decodedParams, encodeHexByteArray(decodedParam))
	}
	return decodedParams, nil
}
