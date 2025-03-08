# Raspberry Pi Extras
This directory contains some helpers and extra features developed for use on a Raspberry Pi. 

The goal here was originally to build a set of kiosk-like displays, one of which could also serve as the Showtimer server.

Included here:
  - `gpio-rcv.py`, `gpio-xmit.py`: Scripts to send and receive GPIO pin status.
  - `osc-rcv.py`: Script to run a simple OSC server to operate GPO pins.
  - `config.py`: Shared configuration.
  - `labwc/autostart`: Kiosk mode autostart script for showtimer.
  - `systemd/`: systemd service configuration for showtimer, GPIO, and OSC modules.

## Raspberry Pi Configuration
### Components
#### Computing Components
These scripts and configuration files were built to utilize Raspberry Pi 5 SBCs, running Raspberry Pi OS 64-bit (release date as of this writing is 2024-11-19).

#### Displays
The displays used for this project are these [10.1" Touchscreen 1280x800 panels](https://www.amazon.com/dp/B0CLQZ48BF). Conveniently, they include a mount on the back in the correct configuration for a Raspberry Pi, and include standoffs and screws to accomplish the mounting. They also include all of the necessary cabling to connect directly to a Raspberry Pi, even down to the micro-HDMI to HDMI adapter. They are USB powered and are fine drawing power directly from the Pi USB port. 

#### HCI
Any USB keyboard and mouse can be helpful for initial configuration and testing, and for rescue if network configuration goes sideways. For an all-in-one solution, the [Logitech K400](https://www.amazon.com/dp/B014EUQOGK) is a reasonable option that combines touchpad and keyboard.

### OS Configuration
The Raspberry Pi boards are imaged with Raspberry Pi OS 64-bit, using a predetermined hostname, username (`showtimer`, for example), and SSH key authentication. One can also use password-based authentication. For the purposes of this effort, it doesn't particularly matter, as direct login is only used for setup. 

Hostnames will be used later, so do use sensible things here. For the purposes of this document:
  - Main Showtimer display: `showtimer-main`
  - Orchestra display: `showtimer-orchestra`
  - SM display: `showtimer-sm`

### Network Configuration
These systems can be set up on WiFi or hardwired ethernet, depending on use case. Raspberry Pi OS uses NetworkManager (and its command-line component `nmcli`) for interface management. This guide won't dive into the details on configuring network interfaces (with the exception of the caveat on link-local addressing noted below). There are lots of helpful guides out there already, and the GUI tool works well enough for most configuration uses. During imaging, it is possible to configure the WiFi during imaging to aid in rapid setup.

#### Link-Local (zeroconf) Addressing
Addressing in these examples use local name resolution (like, `showtimer-main.local`). If one is operating on a network without DHCP, it is possible to use link-local addressing for zeroconfig-type operation. Under the version of NetworkManager running on Raspberry Pi OS, however, link-local IPv4 addresses appear to not be automatically configured when DHCP fails. Some of these components do not bind correctly to IPv6 addresses and require IPv4 (this has not been extensively investigated as of this time).

In order to configure the hardwired interface to try DHCP first, and then fail-over to a zeroconf-style link-local address, the following `nmcli` commands can be entered in a terminal. Be aware that your connection may be lost if you do this while you are using SSH to get into the system, so you may want to do this from a local keyboard/mouse before setting up the kiosk display. 

```
# Add a DHCP-enabled eth0 connection type with a shortened
# retry cycle and shortened timeout.
nmcli conn add type ethernet ifname eth0 con-name eth0-auto
nmcli conn mod eth0-auto connection.autoconnect-priority 100
nmcli conn mod eth0-auto connection.autoconnect-retries 3
nmcli conn mod eth0-auto ipv4.dhcp-timeout 3

# Add a link-local eth0 connection at a lower priority than DHCP
nmcli conn add type ethernet ifname eth0 con-name eth0-ll
nmcli conn mod eth0-ll connection.autoconnect-priority 50 ipv4.method link-local
```

After completing these steps, one can reboot the system and a hardwired connection will first attempt to use DHCP, and if DHCP fails, it will fall-back to a link-local address. 

If you are working on a network where you expect to use link-local addressing all of the time, you can omit the eth0-auto portion, which will result in only a link-local connection that will come up almost immediately. This can be useful if, say, one is hosting these devices on a Dante network. 

## Showtimer Server Configuration
These steps will configure the Showtimer server on one of the RPis. For the purposes of this document, the server will be referred to as `showtimer-main`. 

Download and copy the [latest Linux ARM64 release build of Showtimer](https://github.com/jwoytek/showtimer/releases/latest/download/showtimer-linux-arm64.zip) to `showtimer-main`:

```
curl -o showtimer.zip -L "https://github.com/jwoytek/showtimer/releases/latest/download/showtimer-linux-arm64.zip"
scp showtimer.zip showtimer@showtimer-main:
```

(Note: If your RPi has internet access, you could `curl` it directly onto the RPi instead of fetching and then scp'ing.)

Connect to `showtimer-main` (either via SSH or with a local keyboard and mouse) and use a terminal to unzip and configure Showtimer:

```
unzip -x showtimer.zip

cd showtimer

# use your favorite editor here (nano, pico, vi, etc.)
vi showtimer.yaml
```

Once you are happy with your configuration for this show, save the `showtimer.yaml` file. 

Now it's time for a quick test! Start showtimer in a terminal window and make sure that it starts correctly and says it is listening on one or more interfaces:

```
./showtimer
```

You should see something like this:

```
showtimer@showtimer-sm:~/showtimer $ ./showtimer 
2025/03/08 14:50:40 Creating new timer 'Act 1' with initial duration of 0.000000s
2025/03/08 14:50:40 Creating new timer 'Intermission' with initial duration of 60.000000s
2025/03/08 14:50:40 Creating new timer 'Act 2' with initial duration of 0.000000s
2025/03/08 14:50:40 Webserver listening on 127.0.0.1:8080
2025/03/08 14:50:40 Webserver listening on 192.168.200.173:8080
2025/03/08 14:50:40 OSC server listening on 127.0.0.1:8000
2025/03/08 14:50:40 OSC server listening on 192.168.200.173:8000
```

The IP addresses, names of timers, and durations might be different. That's OK. 

Type [CTRL]-c to quit Showtimer. 

Now it is time to set up the service so that Showtimer starts automatically when the system boots. The next steps require use of `sudo` to operate with root privileges on the RPi. Be careful while doing these steps! 

First, the systemd service configuration file for Showtimer needs to be added to `/etc/systemd/system/` on the RPi. There are multiple ways to do this (scp from another computer to the RPi and then `sudo cp` to put it in place, or edit directly on the RPi using `sudo vi /etc/systemd/system/showtimer.service` (or your favorite editor instead of vi) and copy/paste, or `curl` the file onto the system directly, etc.). For the purposes of this document, we will illustrate using curl, scp, and cp:

```
# On a computer with internet access and SSH:
curl -o showtimer.service -L "https://raw.githubusercontent.com/jwoytek/showtimer/refs/heads/main/rpi/systemd/showtimer.service"

scp showtimer.service showtimer@showtimer-main:

# On the RPi itself:
sudo cp showtimer.service /etc/systemd/system/
```

Now we can start the service to make sure it starts up correctly:

```
# on the RPi
sudo systemctl start showtimer
```

To make sure it started: 
```
sudo systemctl status showtimer
```

You should see that it started correctly and is listed as running. If not, use `journalctl` to see what happened:
```
sudo journalctl -xe
```

If it started correctly, you can now set it to start automaticall on boot:
```
sudo systemctl enable showtimer
```

Finally, point a web browser at http://showtimer-main.local:8080. You should get the Showtimer timer page! If not, it is time to start digging backwards to see what is broken. As a secondary test, start a web browser directly on the RPi and try the same URL. If that fails, also try http://localhost:8080. Debugging past this point is beyond the scope of this document. 

## Showtimer Kiosk Configuration
Configuring the RPis to act as kiosks is easy. You will need to know the name of the system running the Showtimer server (in this example, that will be `showtimer-main.local`). 

The standard options for the Raspberry Pi OS imager has the system auto-login as the configured user. This is already half of the configuration needed for kiosk mode. The second half is to create an autostart configuration file to start a browser in kiosk mode. 

Use your favorite editor (`vi` used in this example, but `nano` or `pico` also work) on the RPi to create the file `.config/labwc/autostart`:

```
# On the RPi
cd ~
vi .config/labwc/autostart
```

Paste or enter in the following into the file:
```
/usr/bin/lwrespawn chromium-browser http://showtimer-main.local:8080/ --kiosk --start-maximized --noerrdialogs --disable-infobars --no-first-run --ozone-platform=wayland
```

If your server is on a different hostname or port other than 8080, change those portions of the file, then save the file. 

Restart the RPi, and the timer should start up automatically. If you initially get a white screen saying that the site can't be reached, it means the browser tried to start before the server was available. Try tapping the blue "reload" button. 

```
# use the following to reboot cleanly
sudo reboot
```

## GPIO Configuration
For the purposes of a recent production, we needed an on-air light for a remote orchestra room, and wanted the ability to add other physical light and buzzer outputs in the future. To accomplish this, the first pass was a broadcast/receiver concept, where a module would run to read the state of GPI pins on the RPi, and then broadcast that state on the network to any listening receivers. This service was modeled in the `gpio-rcv.py` and `gpio-xmit.py` scripts. The initial plan was to use a GPO pin on the Yamaha mixing desk to drive a pin low when the orchestra fader was brought up, activating the light in the remote orchestra room. The idea was later expanded, though, and a more generic solution that did not require that a specific console and GPO interface be used, so `osc-rcv.py` was created to drive the on-air light via OSC. 

The gpio-rcv and xmit scripts are commented and present in the repository for anyone who would like to go down the route of using a broadcast/receiver type of system, driven by buttons or switches or any other sort of physical I/O into a RPi. This document will cover configuring the OSC service, but the other service configuration files are present as a base for further work if desired. 

### GPIO Components
For anything that needs to interface with GPO pins, the [Adafruit Pi-EZConnect Breakout](https://www.adafruit.com/product/2711) was used to make experimentation and interfacing a little easier.

To operate the light, we used an [5V opto-isolated relay with screw terminal blocks](https://www.amazon.com/dp/B07TWH7DZ1) to switch 110VAC line current to run a standard medium base lightbulb. The relay board was attached inside an outdoor-type duplex box with a solid cover, and a lamp socket screwed into one of the knock-outs. Why an outdoor box? It was available and made adding the lamp socket easy. 

The relay board requires 5V power, which was sourced from the RPi. In this case, a length of balanced audio cable was used to connect the RPi GPO pin, 5VDC, and GND from the RPi to the relay board. 

Please consult with an electrician regarding proper connections and safety for wiring the AC side of the system to the relay and light socket. That will not be covered here. 

## OSC Server Setup
Copy the `osc-rcv.py` and `config.py` files to the RPi. For the purposes of this document, the script was copied into the `showtimer` directory. 

As above for the Showtimer service setup, copy or paste the contents of `systemd/osc-rcv.service` into `/etc/systemd/system/osc-rcv.service`. 

The RPi will need to have the python-osc library installed. The method covered here will result in big warnings from the RPi Python subsystem that you may be breaking the world. This does not appear to be an issue, and this method is basically required to install the library on a system without internet access. 

First, download the built [python-osc library wheel](https://pypi.org/project/python-osc/#files) on a system with internet connectivity. As of the time of this writing, the most recent version is 1.9.3. Make sure that you get the `.whl` file.

Next, copy this file to the RPi:

```
scp python_osc-1.9.3-py3-none-any.whl showtimer@showtimer-orchestra.local
```

Finally, install it on the RPi:

```
sudo pip install --no-index --find-links . python-osc --break-system-packages
```

Now it's time to make sure that the OSC server is configured correctly and will start. 

```
# on the RPi, check the configuration for port and pin numbers:
cd showtimer
vi config.py

# now, try to start it:
./osc-rcv.py
```

You should see messages indicating that the server is started and listening for OSC messages. 

Use [CTRL]-c to stop the script, then start the service using systemd:

```
sudo systemctl start osc-rcv

# check to make sure it started successfully
sudo systemctl status osc-rcv
```

If it has started correctly, enable it so that it starts on boot:

```
sudo systemctl enable osc-rcv
```

### Controlling the On-Air Light
To control the on-air light, set up an OSC application (QLab, etc.) to send OSC to `showtimer-orchestra.local` (or whatever system on which you installed the OSC server), port 8001 (make sure this matches what you had in `config.py`). 

Valid commands:
  - `/onair on`: Turn on the light
  - `/onair off`: Turn off the light
  - `/onair blink`: Blink the light five times, once per second, 50% duty cycle (on .5 seconds, off .5 seconds)
