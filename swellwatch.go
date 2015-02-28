package main

import (
	"syscall"
	"time"
	"log"
)

func SwellWatch() {
	for {
		select {
		case tm := <-time.After(time.Second * 60):
			stat := syscall.Statfs_t{}
			err := syscall.Statfs("/", &stat)
			if err == nil {
				total := int(stat.Bsize) * int(stat.Blocks)
				free := int(stat.Bsize) * int(stat.Bfree)
				used := 100.0 - (float64(free) / float64(total) * 100)
				log.Println("used(%):", used, " total:", total, " free:", free, " time:" ,tm)
				//TODO modify more flexible
				noticeToAwsDiskUse(float64(used))
			} else {
				log.Println("Error:", err)
			}
		}
	}
}
