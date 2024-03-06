package app

import (
	"strconv"
	"time"

	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/auth"
	"github.com/golang-jwt/jwt"
)

// JWTHandler ...
func JWTHandler(c *routing.Context, j *jwt.Token) error {
	claims := j.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	roleID := claims["roleid"].(string)
	companyID := claims["companyid"].(string)
	GetRequestScope(c).SetUserID(userID)
	GetRequestScope(c).SetRoleID(roleID)
	GetRequestScope(c).SetCompanyID(companyID)
	return nil
}

// CreateToken create new token
func CreateToken(u *models.AdminUserDetails) (string, error) {
	token, err := auth.NewJWT(jwt.MapClaims{
		"id":         strconv.Itoa(int(u.UserID)),
		"authorized": true,
		"exp":        time.Now().Add(time.Hour * 2).Unix(),
		"companyid":  strconv.Itoa(int(u.CompanyID)),
		"roleid":     strconv.Itoa(int(u.RoleID)),
	}, Config.JWTSigningKey)

	return token, err
}
