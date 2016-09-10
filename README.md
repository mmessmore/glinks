# glinks
Linux Kernel Stats in Go

This tool basically parses and presents all useful data I can find under /proc on Linux.

## Fun with math
This originally started with gathering this sort of data via Sensu for Graphite.  The Sensu plugins 
for checking CPU statistics are seriously flawed 
(eg. [sensu-plugins-cpu-checks](https://github.com/sensu-plugins/sensu-plugins-cpu-checks)).

Basically they screw up sampling royally by polling the values, sleeping a second, polling again, and calculating over 
the 1 second interval. It's a completely accurate measure over that second, but you really want to sample the entire 
time from the last run.

### Whawhat???
OK, so imagine you set the check interval at 10 seconds.  Basically you're measuring only the first second out of every
ten. Worst case you could have a workload that beats every ten seconds.  The check is going to either show the highest 
workload if the check coincides with the work, or lowest workload if they're out of phase.  The further apart the checks
are, the less representative your 1 second sample is.

But it's not just this weird workload, any workload that fluctuates will be inaccurately measured and graphed even 
worse. The fact that it's being graphed assumes that the check needs to cover the duration since the last, not just the 
next second.


## So why not fix their stuff?
So it goes back to the awesome sysstat package.  I really just wanted to use the venerable `sar` and co.  But their use 
is really intended to be interactive.  I don't want to push the data from every host, I want to pull it and be able to 
change intervals on the fly.

So I could rewrite the Sensu plugin and do all this work in Ruby or something.  But that just seems so single purpose.
Plus I don't want to fire up an interpreter for something that I may need to run at a pretty high frequency.

## I'm selfish
I also was just looking for a good excuse to use Go.  I'd written toys, but this seemed like a good task for it.

* I want it to be compiled so I don't have interpreter/JVM overhead
* I don't want to have to worry about the version of Ruby or Python on the host or that ships with Sensu or whatnot
* Statically linked is a real plus for easy distribution, no library hell across distros and versions.
* I don't want to introduce a lot of memory issues with the strings coming from these text files


## TODO

* Support running as a daemon with an REST GET-only interface
* Should I parse `/proc/zoneinfo` or anything else?
