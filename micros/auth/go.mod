module github.com/GMcD/telar-web/micros/auth

replace github.com/red-gold/telar-core v0.1.16 => github.com/GMcD/telar-core v0.1.37

go 1.16

require (
	github.com/GMcD/telar-web v0.1.143
	github.com/alexellis/hmac v0.0.0-20180624211220-5c52ab81c0de
	github.com/gofiber/adaptor/v2 v2.1.3
	github.com/gofiber/fiber/v2 v2.10.0
	github.com/gofiber/template v1.6.9
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/golang-jwt/jwt/v4 v4.0.0
	github.com/google/uuid v1.2.0
	github.com/nbutton23/zxcvbn-go v0.0.0-20210217022336-fa2cb2858354
	github.com/red-gold/telar-core v0.1.16
	github.com/sethvargo/go-password v0.2.0
	golang.org/x/oauth2 v0.0.0-20210220000619-9bb904979d93
)
