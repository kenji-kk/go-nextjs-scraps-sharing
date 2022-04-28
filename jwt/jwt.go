package jwt

import (
	"app/db"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cristalhq/jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var (
  jwtSigner   jwt.Signer
  jwtVerifier jwt.Verifier
)

func JwtSetup() {
  var err error
  key := []byte("jwtSecret123")

  jwtSigner, err = jwt.NewSignerHS(jwt.HS256, key)
  if err != nil {
    fmt.Printf("Error creating JWT signer")
  }

  jwtVerifier, err = jwt.NewVerifierHS(jwt.HS256, key)
  if err != nil {
    fmt.Printf("Error creating JWT verifier")
  }
}

func VerifyJWT(tokenStr string) (int, error) {
  token, err := jwt.Parse([]byte(tokenStr))
  if err != nil {
    log.Error().Err(err).Str("tokenStr", tokenStr).Msg("Error parsing JWT")
    return 0, err
  }

  if err := jwtVerifier.Verify(token.Payload(), token.Signature()); err != nil {
    log.Error().Err(err).Msg("Error verifying token")
    return 0, err
  }

  var claims jwt.StandardClaims
  if err := json.Unmarshal(token.RawClaims(), &claims); err != nil {
    log.Error().Err(err).Msg("Error unmarshalling JWT claims")
    return 0, err
  }

  if notExpired := claims.IsValidAt(time.Now()); !notExpired {
    return 0, errors.New("Token expired.")
  }

  id, err := strconv.Atoi(claims.ID)
  if err != nil {
    log.Error().Err(err).Str("claims.ID", claims.ID).Msg("Error converting claims ID to number")
    return 0, errors.New("ID in token is not valid")
  }
  return id, err
}

func Authorization(ctx *gin.Context) {
  authHeader := ctx.GetHeader("Authorization")
  if authHeader == "" {
    ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing."})
    return
  }
  headerParts := strings.Split(authHeader, " ")
  if len(headerParts) != 2 {
    ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format is not valid."})
    return
  }
  if headerParts[0] != "Bearer" {
    ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing bearer part."})
    return
  }
  userID, err := VerifyJWT(headerParts[1])
  if err != nil {
    ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    return
  }
  user, err := db.FetchUser(int64(userID))
  if err != nil {
    ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    return
  }
  ctx.Set("user", user)
  ctx.Next()
}

func CurrentUser(ctx *gin.Context) (*db.User, error) {
  var err error
  _user, exists := ctx.Get("user")
  if !exists {
    err = errors.New("Current context user not set")
    log.Error().Err(err).Msg("")
    return nil, err
  }
  user, ok := _user.(*db.User)
  if !ok {
    err = errors.New("Context user is not valid type")
    log.Error().Err(err).Msg("")
    return nil, err
  }
  return user, nil
}

func GenerateJWT(user *db.User) string {
  claims := &jwt.RegisteredClaims{
    ID:        fmt.Sprint(user.Id),
    ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
  }
  builder := jwt.NewBuilder(jwtSigner)
  token, err := builder.Build(claims)
  if err != nil {
    fmt.Printf("Error building JWT")
  }
  return token.String()
}

