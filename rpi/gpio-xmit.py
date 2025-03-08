#!/usr/bin/env python3
# 
# gpio-xmit.py
#
# Read RPi GPI pins, build a byte to broadcast to any listening gpio-rcv.py. 
# Also periodically send a heartbeat so that the remote side can indicate a
# loss of comms.
#
# TODO: This was never fully realized due to changing needs for the
#       production. The script is left here to help those who may want to
#       implement this feature in the future. Currently, the script will
#       broadcast the heartbeat, and every five seconds, it will toggle the
#       first bit and send the packet.
#
# TODO: Configure with python logging for more appropriate log output.
# 

import socket
import time
try:
    from rpiozero import LED
except ModuleNotFoundError:
    print("Not running on a RPi!")
import config


def set_bit(byte_val, bit_index):
    """
    Set the bit at position 'bit_index' (0 is least-significant) to 1.
    :param byte_val: The current byte value (0-255)
    :param bit_index: The bit index to set (0-7)
    :return: New byte value with the specified bit set.
    """
    if not (0 <= bit_index <= 7):
        raise ValueError("bit_index must be between 0 and 7")
    return byte_val | (1 << bit_index)

def unset_bit(byte_val, bit_index):
    """
    Clear (set to 0) the bit at position 'bit_index'.
    :param byte_val: The current byte value (0-255)
    :param bit_index: The bit index to clear (0-7)
    :return: New byte value with the specified bit cleared.
    """
    if not (0 <= bit_index <= 7):
        raise ValueError("bit_index must be between 0 and 7")
    return byte_val & ~(1 << bit_index)

def main():
    # Broadcast address and port.
    BROADCAST_IP = '255.255.255.255'  # Alternatively, "255.255.255.255"
    PORT = config.network_port

    # Create a UDP socket and enable broadcasting.
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)

    byte_val = 0
    last_broadcasted_byte = None

    # Time management:
    toggle_interval = 5              # Toggle bit 0 every 5 seconds.
    last_toggle = time.time()         

    next_broadcast_time = time.time() + config.heartbeat_secs  # Next regular broadcast in seconds.

    print("Starting GPIO Xmit. Press Ctrl+C to stop.")

    try:
        while True:
            current_time = time.time()

            # Check if it's time to toggle bit 0.
            if current_time - last_toggle >= toggle_interval:
                if byte_val & 0x01:
                    byte_val = unset_bit(byte_val, 0)
                    print("Bit 0 cleared.")
                else:
                    byte_val = set_bit(byte_val, 0)
                    print("Bit 0 set.")
                last_toggle = current_time

                # Immediately broadcast the updated byte.
                if byte_val != last_broadcasted_byte:
                    data = byte_val.to_bytes(1, byteorder='big')
                    sock.sendto(data, (BROADCAST_IP, PORT))
                    print(f"Broadcasted byte immediately (change): {byte_val:08b}")
                    last_broadcasted_byte = byte_val

                    # Reset the regular broadcast timer.
                    next_broadcast_time = current_time + config.heartbeat_secs

            # Regular broadcast: send once per second.
            if current_time >= next_broadcast_time:
                data = byte_val.to_bytes(1, byteorder='big')
                sock.sendto(data, (BROADCAST_IP, PORT))
                print(f"Broadcasted byte (regular): {byte_val:08b}")
                last_broadcasted_byte = byte_val
                next_broadcast_time = current_time + config.heartbeat_secs

            # Short sleep to prevent busy waiting.
            time.sleep(0.1)
    except KeyboardInterrupt:
        print("Broadcast stopped by user.")
    finally:
        sock.close()

if __name__ == '__main__':
    main()

