module github.com/GMcD/telar-web/micros/storage

replace github.com/red-gold/telar-core v0.1.16 => github.com/GMcD/telar-core v0.1.30

go 1.16

require (
	cloud.google.com/go/storage v1.13.0
	github.com/GMcD/cognito-jwt v0.0.0-20210806015718-8416e465865c
	github.com/GMcD/telar-web v0.1.133
	github.com/aws/aws-sdk-go v1.34.28
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gofiber/adaptor/v2 v2.1.4
	github.com/gofiber/fiber/v2 v2.10.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/onsi/ginkgo v1.15.0 // indirect
	github.com/onsi/gomega v1.10.5 // indirect
	github.com/red-gold/telar-core v0.1.16
	golang.org/x/oauth2 v0.0.0-20210220000619-9bb904979d93
	google.golang.org/api v0.40.0 // indirect
	google.golang.org/genproto v0.0.0-20210222152913-aa3ee6e6a81c // indirect
)
