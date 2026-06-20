#!/bin/sh
# Sidecar init: runs in an Alpine container (which has mkdir/cp/chown/sed — the
# dovecot/dovecot image is too minimal) to bootstrap the Maildir skeleton on the
# shared `dovecot-mail` volume and deliver the sample DMARC reports from
# ./contrib/dovecot/seed into INBOX/new (i.e. unseen — the fetcher searches for
# UNSEEN messages). Exits 0 when done so compose's service_completed_successfully
# gates the dovecot service on this finishing cleanly.
set -eu

MAILROOT=/srv/mail/dmarc/Maildir
SEEDED=/srv/mail/.seeded

mkdir -p "${MAILROOT}/cur" "${MAILROOT}/new" "${MAILROOT}/tmp"

# Sentinel keeps seeding idempotent: a restart must not re-deliver duplicates.
# Wipe the volume (`docker compose down -v`) to re-seed from scratch.
if [ ! -f "${SEEDED}" ]; then
  n=0
  for f in /import/seed/*.eml; do
    [ -e "$f" ] || continue
    n=$((n + 1))
    # RFC 5322 mandates CRLF line endings; seed files are stored LF-only so they
    # stay easy to edit, then converted on delivery.
    sed 's/$/\r/' "$f" > "${MAILROOT}/new/seed-${n}.eml"
  done
  touch "${SEEDED}"
fi

# UID:GID 1000 matches the vmail user inside the dovecot/dovecot image.
chown -R 1000:1000 /srv/mail
