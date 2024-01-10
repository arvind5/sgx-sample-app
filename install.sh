#!/bin/bash

# Install required packages
sudo apt-get install build-essential -y
sudo apt-get install gnupg wget -y

# Configure SGX repos
sudo echo 'deb [arch=amd64] https://download.01.org/intel-sgx/sgx_repo/ubuntu focal main' > /etc/apt/sources.list.d/intel-sgx.list
wget -qO - https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key | apt-key add -

# Update
sudo apt-get update

# Fetch sgxsdk and golang installer
wget https://download.01.org/intel-sgx/sgx-linux/2.19/distro/ubuntu20.04-server/sgx_linux_x64_sdk_2.19.100.3.bin
wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz

# Install sgxsdk and golang
chmod +x sgx_linux_x64_sdk_2.19.100.3.bin
./sgx_linux_x64_sdk_2.19.100.3.bin --prefix=/opt/intel
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz

# Remove downloads
rm sgx_linux_x64_sdk_2.19.100.3.bin
rm go1.21.1.linux-amd64.tar.gz

# Install DCAP QuoteWrapper Package
sudo apt-get install libsgx-dcap-ql=1.16.100.2-focal1 libsgx-dcap-ql-dev=1.16.100.2-focal1 -y