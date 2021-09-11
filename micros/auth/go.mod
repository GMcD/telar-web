module github.com/GMcD/telar-web/micros/auth

replace github.com/red-gold/telar-core v0.1.16 => github.com/GMcD/telar-core v0.1.21
replace github.com/GMcD/telar-web v0.1.89 => github.com/GMcD/telar-web v0.1.95

go 1.15

require (
	github.com/GMcD/cognito-jwt v0.0.0-20210806015718-8416e465865c
	github.com/alexellis/hmac v0.0.0-20180624211220-5c52ab81c0de
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gofiber/adaptor/v2 v2.1.3
	github.com/gofiber/fiber/v2 v2.10.0
	github.com/gofiber/template v1.6.9
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/google/uuid v1.2.0
	github.com/nbutton23/zxcvbn-go v0.0.0-20210217022336-fa2cb2858354
	github.com/red-gold/telar-core v0.1.16
	github.com/GMcD/telar-web v0.1.95
	golang.org/x/oauth2 v0.0.0-20210220000619-9bb904979d93
)
