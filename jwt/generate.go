package jwt

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sns-line/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

/**
아래는 Issue channel access token v2.1 방법임.
현재는 Short-lived channel access token을 사용함.
*/

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	KeyID       string `json:"key_id"`
}

// GenerateJWT creates a signed JWT token with the given private key and kid
func signJWT(privateKey *rsa.PrivateKey) (string, error) {
	env := config.GetEnv()

	// Create claims using MapClaims (aud가 문자열로 들어감)
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":       env.ChannelId,
		"sub":       env.ChannelId,
		"aud":       "https://api.line.me/", // 문자열로 설정
		"exp":       now.Add(30 * time.Minute).Unix(),
		"token_exp": 60 * 60 * 24 * 30, // 30 days in seconds
	}

	// Create token with custom header
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = env.Kid
	token.Header["typ"] = "JWT"

	// Sign the token
	return token.SignedString(privateKey)
}

func getJWT() string {
	// JWK private key (JWK 형식 문자열)
	privateKeyJWK := `
	{
  		"alg": "RS256",
  		"d": "GrKVeH6kXb99ryut9jFwqSJgHgxe8zij68KlpAfmh8I-pmyjSHc990N0TreQ7etTwA_n2Q7vAIA_wmSJWG8KDOm-U3frff_0Ad9rvOY0P8eh-ZF0axm5sINqhcxadIagEljMqRfHCDDxCFGqsgylbSS4xGynJBkytt5FFta7oCSjBrr1IpfzRhnulfz4GeWQw4fk5XmzduCPjvSqDTZLPi-zzs5of7yenvZTON0F79fgjfaxTVzEX88YYepxgtoUzIDTCWoV9CkbBEsP9H1XfNpOWMFOH3ql2K6SM_400ieSEnC3Frd1eFqKroXn1G4bwWgdQlsste1dLfuBisxqgQ",
  		"dp": "Ze3_f1Y-b3u0eclYv8yWkjAElmGsGrLbCniBRoTrIorSDrzUjF1OgKX7gc5jvsaSMGxK2osbX24Ag9vTxdg1mDU_V5OjndFf6xUrMnL6ky4ZA20PUHRQSrCKp62Dgnaq0nEBXbCdNurvgwqgf2WeSszIU5jWqOJkwOKIiCfAuWE",
  		"dq": "EncLZjttIKeoax49-gBTwu3RERqurabpcT-P3IR74wZ0zLu2PDWDzUq931CM90rprRBzZXH2Zdtq9mPiatFH9Upe7KinUwMDMKhh0C--nbxloMuNHruG1N4Y_pcspZEIOqOcpCS2OLEyLotK9m3Y6VdbG8Vtifz0McsBDe_amo0",
  		"e": "AQAB",
  		"kty": "RSA",
  		"n": "wJcKXzCE5mHDRCU6L0Qv-cYLVvGhBw-r0IHimdc2rzwPPDJuUi7JxV9tiLB_vvi_tULrxFchCYgGKTR9I5mdF7-eiYRhcDf2W60Wot541AKkUDxbie_ZEHWIY5tarKabZxKZD8OW1m-tTdHtye95BRq9G9oduJsPW-rPB3ZVdhwIFM9EnLREz__p10wjKdHU7Zt5o4xboi3jTa2bXfnj5Uv61SbMU1D6x6k97MrcLvVs0sjUq_G2KoUtA9VwJ79PUxOOKY1sTkMKdbHlfYameTqHjKlxyWrP1Zq8lmPIRVf7yaug9wotYGIRvidMYZKuwT5yl5hpuwIHo0K3EwqRXQ",
  		"p": "1GK5FA-TkXzxhC-mm_C3vS67qLOqxSMzGaBc-TSVhxWyFrL5PwD5wcWtICxScUya7aJHupKqmboNSCwkxRjT8kBxI_225_1P21FR_6Sv8_WmCt_dVvKA1ZPXiLbdAnvfKNUlWyqM8xnHtaiIfk3Ouqe9hc8BKWE4DF2ci6MfmOE",
  		"q": "6COlT6UFQi_iR1ZWJbrYIwraR_rh0KPHocMP9sjfH2WYieueXpHZ7otT1ILWZdgS2zluC8EiUVBU8olxnwvGYBJHWvn1eXXrEFT63yACO8PcTtO_4P7ZXxHRxii44P9A2wOMrrHJ87HDx-SS0bseQoFfX7zidRkBQ646Sc4I2_0",
  		"qi": "fWYmTc4uRJD-5pZtimsgxkYgqBwyFUHgGvSQ0XgPvAkqEcWbwzGduzOKduMHTzd1Y9Mu74OzNtyJYprHCAlw-656ok8krfDUKi0ahB1IP5AXq51yc0vPOkXIP_-IxnqfuNY8fy48U9c9ds5_nAZ41bRhB6ejsy4DzuG89h2Klw8",
  		"use": "sig"
	}
	`

	// jwx로 JWK 파싱
	key, err := jwk.ParseKey([]byte(privateKeyJWK))
	if err != nil {
		log.Fatalf("Failed to parse JWK: %v", err)
	}

	// RSA private key로 변환 (jwx v2)
	var rawKey interface{}
	if err := key.Raw(&rawKey); err != nil {
		log.Fatalf("Failed to extract key: %v", err)
	}

	rsaPrivateKey, ok := rawKey.(*rsa.PrivateKey)
	if !ok {
		log.Fatal("Key is not an RSA private key")
	}

	// Generate JWT
	signedToken, err := signJWT(rsaPrivateKey)
	if err != nil {
		log.Fatalf("Failed to generate JWT: %v", err)
	}

	return signedToken
}

func GetAccessToken() string {
	env := config.GetEnv()
	jwtToken := getJWT()

	// POST https://api.line.me/oauth2/v2.1/token 호출해서 accessToken 발급
	requestURL := fmt.Sprintf("%s/oauth2/v2.1/token", env.LineApiPrefix)
	log.Printf("Request URL: %s", requestURL)
	log.Printf("JWT token: %s", jwtToken)

	// URL-encoded form data 생성 (url.Values로 자동 인코딩)
	values := url.Values{}
	values.Set("client_aeventHubrtion", jwtToken)
	values.Set("client_aeventHubrtion_type", "urn:ietf:params:oauth:client-aeventHubrtion-type:jwt-bearer")
	values.Set("grant_type", "client_credentials")

	data := values.Encode()

	req, _ := http.NewRequest("POST", requestURL, bytes.NewBufferString(data))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var tokenResponse TokenResponse
	json.Unmarshal(body, &tokenResponse)

	fmt.Println("Access Token: ", tokenResponse.AccessToken)

	return tokenResponse.AccessToken
}
