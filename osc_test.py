#!/usr/bin/env python3
"""
OSC test client for showtimer.
Sends /timer/start and /timer/stop commands using only Python stdlib.
"""

import argparse
import socket
import struct


def osc_string(s: str) -> bytes:
    """Encode a string as an OSC string (null-terminated, padded to 4 bytes)."""
    encoded = s.encode('utf-8') + b'\x00'
    # Pad to multiple of 4 bytes
    padding = (4 - len(encoded) % 4) % 4
    return encoded + b'\x00' * padding


def osc_message(address: str, *args) -> bytes:
    """Build an OSC message with the given address and arguments."""
    # Address pattern
    msg = osc_string(address)

    # Type tag string
    type_tags = ','
    arg_data = b''

    for arg in args:
        if isinstance(arg, str):
            type_tags += 's'
            arg_data += osc_string(arg)
        elif isinstance(arg, int):
            type_tags += 'i'
            arg_data += struct.pack('>i', arg)
        elif isinstance(arg, float):
            type_tags += 'f'
            arg_data += struct.pack('>f', arg)

    msg += osc_string(type_tags)
    msg += arg_data

    return msg


def send_osc(host: str, port: int, address: str, *args):
    """Send an OSC message via UDP."""
    msg = osc_message(address, *args)

    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        sock.sendto(msg, (host, port))
        print(f"Sent OSC: {address} {args if args else ''}")
    finally:
        sock.close()


def main():
    parser = argparse.ArgumentParser(description='Send OSC commands to showtimer')
    parser.add_argument('--host', default='127.0.0.1', help='OSC server host (default: 127.0.0.1)')
    parser.add_argument('--port', type=int, default=8000, help='OSC server port (default: 8000)')
    parser.add_argument('--start', action='store_true', help='Send /timer/start [timer]')
    parser.add_argument('--stop', action='store_true', help='Send /timer/stop [timer]')
    parser.add_argument('--reset', action='store_true', help='Send /timer/reset [timer]')
    parser.add_argument('timer', help='timer name')

    args = parser.parse_args()

    #if not args.message and not args.timers:
    #    parser.error('At least one of --message or --timers must be specified')

    if args.start:
        send_osc(args.host, args.port, '/timer/start', args.timer)
    if args.stop:
        send_osc(args.host, args.port, '/timer/stop', args.timer)
    if args.reset:
        send_osc(args.host, args.port, '/timer/reset', args.timer)


if __name__ == '__main__':
    main()
