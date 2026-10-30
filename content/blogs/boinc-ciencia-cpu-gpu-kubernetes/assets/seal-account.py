#!/usr/bin/env python3
"""Seal a Science United account login without writing plaintext to disk."""
import argparse
import getpass
import hashlib
import json
from pathlib import Path
import subprocess
import xml.etree.ElementTree as ET


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--context', required=True)
    parser.add_argument('--namespace', default='boinc-example')
    parser.add_argument('--controller-name', default='sealed-secrets-controller')
    parser.add_argument('--controller-namespace', default='sealed-secrets')
    parser.add_argument('--output', type=Path, default=Path('account.sealed.yaml'))
    args = parser.parse_args()
    if args.output.exists():
        raise SystemExit('Output already exists; preserve it before rotating credentials.')
    email = input('Science United email: ').strip().lower()
    password = getpass.getpass('Science United password: ')
    if not email or not password:
        raise SystemExit('Email and password must be nonempty.')
    login = ET.Element('acct_mgr_login')
    ET.SubElement(login, 'login').text = email
    # This hash is mandated by the account-manager protocol, not password storage advice.
    ET.SubElement(login, 'password_hash').text = hashlib.md5(
        (password + email).encode(), usedforsecurity=False
    ).hexdigest()
    secret = {'apiVersion': 'v1', 'kind': 'Secret',
              'metadata': {'name': 'boinc-science-united', 'namespace': args.namespace},
              'type': 'Opaque',
              'stringData': {'acct_mgr_login.xml': ET.tostring(login, encoding='unicode')}}
    sealed = subprocess.run(
        ['kubeseal', '--context', args.context, '--controller-name', args.controller_name,
         '--controller-namespace', args.controller_namespace, '--scope', 'strict', '--format', 'yaml'],
        input=json.dumps(secret), text=True, capture_output=True, check=False)
    if sealed.returncode:
        raise SystemExit('kubeseal failed; no credential file written. Check context and controller.')
    with args.output.open('x') as output:
        output.write(sealed.stdout)
    print(f'Sealed login written to {args.output}. No plaintext credential file created.')


if __name__ == '__main__':
    main()
