package security

import (
	jwt "gopkg.in/square/go-jose.v2/jwt"
)

type Claims struct {
	Issuer      string           `json:"iss,omitempty"`
	Subject     string           `json:"sub,omitempty"`
	Audience    jwt.Audience     `json:"aud,omitempty"`
	Expiry      *jwt.NumericDate `json:"exp,omitempty"`
	NotBefore   *jwt.NumericDate `json:"nbf,omitempty"`
	IssuedAt    *jwt.NumericDate `json:"iat,omitempty"`
	ID          string           `json:"jti,omitempty"`
	Teams       []string         `json:"teams,omitempty"`
	Roles       []string         `json:"roles,omitempty"`
	Permissions []string         `json:"permissions,omitempty"`
}

func (c *Claims) ToMap() map[string]interface{} {
	claimsData := make(map[string]interface{})
	claimsData["iss"] = c.Issuer
	claimsData["sub"] = c.Subject
	claimsData["aud"] = c.Audience
	claimsData["exp"] = c.Expiry
	claimsData["nbf"] = c.NotBefore
	claimsData["iat"] = c.IssuedAt
	claimsData["jti"] = c.ID
	claimsData["teams"] = c.Teams
	claimsData["roles"] = c.Roles
	claimsData["permissions"] = c.Permissions
	return claimsData
}

func (c *Claims) FromMap(MapClaims map[string]interface{}) *Claims {
	c.Issuer = MapClaims["iss"].(string)
	c.Subject = MapClaims["sub"].(string)
	c.Audience = MapClaims["aud"].(jwt.Audience)
	c.Expiry = MapClaims["exp"].(*jwt.NumericDate)
	c.NotBefore = MapClaims["nbf"].(*jwt.NumericDate)
	c.IssuedAt = MapClaims["iat"].(*jwt.NumericDate)
	c.ID = MapClaims["jti"].(string)
	c.Teams = MapClaims["teams"].([]string)
	c.Roles = MapClaims["roles"].([]string)
	c.Permissions = MapClaims["permissions"].([]string)
	return c
}
