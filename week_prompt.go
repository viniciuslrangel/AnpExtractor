package main

import (
	"AnpExtractor/anp"
	"fmt"
	"github.com/eiannone/keyboard"
)

func promptWeek(startWeek int, maxWeek int) int {

	err := keyboard.Open()
	if err != nil {
		panic(err)
	}

	weekNum := startWeek

loop:
	for {
		firstDay := anp.WeekToTime(weekNum)
		fmt.Printf(
			"\rSelect week (up/down arrow): %s", weekFormat(firstDay),
		)

		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		switch key {
		case keyboard.KeyArrowUp:
			if weekNum < maxWeek {
				weekNum++
			}
		case keyboard.KeyArrowDown:
			weekNum--
		case keyboard.KeyEnter:
			break loop
		}
	}
	keyboard.Close()

	return weekNum
}
