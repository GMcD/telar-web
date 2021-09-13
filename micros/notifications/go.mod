module github.com/GMcD/telar-web/micros/notifications

replace github.com/red-gold/telar-core v0.1.16 => github.com/GMcD/telar-core v0.1.30

replace github.com/GMcD/telar-web v0.1.103 => github.com/GMcD/telar-web v0.1.108

go 1.16

require (
	github.com/GMcD/cognito-jwt v0.0.0-20210806015718-8416e465865c
	github.com/GMcD/telar-web v0.1.108
	github.com/alexellis/hmac v0.0.0-20180624211220-5c52ab81c0de
	github.com/gofiber/adaptor/v2 v2.1.4
	github.com/gofiber/fiber/v2 v2.10.0
	github.com/gofiber/template v1.6.10
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/red-gold/telar-core v0.1.16
)
