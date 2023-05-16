package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type TimerMap map[string]*Timer

var Timers TimerMap //map[string]*Timer

// for some things, we want the timers in order
func TimersAsSlice(ts TimerMap) []*Timer {
	tSlice := make([]*Timer, 0, len(ts))

	for _, t := range ts {
		tSlice = append(tSlice, t)
	}
	sort.Slice(tSlice, func(i, j int) bool {
		return tSlice[i].index < tSlice[j].index
	})
	return tSlice
}

type timerConfig struct {
	Key      string
	Name     string
	Type     string
	Duration string
}

type config struct {
	Osc      map[string]string
	Web      map[string]string
	Darkmode bool
	Timers   []timerConfig
}

func main() {
	// set up context for daemonizing
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// set up signal handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt) //syscall.SIGHUP to also listen to SIGHUP

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	// start the signal handler
	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				//case syscall.SIGHUP:
				// reload config?
				case os.Interrupt:
					cancel()
					log.Printf("Exiting on signal...")
					os.Exit(1)
				}
			case <-ctx.Done():
				log.Printf("Done.")
				os.Exit(1)
			}
		}
	}()

	// handle config
	var Config config
	var configFile = flag.String("config", "showtimer.yaml", "name of configuration file to read")
	flag.Parse()
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Fatalf("error reading timer config: %s", err)
	}
	//log.Printf("config: %v", Config)

	// timer map
	Timers = make(map[string]*Timer)

	// set up timers
	// assign an automatic index to preserve the as-written order in config
	for i, t := range Config.Timers {
		//log.Printf("Timer %s(%d): %s", t.Key, i, t.Name)
		switch t.Type {
		case "CountUp":
			Timers[t.Key], err = NewTimer(CountUp, t.Name, t.Key, i, time.Duration(0))
			if err != nil {
				log.Fatalf("Unable to create timer %s: %s", t.Key, err)
			}
		case "CountDown":
			d, err := time.ParseDuration(t.Duration)
			if err != nil {
				log.Fatalf("Unable to parse duration for countdown timer %s: %s", t.Key, err)
			}
			Timers[t.Key], err = NewTimer(CountDown, t.Name, t.Key, i, d)
			if err != nil {
				log.Fatalf("Unable to create timer %s: %s", t.Key, err)
			}
		default:
			log.Fatalf("Unknown timer type '%s' for timer %s (must be CountDown or CountUp)", t.Type, t.Key)
		}
	}

	webPort, err := strconv.Atoi(Config.Web["port"])
	if err != nil {
		log.Fatalf("Unable to parse port for web server from '%s': %s", Config.Web["port"], err)
	}
	runWebServer(Config.Web["bind"], webPort, Config.Darkmode, ctx)

	oscPort, err := strconv.Atoi(Config.Osc["port"])
	if err != nil {
		log.Fatalf("Unable to parse port for OSC server from '%s': %s", Config.Osc["port"], err)
	}
	runOSCServer(Config.Osc["bind"], oscPort, ctx)

	// wait until we are somehow told to exit
	<-ctx.Done()
}
