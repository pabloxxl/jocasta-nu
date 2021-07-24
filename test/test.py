#!/usr/bin/env python
import argparse
import json
import pprint
import requests

import dns.resolver

def resolve(fqdn, nameserver, record_type="A"):
    resolver = dns.resolver.Resolver()
    resolver.nameservers = [nameserver]
    try:
        result = [ip.to_text() for ip in resolver.resolve(fqdn, record_type)]
    except dns.resolver.NoNameservers:
        result = None

    print(f"{fqdn} -> {result if result else []}")
    return result

def reset_rules(rest_ip):
    print("Resetting rules")
    response = requests.get(f"http://{rest_ip}:8080/clear")

def add_rule(fqdn, rest_ip, action="BLOCK", record_type="A"):
    print(f"Adding rule: {action} {fqdn} {record_type}")
    response = requests.get(f"http://{rest_ip}:8080/insert?url={fqdn}&action={action}&type={record_type}")

def stats(rest_ip):
    print(f"Getting statistics")
    response = requests.get(f"http://{rest_ip}:8080/stats")
    return json.loads(response.content)

parser = argparse.ArgumentParser(description='Jocasta Nu tester')
parser.add_argument('dns_ip',
    help='DNS server IP' )

args = parser.parse_args()
reset_rules(args.dns_ip)
result = resolve("www.google.pl", args.dns_ip)
assert result

add_rule("www.google.pl", args.dns_ip)
add_rule("www.onet.pl", args.dns_ip, action="LOG")
result = resolve("www.google.pl", args.dns_ip)
assert not result
result = resolve("www.onet.pl", args.dns_ip)
assert result

pprint.pprint(stats(args.dns_ip))
reset_rules(args.dns_ip)

result = resolve("www.google.pl", args.dns_ip)
assert result
