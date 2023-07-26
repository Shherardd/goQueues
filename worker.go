package MQWP

import "sync"

type _Worker struct {
	wg *sync.WaitGroup
}
