package queue

import "github.com/m-mizutani/octovy/pkg/domain/model"

type Queue struct {
	ch chan *model.ScanRepositoryRequest
}

const scanQueueLength = 1024

func New() *Queue {
	return &Queue{
		ch: make(chan *model.ScanRepositoryRequest, scanQueueLength),
	}
}

type Interface interface {
	SendQueue(q *model.ScanRepositoryRequest)
	RecvQueue() *model.ScanRepositoryRequest
}

func (x *Queue) SendQueue(q *model.ScanRepositoryRequest) {
	x.ch <- q
}

func (x *Queue) RecvQueue() *model.ScanRepositoryRequest {
	return <-x.ch
}
