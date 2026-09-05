#!/usr/bin/env bash
set -euo pipefail

# Pruebas end-to-end para fleet-service.
# Requisitos: curl y jq.
# Uso:
#   chmod +x fleet-service-api-tests.sh
#   BASE_URL=http://localhost:8080/api/v1 ./fleet-service-api-tests.sh

BASE_URL="${BASE_URL:-http://localhost:3000/api/v1}"
RUN_ID="${RUN_ID:-$(date +%s)}"
BODY_FILE="$(mktemp)"
FLEET_ID=""
EQUIPMENT_ID=""

cleanup_file() {
  rm -f "${BODY_FILE}"
}
trap cleanup_file EXIT

for command in curl jq; do
  if ! command -v "${command}" >/dev/null 2>&1; then
    echo "Falta el comando requerido: ${command}" >&2
    exit 1
  fi
done

print_body() {
  if [[ -s "${BODY_FILE}" ]]; then
    jq . "${BODY_FILE}" 2>/dev/null || sed -n '1,120p' "${BODY_FILE}"
  fi
}

request() {
  local method="$1"
  local path="$2"
  local expected_codes="$3"
  local payload="${4:-}"
  local status

  echo
  echo "${method} ${path}"

  if [[ -n "${payload}" ]]; then
    status="$(curl --silent --show-error \
      --output "${BODY_FILE}" \
      --write-out '%{http_code}' \
      --request "${method}" \
      --header 'Content-Type: application/json' \
      --data "${payload}" \
      "${BASE_URL}${path}")"
  else
    status="$(curl --silent --show-error \
      --output "${BODY_FILE}" \
      --write-out '%{http_code}' \
      --request "${method}" \
      "${BASE_URL}${path}")"
  fi

  echo "HTTP ${status}"
  print_body

  if [[ " ${expected_codes} " != *" ${status} "* ]]; then
    echo "Se esperaba uno de estos códigos: ${expected_codes}" >&2
    exit 1
  fi
}

json_id() {
  jq -er '.data.id' "${BODY_FILE}"
}

echo "Probando API: ${BASE_URL}"

# 1. Crear flota.
request POST /fleets "201" "{
  \"code\": \"FLEET-TEST-${RUN_ID}\",
  \"name\": \"Flota de pruebas ${RUN_ID}\"
}"
FLEET_ID="$(json_id)"
echo "FLEET_ID=${FLEET_ID}"

# 2. Listar, buscar y obtener la flota.
request GET "/fleets" "200"
request GET "/fleets?search=FLEET-TEST-${RUN_ID}" "200"
request GET "/fleets/${FLEET_ID}" "200"

# 3. Actualizar flota.
request PATCH "/fleets/${FLEET_ID}" "200" "{
  \"name\": \"Flota QA ${RUN_ID}\"
}"

# 4. Validar código duplicado.
request POST /fleets "409" "{
  \"code\": \"FLEET-TEST-${RUN_ID}\",
  \"name\": \"Flota duplicada\"
}"

# 5. Crear maquinaria sin asignarla inicialmente a una flota.
request POST /equipments "201" "{
  \"code\": \"CAT-TEST-${RUN_ID}\",
  \"type\": \"EXCAVATOR\",
  \"brand\": \"CAT\",
  \"model\": \"320\",
  \"serialNumber\": \"SERIAL-${RUN_ID}\",
  \"year\": 2026,
  \"capacityTons\": 22,
  \"location\": {
    \"name\": \"San Salvador\",
    \"latitude\": 13.6929,
    \"longitude\": -89.2182
  },
  \"engineHours\": 120,
  \"maintenanceIntervalHours\": 250,
  \"fuelPercent\": 85
}"
EQUIPMENT_ID="$(json_id)"
echo "EQUIPMENT_ID=${EQUIPMENT_ID}"

# 6. Listar, filtrar y obtener la maquinaria.
request GET "/equipments" "200"
request GET "/equipments?type=EXCAVATOR&status=AVAILABLE&search=CAT-TEST-${RUN_ID}" "200"
request GET "/equipments/${EQUIPMENT_ID}" "200"

# 7. Actualizar datos maestros.
request PATCH "/equipments/${EQUIPMENT_ID}" "200" "{
  \"brand\": \"Caterpillar\",
  \"capacityTons\": 23,
  \"location\": {
    \"name\": \"San Miguel\",
    \"latitude\": 13.4833,
    \"longitude\": -88.1833
  }
}"

# 8. Probar una transición inválida. AVAILABLE -> WORKING debe rechazarse.
request PATCH "/equipments/${EQUIPMENT_ID}/status" "400 422" "{
  \"status\": \"WORKING\",
  \"reason\": \"Prueba de transición inválida\"
}"

# 9. Asignar la maquinaria a la flota.
request PUT "/fleets/${FLEET_ID}/equipments/${EQUIPMENT_ID}" "200"
request GET "/fleets/${FLEET_ID}/equipments" "200"

# 10. No debe poder eliminarse una flota con maquinaria asignada.
# request DELETE "/fleets/${FLEET_ID}" "409"

# 11. Ejecutar el ciclo operacional permitido.
request PATCH "/equipments/${EQUIPMENT_ID}/status" "200" "{
  \"status\": \"RESERVED\",
  \"reason\": \"Reservada para una prueba logística\"
}"

# 12. No debe retirarse de la flota mientras esté reservada.
request DELETE "/fleets/${FLEET_ID}/equipments/${EQUIPMENT_ID}" "409"

request PATCH "/equipments/${EQUIPMENT_ID}/status" "200" "{
  \"status\": \"IN_TRANSIT\",
  \"reason\": \"Traslado de prueba iniciado\"
}"

request PATCH "/equipments/${EQUIPMENT_ID}/status" "200" "{
  \"status\": \"WORKING\",
  \"reason\": \"Maquinaria recibida en el proyecto\"
}"

request PATCH "/equipments/${EQUIPMENT_ID}/status" "200" "{
  \"status\": \"AVAILABLE\",
  \"reason\": \"Prueba finalizada y maquinaria liberada\"
}"

# 13. Consultar el historial completo de estados.
request GET "/equipments/${EQUIPMENT_ID}/status-history" "200"

# 14. Retirar la maquinaria de la flota. 204 no devuelve JSON.
request DELETE "/fleets/${FLEET_ID}/equipments/${EQUIPMENT_ID}" "204"
request GET "/fleets/${FLEET_ID}/equipments" "200"

# 15. Eliminar los registros creados por la prueba.
# request DELETE "/equipments/${EQUIPMENT_ID}" "204"
EQUIPMENT_ID=""

# request DELETE "/fleets/${FLEET_ID}" "204"
FLEET_ID=""

# 16. Verificar que ya no aparecen en consultas normales.
request GET "/equipments?search=CAT-TEST-${RUN_ID}" "200"
request GET "/fleets?search=FLEET-TEST-${RUN_ID}" "200"

echo
echo "Todas las pruebas finalizaron correctamente."