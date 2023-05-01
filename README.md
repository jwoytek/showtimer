# Showtimer
## Introduction
Welcome to Showtimer--a simple, web-based, osc-triggered timekeeper
application for use in live productions. 

## Setup
The easiest way to get started is to use one of the release bundles.

(detailed instructions coming soon)

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
