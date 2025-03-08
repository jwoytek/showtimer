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
	t, ok := Timers[name]
	if !ok {
		log.Printf("Asked to start unknown timer: %s", name)
	}
	log.Printf("OSC timer start for %s", name)
	t.Start()
}

func oscHandleTimerStop(msg *osc.Message) {
	if msg.CountArguments() != 1 {
		log.Printf("Bad OSC /timer/stop message: %s", msg)
		return
	}
	name := msg.Arguments[0].(string)
	t, ok := Timers[name]
	if !ok {
		log.Printf("Asked to stop unknown timer: %s", name)
	}
	log.Printf("OSC timer stop for %s", name)
	t.Stop()
}

func oscHandleTimerReset(msg *osc.Message) {
	if msg.CountArguments() != 1 {
		log.Printf("Bad OSC /timer/reset message: %s", msg)
		return
	}
	name := msg.Arguments[0].(string)
	t, ok := Timers[name]
	if !ok {
		log.Printf("Asked to reset unknown timer: %s", name)
	}
	log.Printf("OSC timer reset for %s", name)
	t.Reset()
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
