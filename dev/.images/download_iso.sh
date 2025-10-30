#!/usr/bin/env bash
set -euo pipefail

# Download a Debian ISO (latest or specific version) and verify its SHA256 checksum.

VERSION="latest"
DEST_DIR="$PWD"
MEDIA_TYPE="cd"
DEFAULT_VARIANT_CD="netinst"
DEFAULT_VARIANT_DVD="DVD-1"
ISO_VARIANT="$DEFAULT_VARIANT_CD"
VARIANT_OVERRIDDEN=0

usage() {
	cat <<'EOF'
Usage: download_iso.sh [--version <version>|latest] [--dest <directory>] [--media cd|dvd] [--variant <iso-suffix>]

Downloads the requested Debian amd64 ISO (defaults to latest netinst CD) and verifies its integrity.
Examples:
  download_iso.sh
  download_iso.sh --version 12.5.0 --dest /tmp/debian
	download_iso.sh --media dvd
	download_iso.sh --media dvd --variant DVD-2
EOF
}

log() {
	printf '[%s] %s\n' "$(date -u +'%Y-%m-%dT%H:%M:%SZ')" "$*"
}

require_command() {
	if ! command -v "$1" >/dev/null 2>&1; then
		printf 'Error: required command "%s" not found.\n' "$1" >&2
		exit 1
	fi
}

cleanup() {
	rm -rf "${TMP_DIR:-}"
}

parse_args() {
	while [[ $# -gt 0 ]]; do
		case "$1" in
			-h|--help)
				usage
				exit 0
				;;
			-v|--version)
				VERSION="${2:-}"
				if [[ -z "$VERSION" ]]; then
					printf 'Error: --version expects a value.\n' >&2
					exit 1
				fi
				shift 2
				;;
			--dest)
				DEST_DIR="${2:-}"
				if [[ -z "$DEST_DIR" ]]; then
					printf 'Error: --dest expects a value.\n' >&2
					exit 1
				fi
				shift 2
				;;
			--variant)
				ISO_VARIANT="${2:-}"
				if [[ -z "$ISO_VARIANT" ]]; then
					printf 'Error: --variant expects a value.\n' >&2
					exit 1
				fi
				VARIANT_OVERRIDDEN=1
				shift 2
				;;
			--media)
				MEDIA_TYPE="${2:-}"
				if [[ -z "$MEDIA_TYPE" ]]; then
					printf 'Error: --media expects a value.\n' >&2
					exit 1
				fi
				shift 2
				;;
			*)
				printf 'Error: unknown argument "%s".\n' "$1" >&2
				usage >&2
				exit 1
				;;
		esac
	done
}

main() {
	parse_args "$@"

	require_command curl
	require_command sha256sum

	TMP_DIR="$(mktemp -d)"
	trap cleanup EXIT

	mkdir -p "$DEST_DIR"

	local base_url="https://cdimage.debian.org/debian-cd"
	local version_path
	MEDIA_TYPE="${MEDIA_TYPE,,}"
	if [[ "$MEDIA_TYPE" != "cd" && "$MEDIA_TYPE" != "dvd" ]]; then
		printf 'Error: --media must be either "cd" or "dvd".\n' >&2
		exit 1
	fi

	if [[ "$MEDIA_TYPE" == "dvd" && "$VARIANT_OVERRIDDEN" -eq 0 ]]; then
		ISO_VARIANT="$DEFAULT_VARIANT_DVD"
	elif [[ "$MEDIA_TYPE" == "cd" && "$VARIANT_OVERRIDDEN" -eq 0 ]]; then
		ISO_VARIANT="$DEFAULT_VARIANT_CD"
	fi

	if [[ "$VERSION" == "latest" ]]; then
		version_path="current"
	else
		version_path="$VERSION"
	fi

	local sha_url="${base_url}/${version_path}/amd64/iso-${MEDIA_TYPE}/SHA256SUMS"
	local sha_local="${TMP_DIR}/SHA256SUMS"
	log "Fetching checksum manifest from ${sha_url}"
	curl --fail --location --silent --show-error "$sha_url" --output "$sha_local"

	local iso_name
	iso_name="$(awk -v variant="$ISO_VARIANT" 'BEGIN{IGNORECASE=1} $2 ~ variant".iso$" {print $2; exit}' "$sha_local")"
	if [[ -z "$iso_name" ]]; then
		printf 'Error: could not find an ISO matching variant "%s" in SHA256SUMS.\n' "$ISO_VARIANT" >&2
		exit 1
	fi

	local iso_url="${base_url}/${version_path}/amd64/iso-${MEDIA_TYPE}/${iso_name}"
	local iso_path="${DEST_DIR}/${iso_name}"

	log "Downloading ISO ${iso_name}"
	curl --fail --location --continue-at - --show-error "$iso_url" --output "$iso_path"

	local sha_subset="${TMP_DIR}/SHA256SUMS.selected"
	grep " ${iso_name}$" "$sha_local" > "$sha_subset"

	log "Verifying checksum"
	(cd "$DEST_DIR" && sha256sum --check "$sha_subset")

	log "ISO stored at ${iso_path}"
	log "Checksum verification successful"

	if command -v gpg >/dev/null 2>&1; then
		local sig_url="${sha_url}.sign"
		local sig_local="${TMP_DIR}/SHA256SUMS.sign"
		log "Optional: verifying signature of checksum (requires imported Debian signing key)"
		if curl --fail --location --silent --show-error "$sig_url" --output "$sig_local"; then
			if gpg --verify "$sig_local" "$sha_local" >/dev/null 2>&1; then
				log "GPG signature verification succeeded"
			else
				log "Warning: GPG signature verification failed; ensure Debian CD signing key is imported"
			fi
		else
			log "Warning: could not download signature file from ${sig_url}"
		fi
	else
		log "Tip: install gpg to verify checksum signatures"
	fi
}

main "$@"
