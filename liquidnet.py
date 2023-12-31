#!/bin/python3

import ipaddress
import random
import argparse

def gather_random_ipv4():
    return '.'.join(str(random.randint(0, 255)) for _ in range(4))

def gather_valid_ipv4_addresses(num_addresses):
    valid_addresses = set()
    while len(valid_addresses) < num_addresses:
        ip = gather_random_ipv4()
        try:
            ipaddress.IPv4Address(ip)
            valid_addresses.add(ip)
        except ipaddress.AddressValueError:
            pass
    return list(valid_addresses)

def main():
    parser = argparse.ArgumentParser(description='Gather random valid IPv4 addresses')
    parser.add_argument('-n', dest='num_addresses', type=int, required=True, help='Number of IP addresses to gather')
    args = parser.parse_args()

    num_addresses = args.num_addresses
    valid_ips = gather_valid_ipv4_addresses(num_addresses)
    for ip in valid_ips:
        print(ip)

if __name__ == "__main__":
    main()
