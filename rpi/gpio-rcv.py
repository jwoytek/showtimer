#!/usr/bin/env python3
# 
# gpio-rcv.py
#
# Receive a byte sent by gpio-xmit, and enable/disable RPi GPO pins according
# to the configuration in config.py
#
# TODO: Configure with python logging for more appropriate log output.
# 

import socket
import time
from gpiozero import LED, Buzzer
import config


def is_bit_set(byte_val, bit_index):
    """
    Check if the bit at 'bit_index' (0 is least-significant) is set in 'byte_val'.
    
    :param byte_val: The byte value (0-255).
    :param bit_index: The bit index to check (0-7).
    :return: True if the bit is set; False otherwise.
    """
    if not (0 <= bit_index <= 7):
        raise ValueError("bit_index must be between 0 and 7")
    return bool(byte_val & (1 << bit_index))

def main():
    PORT = config.network_port
    BUFFER_SIZE = 1024
    pins = {}

    # configure RPi pins
    for bit, bit_config in enumerate(config.bit_config):
        if 'type' not in bit_config:
            print(f"Bit with no type at position {bit}: skipping")
            continue

        if bit_config['type'] == 'light':
            pins[bit] = LED(bit_config['gpo'], active_high=(True & bit_config['gpo_mode']), initial_value=(False & bit_config['gpo_mode']))
        elif bit_config['type'] == 'buzzer':
            pins[bit] = Buzzer(bit_config['gpo'], active_high=(True & bit_config['gpo_mode']), initial_value=(False & bit_config['gpo_mode']))
        else:
            print(f"Unknown bit type at position {bit}: {bit_config}")

    # Create a UDP socket and bind to all network interfaces.
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.bind(('', PORT))
    
    # Set a timeout to periodically check for heartbeat (in seconds).
    sock.settimeout(1)
    print(f"Receiver listening on port {PORT}...")

    last_received_time = time.time()
    connection_lost = False  # Flag to track heartbeat state
    reset_all = False # Reset all pins to received state regardless of current value

    try:
        while True:
            try:
                data, addr = sock.recvfrom(BUFFER_SIZE)
                if data:
                    # Extract the broadcast byte (assuming one byte is sent)
                    byte_val = data[0]

                    # If the connection was previously flagged as lost, announce reestablishment.
                    if connection_lost:
                        print(f"[INFO] Packet received from {addr}. Connection reestablished.")
                        connection_lost = False
                        reset_all = True

                    # Update the last received timestamp.
                    last_received_time = time.time()

                    ## Print the received byte in binary.
                    #print(f"Received byte: {byte_val:08b} from {addr}")

                    for bit, bit_config in enumerate(config.bit_config):
                        if is_bit_set(byte_val, bit):
                            if pins[bit].value == 1 and not reset_all:
                                continue
                            print(f"Enable {bit_config['name']}")
                            pins[bit].on()
                            #GPIO.output(bit_config['gpo'], True & bit_config['gpo_mode'])
                        else:
                            if pins[bit].value == 0 and not reset_all:
                                continue
                            print(f"Disable {bit_config['name']}")
                            pins[bit].off()
                            #GPIO.output(bit_config['gpo'], False & bit_config['gpo_mode'])

                    reset_all = False

                    ## Optionally, check and print the state of each bit.
                    #for bit in range(8):
                    #    state = "set" if is_bit_set(byte_val, bit) else "not set"
                    #    print(f"  Bit {bit}: {state}")
                    #print("-" * 40)

            except socket.timeout:
                # No packet received in this cycle. Check if we've exceeded 5 seconds.
                if time.time() - last_received_time >= config.timeout_secs and not connection_lost:
                    print("[WARNING] No packet received for 5 seconds or more. Connection appears lost.")
                    connection_lost = True
                    pins[0].blink(.25, 5, background=True)

    except KeyboardInterrupt:
        print("Receiver stopped by user.")

    finally:
        sock.close()

if __name__ == '__main__':
    main()
