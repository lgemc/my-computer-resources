# Monitor de Recursos del Servidor - GoLang 🚀

Este proyecto es una aplicación escrita en **Go** para monitorear el uso de recursos en un servidor baremetal en casa. La aplicación expone endpoints HTTP para obtener el estado actual del servidor y almacenar un historial del uso de recursos en una base de datos **PostgreSQL**.

---

## 🛠 Características

- 📊 **Monitoreo en tiempo real** del uso de CPU, RAM y disco.
- 📝 **Almacenamiento histórico** de los recursos en una base de datos PostgreSQL.
- 📂 **Exportación del historial** en formato JSON.
- 🌍 **Acceso remoto** a través de un servidor HTTP.

---

## 📡 Endpoints

### 1️⃣ Obtener el estado actual del servidor
- **Método:** `GET /resources/now`
- **Descripción:** Retorna el estado en tiempo real del servidor.

**Ejemplo de respuesta:**
```json
{
  "cpu_usage": 0.0743,
  "ram_usage": 0.6856,
  "disk_usage": 0.7435,
  "total_disk": 460.43
}
```
### 2️⃣ Obtener el historial de uso
- **Método:** `GET /resources/history`
- **Descripción:** Devuelve el historial de uso en formato JSON con timestamp.

**Ejemplo de respuesta:**
```json
[
  {
    "date": "2025-03-31 14:24:34.13440",
    "cpu_usage": 0.0743,
    "ram_usage": 0.6856,
    "disk_usage": 0.7435,
    "total_disk": 460.43
  }
]
```

## 🛠 Tecnologías utilizadas

- **Golang** → Lenguaje principal de desarrollo.
- **Echo Framework (opcional)** → Para manejar las rutas HTTP.
- **PostgreSQL** → Base de datos para almacenar el historial.
- **gopsutil** → Librería para obtener información del sistema.

---

## ⚡ Instalación y ejecución

### 1️⃣ Clonar el repositorio:

```bash
git clone https://github.com/tu-usuario/tu-repositorio.git
cd tu-repositorio
```

### 2️⃣ Instalar dependencias:

```bash
go mod tidy
```

### 3️⃣ Configurar la base de datos en PostgreSQL:

```sql
CREATE TABLE resource_usage (
    date TIMESTAMP PRIMARY KEY,
    cpu_usage FLOAT,
    ram_usage FLOAT,
    disk_usage FLOAT,
    total_disk FLOAT
);
```

### 4️⃣ Ejecutar la aplicación:

```bash
go run .
```

### 5️⃣ Acceder desde el navegador o herramientas como curl o Postman:

```bash
curl http://localhost:8080/resources/now
```

---

## 📌 Mejoras futuras

- ✅ Agregar autenticación para proteger los endpoints.
- ✅ Implementar gráficos en una UI para visualizar datos históricos.
- ✅ Optimizar almacenamiento y limpieza automática del historial.

---

Este proyecto es ideal para quienes desean aprender Go, monitoreo de servidores y bases de datos en tiempo real. 🚀 ¡Pull requests y sugerencias son bienvenidas!