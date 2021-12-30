#!/usr/bin/env python
import argparse
import json
import requests

import dns.resolver

TC_NUMBER = 0

def resolve(fqdn, nameserver, record_type="A"):
    resolver = dns.resolver.Resolver()
    resolver.nameservers = [nameserver]
    try:
        result = [ip.to_text() for ip in resolver.resolve(fqdn, record_type)]
    except dns.resolver.NoNameservers:
        result = None

    print(f"\t{fqdn} -> {result if result else []}")
    return result

def reset_rules(rest_ip):
    print("\tResetting rules")
    response = requests.get(f"http://{rest_ip}:8080/clear")

def add_rule(fqdn, rest_ip, action="BLOCK", record_type="A"):
    print(f"\tAdding rule: {action} {fqdn} {record_type}")
    response = requests.get(f"http://{rest_ip}:8080/insert?url={fqdn}&action={action}&type={record_type}")

def get_statistics(rest_ip):
    print(f"\tGetting statistics")
    response = requests.get(f"http://{rest_ip}:8080/stats")
    output = json.loads(response.content)
    print(f"\t{output}")
    return output

def print_header(description):
    global TC_NUMBER
    TC_NUMBER += 1
    print(f"========== TEST CASE {TC_NUMBER} ({description}) ==========")

parser = argparse.ArgumentParser(description='Jocasta Nu tester')
parser.add_argument('dns_ip', help='DNS server IP' )

args = parser.parse_args()
print_header("Blocked URL")
reset_rules(args.dns_ip)
result = resolve("www.google.pl", args.dns_ip)
assert result

print_header("One blocked URL and one logget URL")
add_rule("www.google.pl", args.dns_ip)
add_rule("www.onet.pl", args.dns_ip, action="LOG")
assert not resolve("www.google.pl", args.dns_ip)
assert resolve("www.onet.pl", args.dns_ip)

print_header("Blocked URL cleared by rest api")
reset_rules(args.dns_ip)
assert resolve("www.google.pl", args.dns_ip)

print_header("Mixing A and AAAA records")
reset_rules(args.dns_ip)
add_rule("www.google.pl", args.dns_ip, record_type="AAAA")
assert resolve("www.google.pl", args.dns_ip)
assert not resolve("www.google.pl", args.dns_ip, record_type="AAAA")
add_rule("www.google.pl", args.dns_ip, record_type="A")
assert not resolve("www.google.pl", args.dns_ip)
assert not resolve("www.google.pl", args.dns_ip, record_type="AAAA")

print_header("Generation of statistics")
reset_rules(args.dns_ip)
add_rule("www.google.pl", args.dns_ip)
resolve("www.google.pl", args.dns_ip)
stats = get_statistics(args.dns_ip)
assert stats
assert len(stats.get("requests", {})) == 1

print_header("Logging queries")
reset_rules(args.dns_ip)
add_rule("www.google.pl", args.dns_ip, action="LOG")
assert resolve("www.google.pl", args.dns_ip)
assert resolve("www.google.pl", args.dns_ip)
stats = get_statistics(args.dns_ip)
assert stats
assert len(stats.get("requests", {})) == 2
