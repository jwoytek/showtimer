# Showtimer
## Introduction
Welcome to Showtimer--a simple, web-based, osc-triggered timekeeper
application for use in live productions. 

## Setup
The easiest way to get started is to use one of the release bundles.

### Download
Pick the binary bundle appropriate for your operating system and
architecture:
- darwin/arm64 for: MacOS on Apple Silicon
- darwin/amd64 for: MacOS on Intel processors
- windows/amd64 for: Windows on Intel processors
- linux/amd64 for: Linux on Intel processors
- linux/arm64 for: Linux on Apple Silicon or other ARM CPU

### Setup and Use:
Currently, most setup and startup of the program must be done from the
command line (Windows CMD.EXE, MacOS Terminal, etc.). 

1. unzip the bundle in a useful location, like your home directory
2. Edit the `showtimer.yaml` file with your favorite editor
3. Set the port values for the webserver and OSC (or note the defaults)
4. Optionally, if you would like to only listen on a certain interface,
   enter the IP address of that interface. This is probably not common.
   Leaving the bind address commented out will cause the services to listen
   on all available interfaces. This is the default.
5. Edit the timer list if necessary. 
    - `name`: Friendly name of the timer, can contain spaces and punctuation
    - `key`: Short name used in OSC commands. Must not contain punctuation or spaces.
    - `type`: "CountDown" for a countdown timer; "CountUp" for an elapsed timer
    - `duration`: __REQUIRED__ for CountDown timers, this specifies the 
      initial value of the timer. This is set in human readable form, like 
      "15m30s" is 15 minutes and 30 seconds; "1h30m" is one hour and 30 minutes, etc.
      **Note** that duration is ignored for CountUp timers, which always start 
      at 00:00:00.
6. Start the showtimer by running `showtimer` in the showtimer directory.
7. Point a web browser at one of the web server addresses printed to the 
   terminal during startup. You should see the timers specified in the
   configuration file.
8. Configure your OSC client to send a message to one of the address:port
   combinations printed to the terminal during startup for the OSC server.
   The message should be in the format `/timer/start [key]` where "[key]" is
   the key you entered in the config file. You should see the timer start
   to run in your web browser. 

### MACOS USERS PLEASE READ
MacOS does some smart things and will by default not let you run random
things you download from the internet. This is usually good. However, if
you are trying to run `showtimer`, you'll get a dialog saying that Apple
cannot check the program for malicious software. Your only choices are to
quit or show the program in the Finder. 

In order to get past this warning and tell MacOS that you think the program
is safe to run, click the option to show the program in the finder. Then
right-click (or option-click or two-finger-click) the `showtimer` binary, 
and select "Open" from the menu. MacOS will ask you if you are sure you 
want to do this. If you are sure, click OK and the program will run and
immediately exit because it couldn't find the configuration file. This is
OK. Close the window where the program started, and now you can start it
from a regular terminal window. 

## Available OSC commands
All commands take a `key` argument, which is the key for the timer in the
configuration file. 

* `/timer/start [key]`
    * Start the specified timer
* `/timer/stop [key]`
    * Stop the specified timer
* `/timer/reset [key]`
    * Reset the specified timer

## Available API commands
The web API can be used to retrieve timer values directly to use them in 
other places.
* `GET /timer/[key]`
    * Returns a JSON structure with timer data for the specified timer

## Background
Showtimer arose from an offhand request from one of the finest techies
with whom I work when he saw me messing about with some script magic
in QLab to automatically stop and stop multitrack recording sessions
during rehearsals. During a run of shows we did together in 2022, he 
asked if I could come up with something to help track runtimes for 
shows that would trigger on my QLab board cues that I fire at the 
beginning and end of every act and intermission. 

I poked around and put together some QLab script cues based on work
by a few others trying to do similar things, and used Open Stage
Control as a web front-end to share the timers to production 
departments. It worked and was quite useful during the show, but it
depended on my QLab instance being up and running, and times would
not survive a QLab restart. This was sub-optimal. 

A new solution was needed, and I couldn't find much out there that
would answer the need. Specifically, we wanted something simple
that was standalone, could be controlled via OSC, and had the ability
to handle both count-up (runtime/elapsed) and count-down timers. It
had to have a web front end to make it easy for others to get the 
timers. Sub-second synchronization among viewers was not necessary,
though it would be nice to get there eventually. 

"Showtimer" was born, written in Go, with a simple configuration 
file.
