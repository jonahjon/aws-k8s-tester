module github.com/aws/aws-k8s-tester/utils

go 1.16

require (
	github.com/briandowns/spinner v1.12.0
	github.com/dustin/go-humanize v1.0.0
	github.com/mitchellh/ioprogress v0.0.0-20180201004757-6a23b12fa88e
	go.uber.org/zap v1.16.0
	k8s.io/utils v0.0.0-20210305010621-2afb4311ab10
)

replace (
	github.com/aws/aws-k8s-tester/utils/file => ./file
	github.com/aws/aws-k8s-tester/utils/log => ./log
	github.com/aws/aws-k8s-tester/utils/rand => ./rand
	github.com/aws/aws-k8s-tester/utils/spinner => ./spinner
	github.com/aws/aws-k8s-tester/utils/terminal => ./terminal
)
