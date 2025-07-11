module github.com/robrt95x/godops/services/order

go 1.24.0

require (
	github.com/go-chi/chi/v5 v5.0.12
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/robrt95x/godops/pkg v0.0.0-00010101000000-000000000000
)

require (
	github.com/lib/pq v1.10.9
	github.com/sirupsen/logrus v1.9.3
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect

replace github.com/robrt95x/godops/pkg => ../../pkg
