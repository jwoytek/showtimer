#!/usr/bin/env python3
# 
# osc-rcv.py
#
# Start a simple OSC server, listening on the port configured in config.py.
# Currently, we respond only to onair. This could be extended with a little
# additional machinery to use a key in the bit config in config.py as the
# OSC command string. 
#
# TODO: Extend to honor all bits and auto-config OSC command endpoints based
#       on an additional key in the bit config.

# TODO: Configure with python logging for more appropriate log output.
# 

import socket
import time
from gpiozero import LED, Buzzer
import config
from pythonosc.dispatcher import Dispatcher
from pythonosc.osc_server import BlockingOSCUDPServer

# pins moved to a global here, unlike gpio-xmit and gpio-rcv.
pins = {}

def onair_handler(address, *args):
    if args[0] == 'on':
        print('On-air ON')
        pins[0].on()
    elif args[0] == 'off':
        print('On-air OFF')
        pins[0].off()
    elif args[0] == 'blink':
        print('On-air BLINK')
        pins[0].blink(.5, 1, n=5, background=True)
    else:
        print(f"*** Unknown onair command: {args}")


def default_handler(address, *args):
    print(f"*** Unknown message received: {address}: {args}")


def main():
    PORT = config.onair_osc_port

    # configure RPi pins
    for bit, bit_config in enumerate(config.bit_config):
        if bit_config['type'] == 'light':
            pins[bit] = LED(bit_config['gpo'], active_high=(True & bit_config['gpo_mode']), initial_value=(False & bit_config['gpo_mode']))
        elif bit_config['type'] == 'buzzer':
            pins[bit] = Buzzer(bit_config['gpo'], active_high=(True & bit_config['gpo_mode']), initial_value=(False & bit_config['gpo_mode']))
        else:
            print(f"Unknown bit type at position {bit}: {bit_config}")

    dispatcher = Dispatcher()
    dispatcher.map("/onair", onair_handler)
    dispatcher.set_default_handler(default_handler)

    server = BlockingOSCUDPServer(('', PORT), dispatcher)

    print(f"Starting on-air OSC listener on {PORT}")
    server.serve_forever()


if __name__ == '__main__':
    main()
