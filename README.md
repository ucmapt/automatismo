# **Procesador Topológico de la UCM-CFE** #
## Estado CARGANDO ##
- Se lee archivo de configuración
- Se prueba conectar a las fuentes de datos
- Se valida el estado del módulo, si está desactivado, se procede a cambiar a estado **FueraLinea**
- Se carga el grafo con la información geográfica de la BD (231)
  > **FALTA** - Como trabajo pendiente, está subir la información a REDIS para acelerar los refrescos de información

- Cargar estados de los elementos de conmutación involucrados (224)
  > **NOTA** - Debido a problemas de consistencia de datos, provisionalemente se está asociando la información basado em una tabla que contiene los puntos por señal de cada elemento, se requiere adecuar esto, mapeando los puntos correspondientes sin este registro.

  > **NOTA** - Raúl tiene APIs que pueden brindar alternativas a esta situación y evitar tener duplicidad de información y configuraciones.

  > **FALTA** - Completar el manejo de pseudopuntos, pero básicamente se incorpora el mismo concepto, un elemento de conmutación cuyo estado se regista en la UCM y afecta un nodo o una línea.
 
- Se actualiza el grafo, aplicando algoritmos de coloreo del grafo
- Refrescar mapa a través de SQL, actualizando circuito y estado en los elementos
> **FALTA** - Incluir el estado del módulo de Procesador Topológico en la IHM del Mapa
 
- Cualquier problema en la carga, conduce al estado **SinConfiguracion**
- Si termina exitosamente la carga, se procede a cambiar a estado **Operando**
 

## Estado OPERANDO ##

- Preparar recepción de mensajes de RABBITMQ
  - Inicio de Atencion a falla franca
  - Respuesta del Operador ante atenciones
  - Solicitud de Suspensión de procesos desde el editor topológico

- Otros eventos a monitorear
  - Cambios en la configuración del módulo
  - Cambios en los catálogos de la UCM 
  - Cambio de estado de los elementos de conmutación de la UCM
  - Cambios en registro de licencias
  - Se superó el tiempo de espera para recibir respuesta del operador

  > NOTA - Los cambios en la configuración requieren ejecutar una consulta cada 150 segundos, podrían adecuarse a un registro en RabbitMQ

  > NOTA - Los cambios en los catálogos que se han estado considerando son los relacionados con el manejo de puntos de los alimentadores y elementos de conmutación.

  > FALTA - Incorporar manejo de pseudopuntos.

  > NOTA - Las licencias se consultan cada 150 segundos a través de una API.

  > FALTA -Agregar mejor manejo de licencias. 

- Al recibir mensaje de RabbitMQ iniciando _*Atencion a Falla Franca*_ se procede a **Atención a Falla Franca**.
- Al recibir mensaje de RabbitMQ iniciando _*Respuesta del Operador*_ se procede a **Respuesta del Operador**
- Al recibir mensaje de RabbitMQ iniciando _*Solicitud de Suspensión*_ de procesos desde el editor topológico, se procede a **Suspensión de procesos**
- Al detectar cambios en la configuración del módulo, se procede a **Analizar nueva configuración**
- Al detectar cambios en los catálogos de la UCM, cambio de estado de los elementos de conmutación de la UCM o cambios en registro de licencias, se procede a **Refrescar información topológica**

  
### Atención a Falla Franca ###
La atención a falla franca requiere procesar el mensaje recibido, complementar la información del evento de falla, calcular el impacto de falla y notificar al operador. 
- Recolectar información de la falla
  1) Elemento que sintió la falla 
     - Incluido en el mensaje, completar con SQL a las tablas del elemento eléctrico en la UCM
  2) Puntos para Estado, Corrientes, Voltajes y Potencias del elemento registrados en la UCM
     - Extaer mediante API 
  3) Corriente de falla 
     - Incluido en el mensaje
  4) Momento de la falla 
     - Incluido en el mensaje
 
- Descartar fallas caducadas, es decir que sucedieron hace más de 10 minutos
 
- Completar contexto de la falla:
  1) Reconstruir grafo general
     - Cargar información topológica desde BD (Podría trabajarse desde REDIS, de momento se vuelve a cargar desde POSTGRES)
     - Aplicando la API correspondiente, completar con información de estado de los elementos de los circuitos
     - Actualizar grafo general
 
  2) Extraer subgrafo que contenga el elemento que sintió la falla
 
  3) Obtener parámetros eléctricos de la ventana de análisis, incluyendo voltajes ABC, corrientes ABC y potencias efectivas
     - Obtener parámetros eléctricos de T-2 (2 minutos previos a la falla aproximadamente)
     - Obtener parámetros eléctricos de T-0 (reciente)
 
- Descartar información nula o caduca, registrada hace más de 10 minutos
 
- Dependiendo que información se obtuvo, calcular el impacto de la falla
  - Si se tienen potencias, calcular como la resta del dato anterior y el más reciente
  - Si se tienen corrientes, calcular como la resta del dato RTS anterior y el RTS más reciente
  - Sin datos adecuados, no se calcula el impacto, pero se registra como "imposible de estimar"

- Obtener listado de elementos desconectados tras la falla

- Refrescar mapa a través de SQL, actualizando circuito y estado en los elementos
- A través de API, mostrar mensaje al operador, incluir impacto de la falla
 
- Registrar evento de falla franca
- Registrar nuevo proceso de atención a falla en curso
 
**TERMINA ATENCIÓN A FALLA FRANCA**

### Respuesta del Operador ###
 
- El operador debe responder si se procede a ejecutar el análisis de corto circuito o no
 - **En caso negativo,**
   - Registrar evento de respuesta del operador
   - Registrar proceso de atención a falla como terminado
 
 - **En caso positivo,**
   - Registrar evento de respuesta del operador
   - Extraer subgrafo del circuito dinámico
   - Extraer la forma matricial del subgrafo
   - Extraer la información eléctrica de los transformadores involucrados en el circuito dinámico
   - Registrar evento de análisis iniciado
   - Ejecutar el proceso de análisis conforme al algortimo definido por el Dr. Dionisio Antonio Suarez Cerda, tomando como base el circuito dinámico en cuestión, con la corriente de corto circuito reportada. 
   - Esperar resultado, habilitar terminación del proceso
   - Registrar los resultados del proceso de análisis
   - Registrar evento de análisis terminado
   - Refrescar mapa a través de SQL, actualizando circuito, estado en los elementos y marcas del polígono de falla
   - Registrar proceso de atención a falla como terminado
 
**TERMINA RESPUESTA DEL OPERADOR**

### Suspensión de procesos ###
- Registrar evento de solicitud de suspensión.
- Descartar que exista otra suspensión vigente. De existir, ignorar esperando que se concluya la solicitud anterior. 
- Registrar todos los procesos en curso como terminados.
- Terminar procesos de análisis en curso.
- Terminar procesos de actualización de eventos.
- Refrescar mapa a través de SQL, restableciendo circuito y estado en los elementos
- A través de API, reportar módulo suspendido
- Registrar estado del módulo como inactivo
- Terminar recepción de mensajes de RABBITMQ
  - Inicio de Atencion a falla franca
  - Respuesta del Operador ante atenciones
  - Solicitud de Suspensión de procesos desde el editor topológico
 
- Otros eventos a dejar de monitorear
  - Cambios en los catálogos de la UCM
  - Cambio de estado de los elementos de conmutación de la UCM
  - Cambios en registro de licencias
  - Registrar proceso de suspensión en curso
 
- Si todo sucede exitosamente, se procede a cambiar a estado **Bloqueado**
 
**TERMINA SUSPENSIÓN DE PROCESOS**

### Refrescar información topológica ###
- Registrar todos los procesos en curso como terminados.
- Terminar procesos de análisis en curso.
- Se carga el grafo con la información geográfica de la BD (231)
- Cargar estados de los elementos de conmutación involucrados (224)
- Se actualiza el grafo, aplicando algoritmos de coloreo del grafo
- Refrescar mapa a través de SQL, actualizando circuito y estado en los elementos
- Registrar evento de refresco de información topológica

**TERMINA REFRESCAR INFORMACIÓN TOPOLÓGICA**

### Analizar nueva configuración ###
- Registrar evento de análisis de configuración
- Si cambia el estado, de activado a desactivado, se procede a cambiar a estado **FueraLinea**
- Si cambia el estado, de desactivado a activado, se procede a cambiar a estado **Cargando**
- Si no hubo cambios de estado, actualizar parámetros.

**TERMINA ANALIZAR NUEVA CONFIGURACIÓN**

## Estado BLOQUEADO ##
- Preparar recepción de mensajes de RABBITMQ
  - Solicitud de Reactivación de procesos desde el editor topológico

- Al recibir mensaje de RabbitMQ iniciando _*Solicitud de Reactivacon*_ de procesos desde el editor topológico, se procede a **Reactivacion de procesos**

### Reactivación de procesos ###
- Registrar evento de solicitud de reactivación
- Descartar que no exista una reactivación vigente
- Registrar todos los procesos en curso como terminados
  > FALTA - Validar consistencia de datos georrefereciados y registrados en la UCM, incluyendo un mensaje en la IHM
- A través de API, reportar módulo reactivado
- Registrar estado del módulo como activo
- Registrar proceso de reactivación en curso
 
**TERMINA REACTIVACIÓN DE PROCESOS**

## Estado FUERALINEA ##
- Preparar eventos a monitorear
  - Cambios en la configuración del módulo

- Al detectar cambios en la configuración del módulo, se procede a **Analizar reactivación**

### Analizar reactivación ###
- Registrar evento de análisis de configuración
- Si cambia el estado, de desactivado a activado, se procede a cambiar a estado **Cargando**

**TERMINA ANALIZAR REACTIVACIÓN**

## Estado SINCONFIGURACIÓN ##
- Registrar evento de problema de configuración
- Es un estado terminal