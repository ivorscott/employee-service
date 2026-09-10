#!/bin/sh
# Traffic generator for the observability lab.
#
# Produces a deliberate mix of response codes using only GET requests:
#   200 - a seeded employee id
#   404 - a well-formed uuid that does not exist
#   400 - a malformed id (fails uuid validation)
#
# GET is used on purpose. The PATCH path publishes to RabbitMQ and
# pkg/adapter/rabbitmq.go calls log.Fatal() if publishing fails, which would
# kill the service mid-lab.

BASE="http://employee-service:8080"

# Seeded ids from res/seed/data.sql (run `docker compose run --rm seed` first,
# otherwise these return 404 too).
SEEDED_A="bc4cd1a1-4f0e-4e39-9960-e6b1cfe388db"
SEEDED_B="35d27bb4-c39e-4c10-9a64-aabc2490ec4d"

# A syntactically valid uuid that is not in the database.
MISSING="00000000-0000-4000-8000-000000000000"

# Not a uuid at all.
MALFORMED="not-a-valid-uuid"

echo "loadgen: starting against ${BASE}"

i=0
while true; do
  i=$((i + 1))

  # Weight the mix so 200s dominate, with a steady trickle of 404s and 400s.
  case $((i % 10)) in
    0|1|2|3) TARGET="$SEEDED_A" ;;
    4|5|6)   TARGET="$SEEDED_B" ;;
    7|8)     TARGET="$MISSING" ;;
    *)       TARGET="$MALFORMED" ;;
  esac

  curl -s -o /dev/null -w "%{http_code} /employees/${TARGET}\n" \
    "${BASE}/employees/${TARGET}" || true

  sleep 2
done
