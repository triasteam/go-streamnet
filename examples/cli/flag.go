package main

import (
	"fmt"
	"flag"
	"time"
	"errors"
	"strings"
)

type interval []time.Duration
func (i *interval) String() string {
	return fmt.Sprint(*i)
}

func (i *interval) Set(value string) error {
	if len(*i) > 0 {
		return errors.New("interval flag already set")
	}
	for _, dt := range strings.Split(value, ",") {
		duration, err := time.ParseDuration(dt)
		if err != nil {
			return err
		}
		*i = append(*i, duration)
	}
	return nil
}

func main() {
	/* usage: flag -flagname 1 --flagname2 2 --flagname3=1s,2s
	 *        flag --help
	 *
	 */
	var ip = flag.Int("flagname", 1234, "help message for flagname")

	var flagvar int
	flag.IntVar(&flagvar, "flagname2", 1234, "help message for flagname2")

	
	var flagvar2 interval
	flag.Var(&flagvar2, "flagname3", "help message for flagname3")


	flag.Parse()

	//fmt.Println(*ip, flagvar)
	fmt.Println(*ip, flagvar, flagvar2)

	//fmt.Println(flag.Args())
}
