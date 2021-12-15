package toolib

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func JwtString(jwtKey string, expired time.Duration) (string, error) {
	return JwtNewWithClaims(jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expired).Unix(), //到期时间
		IssuedAt:  time.Now().Unix(),              //发布时间
		NotBefore: time.Now().Unix(),              //在此之前不可用
	}, jwtKey)
}

func JwtSimple(expired time.Duration, auc, id, issuer, subject, jwtKey string) (string, error) {
	return JwtNewWithClaims(jwt.StandardClaims{
		Audience:  auc,                            //用户
		ExpiresAt: time.Now().Add(expired).Unix(), //到期时间
		Id:        id,                             //jwt标识
		IssuedAt:  time.Now().Unix(),              //发布时间
		Issuer:    issuer,                         //发行人
		NotBefore: time.Now().Unix(),              //在此之前不可用
		Subject:   subject,                        //主题
	}, jwtKey)
}

func JwtNewWithClaims(claims jwt.StandardClaims, jwtKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

func JwtVerify(tokenString, jwtKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt parse:%v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, fmt.Errorf("verify fail")
	} else {
		return claims, nil
	}
}
