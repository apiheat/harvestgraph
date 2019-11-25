#!/bin/bash

set -e

AKAMAI_EDGERC_SECTION="default"
AKAMAI_EDGERC_CONFIG="~/.edgerc"
DESTINATION="~/myrepo/netlist-usage/"
SOURCE_DIR="/tmp"

# Creating Security Configurations Data
echo "Getting all Security Configurations data."
akamai appsec configs --json --edgerc ${AKAMAI_EDGERC_CONFIG} --section ${AKAMAI_EDGERC_SECTION} | jq '[.configurations[] | {configId: .id, configName: .name}]' > ${SOURCE_DIR}/configurations_map.json

# Export all configuration data
echo "Exporting Security Configuration data..."
for CONFIGURATION in $(akamai appsec --edgerc ${AKAMAI_EDGERC_CONFIG} --section ${AKAMAI_EDGERC_SECTION} configs)
do
  akamai appsec --json --edgerc ${AKAMAI_EDGERC_CONFIG} --section ${AKAMAI_EDGERC_SECTION} export --config ${CONFIGURATION} | jq . > ${SOURCE_DIR}/${CONFIGURATION}.json
done

for NETLIST_ID in $(akamai netlist --config ${AKAMAI_EDGERC_CONFIG} --section ${AKAMAI_EDGERC_SECTION} get all | jq '.[].uniqueId' | tr -d '"')
do
  NETLIST_NAME=$(akamai netlist --config ${AKAMAI_EDGERC_CONFIG} --section ${AKAMAI_EDGERC_SECTION} get all | jq '.[].uniqueId' | tr -d '"')
  echo "Creating json metadata file for ${NETLIST_NAME}"
  harvestgraph --id "${NETLIST_ID}" --name "${NETLIST_NAME}" -m ${SOURCE_DIR}/configurations_map.json -s ${SOURCE_DIR}/appsecConfigs -d ${DESTINATION}
  echo "Creating dot graph file for ${NETLIST_NAME}"
  harvestgraph --id "${NETLIST_ID}" --name "${NETLIST_NAME}" -m ${SOURCE_DIR}/configurations_map.json -s ${SOURCE_DIR}/appsecConfigs -d ${DESTINATION} -o "dot"
done
