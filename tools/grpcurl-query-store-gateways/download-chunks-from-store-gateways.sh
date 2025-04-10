#!/usr/bin/env bash
# SPDX-License-Identifier: AGPL-3.0-only

# Begin of configuration.
K8S_CONTEXT=""
K8S_NAMESPACE=""
MIMIR_TENANT_ID=""
# End of configuration.

SCRIPT_DIR=$(realpath "$(dirname "${0}")")
OUTPUT_DIR="chunks-dump"
# Randomize port and start from higher ports, so we don't collide with random apps on your laptop.
NEXT_PORT=$(( ((RANDOM % 49152) + 16384) ))

# File used to keep track of the list of store-gateways failed to be queried
# Reset it each time this script is called.
FAILURES_TRACKING_FILE="${SCRIPT_DIR}/${OUTPUT_DIR}/.failures"
echo -n "" > "${FAILURES_TRACKING_FILE}"

mkdir -p "$OUTPUT_DIR"

# Print a message in green.
print_success() {
  echo -e "\033[0;32m${1}\033[0m"
}

# Print a message in red.
print_failure() {
  echo -e "\033[0;31m${1}\033[0m"
}

# Utility function to query a single store-gateway
#
# Parameters:
# - $1: The pod ID
# - $2: The local port to use
query_store_gateway() {
  POD=$1
  LOCAL_PORT=$2

  echo "Querying $POD"

  # Open port-forward
  kubectl port-forward --context "$K8S_CONTEXT" -n "$K8S_NAMESPACE" "$POD" ${LOCAL_PORT}:9095 > /dev/null &
  KUBECTL_PID=$!

  # Wait some time
  sleep 5

  # HACK
  # If you get an error resolving the reference to "github.com/grafana/mimir/pkg/mimirpb/mimir.proto" in
  # pkg/ingester/client/ingester.proto, you need to manually modify the import statement to be just
  # "pkg/mimirpb/mimir.proto".
  cat "$SCRIPT_DIR/download-chunks-from-store-gateways-query.json" | grpcurl \
    -d @ \
    -rpc-header "x-scope-orgid: $MIMIR_TENANT_ID" \
    -rpc-header "__org_id__: $MIMIR_TENANT_ID" \
    -proto pkg/storegateway/storegatewaypb/gateway.proto \
    -import-path "$SCRIPT_DIR/../.." \
    -import-path "$SCRIPT_DIR/../../pkg/storegateway/storepb" \
    -import-path "$SCRIPT_DIR/../../vendor" \
    -plaintext \
    localhost:${LOCAL_PORT} "gatewaypb.StoreGateway/Series" > "$OUTPUT_DIR/$POD"
  STATUS_CODE=$?

  kill $KUBECTL_PID > /dev/null
  wait $KUBECTL_PID > /dev/null 2> /dev/null

  if [ $STATUS_CODE -eq 0 ]; then
    print_success "Successfully queried $POD"
  else
    print_failure "Failed to query $POD"

    # Keep track of the failure.
    echo "$POD" >> "${FAILURES_TRACKING_FILE}"
  fi
}

# Get list of store-gateway pods
PODS=$(kubectl --context "$K8S_CONTEXT" -n "$K8S_NAMESPACE" get pods --no-headers | grep store-gateway | awk '{print $1}')

# Concurrently query store-gateways
for POD in $PODS; do
  query_store_gateway "${POD}" "${NEXT_PORT}" &

  NEXT_PORT=$((NEXT_PORT+1))

  # Throttle to reduce the likelihood of networking issues and K8S rate limiting.
  sleep 0.25
done

# Wait for all background jobs to finish
wait

# Print final report.
echo ""
echo ""

if [ ! -s "${FAILURES_TRACKING_FILE}" ]; then
  print_success "Successfully queried all store-gateways"
  exit 0
else
  # Count the number of failed store-gateways.
  FAILURES_COUNT=$(wc -l "${FAILURES_TRACKING_FILE}" | awk '{print $1}')

  print_failure "Failed to query $FAILURES_COUNT store-gateways:"

  # Print the list of failed store-gateways.
  sort < "${FAILURES_TRACKING_FILE}" | sed 's/^/- /g'

  exit 1
fi
