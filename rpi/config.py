# define pins and ports for gpio

# network port for GPIO xmit/rcv
network_port = 8817
# network port for on-air OSC receiver
onair_osc_port = 8001
# timeout for GPIO rcv, after which it will indicate loss of comms
timeout_secs = 5
# delay between heartbeat packets sent from GPIO xmit
heartbeat_secs = 1

# bit config is in order 0 - 7
# currently, there is no catch to skip empty bits
bit_config = [
    {
        'name': 'On Air',
        'type': 'light',
        'gpi': 4, # BCM pin number for gpio-xmit to monitor
        'gpo': 4, # BCM pin number for gpio-rcv to monitor
        'gpi_mode': 0, # goes LOW when active
        'gpo_mode': 1  # goes HIGH when active

    },
    {
        'name': 'Orchestra Buzzer',
        'type': 'buzzer',
        'gpo': 17,
        'gpo_mode': 1
    },
    {
        'name': 'SM Buzzer',
        'type': 'buzzer',
        'gpo': 22,
        'gpo_mode': 1
    }
]