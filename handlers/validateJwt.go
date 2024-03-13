package handlers

import (
	"bytes"
	"errors"

	"crypto"
	"crypto/rsa"
	"crypto/sha256"

	"encoding/base64"
	"encoding/binary"
	"encoding/json"

	"fmt"
	"log"

	"math/big"
	"time"

	"github.com/IntuitDeveloper/OAuth2-Go/cache"
	"github.com/IntuitDeveloper/OAuth2-Go/config"

	"strings"
)

/*
 * Method to validate IDToken
 */
func ValidateIDToken(idToken string) bool {

	log.Println("Ending ValidateIDToken")
	if idToken != "" {
		parts := strings.Split(idToken, ".")

		if len(parts) < 3 {
			log.Fatalln("Malformed ID token")
			return false
		}

		idTokenHeader, err := base64.StdEncoding.DecodeString(parts[0])
		if err != nil {
			log.Fatalln("error parsing idTokenHeader:", err)
			return false
		}
		idTokenPayload, errr := base64.RawStdEncoding.DecodeString(parts[1])
		if errr != nil {
			log.Fatalln("error parsing idTokenPayload:", errr)
			return false
		}

		var payload = new(Claims)
		error := json.Unmarshal(idTokenPayload, &payload)
		if error != nil {
			log.Fatalln("error parsing payload:", error)
			return false
		}

		var header = new(Header)
		error1 := json.Unmarshal(idTokenHeader, &header)
		if error1 != nil {
			log.Fatalln("error parsing idTokenHeader:", error1)
			return false
		}

		//Step 1 : First check if the issuer is as mentioned in "issuer" in the discovery doc
		issuer := payload.ISS
		if issuer != cache.GetFromCache("issuer") {
			log.Fatalln("issuer value mismtach")
			return false
		}

		//Step 2 : check if the aud field in idToken is same as application's clientId
		audArray := payload.AUD
		aud := audArray[0]
		if aud != config.OAuthConfig.ClientId {
			log.Fatalln("incorrect client id")
			return false
		}

		//Step 3 : ensure the timestamp has not elapsed
		expirationTimestamp := payload.EXP
		now := time.Now().Unix()
		if (expirationTimestamp - now) <= 0 {
			log.Fatalln("expirationTimestamp has elapsed")
			return false
		}

		//Step 4: Verify that the ID token is properly signed by the issuer
		jwksResponse, err := CallJWKSAPI()
		if err != nil {
			log.Fatalln("error calling jwks", err)
			return false
		}

		//check if keys[0] belongs to the right kid
		headerKid := header.KID
        key, err := getKeyByKeyID(jwksResponse.KEYS, headerKid)

        if err != nil {
		    log.Fatalln("no keys found for the header ", err)
			return false
        }
        //get the exponent (e) and modulo (n) to form the PublicKey
		e := key.E
		n := key.N

		//build the public key
		pubKey, err := getPublicKey(n, e)
		if err != nil {
			log.Fatalln("unable to get public key", err)
			return false
		}

		//verify token using public key
		data := []byte(parts[0] + "." + parts[1])
		signature, err := base64.RawURLEncoding.DecodeString(parts[2])
		if err != nil {
			log.Fatalln("error decoding tokensignature:", err)
			return false
		}
		if err := verify(signature, data, pubKey); err != nil {
			log.Fatalln("unable to verify signature", err)
			return false
		}

		log.Println("Token Signature validated")
		return true
	}
	log.Println("Exiting ValidateIDToken")
	return false
}

type Header struct {
	ALG string `json:"alg"`
	KID string `json:"kid"`
}
type Claims struct {
	AUD       []string `json:"aud"`
	EXP       int64    `json:"exp"`
	IAT       int      `json:"iat"`
	ISS       string   `json:"iss"`
	REALMID   string   `json:"realmid"`
	SUB       string   `json:"sub"`
	AUTH_TIME int      `json:"auth_time"`
}

/*
 * Build public key
 */
func getPublicKey(modulus, exponent string) (*rsa.PublicKey, error) {

	decN, err := base64.RawURLEncoding.DecodeString(modulus)
	if err != nil {
		log.Fatalln("error decoding modulus", err)
	}
	n := big.NewInt(0)
	n.SetBytes(decN)

	decE, err := base64.RawURLEncoding.DecodeString(exponent)
	if err != nil {
		log.Fatalln("error decoding exponent", err)
	}
	var eBytes []byte
	if len(decE) < 8 {
		eBytes = make([]byte, 8-len(decE), 8)
		eBytes = append(eBytes, decE...)
	} else {
		eBytes = decE
	}
	eReader := bytes.NewReader(eBytes)
	var e uint64
	err = binary.Read(eReader, binary.BigEndian, &e)
	if err != nil {
		log.Fatalln("error reading exponent", err)
	}

	s := rsa.PublicKey{N: n, E: int(e)}
	var pKey *rsa.PublicKey = &s
	return pKey, err

}

/*
 * verify token using public key
 */
func verify(signature, data []byte, pubKey *rsa.PublicKey) error {
	hash := sha256.New()
	hash.Write(data)
	digest := hash.Sum(nil)

	if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, digest, signature); err != nil {
		return fmt.Errorf("unable to verify signature, %s", err.Error())
	}

	return nil
}

/*
 * retrieve keys matching the token kid
 */
func getKeyByKeyID(a []Keys, tknkid string) (Keys, error) {

	for i := 0; i < len(a); i++ {
        if a[i].KID == tknkid {
			return a[i], nil
		}
    }

	err := errors.New("Token is not valid, kid from token and certificate don't match")
	var b Keys
	return b, err
}
