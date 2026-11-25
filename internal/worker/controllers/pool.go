package controllers

import (
	"context"
	"fmt"
	"time"
)

func FillPool(ctx context.Context, ctrl *Controller) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			newVm, err := ctrl.create()
			fmt.Printf("%+v\n", newVm)
			if err != nil {
				time.Sleep(1 * time.Second)
			} else {
				ctrl.WarmVMPool <- newVm
			}
		}
	}
}
