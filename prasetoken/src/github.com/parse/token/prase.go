package token

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

func NewToken(rawToken string) (*Token, error) {
	parts := strings.Split(rawToken, TokenSeparator)
	if len(parts) != 3 {
		return nil, errors.New("t format error")
	}

	var (
		rawHeader, rawClaims   = parts[0], parts[1]
		headerJSON, claimsJSON []byte
		err                    error
	)

	defer func() {
		if err != nil {
			log.Printf("error while unmarshalling raw t: %s", err)
		}
	}()

	if headerJSON, err = joseBase64UrlDecode(rawHeader); err != nil {
		err = fmt.Errorf("unable to decode header: %s", err)
		return nil, err
	}

	if claimsJSON, err = joseBase64UrlDecode(rawClaims); err != nil {
		err = fmt.Errorf("unable to decode claims: %s", err)
		return nil, err
	}

	t := new(Token)
	t.Header = new(Header)
	t.Claims = new(ClaimSet)

	t.Raw = strings.Join(parts[:2], TokenSeparator)
	if t.Signature, err = joseBase64UrlDecode(parts[2]); err != nil {
		err = fmt.Errorf("unable to decode signature: %s", err)
		return nil, err
	}

	if err = json.Unmarshal(headerJSON, t.Header); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(claimsJSON, t.Claims); err != nil {
		return nil, err
	}

	return t, nil
}

func joseBase64UrlDecode(s string) ([]byte, error) {
	switch len(s) % 4 {
	case 0:
	case 2:
		s += "=="
	case 3:
		s += "="
	default:
		return nil, errors.New("illegal base64url string")
	}
	return base64.URLEncoding.DecodeString(s)
}
