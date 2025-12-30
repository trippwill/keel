#!/usr/bin/env bash
set -euo pipefail

cpu=$(grep -m1 'model name' /proc/cpuinfo | cut -d: -f2 | xargs)
mem=$(grep MemTotal /proc/meminfo | awk '{print $2}')
{
  printf '%s\n' "# Keel Benchmarks"
  printf '\n'
  printf '%s\n' "## Environment"
  printf '%s\n' "- $(go version)"
  printf '%s\n' "- cpu: $cpu"
  printf '%s\n' "- mem: $((mem / 1024)) MB"
} >BENCHMARKS.md
