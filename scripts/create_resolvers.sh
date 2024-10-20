#!/bin/sh

# List of TLDs and custom domains to support
supportedTLDs="ao ara epic fuck internal ki local localhost lore mail mi myth neko os pwn root test thc waifu"
customDomains="domains.internal"

# Check if the script is being run with sudo
if [ "$(id -u)" -ne 0 ]; then
  echo "This script must be run as root. Please enter your password."
  exec sudo "$0" "$@"
fi

# Create the /etc/resolver directory if it doesn't exist
if [ ! -d /etc/resolver ]; then
  echo "Creating /etc/resolver directory..."
  mkdir -p /etc/resolver
fi

# Loop through each TLD and create a resolver file
for tld in $supportedTLDs; do
  resolver_file="/etc/resolver/$tld"
  echo "Creating resolver for .$tld TLD..."
  
  echo "nameserver 127.0.0.1" > "$resolver_file"
  chmod 644 "$resolver_file"
done

# Create a resolver file for domains.internal
resolver_file="/etc/resolver/domains.internal"
echo "Creating resolver for domains.internal..."
echo "nameserver 127.0.0.1" > "$resolver_file"
chmod 644 "$resolver_file"

echo "Resolvers for the following TLDs have been created:"
echo $supportedTLDs
echo "Resolver for domains.internal has been created."
