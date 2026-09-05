# Fleet Service

Microservicio REST para la gestión de **maquinaria pesada y flotas**. Permite registrar maquinaria, consultar y actualizar sus datos, controlar su ciclo de estados, mantener un historial de cambios y asignar o retirar maquinaria de una flota.

El servicio está desarrollado en **Go**, expone una API HTTP con **Gin**, utiliza **PostgreSQL** mediante **GORM** y cuenta con instrumentación de observabilidad basada en **OpenTelemetry**.

---

## Contenido

- [Características](#características)
- [Tecnologías](#tecnologías)
- [Arquitectura del proyecto](#arquitectura-del-proyecto)
- [Requisitos](#requisitos)
- [Variables de entorno](#variables-de-entorno)
- [Levantar el proyecto localmente](#levantar-el-proyecto-localmente)
- [Health check](#health-check)
- [Modelo de maquinaria](#modelo-de-maquinaria)
- [Tipos de maquinaria](#tipos-de-maquinaria)
- [Estados de maquinaria](#estados-de-maquinaria)
- [Transiciones de estado](#transiciones-de-estado)
- [Endpoints](#endpoints)
- [Casos de uso](#casos-de-uso)
- [Pruebas end-to-end](#pruebas-end-to-end)
- [Manejo de errores](#manejo-de-errores)
- [Observabilidad](#observabilidad)
- [Notas conocidas](#notas-conocidas)

---

## Características

El microservicio actualmente permite:

- Crear y consultar flotas.
- Actualizar información básica de una flota.
- Registrar maquinaria pesada.
- Consultar maquinaria con paginación y filtros.
- Consultar una maquinaria por UUID.
- Actualizar datos maestros de maquinaria.
- Administrar el estado operativo de la maquinaria.
- Validar las transiciones permitidas entre estados.
- Mantener un historial de cambios de estado.
- Asignar maquinaria a una flota.
- Retirar maquinaria de una flota.
- Consultar toda la maquinaria perteneciente a una flota.
- Evitar que una maquinaria sea movida de flota cuando se encuentra en operación.
- Persistir la información en PostgreSQL.
- Ejecutar migraciones automáticamente al iniciar el servicio.
- Emitir trazas mediante OpenTelemetry a `STDOUT`, OTLP o Google Cloud.
- Generar logs locales o mediante Google Cloud Logging.

---

## Tecnologías

| Tecnología | Uso |
|---|---|
| Go 1.27 | Lenguaje principal |
| Gin 1.12 | API REST y routing HTTP |
| PostgreSQL 17 | Base de datos relacional |
| GORM 1.31 | ORM y migraciones |
| Google UUID | Identificadores UUID |
| OpenTelemetry | Trazas distribuidas |
| OTLP/gRPC | Exportación de trazas |
| Google Cloud Logging | Logging opcional en GCP |
| Google Cloud Error Reporting | Reporte de errores en GCP |
| Docker | Construcción del contenedor |
| Docker Compose | Entorno local con API + PostgreSQL |
| curl + jq | Pruebas end-to-end incluidas |

Dependencias principales definidas en `go.mod`:

```text
github.com/gin-gonic/gin
github.com/gin-contrib/cors
github.com/google/uuid
gorm.io/gorm
gorm.io/driver/postgres
go.opentelemetry.io/otel
go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
cloud.google.com/go/logging
cloud.google.com/go/errorreporting
```

---

## Arquitectura del proyecto

```text
fleet-service/
├── cmd/
│   └── services/
│       └── main.go
├── internal/
│   ├── database/
│   │   └── postgres/
│   ├── equipments/
│   │   ├── db/
│   │   ├── domain/
│   │   ├── dto/
│   │   └── handlers/
│   ├── fleets/
│   │   ├── db/
│   │   ├── domain/
│   │   ├── dto/
│   │   └── handlers/
│   ├── interfaces/
│   ├── logging/
│   ├── router/
│   ├── trace/
│   └── validations/
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── fleet-service-test.sh
├── go.mod
└── go.sum
```

### Flujo general

```text
Cliente HTTP
    │
    ▼
Gin Router
    │
    ▼
Handlers
    │
    ├── Validaciones / DTO
    │
    ▼
Interfaces de base de datos
    │
    ▼
GORM
    │
    ▼
PostgreSQL
```

Las trazas y logs se generan de forma transversal durante las operaciones HTTP y de base de datos.

---

## Requisitos

### Usando Docker Compose

Se requiere:

- Docker.
- Docker Compose v2.

Para ejecutar las pruebas incluidas también se requiere:

- `curl`.
- `jq`.

### Ejecutando Go directamente

Se requiere:

- Go 1.27 o compatible con el `go.mod` del proyecto.
- PostgreSQL accesible desde la máquina local.
- Las variables de entorno requeridas por el servicio.

---

## Variables de entorno

El repositorio incluye `.env.example`.

```env
DB_USER="mongo"
DB_PASS="1234"
DB_NAME="backend_golang_gin"
DB_STRING="host=postgres user=mongo password=1234 dbname=backend_golang_gin port=5432 sslmode=disable"
DB_CONTEXT="postgresql"
TRACE_TYPE="STDOUT"
SERVICE_NAME="fleet-service"
GCP_PROJECT_ID="gcp-project-id"
```

### Variables utilizadas por la aplicación

| Variable | Requerida | Descripción | Ejemplo |
|---|---:|---|---|
| `DB_CONTEXT` | Sí | Backend de base de datos. Actualmente solo soporta `postgresql`. | `postgresql` |
| `DB_STRING` | Sí | DSN de conexión a PostgreSQL. | `host=postgres user=mongo password=1234 dbname=backend_golang_gin port=5432 sslmode=disable` |
| `TRACE_TYPE` | Sí | Exportador de trazas. | `STDOUT` |
| `SERVICE_NAME` | Sí | Nombre utilizado por OpenTelemetry. | `fleet-service` |
| `SERVICE_VERSION` | No | Versión del servicio para telemetría. | `0.1.0` |
| `ENVIRONMENT` | No | Ambiente reportado en telemetría. | `local` |
| `LOGGING_TYPE` | No | Tipo de logging. Si no existe, usa logging local. | `GCP` |
| `GCP_PROJECT_ID` | Condicional | Requerido cuando se utiliza tracing o logging de GCP. | `my-project` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | Condicional | Endpoint OTLP personalizado. | `otel-collector:4317` |
| `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` | Condicional | Endpoint específico para trazas OTLP. | `otel-collector:4317` |

### Valores admitidos para `TRACE_TYPE`

```text
STDOUT
OTLP
GCP
NONE
DISABLED
```

Para desarrollo local se recomienda:

```env
TRACE_TYPE=STDOUT
SERVICE_NAME=fleet-service
LOGGING_TYPE=local
```

---

## Levantar el proyecto localmente

### Opción 1: Docker Compose

1. Crear el archivo `.env` a partir del ejemplo:

```bash
cp .env.example .env
```

2. Construir y levantar los servicios:

```bash
docker compose up --build
```

O en segundo plano:

```bash
docker compose up -d --build
```

3. Verificar el estado de los contenedores:

```bash
docker compose ps
```

4. Consultar el health check:

```bash
curl http://localhost:3000/health
```

Respuesta esperada:

```json
{
  "Status": "Up and Running"
}
```

5. Ver logs:

```bash
docker compose logs -f fleet-service
```

6. Detener el entorno:

```bash
docker compose down
```

Para eliminar también el volumen local de PostgreSQL:

```bash
docker compose down -v
```

> Al iniciar, el servicio ejecuta `AutoMigrate` de GORM para crear o actualizar las tablas requeridas.

### Opción 2: Go + PostgreSQL local

Si PostgreSQL se ejecuta directamente en la máquina, configurar por ejemplo:

```env
DB_CONTEXT=postgresql
DB_STRING="host=localhost user=mongo password=1234 dbname=backend_golang_gin port=5432 sslmode=disable"
TRACE_TYPE=STDOUT
SERVICE_NAME=fleet-service
```

Exportar las variables:

```bash
set -a
source .env
set +a
```

Instalar dependencias:

```bash
GODEBUG=http2client=0 go mod tidy
```

Ejecutar:

```bash
go run ./cmd/services/main.go
```

La API escucha en:

```text
http://localhost:3000
```

Base URL de la API:

```text
http://localhost:3000/api/v1
```

---

## Health check

### `GET /health`

Permite comprobar si el proceso HTTP está disponible.

```bash
curl http://localhost:3000/health
```

**200 OK**

```json
{
  "Status": "Up and Running"
}
```

---

## Modelo de maquinaria

Ejemplo de una maquinaria registrada:

```json
{
  "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
  "code": "CAT-320-001",
  "type": "EXCAVATOR",
  "brand": "CAT",
  "model": "320",
  "serialNumber": "CAT320-0001",
  "year": 2026,
  "capacityTons": 22,
  "status": "AVAILABLE",
  "location": {
    "name": "San Salvador",
    "latitude": 13.6929,
    "longitude": -89.2182
  },
  "engineHours": 120,
  "nextMaintenanceHours": 370,
  "fuelPercent": 85,
  "createdAt": "2026-09-05T20:00:00Z",
  "updatedAt": "2026-09-05T20:00:00Z"
}
```

Si la maquinaria pertenece a una flota también contiene:

```json
{
  "fleetId": "9708aeb7-ad5b-4eac-81af-e19096991e89"
}
```

`nextMaintenanceHours` se calcula al crear la maquinaria como:

```text
engineHours + maintenanceIntervalHours
```

---

## Tipos de maquinaria

Los tipos aceptados son:

| Valor | Descripción |
|---|---|
| `EXCAVATOR` | Excavadora |
| `BACKHOE` | Retroexcavadora |
| `BULLDOZER` | Bulldozer |
| `LOADER` | Cargador |
| `CRANE` | Grúa |
| `TRUCK` | Camión |

---

## Estados de maquinaria

| Estado | Significado general |
|---|---|
| `AVAILABLE` | Disponible para asignación u operación |
| `RESERVED` | Reservada para un trabajo |
| `IN_TRANSIT` | En traslado |
| `WORKING` | En operación |
| `MAINTENANCE` | En mantenimiento |
| `INACTIVE` | Inactiva |
| `RETIRED` | Retirada definitivamente |

---

## Transiciones de estado

El dominio valida las transiciones de estado. No cualquier estado puede cambiar directamente a otro.

| Estado actual | Estados siguientes permitidos |
|---|---|
| `AVAILABLE` | `RESERVED`, `MAINTENANCE`, `INACTIVE` |
| `RESERVED` | `AVAILABLE`, `IN_TRANSIT`, `WORKING`, `MAINTENANCE` |
| `IN_TRANSIT` | `AVAILABLE`, `WORKING`, `MAINTENANCE` |
| `WORKING` | `AVAILABLE`, `MAINTENANCE` |
| `MAINTENANCE` | `AVAILABLE`, `INACTIVE` |
| `INACTIVE` | `AVAILABLE`, `RETIRED` |
| `RETIRED` | Ninguno |

Ejemplo de transición válida:

```text
AVAILABLE → RESERVED → IN_TRANSIT → WORKING → AVAILABLE
```

Ejemplo de transición inválida:

```text
AVAILABLE → WORKING
```

La API responde `422 Unprocessable Entity` cuando la transición no es válida.

---

# Endpoints

## Resumen

### Maquinaria

| Método | Endpoint | Descripción |
|---|---|---|
| `POST` | `/api/v1/equipments` | Registrar maquinaria |
| `GET` | `/api/v1/equipments` | Listar y filtrar maquinaria |
| `GET` | `/api/v1/equipments/:id` | Obtener maquinaria por UUID |
| `PATCH` | `/api/v1/equipments/:id` | Actualizar datos de maquinaria |
| `PATCH` | `/api/v1/equipments/:id/status` | Cambiar estado operativo |
| `GET` | `/api/v1/equipments/:id/status-history` | Consultar historial de estados |

### Flotas

| Método | Endpoint | Descripción |
|---|---|---|
| `POST` | `/api/v1/fleets` | Crear flota |
| `GET` | `/api/v1/fleets` | Listar flotas |
| `GET` | `/api/v1/fleets/:fleetID` | Obtener flota por UUID |
| `PATCH` | `/api/v1/fleets/:fleetID` | Actualizar flota |
| `PUT` | `/api/v1/fleets/:fleetID/equipments/:equipmentID` | Asignar maquinaria |
| `DELETE` | `/api/v1/fleets/:fleetID/equipments/:equipmentID` | Retirar maquinaria |
| `GET` | `/api/v1/fleets/:fleetID/equipments` | Listar maquinaria de una flota |

---

# API de maquinaria

## Crear maquinaria

### `POST /api/v1/equipments`

Registra una nueva maquinaria. El estado inicial es automáticamente `AVAILABLE`.

### Request

```json
{
  "code": "CAT-320-001",
  "type": "EXCAVATOR",
  "brand": "CAT",
  "model": "320",
  "serialNumber": "CAT320-0001",
  "year": 2026,
  "capacityTons": 22,
  "location": {
    "name": "San Salvador",
    "latitude": 13.6929,
    "longitude": -89.2182
  },
  "engineHours": 120,
  "maintenanceIntervalHours": 250,
  "fuelPercent": 85
}
```

También se puede crear directamente dentro de una flota utilizando:

```json
{
  "fleetId": "9708aeb7-ad5b-4eac-81af-e19096991e89"
}
```

### Validaciones relevantes

- `code`: 3 a 50 caracteres.
- `type`: uno de los tipos soportados.
- `brand`: 2 a 100 caracteres.
- `model`: 1 a 100 caracteres.
- `serialNumber`: 3 a 100 caracteres.
- `year`: entre 1950 y 2100.
- `capacityTons`: mayor que 0.
- `latitude`: entre -90 y 90.
- `longitude`: entre -180 y 180.
- `engineHours`: mayor o igual que 0.
- `maintenanceIntervalHours`: mayor que 0.
- `fuelPercent`: entre 0 y 100.

### Respuesta esperada

**201 Created**

```json
{
  "message": "Maquinaria registrada correctamente",
  "data": {
    "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
    "code": "CAT-320-001",
    "type": "EXCAVATOR",
    "brand": "CAT",
    "model": "320",
    "serialNumber": "CAT320-0001",
    "year": 2026,
    "capacityTons": 22,
    "status": "AVAILABLE",
    "location": {
      "name": "San Salvador",
      "latitude": 13.6929,
      "longitude": -89.2182
    },
    "engineHours": 120,
    "nextMaintenanceHours": 370,
    "fuelPercent": 85
  }
}
```

### Errores frecuentes

**400 Bad Request**

```json
{
  "error": "invalid_request",
  "message": "Los datos enviados no son válidos"
}
```

**409 Conflict**

```json
{
  "error": "equipment_already_exists",
  "message": "Ya existe una maquinaria con ese código o número de serie"
}
```

---

## Listar maquinaria

### `GET /api/v1/equipments`

Soporta filtros y paginación.

### Query parameters

| Parámetro | Obligatorio | Descripción |
|---|---:|---|
| `page` | No | Página. Default: `1` |
| `pageSize` | No | Registros por página. Default: `20`; máximo efectivo: `100` |
| `type` | No | Filtrar por tipo |
| `status` | No | Filtrar por estado |
| `brand` | No | Filtrar por marca |
| `search` | No | Busca por código, marca, modelo o número de serie |

Ejemplo:

```bash
curl "http://localhost:3000/api/v1/equipments?type=EXCAVATOR&status=AVAILABLE&search=CAT&page=1&pageSize=20"
```

### Respuesta esperada

**200 OK**

```json
{
  "data": [
    {
      "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
      "code": "CAT-320-001",
      "type": "EXCAVATOR",
      "brand": "CAT",
      "model": "320",
      "status": "AVAILABLE"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "total": 1,
    "totalPages": 1
  }
}
```

---

## Consultar maquinaria por ID

### `GET /api/v1/equipments/:id`

```bash
curl http://localhost:3000/api/v1/equipments/68dd1b50-d732-4fe3-88e9-fb980e61bdd5
```

**200 OK**

```json
{
  "data": {
    "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
    "code": "CAT-320-001",
    "type": "EXCAVATOR",
    "status": "AVAILABLE"
  }
}
```

**400 Bad Request** si el ID no es UUID:

```json
{
  "error": "invalid_id",
  "message": "El identificador no es un UUID válido"
}
```

**404 Not Found**:

```json
{
  "error": "equipment_not_found",
  "message": "La maquinaria no existe"
}
```

---

## Actualizar maquinaria

### `PATCH /api/v1/equipments/:id`

Permite modificar parcialmente los datos maestros.

### Request

```json
{
  "brand": "Caterpillar",
  "capacityTons": 23,
  "location": {
    "name": "San Miguel",
    "latitude": 13.4833,
    "longitude": -88.1833
  }
}
```

Campos soportados:

```text
fleetId
type
brand
model
serialNumber
year
capacityTons
location.name
location.latitude
location.longitude
```

### Respuesta esperada

**200 OK**

```json
{
  "message": "Maquinaria actualizada correctamente",
  "data": {
    "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
    "brand": "Caterpillar",
    "capacityTons": 23
  }
}
```

Si no se envía ningún campo:

**400 Bad Request**

```json
{
  "error": "empty_update",
  "message": "Debe enviar al menos un campo para actualizar"
}
```

---

## Cambiar estado de maquinaria

### `PATCH /api/v1/equipments/:id/status`

### Request

```json
{
  "status": "RESERVED",
  "reason": "Reservada para una solicitud logística"
}
```

### Respuesta funcional esperada

**200 OK**

```json
{
  "message": "Estado actualizado correctamente",
  "data": {
    "equipment": {
      "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
      "status": "RESERVED"
    },
    "transition": {
      "equipmentId": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
      "fromStatus": "AVAILABLE",
      "toStatus": "RESERVED",
      "reason": "Reservada para una solicitud logística"
    }
  }
}
```

Una transición no permitida responde:

**422 Unprocessable Entity**

```json
{
  "error": "invalid_status_transition",
  "message": "invalid equipment status transition: AVAILABLE -> WORKING"
}
```

> Revisar la sección [Notas conocidas](#notas-conocidas): la implementación actual del handler/DB presenta una inconsistencia en el mapeo de `equipment`, `transition` y `reason`.

---

## Consultar historial de estados

### `GET /api/v1/equipments/:id/status-history`

```bash
curl http://localhost:3000/api/v1/equipments/68dd1b50-d732-4fe3-88e9-fb980e61bdd5/status-history
```

**200 OK**

```json
{
  "data": [
    {
      "id": "bb91f38f-75a4-4b42-88ce-5d097922a1d3",
      "equipmentId": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
      "fromStatus": "WORKING",
      "toStatus": "AVAILABLE",
      "reason": "Trabajo finalizado",
      "changedAt": "2026-09-05T21:00:00Z"
    }
  ]
}
```

El historial se devuelve ordenado por `changedAt` descendente.

---

# API de flotas

## Crear flota

### `POST /api/v1/fleets`

### Request

```json
{
  "code": "FLEET-001",
  "name": "Flota Central"
}
```

Validaciones:

- `code`: 3 a 50 caracteres.
- `name`: 3 a 150 caracteres.
- El código debe ser único.

### Respuesta esperada

**201 Created**

```json
{
  "message": "Flota creada correctamente",
  "data": {
    "id": "9708aeb7-ad5b-4eac-81af-e19096991e89",
    "code": "FLEET-001",
    "name": "Flota Central"
  }
}
```

Código duplicado:

**409 Conflict**

```json
{
  "error": "fleet_already_exists",
  "message": "Ya existe una flota con ese código"
}
```

---

## Listar flotas

### `GET /api/v1/fleets`

Query parameters:

| Parámetro | Obligatorio | Descripción |
|---|---:|---|
| `page` | No | Página. Default: `1` |
| `pageSize` | No | Tamaño de página. Default: `20`; máximo efectivo: `100` |
| `search` | No | Busca en código o nombre |

Ejemplo:

```bash
curl "http://localhost:3000/api/v1/fleets?search=central&page=1&pageSize=20"
```

**200 OK**

```json
{
  "data": [
    {
      "id": "9708aeb7-ad5b-4eac-81af-e19096991e89",
      "code": "FLEET-001",
      "name": "Flota Central"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "total": 1,
    "totalPages": 1
  }
}
```

---

## Consultar flota por ID

### `GET /api/v1/fleets/:fleetID`

```bash
curl http://localhost:3000/api/v1/fleets/9708aeb7-ad5b-4eac-81af-e19096991e89
```

**200 OK**

```json
{
  "data": {
    "id": "9708aeb7-ad5b-4eac-81af-e19096991e89",
    "code": "FLEET-001",
    "name": "Flota Central"
  }
}
```

---

## Actualizar flota

### `PATCH /api/v1/fleets/:fleetID`

### Request

```json
{
  "name": "Flota Región Central"
}
```

También se puede actualizar `code`.

### Respuesta esperada

**200 OK**

```json
{
  "message": "Flota actualizada correctamente",
  "data": {
    "id": "9708aeb7-ad5b-4eac-81af-e19096991e89",
    "code": "FLEET-001",
    "name": "Flota Región Central"
  }
}
```

---

## Asignar maquinaria a una flota

### `PUT /api/v1/fleets/:fleetID/equipments/:equipmentID`

No requiere body.

```bash
curl -X PUT \
  http://localhost:3000/api/v1/fleets/9708aeb7-ad5b-4eac-81af-e19096991e89/equipments/68dd1b50-d732-4fe3-88e9-fb980e61bdd5
```

### Respuesta esperada

**200 OK**

```json
{
  "message": "Maquinaria asignada correctamente",
  "data": {
    "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
    "fleetId": "9708aeb7-ad5b-4eac-81af-e19096991e89",
    "status": "AVAILABLE"
  }
}
```

La operación es idempotente si la maquinaria ya pertenece a la misma flota.

No puede asignarse si ya pertenece a otra flota:

**409 Conflict**

```json
{
  "error": "equipment_already_assigned",
  "message": "La maquinaria ya pertenece a otra flota"
}
```

Tampoco puede cambiar de flota cuando se encuentra en un estado operativo no permitido.

---

## Retirar maquinaria de una flota

### `DELETE /api/v1/fleets/:fleetID/equipments/:equipmentID`

```bash
curl -i -X DELETE \
  http://localhost:3000/api/v1/fleets/9708aeb7-ad5b-4eac-81af-e19096991e89/equipments/68dd1b50-d732-4fe3-88e9-fb980e61bdd5
```

### Respuesta esperada

**204 No Content**

No devuelve body.

La maquinaria solamente puede cambiar de flota cuando está en uno de estos estados:

```text
AVAILABLE
MAINTENANCE
INACTIVE
```

Ejemplo de rechazo:

**409 Conflict**

```json
{
  "error": "equipment_in_operation",
  "message": "No se puede mover la maquinaria mientras está en operación"
}
```

---

## Listar maquinaria de una flota

### `GET /api/v1/fleets/:fleetID/equipments`

```bash
curl http://localhost:3000/api/v1/fleets/9708aeb7-ad5b-4eac-81af-e19096991e89/equipments
```

### Respuesta esperada

**200 OK**

```json
{
  "fleet": {
    "id": "9708aeb7-ad5b-4eac-81af-e19096991e89",
    "code": "FLEET-001",
    "name": "Flota Central"
  },
  "data": [
    {
      "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
      "code": "CAT-320-001",
      "fleetId": "9708aeb7-ad5b-4eac-81af-e19096991e89",
      "status": "AVAILABLE"
    }
  ]
}
```

Si la flota no existe:

**404 Not Found**

```json
{
  "error": "fleet_not_found",
  "message": "La flota no existe"
}
```

---

# Casos de uso

## Caso 1: Registrar una nueva flota

Un administrador necesita crear una agrupación lógica de maquinaria.

```text
POST /api/v1/fleets
        │
        ▼
Se valida code + name
        │
        ▼
Se registra la flota
        │
        ▼
201 Created
```

Ejemplo:

```bash
curl -X POST http://localhost:3000/api/v1/fleets \
  -H 'Content-Type: application/json' \
  -d '{
    "code": "FLEET-001",
    "name": "Flota Central"
  }'
```

---

## Caso 2: Registrar maquinaria disponible

Al ingresar una nueva máquina al inventario:

```text
POST /api/v1/equipments
        │
        ▼
Estado inicial AVAILABLE
        │
        ▼
Se calcula nextMaintenanceHours
        │
        ▼
201 Created
```

---

## Caso 3: Asignar maquinaria a una flota

Flujo recomendado:

```text
Crear/consultar flota
        │
        ▼
Crear/consultar maquinaria
        │
        ▼
PUT /fleets/{fleetID}/equipments/{equipmentID}
        │
        ▼
La maquinaria queda asociada a la flota
```

Antes de asignarla, el servicio verifica:

- que la flota exista;
- que la maquinaria exista;
- que no pertenezca a otra flota;
- que su estado permita modificar la pertenencia a una flota.

---

## Caso 4: Ejecutar el ciclo de una operación logística

Un flujo de trabajo típico puede ser:

```text
AVAILABLE
   │
   ▼
RESERVED
   │
   ▼
IN_TRANSIT
   │
   ▼
WORKING
   │
   ▼
AVAILABLE
```

Cada cambio se realiza mediante:

```text
PATCH /api/v1/equipments/:id/status
```

Después puede consultarse la trazabilidad completa con:

```text
GET /api/v1/equipments/:id/status-history
```

---

## Caso 5: Enviar maquinaria a mantenimiento

Desde varios estados operativos se permite pasar a `MAINTENANCE`.

Ejemplo:

```json
{
  "status": "MAINTENANCE",
  "reason": "Mantenimiento preventivo de 500 horas"
}
```

Cuando la maquinaria está en mantenimiento también puede ser movida entre flotas según las reglas actuales del dominio.

---

## Caso 6: Buscar maquinaria disponible para otro servicio

El endpoint de listado permite filtrar el inventario antes de asignarlo:

```bash
curl "http://localhost:3000/api/v1/equipments?type=EXCAVATOR&status=AVAILABLE&brand=CAT"
```

Esto permite que otros microservicios, por ejemplo un servicio logístico o de recomendaciones, consulten maquinaria candidata sin conocer directamente la base de datos del `fleet-service`.

---

# Pruebas end-to-end

El repositorio incluye:

```text
fleet-service-test.sh
```

El script prueba, entre otras cosas:

- creación de flota;
- listado y búsqueda;
- actualización de flota;
- validación de código duplicado;
- creación de maquinaria;
- filtros de maquinaria;
- actualización de maquinaria;
- transición de estado inválida;
- asignación a flota;
- flujo `AVAILABLE → RESERVED → IN_TRANSIT → WORKING → AVAILABLE`;
- bloqueo de retiro mientras la maquinaria está en operación;
- historial de estados;
- retiro de maquinaria de una flota.

Requisitos:

```bash
curl --version
jq --version
```

Dar permisos:

```bash
chmod +x fleet-service-test.sh
```

Ejecutar con la URL por defecto:

```bash
./fleet-service-test.sh
```

Por defecto utiliza:

```text
http://localhost:3000/api/v1
```

También puede indicarse otra URL:

```bash
BASE_URL=http://localhost:3001/api/v1 ./fleet-service-test.sh
```

Al finalizar correctamente muestra:

```text
Todas las pruebas finalizaron correctamente.
```

---

# Manejo de errores

Formato general:

```json
{
  "error": "error_code",
  "message": "Descripción legible del error"
}
```

En algunas validaciones también puede incluirse:

```json
{
  "detail": "detalle técnico de validación"
}
```

Errores relevantes:

| HTTP | Código | Descripción |
|---:|---|---|
| `400` | `invalid_request` | Payload inválido |
| `400` | `invalid_id` | UUID inválido |
| `400` | `invalid_fleet_id` | `fleetId` inválido |
| `400` | `empty_update` | PATCH sin campos |
| `404` | `equipment_not_found` | Maquinaria inexistente |
| `404` | `fleet_not_found` | Flota inexistente |
| `404` | `equipment_not_in_fleet` | La maquinaria no pertenece a esa flota |
| `409` | `equipment_already_exists` | Código o serie duplicados |
| `409` | `fleet_already_exists` | Código de flota duplicado |
| `409` | `equipment_already_assigned` | La maquinaria ya pertenece a otra flota |
| `409` | `equipment_in_operation` | No puede cambiarse de flota en el estado actual |
| `409` | `concurrent_status_change` | Otro proceso modificó el estado durante la operación |
| `422` | `invalid_status_transition` | Transición de estado no permitida |
| `500` | `database_error` | Error interno de persistencia |

---

# Observabilidad

## Tracing

El servicio utiliza OpenTelemetry.

### STDOUT

Ideal para desarrollo:

```env
TRACE_TYPE=STDOUT
```

Las trazas se imprimen en consola.

### Deshabilitado

```env
TRACE_TYPE=NONE
```

O:

```env
TRACE_TYPE=DISABLED
```

### OTLP

```env
TRACE_TYPE=OTLP
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
SERVICE_NAME=fleet-service
```

### Google Cloud

```env
TRACE_TYPE=GCP
GCP_PROJECT_ID=my-gcp-project
SERVICE_NAME=fleet-service
```

El exportador de GCP utiliza Application Default Credentials.

En un entorno local autenticado con Google Cloud puede utilizarse, por ejemplo:

```bash
gcloud auth application-default login
```

## Logging

Por defecto se utiliza logging local.

Para Google Cloud:

```env
LOGGING_TYPE=GCP
GCP_PROJECT_ID=my-gcp-project
```

---

# Notas conocidas

La documentación anterior describe la intención funcional del servicio y sus rutas actuales. Durante la revisión del código se detectaron algunos puntos que conviene corregir.

## Ejemplo rápido de flujo completo

```bash
# 1. Crear una flota
curl -X POST http://localhost:3000/api/v1/fleets \
  -H 'Content-Type: application/json' \
  -d '{"code":"FLEET-001","name":"Flota Central"}'

# 2. Crear maquinaria
curl -X POST http://localhost:3000/api/v1/equipments \
  -H 'Content-Type: application/json' \
  -d '{
    "code":"CAT-320-001",
    "type":"EXCAVATOR",
    "brand":"CAT",
    "model":"320",
    "serialNumber":"CAT320-0001",
    "year":2026,
    "capacityTons":22,
    "location":{
      "name":"San Salvador",
      "latitude":13.6929,
      "longitude":-89.2182
    },
    "engineHours":120,
    "maintenanceIntervalHours":250,
    "fuelPercent":85
  }'

# 3. Asignar maquinaria a una flota
curl -X PUT \
  http://localhost:3000/api/v1/fleets/<FLEET_ID>/equipments/<EQUIPMENT_ID>

# 4. Reservar la maquinaria
curl -X PATCH \
  http://localhost:3000/api/v1/equipments/<EQUIPMENT_ID>/status \
  -H 'Content-Type: application/json' \
  -d '{
    "status":"RESERVED",
    "reason":"Reservada para solicitud logística"
  }'

# 5. Consultar historial
curl \
  http://localhost:3000/api/v1/equipments/<EQUIPMENT_ID>/status-history
```

---