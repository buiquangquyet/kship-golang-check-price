package middlewares

import (
	"check-price/src/common"
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"check-price/src/core/dto"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
)

type AuthMiddleware struct {
	publicKeys map[string]*rsa.PublicKey
}

func loadKeyMap() map[string]*rsa.PublicKey {
	mapPublicRsaKeys := make(map[string]*rsa.PublicKey)
	cf := configs.Get().Token.PublicKeys
	for key, data := range cf {
		keyData, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			log.Fatal("err decode public key, data:[%s], err:[%s]", data, err.Error())
		}

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
		if err != nil {
			log.Fatal("validate: parse key, err:[%s]", err.Error())
		}
		mapPublicRsaKeys[key] = publicKey
	}
	return mapPublicRsaKeys
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		publicKeys: loadKeyMap(),
	}
}

func (a *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		ierr := a.ValidateTokenClient(c, token)
		if ierr != nil {
			log.IErr(c, ierr)
			c.JSON(ierr.GetHttpStatus(), dto.ConvertErrorToResponse(ierr))
			c.Abort()
			return
		}
		c.Next()
	}
}

func (a *AuthMiddleware) ValidateTokenClient(c *gin.Context, token string) *common.Error {
	ierr := common.ErrUnauthorized(c).SetSource(common.SourceAPIService)

	tok, err := jwt.Parse(token, a.keyFunc)
	if err != nil {
		return ierr.SetDetail(err.Error())
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return ierr.SetDetail("token invalid")
	}

	kvrid, ok := claims["kvrid"]
	if !ok {
		return ierr.SetDetail("kvrid not found")
	}
	retailerId := kvrid.(int64)
	source, ok := claims["kvsource"]
	if !ok {
		return ierr.SetDetail("kvsource not found")
	}
	c.Set(constant.MerchantCodeKey, strings.ToUpper(source.(string)))
	c.Set(constant.MerchantIdKey, retailerId)
	versionLocation := 1
	versionLocationString := c.GetHeader("version-location")
	if versionLocationString == "" {
		versionLocation, _ = strconv.Atoi(versionLocationString)
	}
	c.Set(constant.VersionLocation, versionLocation)
	return nil
}

func (a *AuthMiddleware) keyFunc(jwtToken *jwt.Token) (interface{}, error) {
	if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if !ok {
		return nil, fmt.Errorf("invalid JWT Token")
	}

	source, ok := claims["kvsource"].(string)
	if !ok {
		return nil, fmt.Errorf("kvsource wrong format")
	}
	publicKey, ok := a.publicKeys[strings.ToLower(source)]
	if !ok {
		return nil, fmt.Errorf("kvsource not support")
	}
	return publicKey, nil
}
