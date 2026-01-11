package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hypebeast/go-osc/osc"
)

func runOSCServer(bindAddr string, bindPort int, ctx context.Context) {
	oscDispatcher := osc.NewStandardDispatcher()
	oscDispatcher.AddMsgHandler("/timer/start", oscHandleTimerStart)
	oscDispatcher.AddMsgHandler("/timer/stop", oscHandleTimerStop)
	oscDispatcher.AddMsgHandler("/timer/reset", oscHandleTimerReset)
	oscDispatcher.AddMsgHandler("/timer/stop_and_reset", oscHandleTimerStopAndReset)
	oscDispatcher.AddMsgHandler("/timer/addmsg", oscHandleAddMsg)
	if bindAddr == "" {
		// binding to all addresses
		/*
			ips, err := findMyIPs()
			if err != nil {
				log.Fatalf("Unable to get my IP addresses: %s", err)
			}
		*/
		for _, ip := range IPAddrs {
			log.Printf("OSC server listening on %s:%d", ip, bindPort)
		}
	} else {
		log.Printf("OSC server listening on %s:%d", bindAddr, bindPort)
	}

	oscServer := &osc.Server{
		Addr:       fmt.Sprintf("%s:%d", bindAddr, bindPort),
		Dispatcher: oscDispatcher,
	}

	go func() {
		log.Fatal(oscServer.ListenAndServe())
	}()
}

func oscHandleTimerStart(msg *osc.Message) {
	if msg.CountArguments() != 1 {
		log.Printf("Bad OSC /timer/start message: %s", msg)
		return
	}
	name := msg.Arguments[0].(string)
	err := Timers.Start(name)
	if err != nil {
		log.Print(err)
		return
	}
}

func oscHandleTimerStop(msg *osc.Message) {
	if msg.CountArguments() != 1 {
		log.Printf("Bad OSC /timer/stop message: %s", msg)
		return
	}
	name := msg.Arguments[0].(string)
	err := Timers.Stop(name)
	if err != nil {
		log.Print(err)
		return
	}
}

func oscHandleTimerReset(msg *osc.Message) {
	if msg.CountArguments() != 1 {
		log.Printf("Bad OSC /timer/reset message: %s", msg)
		return
	}
	name := msg.Arguments[0].(string)
	err := Timers.Reset(name)
	if err != nil {
		log.Print(err)
		return
	}
}

func oscHandleTimerStopAndReset(msg *osc.Message) {
	if msg.CountArguments() != 1 {
		log.Printf("Bad OSC /timer/stop_and_reset message: %s", msg)
		return
	}
	name := msg.Arguments[0].(string)
	err := Timers.Stop(name)
	if err != nil {
		log.Printf("Stop and reset failed: %s", err)
		return
	}
	err = Timers.Reset(name)
	if err != nil {
		log.Printf("Stop and reset failed: %s", err)
		return
	}
}

func oscHandleAddMsg(msg *osc.Message) {
	if msg.CountArguments() != 1 {
		log.Printf("Bad OSC /timer/addmsg message: %s", msg)
		return
	}
	m := msg.Arguments[0].(string)
	AddMessage(m)
	log.Printf("OSC message added: %s", m)
}
