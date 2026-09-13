# Fleet Service

Microservicio REST que expone la información de **flota, tareas, mantenimientos y geocercas** del proveedor de rastreo **Startrack** dentro de la plataforma Entropy.

Permite registrar vehículos y maquinaria, consultarlos con filtros y paginación, cambiar su estado operativo dejando historial, registrar órdenes de mantenimiento, programar tareas de traslado y administrar geocercas.

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
- [Modelos de datos](#modelos-de-datos)
- [Endpoints](#endpoints)
- [API de vehículos](#api-de-vehículos)
- [API de tareas](#api-de-tareas)
- [API de mantenimientos](#api-de-mantenimientos)
- [API de geocercas](#api-de-geocercas)
- [API de flotas](#api-de-flotas)
- [Casos de uso](#casos-de-uso)
- [Carga de datos sintéticos](#carga-de-datos-sintéticos)
- [Pruebas end-to-end](#pruebas-end-to-end)
- [Manejo de errores](#manejo-de-errores)
- [Observabilidad](#observabilidad)
- [Notas conocidas](#notas-conocidas)

---

## Características

El microservicio actualmente permite:

- Registrar vehículos y maquinaria con la nomenclatura de Startrack.
- Consultar vehículos con paginación y filtros por tipo, estado, marca y búsqueda libre.
- Consultar un vehículo por UUID.
- Actualizar datos maestros de un vehículo.
- Cambiar el estado operativo de un vehículo y dejar registro en el historial.
- Consultar el historial de cambios de estado.
- Eliminar vehículos mediante borrado lógico.
- Registrar, consultar, actualizar y eliminar órdenes de mantenimiento.
- Registrar, consultar, actualizar y eliminar tareas de traslado.
- Registrar, consultar, actualizar y eliminar geocercas.
- Crear, listar, consultar y actualizar flotas como agrupación lógica.
- Persistir la información en PostgreSQL.
- Ejecutar migraciones automáticamente al iniciar el servicio.
- Emitir trazas mediante OpenTelemetry a `STDOUT`, OTLP o Google Cloud.
- Generar logs locales o mediante Google Cloud Logging.

> El agrupamiento operativo de un vehículo se almacena directamente en su campo `group`, tal como lo entrega Startrack. Las flotas se mantienen como recurso de compatibilidad y ya **no** administran la pertenencia de vehículos.

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
| Logrus | Logging local estructurado |
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
github.com/sirupsen/logrus
gorm.io/gorm
gorm.io/driver/postgres
go.opentelemetry.io/otel
go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
go.opentelemetry.io/otel/exporters/stdout/stdouttrace
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
│   ├── equipments/          # vehículos Startrack
│   │   ├── db/
│   │   ├── domain/
│   │   ├── dto/
│   │   └── handlers/
│   ├── fleets/              # flotas (compatibilidad)
│   │   ├── db/
│   │   ├── domain/
│   │   ├── dto/
│   │   └── handlers/
│   ├── geofences/
│   │   ├── db/
│   │   ├── domain/
│   │   ├── dto/
│   │   └── handlers/
│   ├── maintenance/
│   │   ├── db/
│   │   ├── domain/
│   │   ├── dto/
│   │   └── handlers/
│   ├── tasks/
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

Cada módulo de dominio sigue la misma estructura de cuatro capas:

```text
domain/    Modelo GORM y nombre de tabla
dto/       Contratos de entrada y sus validaciones
db/        Interfaz de persistencia + implementación PostgreSQL
handlers/  Adaptadores HTTP de Gin
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

### Tablas creadas

| Módulo | Tabla |
|---|---|
| Vehículos | `startrack_vehicles` |
| Historial de estados | `startrack_vehicle_status_history` |
| Tareas | `startrack_tasks` |
| Mantenimientos | `startrack_maintenance` |
| Geocercas | `startrack_geofences` |
| Flotas | `fleets` |

Todas las tablas de dominio usan borrado lógico (`deleted_at`).

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
| `PORT` | No | Puerto HTTP. Si no existe, usa `3000`. | `8080` |
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

> Al iniciar, el servicio ejecuta `AutoMigrate` de GORM para crear o actualizar las seis tablas requeridas.

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

## Modelos de datos

### Vehículo

Tabla `startrack_vehicles`. Ejemplo de un vehículo registrado:

```json
{
  "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
  "description": "SYN-DEMO-01",
  "status": "Normal",
  "type": "Excavadora",
  "year": 2024,
  "color": "Amarillo",
  "brand": "Caterpillar",
  "model": "320",
  "group": "Zona Central",
  "tags": "excavacion,movimiento-tierra",
  "driver": "Carlos Méndez",
  "remoteId": "REMOTE-DEMO-01",
  "location": "",
  "latitude": 0,
  "longitude": 0,
  "createdAt": "2026-09-13T20:00:00Z",
  "updatedAt": "2026-09-13T20:00:00Z"
}
```

| Campo | Tipo | Notas |
|---|---|---|
| `id` | UUID | Generado por el servicio |
| `description` | string | Identificador legible del activo. **Único** |
| `status` | string | Estado operativo. Texto libre, sin catálogo fijo |
| `type` | string | Tipo de maquinaria. Texto libre |
| `year` | int | Entre 1900 y 2100 |
| `color` | string | Opcional |
| `brand` | string | Opcional |
| `model` | string | Opcional |
| `group` | string | Agrupación operativa de Startrack |
| `tags` | string | Etiquetas separadas por coma |
| `driver` | string | Operador asignado |
| `remoteId` | string | Identificador en el sistema remoto. **Único** |
| `location` | string | Ubicación textual. Ver [Notas conocidas](#notas-conocidas) |
| `latitude` | float | Grados decimales. Ver [Notas conocidas](#notas-conocidas) |
| `longitude` | float | Grados decimales. Ver [Notas conocidas](#notas-conocidas) |
| `lastPositionAt` | fecha | Se omite cuando es nulo |

> `status` y `type` son cadenas libres. El servicio ya **no** valida un catálogo de tipos ni una máquina de estados: acepta los valores que envía el origen (`Normal`, `Alerta`, `Mantenimiento`, `Inactivo`, `Disponible`, etc.).

### Historial de estado

Tabla `startrack_vehicle_status_history`:

```json
{
  "id": "bb91f38f-75a4-4b42-88ce-5d097922a1d3",
  "equipmentId": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
  "fromStatus": "Normal",
  "toStatus": "Mantenimiento",
  "reason": "Ingreso a taller",
  "changedAt": "2026-09-13T21:00:00Z"
}
```

### Tarea

Tabla `startrack_tasks`:

```json
{
  "id": "3f0ad7a2-9c21-4a60-9d3f-1b2c3d4e5f60",
  "taskId": "SYN-DEMO-TASK-ACT01",
  "title": "Traslado de excavadora a San Miguel",
  "description": "Tarea de traslado entre planteles.",
  "type": "Traslado",
  "status": "Pendiente",
  "scheduledDate": "2030-07-15T08:00:00Z",
  "origin": "Plantel central ECON",
  "destination": "Parque industrial San Miguel",
  "latitude": 134833000,
  "longitude": -881833000,
  "assignee": "Adriana Steiner",
  "createdAt": "2026-09-13T20:00:00Z",
  "updatedAt": "2026-09-13T20:00:00Z"
}
```

> `latitude` y `longitude` son **enteros**, no grados decimales. La carga sintética usa grados × 10⁷ (`13.4833` → `134833000`). El servicio almacena el entero tal cual, sin convertirlo.

### Mantenimiento

Tabla `startrack_maintenance`:

```json
{
  "id": "9d1e6c34-7a42-4c8b-9f10-5b6a7c8d9e01",
  "vehicle": "SYN-DEMO-01",
  "reference": "Orden SYN-DEMO-01: Preventivo",
  "serviceDate": "2026-08-14T00:00:00Z",
  "odometer": 18200,
  "serviceTime": "08:30",
  "hourMeter": 850,
  "repairReason": "Cambio de filtros y revisión general",
  "provider": "Taller Central ECON",
  "mechanic": "Equipo de mantenimiento",
  "serviceType": "Preventivo",
  "createdAt": "2026-09-13T20:00:00Z",
  "updatedAt": "2026-09-13T20:00:00Z"
}
```

> `vehicle` es la **descripción** del vehículo en texto, no su UUID. No existe llave foránea: el mantenimiento no valida que el vehículo exista.

### Geocerca

Tabla `startrack_geofences`:

```json
{
  "id": "c4e5f6a7-b8c9-4d0e-9f1a-2b3c4d5e6f70",
  "geofenceId": "SYN-DEMO-GEO-SS",
  "name": "Proyecto Centro San Salvador",
  "group": "Proyectos centrales",
  "additionalMargin": 250,
  "latitude": 136929000,
  "longitude": -892182000,
  "createdAt": "2026-09-13T20:00:00Z",
  "updatedAt": "2026-09-13T20:00:00Z"
}
```

> Igual que en tareas, `latitude` y `longitude` son enteros escalados.

### Flota

Tabla `fleets`:

```json
{
  "id": "9708aeb7-ad5b-4eac-81af-e19096991e89",
  "code": "FLEET-001",
  "name": "Flota Central",
  "description": "",
  "branch": "",
  "active": true,
  "createdAt": "2026-09-13T20:00:00Z",
  "updatedAt": "2026-09-13T20:00:00Z"
}
```

> El modelo incluye `description`, `branch` y `active`, pero la API de creación y actualización solo recibe `code` y `name`.

---

# Endpoints

## Resumen

### Vehículos

Disponibles bajo dos prefijos equivalentes: `/api/v1/vehicles` (nomenclatura Startrack, recomendado) y `/api/v1/equipments` (compatibilidad con el flujo n8n existente). Ambos apuntan a los mismos handlers.

| Método | Endpoint | Descripción |
|---|---|---|
| `POST` | `/api/v1/vehicles` | Registrar vehículo |
| `GET` | `/api/v1/vehicles` | Listar y filtrar vehículos |
| `GET` | `/api/v1/vehicles/:id` | Obtener vehículo por UUID |
| `PATCH` | `/api/v1/vehicles/:id` | Actualizar datos del vehículo |
| `DELETE` | `/api/v1/vehicles/:id` | Eliminar vehículo (borrado lógico) |
| `PATCH` | `/api/v1/vehicles/:id/status` | Cambiar estado operativo |
| `GET` | `/api/v1/vehicles/:id/status-history` | Consultar historial de estados |

### Tareas

| Método | Endpoint | Descripción |
|---|---|---|
| `POST` | `/api/v1/tasks` | Registrar tarea |
| `GET` | `/api/v1/tasks` | Listar y filtrar tareas |
| `GET` | `/api/v1/tasks/:id` | Obtener tarea por UUID |
| `PATCH` | `/api/v1/tasks/:id` | Actualizar tarea |
| `DELETE` | `/api/v1/tasks/:id` | Eliminar tarea |

### Mantenimientos

| Método | Endpoint | Descripción |
|---|---|---|
| `POST` | `/api/v1/maintenance` | Registrar mantenimiento |
| `GET` | `/api/v1/maintenance` | Listar y filtrar mantenimientos |
| `GET` | `/api/v1/maintenance/:id` | Obtener mantenimiento por UUID |
| `PATCH` | `/api/v1/maintenance/:id` | Actualizar mantenimiento |
| `DELETE` | `/api/v1/maintenance/:id` | Eliminar mantenimiento |

### Geocercas

| Método | Endpoint | Descripción |
|---|---|---|
| `POST` | `/api/v1/geofences` | Registrar geocerca |
| `GET` | `/api/v1/geofences` | Listar y buscar geocercas |
| `GET` | `/api/v1/geofences/:id` | Obtener geocerca por UUID |
| `PATCH` | `/api/v1/geofences/:id` | Actualizar geocerca |
| `DELETE` | `/api/v1/geofences/:id` | Eliminar geocerca |

### Flotas

| Método | Endpoint | Descripción |
|---|---|---|
| `POST` | `/api/v1/fleets` | Crear flota |
| `GET` | `/api/v1/fleets` | Listar flotas |
| `GET` | `/api/v1/fleets/:fleetID` | Obtener flota por UUID |
| `PATCH` | `/api/v1/fleets/:fleetID` | Actualizar flota |

> Los endpoints de asignación y retiro de maquinaria de una flota (`PUT`/`DELETE /fleets/:fleetID/equipments/:equipmentID`) y el listado `GET /fleets/:fleetID/equipments` **ya no existen**. La agrupación se lee del campo `group` del vehículo.

---

# API de vehículos

Los ejemplos usan `/api/v1/vehicles`; `/api/v1/equipments` acepta exactamente las mismas peticiones.

## Crear vehículo

### `POST /api/v1/vehicles`

### Request

```json
{
  "description": "SYN-DEMO-01",
  "status": "Normal",
  "type": "Excavadora",
  "year": 2024,
  "color": "Amarillo",
  "brand": "Caterpillar",
  "model": "320",
  "group": "Zona Central",
  "tags": "excavacion,movimiento-tierra",
  "driver": "Carlos Méndez",
  "remoteId": "REMOTE-DEMO-01"
}
```

### Validaciones

| Campo | Obligatorio | Regla |
|---|---:|---|
| `description` | Sí | 2 a 200 caracteres. Único |
| `status` | Sí | 2 a 50 caracteres |
| `type` | Sí | 2 a 100 caracteres |
| `year` | Sí | Entre 1900 y 2100 |
| `remoteId` | Sí | Máximo 100 caracteres. Único |
| `color` | No | Máximo 80 caracteres |
| `brand` | No | Máximo 100 caracteres |
| `model` | No | Máximo 100 caracteres |
| `group` | No | Máximo 150 caracteres |
| `tags` | No | Máximo 500 caracteres |
| `driver` | No | Máximo 200 caracteres |

### Respuesta esperada

**201 Created**

```json
{
  "message": "Vehículo registrado correctamente",
  "data": {
    "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
    "description": "SYN-DEMO-01",
    "status": "Normal",
    "type": "Excavadora",
    "year": 2024,
    "color": "Amarillo",
    "brand": "Caterpillar",
    "model": "320",
    "group": "Zona Central",
    "tags": "excavacion,movimiento-tierra",
    "driver": "Carlos Méndez",
    "remoteId": "REMOTE-DEMO-01",
    "location": "",
    "latitude": 0,
    "longitude": 0,
    "createdAt": "2026-09-13T20:00:00Z",
    "updatedAt": "2026-09-13T20:00:00Z"
  }
}
```

### Errores frecuentes

**400 Bad Request**

```json
{
  "error": "invalid_request",
  "message": "Los datos enviados no son válidos",
  "detail": "Key: 'CreateEquipmentRequest.RemoteID' Error:Field validation for 'RemoteID' failed on the 'required' tag"
}
```

**409 Conflict**

```json
{
  "error": "vehicle_already_exists",
  "message": "Ya existe un vehículo con esa descripción o ID remoto"
}
```

---

## Listar vehículos

### `GET /api/v1/vehicles`

### Query parameters

| Parámetro | Obligatorio | Descripción |
|---|---:|---|
| `page` | No | Página. Default: `1` |
| `pageSize` | No | Registros por página. Default: `20`; máximo efectivo: `100` |
| `type` | No | Coincidencia exacta, sin distinguir mayúsculas |
| `status` | No | Coincidencia exacta, sin distinguir mayúsculas |
| `brand` | No | Coincidencia exacta, sin distinguir mayúsculas |
| `search` | No | Busca en `description`, `type`, `brand`, `model`, `remoteId` y `tags` |

Ejemplo:

```bash
curl "http://localhost:3000/api/v1/vehicles?type=Excavadora&status=Normal&search=SYN-DEMO&page=1&pageSize=20"
```

### Respuesta esperada

**200 OK**

```json
{
  "data": [
    {
      "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
      "description": "SYN-DEMO-01",
      "status": "Normal",
      "type": "Excavadora",
      "brand": "Caterpillar",
      "model": "320",
      "group": "Zona Central"
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

Los resultados se ordenan por `createdAt` descendente.

---

## Consultar vehículo por ID

### `GET /api/v1/vehicles/:id`

```bash
curl http://localhost:3000/api/v1/vehicles/68dd1b50-d732-4fe3-88e9-fb980e61bdd5
```

**200 OK**

```json
{
  "data": {
    "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
    "description": "SYN-DEMO-01",
    "status": "Normal",
    "type": "Excavadora"
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

> Este endpoint conserva el código heredado `equipment_not_found`. Las demás operaciones devuelven `vehicle_not_found`. Ver [Notas conocidas](#notas-conocidas).

---

## Actualizar vehículo

### `PATCH /api/v1/vehicles/:id`

Permite modificar parcialmente los datos maestros.

### Request

```json
{
  "brand": "Caterpillar",
  "driver": "Ana López",
  "group": "Zona Occidente"
}
```

Campos soportados:

```text
description
type
year
color
brand
model
group
tags
driver
remoteId
```

### Respuesta esperada

**200 OK**

```json
{
  "message": "Vehículo actualizado correctamente",
  "data": {
    "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
    "brand": "Caterpillar",
    "driver": "Ana López",
    "group": "Zona Occidente"
  }
}
```

Si no se envía ningún campo:

**400 Bad Request**

```json
{
  "error": "empty_update",
  "message": "Debe enviar al menos un campo"
}
```

---

## Eliminar vehículo

### `DELETE /api/v1/vehicles/:id`

```bash
curl -i -X DELETE \
  http://localhost:3000/api/v1/vehicles/68dd1b50-d732-4fe3-88e9-fb980e61bdd5
```

**204 No Content**, sin body. El borrado es lógico: el registro conserva su fila con `deleted_at` y deja de aparecer en las consultas.

**404 Not Found**

```json
{
  "error": "vehicle_not_found",
  "message": "El vehículo no existe"
}
```

---

## Cambiar estado de un vehículo

### `PATCH /api/v1/vehicles/:id/status`

El cambio se ejecuta dentro de una transacción con bloqueo `FOR UPDATE` sobre el vehículo y registra una fila en el historial.

### Request

```json
{
  "status": "Mantenimiento",
  "reason": "Ingreso a taller por falla hidráulica"
}
```

Ambos campos son obligatorios: `status` de 2 a 50 caracteres y `reason` de 2 a 250.

### Respuesta esperada

**200 OK**

```json
{
  "message": "Estado actualizado correctamente",
  "data": {
    "vehicle": {
      "id": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
      "description": "SYN-DEMO-01",
      "status": "Mantenimiento"
    },
    "transition": {
      "id": "bb91f38f-75a4-4b42-88ce-5d097922a1d3",
      "equipmentId": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
      "fromStatus": "Normal",
      "toStatus": "Mantenimiento",
      "reason": "Ingreso a taller por falla hidráulica",
      "changedAt": "2026-09-13T21:00:00Z"
    }
  }
}
```

No existe una máquina de estados: **cualquier valor de `status` es aceptado**. La única restricción es que el estado nuevo sea distinto del actual.

**409 Conflict** al enviar el mismo estado que ya tiene el vehículo:

```json
{
  "error": "invalid_status_change",
  "message": "vehicle already has the requested status"
}
```

**404 Not Found**

```json
{
  "error": "vehicle_not_found",
  "message": "El vehículo no existe"
}
```

---

## Consultar historial de estados

### `GET /api/v1/vehicles/:id/status-history`

```bash
curl http://localhost:3000/api/v1/vehicles/68dd1b50-d732-4fe3-88e9-fb980e61bdd5/status-history
```

**200 OK**

```json
{
  "data": [
    {
      "id": "bb91f38f-75a4-4b42-88ce-5d097922a1d3",
      "equipmentId": "68dd1b50-d732-4fe3-88e9-fb980e61bdd5",
      "fromStatus": "Normal",
      "toStatus": "Mantenimiento",
      "reason": "Ingreso a taller por falla hidráulica",
      "changedAt": "2026-09-13T21:00:00Z"
    }
  ]
}
```

El historial se devuelve ordenado por `changedAt` descendente. Si el vehículo no existe responde `404`.

---

# API de tareas

## Crear tarea

### `POST /api/v1/tasks`

### Request

```json
{
  "taskId": "SYN-DEMO-TASK-ACT01",
  "title": "Traslado de excavadora a San Miguel",
  "description": "Tarea de traslado entre planteles.",
  "type": "Traslado",
  "status": "Pendiente",
  "scheduledDate": "2030-07-15T08:00:00Z",
  "origin": "Plantel central ECON",
  "destination": "Parque industrial San Miguel",
  "latitude": 134833000,
  "longitude": -881833000,
  "assignee": "Adriana Steiner"
}
```

### Validaciones

| Campo | Obligatorio | Regla |
|---|---:|---|
| `taskId` | Sí | Máximo 100 caracteres. Único |
| `title` | Sí | Máximo 200 caracteres |
| `type` | Sí | Máximo 100 caracteres |
| `status` | Sí | Máximo 100 caracteres |
| `scheduledDate` | Sí | Fecha RFC 3339. Se almacena en UTC |
| `origin` | Sí | Máximo 250 caracteres |
| `destination` | Sí | Máximo 250 caracteres |
| `latitude` | Sí | Entero. No admite `0` |
| `longitude` | Sí | Entero. No admite `0` |
| `assignee` | Sí | Máximo 200 caracteres |
| `description` | No | Máximo 2000 caracteres |

**201 Created**

```json
{
  "message": "Tarea registrada correctamente",
  "data": {
    "id": "3f0ad7a2-9c21-4a60-9d3f-1b2c3d4e5f60",
    "taskId": "SYN-DEMO-TASK-ACT01",
    "status": "Pendiente"
  }
}
```

**409 Conflict** con `taskId` duplicado:

```json
{
  "error": "task_already_exists",
  "message": "Ya existe una tarea con ese identificador"
}
```

---

## Listar tareas

### `GET /api/v1/tasks`

| Parámetro | Obligatorio | Descripción |
|---|---:|---|
| `page` | No | Página. Default: `1` |
| `pageSize` | No | Default: `20`; máximo efectivo: `100` |
| `status` | No | Coincidencia exacta, sin distinguir mayúsculas |
| `type` | No | Coincidencia exacta, sin distinguir mayúsculas |
| `search` | No | Busca en `taskId`, `title`, `description`, `origin`, `destination` y `assignee` |

```bash
curl "http://localhost:3000/api/v1/tasks?status=Pendiente&search=SYN-DEMO"
```

Devuelve `data` + `pagination`, ordenado por `scheduledDate` descendente.

---

## Consultar, actualizar y eliminar tareas

```text
GET    /api/v1/tasks/:id
PATCH  /api/v1/tasks/:id
DELETE /api/v1/tasks/:id
```

`PATCH` acepta cualquier subconjunto de los campos de creación y responde:

```json
{
  "message": "Tarea actualizada correctamente",
  "data": { "id": "3f0ad7a2-9c21-4a60-9d3f-1b2c3d4e5f60", "status": "Completada" }
}
```

`DELETE` responde **204 No Content**.

**404 Not Found**

```json
{
  "error": "task_not_found",
  "message": "La tarea no existe"
}
```

---

# API de mantenimientos

## Crear mantenimiento

### `POST /api/v1/maintenance`

### Request

```json
{
  "vehicle": "SYN-DEMO-01",
  "reference": "Orden SYN-DEMO-01: Preventivo",
  "serviceDate": "2026-08-14T00:00:00Z",
  "odometer": 18200,
  "serviceTime": "08:30",
  "hourMeter": 850,
  "repairReason": "Cambio de filtros y revisión general",
  "provider": "Taller Central ECON",
  "mechanic": "Equipo de mantenimiento",
  "serviceType": "Preventivo"
}
```

### Validaciones

| Campo | Obligatorio | Regla |
|---|---:|---|
| `vehicle` | Sí | Máximo 200 caracteres. Es la `description` del vehículo, no su UUID |
| `reference` | Sí | Máximo 2000 caracteres |
| `serviceDate` | Sí | Fecha RFC 3339. Se almacena en UTC |
| `serviceTime` | Sí | Exactamente 5 caracteres, formato `HH:MM` |
| `repairReason` | Sí | Máximo 200 caracteres |
| `odometer` | No | Mayor o igual que 0 |
| `hourMeter` | No | Mayor o igual que 0 |
| `provider` | No | Máximo 200 caracteres |
| `mechanic` | No | Máximo 200 caracteres |
| `serviceType` | No | Máximo 200 caracteres |

**201 Created**

```json
{
  "message": "Mantenimiento registrado correctamente",
  "data": {
    "id": "9d1e6c34-7a42-4c8b-9f10-5b6a7c8d9e01",
    "vehicle": "SYN-DEMO-01",
    "serviceType": "Preventivo"
  }
}
```

No hay restricción de unicidad: pueden registrarse varias órdenes idénticas para el mismo vehículo.

---

## Listar mantenimientos

### `GET /api/v1/maintenance`

| Parámetro | Obligatorio | Descripción |
|---|---:|---|
| `page` | No | Página. Default: `1` |
| `pageSize` | No | Default: `20`; máximo efectivo: `100` |
| `vehicle` | No | Coincidencia exacta, sin distinguir mayúsculas |
| `search` | No | Busca en `vehicle`, `reference`, `repairReason`, `provider`, `mechanic` y `serviceType` |

```bash
curl "http://localhost:3000/api/v1/maintenance?vehicle=SYN-DEMO-01"
```

Devuelve `data` + `pagination`, ordenado por `serviceDate` descendente.

---

## Consultar, actualizar y eliminar mantenimientos

```text
GET    /api/v1/maintenance/:id
PATCH  /api/v1/maintenance/:id
DELETE /api/v1/maintenance/:id
```

`DELETE` responde **204 No Content**.

**404 Not Found**

```json
{
  "error": "maintenance_not_found",
  "message": "El mantenimiento no existe"
}
```

---

# API de geocercas

## Crear geocerca

### `POST /api/v1/geofences`

### Request

```json
{
  "geofenceId": "SYN-DEMO-GEO-SS",
  "name": "Proyecto Centro San Salvador",
  "group": "Proyectos centrales",
  "additionalMargin": 250,
  "latitude": 136929000,
  "longitude": -892182000
}
```

### Validaciones

| Campo | Obligatorio | Regla |
|---|---:|---|
| `geofenceId` | Sí | Máximo 100 caracteres. Único |
| `name` | Sí | Máximo 200 caracteres |
| `latitude` | Sí | Entero. No admite `0` |
| `longitude` | Sí | Entero. No admite `0` |
| `group` | No | Máximo 150 caracteres |
| `additionalMargin` | No | Mayor o igual que 0. Default `0` |

**201 Created**

```json
{
  "message": "Geocerca registrada correctamente",
  "data": {
    "id": "c4e5f6a7-b8c9-4d0e-9f1a-2b3c4d5e6f70",
    "geofenceId": "SYN-DEMO-GEO-SS",
    "name": "Proyecto Centro San Salvador"
  }
}
```

**409 Conflict**

```json
{
  "error": "geofence_already_exists",
  "message": "Ya existe una geocerca con ese identificador"
}
```

---

## Listar geocercas

### `GET /api/v1/geofences`

| Parámetro | Obligatorio | Descripción |
|---|---:|---|
| `page` | No | Página. Default: `1` |
| `pageSize` | No | Default: `20`; máximo efectivo: `100` |
| `search` | No | Busca en `geofenceId`, `name` y `group` |

Devuelve `data` + `pagination`, ordenado por `createdAt` descendente.

---

## Consultar, actualizar y eliminar geocercas

```text
GET    /api/v1/geofences/:id
PATCH  /api/v1/geofences/:id
DELETE /api/v1/geofences/:id
```

`DELETE` responde **204 No Content**.

**404 Not Found**

```json
{
  "error": "geofence_not_found",
  "message": "La geocerca no existe"
}
```

---

# API de flotas

Las flotas se conservan como agrupación lógica y como endpoint de compatibilidad. No administran la pertenencia de vehículos.

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

- `code`: 3 a 50 caracteres. Se normaliza a mayúsculas. Único.
- `name`: 3 a 150 caracteres.

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

| Parámetro | Obligatorio | Descripción |
|---|---:|---|
| `page` | No | Página. Default: `1` |
| `pageSize` | No | Default: `20`; máximo efectivo: `100` |
| `search` | No | Busca en `code` o `name` |

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

## Consultar y actualizar flotas

```text
GET   /api/v1/fleets/:fleetID
PATCH /api/v1/fleets/:fleetID
```

`PATCH` acepta `code` y/o `name`:

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

# Casos de uso

## Caso 1: Ingresar un vehículo del rastreo satelital

```text
POST /api/v1/vehicles
        │
        ▼
Se validan description, status, type, year y remoteId
        │
        ▼
Se registra el vehículo con su group y tags
        │
        ▼
201 Created
```

```bash
curl -X POST http://localhost:3000/api/v1/vehicles \
  -H 'Content-Type: application/json' \
  -d '{
    "description":"SYN-DEMO-01",
    "status":"Normal",
    "type":"Excavadora",
    "year":2024,
    "brand":"Caterpillar",
    "model":"320",
    "group":"Zona Central",
    "tags":"excavacion,movimiento-tierra",
    "driver":"Carlos Méndez",
    "remoteId":"REMOTE-DEMO-01"
  }'
```

---

## Caso 2: Enviar un vehículo a mantenimiento y registrar la orden

```text
PATCH /api/v1/vehicles/{id}/status   → status "Mantenimiento"
        │
        ▼
POST /api/v1/maintenance             → vehicle = description del vehículo
        │
        ▼
GET /api/v1/vehicles/{id}/status-history
```

El historial queda disponible para auditar cuándo entró y salió de taller.

---

## Caso 3: Programar un traslado

```text
GET /api/v1/vehicles?status=Normal&type=Excavadora
        │
        ▼
POST /api/v1/tasks   → origen, destino, fecha y responsable
        │
        ▼
PATCH /api/v1/tasks/{id}  → status "Completada" o "Cancelada"
```

---

## Caso 4: Delimitar zonas de proyecto

```text
POST /api/v1/geofences
        │
        ▼
Se registra la zona con su margen adicional
        │
        ▼
GET /api/v1/geofences?search=San Salvador
```

---

## Caso 5: Buscar vehículos disponibles para otro servicio

```bash
curl "http://localhost:3000/api/v1/vehicles?type=Excavadora&status=Normal&brand=Caterpillar"
```

Esto permite que otros microservicios de la plataforma, por ejemplo el servicio logístico o el MCP, consulten el inventario sin conocer directamente la base de datos del `fleet-service`.

---

# Carga de datos sintéticos

El script `seed-entropy-synthetic-data.sh` (mantenido fuera de este repositorio, junto a los demás servicios de Entropy) carga un conjunto correlacionado de datos de prueba usando únicamente HTTP.

Contra este servicio ejecuta:

| Recurso | Cantidad aproximada |
|---|---|
| `POST /api/v1/vehicles` | 25 vehículos, incluido uno sin pareja en Logistic |
| `POST /api/v1/maintenance` | 12 órdenes preventivas y correctivas |
| `POST /api/v1/geofences` | 5 geocercas regionales |
| `POST /api/v1/tasks` | 4 tareas de traslado en distintos estados |

Además llama a `GET /health` antes de empezar, y coordina cada vehículo con su equivalente en `logistic-service` y con la proyección unificada del `mcp-server`.

Configuración por variables de entorno, sin valores dentro del repositorio:

```bash
FLEET_URL=http://localhost:3000 \
LOGISTIC_URL=http://localhost:3001 \
MCP_URL=http://localhost:3002 \
SEED_RUN=SYNTH-DEMO \
./seed-entropy-synthetic-data.sh
```

`SEED_RUN` identifica la corrida y se incrusta en `description`, `remoteId`, `taskId` y `geofenceId`, de modo que todo el conjunto puede aislarse después con el parámetro `search`.

> El script envía una cabecera `Authorization: Bearer`. Este servicio **no** valida credenciales: la cabecera se ignora. No incrustes tokens reales en scripts versionados.

---

# Pruebas end-to-end

El repositorio incluye:

```text
fleet-service-test.sh
```

> **Este script está desactualizado.** Fue escrito para el modelo anterior de maquinaria (`code`, `serialNumber`, `capacityTons`, `engineHours`, estados `AVAILABLE`/`RESERVED`/`WORKING`) y para los endpoints de pertenencia a flota que ya se eliminaron. Falla en el paso 5, al crear maquinaria con el payload antiguo. Ver [Notas conocidas](#notas-conocidas).

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

Mientras tanto, el flujo mínimo puede verificarse a mano:

```bash
# 1. Crear un vehículo
curl -X POST http://localhost:3000/api/v1/vehicles \
  -H 'Content-Type: application/json' \
  -d '{"description":"SMOKE-01","status":"Normal","type":"Excavadora","year":2024,"remoteId":"REMOTE-SMOKE-01"}'

# 2. Cambiar su estado
curl -X PATCH http://localhost:3000/api/v1/vehicles/<VEHICLE_ID>/status \
  -H 'Content-Type: application/json' \
  -d '{"status":"Mantenimiento","reason":"Prueba de humo"}'

# 3. Consultar el historial
curl http://localhost:3000/api/v1/vehicles/<VEHICLE_ID>/status-history

# 4. Registrar un mantenimiento
curl -X POST http://localhost:3000/api/v1/maintenance \
  -H 'Content-Type: application/json' \
  -d '{"vehicle":"SMOKE-01","reference":"Orden de prueba","serviceDate":"2026-09-13T00:00:00Z","serviceTime":"08:30","repairReason":"Prueba de humo","odometer":100,"hourMeter":10}'

# 5. Eliminar el vehículo
curl -i -X DELETE http://localhost:3000/api/v1/vehicles/<VEHICLE_ID>
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

En los errores de binding también se incluye el detalle del validador:

```json
{
  "detail": "detalle técnico de validación"
}
```

Errores relevantes:

| HTTP | Código | Descripción |
|---:|---|---|
| `400` | `invalid_request` | Payload inválido |
| `400` | `invalid_id` | UUID inválido en la ruta |
| `400` | `empty_update` | PATCH sin campos |
| `404` | `vehicle_not_found` | Vehículo inexistente |
| `404` | `equipment_not_found` | Vehículo inexistente, devuelto por `GET /vehicles/:id` |
| `404` | `task_not_found` | Tarea inexistente |
| `404` | `maintenance_not_found` | Mantenimiento inexistente |
| `404` | `geofence_not_found` | Geocerca inexistente |
| `409` | `vehicle_already_exists` | `description` o `remoteId` duplicados |
| `409` | `task_already_exists` | `taskId` duplicado |
| `409` | `geofence_already_exists` | `geofenceId` duplicado |
| `409` | `fleet_already_exists` | Código de flota duplicado |
| `409` | `invalid_status_change` | El vehículo ya tiene ese estado |
| `500` | `database_error` | Error interno de persistencia |
| `500` | `fleet_not_found` | Flota inexistente al actualizar. El código HTTP es incorrecto, ver [Notas conocidas](#notas-conocidas) |

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

> La instrumentación con spans explícitos está completa en los módulos `equipments` y `fleets`. Los módulos `tasks`, `maintenance` y `geofences` heredan las trazas automáticas de Gin, pero aún no abren spans propios.

## Logging

Por defecto se utiliza logging local.

Para Google Cloud:

```env
LOGGING_TYPE=GCP
GCP_PROJECT_ID=my-gcp-project
```

---

# Notas conocidas

Puntos detectados durante la revisión del código que conviene tener presentes o corregir.

## Coordenadas del vehículo nunca se guardan

El modelo `Equipment` tiene `location`, `latitude`, `longitude` y `lastPositionAt`, pero los DTO `CreateEquipmentRequest` y `UpdateEquipmentRequest` **no incluyen esos campos**. Si el cliente los envía, Gin los descarta en silencio y el vehículo se persiste con `location: ""`, `latitude: 0`, `longitude: 0` y `lastPositionAt` nulo. El script de datos sintéticos los envía en cada `POST /api/v1/vehicles` y ninguno se almacena.

## Códigos de error inconsistentes

- `GET /api/v1/vehicles/:id` devuelve `equipment_not_found` con el mensaje "La maquinaria no existe"; las demás operaciones usan `vehicle_not_found` y "El vehículo no existe".
- `GET /api/v1/fleets/:fleetID` devuelve el código `equipment_not_found` con el mensaje "La flota no existe".

## `PATCH /fleets/:fleetID` responde 500 cuando la flota no existe

En `internal/fleets/db/postgres/updateFleets.go`, el caso `RowsAffected == 0` construye el error `fleet_not_found` con `http.StatusInternalServerError` en lugar de `http.StatusNotFound`.

## Latitud y longitud cero son rechazadas

En tareas y geocercas, `latitude` y `longitude` son `int64` con la regla `required`. El validador de Gin trata el `0` como valor ausente, así que no es posible registrar un punto exactamente sobre el ecuador o el meridiano de Greenwich.

## Escala de coordenadas no documentada en el contrato

Tareas y geocercas guardan coordenadas como enteros, mientras que el vehículo las guarda como decimales. La carga sintética usa grados × 10⁷, pero el servicio no valida ni convierte la escala: cualquier entero es aceptado.

## Código muerto de la etapa anterior

Quedaron sin uso tras eliminar la pertenencia de maquinaria a flotas:

```text
internal/validations/fleetMembershipError.go
internal/validations/fleetId.go
internal/equipments/dto/createLocations.go
internal/equipments/dto/updateLocations.go
```

## `fleet-service-test.sh` desactualizado

El script todavía prueba el modelo anterior (`code`, `serialNumber`, `capacityTons`, transiciones `AVAILABLE → RESERVED → IN_TRANSIT → WORKING`) y los endpoints `PUT`/`DELETE /fleets/:fleetID/equipments/:equipmentID`, que ya no existen. Necesita reescribirse contra los cinco recursos actuales.

## El servicio no valida autenticación

No hay middleware de autenticación ni de API key. Las cabeceras `Authorization` que envían los clientes de integración se ignoran. El control de acceso depende por completo de la capa de red o del gateway que exponga el servicio.

## CORS restringido a un solo origen

`internal/router/cors.go` permite únicamente `http://localhost:8080`. Cualquier frontend desplegado en otro origen será bloqueado por el navegador hasta que se agregue a la lista.
