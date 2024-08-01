#!/usr/bin/env bash

export SSH_ASKPASS=./sshpass.sh
export SSH_ASKPASS_REQUIRE=force
export GIT_TRACE_PACKET=true
# Test transport
git clone foo://git@git-server:git/repo
