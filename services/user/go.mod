module github.com/robrt95x/godops/services/user

go 1.21

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/robrt95x/godops/pkg/errors v0.0.0-00010101000000-000000000000
	github.com/robrt95x/godops/pkg/logger v0.0.0-00010101000000-000000000000
	github.com/robrt95x/godops/pkg/middleware v0.0.0-00010101000000-000000000000
)

replace github.com/robrt95x/godops/pkg/errors => ../../pkg/errors

replace github.com/robrt95x/godops/pkg/logger => ../../pkg/logger

replace github.com/robrt95x/godops/pkg/middleware => ../../pkg/middleware
