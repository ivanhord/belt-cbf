module github.com/ivanhord/belt-cbf/verify_service

go 1.24.4

require (
	github.com/ivanhord/belt-cbf/shared/bee2 v0.0.0
	golang.org/x/text v0.27.0
)

replace github.com/ivanhord/belt-cbf/shared/bee2 => ../shared/bee2
