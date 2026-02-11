package handlers

import (
	"github.com/gin-gonic/gin"
)

// AuthCookieMaxAge defines the lifetime of authentication cookies in seconds (7 days).
// Used for access token, refresh token, and CSRF token cookies.
const AuthCookieMaxAge = 7 * 24 * 60 * 60 

// setAuthCookies sets the access token, refresh token, and CSRF token as cookies for the client.
// - Access and refresh tokens are set as HttpOnly cookies to prevent XSS attacks (JS cannot read them).
// - CSRF token is set as a regular cookie so the frontend can read it and send it in headers, preventing CSRF attacks.
// - All cookies use a predefined max age for consistency.
func setAuthCookies(c *gin.Context, accessToken, refreshToken, csrfToken string) {
	c.SetCookie("access_token", accessToken, AuthCookieMaxAge, "/", "", true, true)
	c.SetCookie("refresh_token", refreshToken, AuthCookieMaxAge, "/", "", true, true)
	c.SetCookie("csrf_token", csrfToken, AuthCookieMaxAge, "/", "", false, false)
}
