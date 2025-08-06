#!/usr/bin/env bash
# a simple install script for ctop

KERNEL=$(uname -s)

# Color definitions (BLUE is defined but unused in the original script)
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

function log_warning() { 
    echo -e "${YELLOW}[WARNING]${NC} $1" 
}
function log_error() { 
    echo -e "${RED}[ERROR]${NC} $1"
}
function output() { 
    echo -e "${GREEN}[ctop-install]${NC} $*"
}

function command_exists() {
  command -v "$@" > /dev/null 2>&1
}

# extract github download url matching pattern
function extract_url() {
  local match=$1
  shift
  while read -r line; do
    case $line in
      *browser_download_url*"${match}"*)
        url=$(echo "$line" | sed -e 's/^.*"browser_download_url":[ ]*"//' -e 's/".*//;s/\ //g')
        echo "$url"
        break
      ;;
    esac
  done <<< "$*"
}

case $KERNEL in
  Linux) MATCH_BUILD="linux-amd64" ;;
  Darwin) MATCH_BUILD="darwin-amd64" ;;
  *)
    log_error "platform not supported by this install script"
    exit 1
    ;;
esac

for req in curl wget; do
  command_exists "$req" || {
    output "missing required $req binary"
    req_failed=1
  }
done
[ "$req_failed" = 1 ] && exit 1

sh_c='sh -c'
if [[ $EUID -ne 0 ]]; then
  if command_exists sudo; then
    log_warning "sudo is required to install ctop"
    sh_c='sudo -E sh -c'
  elif command_exists su; then
    log_warning "su is required to install ctop"
    sh_c='su -c'
  else
    log_error "This installer needs the ability to run commands as root. We are unable to find either sudo or su available to make this happen."
    exit 1
  fi
fi

TMP=$(mktemp -d "${TMPDIR:-/tmp}/ctop.XXXXX")
cd "${TMP}" || exit

output "fetching latest release info"
resp=$(curl -s https://api.github.com/repos/Betzalel75/ctop/releases/latest)

output "fetching release checksums"
checksum_url=$(extract_url sha256sums.txt "$resp")
wget -q "$checksum_url" -O sha256sums.txt

# skip if latest already installed
cur_ctop=$(command -v ctop 2> /dev/null)
if [[ -n "$cur_ctop" ]]; then
  cur_sum=$(sha256sum "$cur_ctop" | sed 's/ .*//')
  (grep -q "$cur_sum" sha256sums.txt) && {
    output "already up-to-date"
    exit 0
  }
fi

output "fetching latest ctop"
url=$(extract_url "$MATCH_BUILD" "$resp")
wget -q --show-progress "$url"
(sha256sum -c --quiet --ignore-missing sha256sums.txt) || exit 1

output "installing to /usr/local/bin"
chmod +x ctop-*
$sh_c "mv ctop-* /usr/local/bin/ctop"

output "done!"
