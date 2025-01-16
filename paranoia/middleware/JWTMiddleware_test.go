package middleware

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
)

// Утилита для создания тестового токена
func createTestToken(privateKey *rsa.PrivateKey, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, _ := token.SignedString(privateKey)
	return tokenString
}

// Тесты
//func TestJWTMiddleware(t *testing.T) {
//	// Генерация тестовых ключей RSA
//	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
//	assert.NoError(t, err)
//	publicKey := &privateKey.PublicKey
//
//	// Настройка middleware
//	middleware := &JWTMiddleware{
//		name: "TestJWTMiddleware",
//		config: JWTMiddlewareConfig{
//			CtxKey: "user",
//		},
//		publicKey: publicKey,
//	}
//
//	// Подготовка тестового токена
//	validClaims := jwt.MapClaims{
//		"user_id": "12345",
//		"exp":     time.Now().Add(1 * time.Hour).Unix(),
//	}
//	validToken := createTestToken(privateKey, validClaims)
//
//	expiredClaims := jwt.MapClaims{
//		"user_id": "12345",
//		"exp":     time.Now().Add(-1 * time.Hour).Unix(),
//	}
//	expiredToken := createTestToken(privateKey, expiredClaims)
//
//	tests := []struct {
//		name           string
//		authHeader     string
//		expectedStatus int
//		expectedBody   string
//	}{
//		{
//			name:           "No Authorization Header",
//			authHeader:     "",
//			expectedStatus: http.StatusUnauthorized,
//			expectedBody:   "Authorization header is missing",
//		},
//		{
//			name:           "Invalid Authorization Format",
//			authHeader:     "InvalidFormat",
//			expectedStatus: http.StatusUnauthorized,
//			expectedBody:   "Invalid Authorization header format",
//		},
//		{
//			name:           "Invalid Token",
//			authHeader:     "Bearer invalid.token",
//			expectedStatus: http.StatusUnauthorized,
//			expectedBody:   "Invalid or expired token",
//		},
//		{
//			name:           "Expired Token",
//			authHeader:     "Bearer " + expiredToken,
//			expectedStatus: http.StatusUnauthorized,
//			expectedBody:   "Invalid or expired token",
//		},
//		{
//			name:           "Valid Token",
//			authHeader:     "Bearer " + validToken,
//			expectedStatus: http.StatusOK,
//			expectedBody:   "",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			ctx := srvUtils.HttpCtxPool.Get().(*srvUtils.HttpCtx)
//			ctx.Fill(&http.Request{
//				Header: http.Header{"Authorization": []string{tt.authHeader}},
//			})
//
//			nextCalled := false
//			next := func(_ context.Context, _ interfaces.ICtx) {
//				nextCalled = true
//			}
//			middleware.Invoke(next)(context.Background(), ctx)
//
//			assert.Equal(t, tt.expectedStatus, ctx.GetResponse().GetStatus())
//			assert.Equal(t, tt.expectedBody, string(ctx.GetResponse().GetBody()))
//
//			if tt.expectedStatus == http.StatusOK {
//				assert.True(t, nextCalled)
//			} else {
//				assert.False(t, nextCalled)
//			}
//		})
//	}
//}
