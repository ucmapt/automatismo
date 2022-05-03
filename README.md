# Procesador Topológico de la UCM-CFE

## Estado CARGANDO

- Se lee archivo de configuración
- Se prueba conectar a las fuentes de datos
- Se carga el grafo con la información geográfica de la BD (231)
> OJO - Como trabajo pendiente, está subir la información a REDIS para acelerar los refrescos de información

- Cargar estados de los elementos de conmutación involucrados
 
> OJO - Debido a problemas de consistencia de datos al usar la API, provisionalemente se está asociando la información
       basado em una tabla que contiene los puntos por señal de cada elemento, se requiere adecuar esto, mapeando
       los puntos correspondientes sin este registro.
       Falta completar el manejo de pseudopuntos, pero básicamente se incorpora el mismo concepto, un elemento de conmutación 
       cuyo estado se regista en la UCM y afecta un nodo o una línea.
 
- Se actualiza el grafo, aplicando algoritmos de coloreo del grafo
- Refrescar mapa a través de SQL, actualizando circuito y estado en los elementos
> FALTA - Incluir el estado del módulo de Procesador Topológico en la IHM del Mapa
 
- Cualquier problema en la carga, conduce al estado SinConfiguracion

- Si termina exitosamente la carga, se procede a cambiar a estado Operando
 

## Estado OPERANDO 

- Preparar recepción de mensajes de RABBITMQ
  - Inicio de Atencion a falla franca
  - Respuesta del Operador ante atenciones
  - Solicitud de Suspensión de procesos desde el editor topológico

- Otros eventos a monitorear
  - Cambios en los catálogos de la UCM
  - Cambios en la configuración del módulo
  - Cambio de estado de los elementos de conmutación de la UCM
  - Cambios en registro de licencias
  - Se superó el tiempo de espera para recibir respuesta del operador
 
 
 
  Al recibir mensaje de RabbitMQ iniciando Atencion a Falla Franca
 
  ### INICIA ATENCIÓN A FALLA FRANCA
 
  Recolectar información de la falla
  1) Elemento que sintió la falla - incluido en el mensaje, completar con SQL a las tablas del elemento eléctrico en la UCM
  2) Puntos para Estado, Corrientes, Voltajes y Potencias del elemento registrados en la UCM
  2) Corriente de falla - incluido en el mensaje
  3) Momento de la falla - incluido en el mensaje
 
  Descartar fallas caducadas, es decir que sucedieron hace más de 10 minutos
 
  Completar contexto de la falla:
 
  A) Reconstruir grafo general
  - Cargar información topológica desde BD (Podría trabajarse desde REDIS, de momento se vuelve a cargar desde POSTGRES)
  - Aplicando la API correspondiente, completar con información de estado de los elementos de los circuitos
  - Actualizar grafo general
 
  B) Extraer subgrafo que contenga el elemento que sintió la falla
 
  C) Obtener parámetros eléctricos de la ventana de análisis, incluyendo voltajes ABC, corrientes ABC y potencias efectivas
  - Obtener parámetros eléctricos de T-2 (2 minutos previos a la falla aproximadamente)
  - Obtener parámetros eléctricos de T-0 (reciente)
 
  Descartar información nula o caduca, registrada hace más de 10 minutos
 
  Dependiendo que información se obtuvo, calcular el impacto de la falla
  - Si se tienen potencias, calcular como la resta del dato anterior y el más reciente
  - Si se tienen corrientes, calcular como la resta del dato RTS anterior y el RTS más reciente
  - Sin datos adecuados, no se calcula el impacto, pero se registra como "imposible de estimar"
 
  Obtener listado de elementos desconectados tras la falla
 
  Refrescar mapa a través de SQL, actualizando circuito y estado en los elementos
  A través de API, mostrar mensaje al operador, incluir impacto de la falla
 
  Registrar evento de falla franca
  Registrar nuevo proceso de atención a falla en curso
 
  TERMINA ATENCIÓN A FALLA FRANCA
 
 
 
  Al recibir mensaje de RabbitMQ iniciando Respuesta del operador
 
 ### INICIA RESPUESTA DEL OPERADOR
 
  El operador debe responder si se procede a ejecutar el análisis de corto circuito o no
 
  En caso negativo,
  Registrar evento de respuesta del operador
  Registrar proceso de atención a falla como terminado
 
  En caso positivo,
  Registrar evento de respuesta del operador
 
  Extraer la forma matricial del subgrafo
  Extraer la información eléctrica de los transformadores involucrados en el circuito dinámico
 
  Registrar evento de análisis iniciado
  Ejecutar el proceso de análisis conforme al algortimo definido en MATLAB por el
  Dr. Dionisio Antonio Suarez Cerda, tomando como base el circuito dinámico en cuestión, con la
  corriente de corto circuito reportada
 
  Registrar los resultados del proceso de análisis
  Registrar evento de análisis terminado
  Refrescar mapa a través de SQL, actualizando circuito, estado en los elementos y marcas del polígono de falla
  Registrar proceso de atención a falla como terminado
 
  TERMINA RESPUESTA DEL OPERADOR
 
 
 
  Al recibir mensaje de RabbitMQ iniciando Solicitud de Suspensión de procesos desde el editor topológico
 
  ### INICIA SUSPENSIÓN DE PROCESOS
 
  Registrar evento de solicitud de suspensión
  Descartar que exista otra suspensión vigente
 
  Registrar todos los procesos en curso como terminados
  Terminar procesos de análisis en curso
  Terminar procesos de actualización de eventos
 
  Refrescar mapa a través de SQL, restableciendo circuito y estado en los elementos
  A través de API, reportar módulo suspendido
  Registrar estado del módulo como inactivo
 
  Terminar recepción de mensajes de RABBITMQ
  - Inicio de Atencion a falla franca
  - Respuesta del Operador ante atenciones
  - Solicitud de Suspensión de procesos desde el editor topológico
 
  Otros eventos a dejar de monitorear
  - Cambios en los catálogos de la UCM
  - Cambio de estado de los elementos de conmutación de la UCM
  - Cambios en registro de licencias
 
  Registrar proceso de suspensión en curso
 
  Si todo sucede exitosamente, se procede a cambiar a estado Bloqueado
 
  TERMINA SUSPENSIÓN DE PROCESOS
 
  - Cambios en los catálogos de la UCM
  - Cambios en la configuración del módulo
  - Cambio de estado de los elementos de conmutación de la UCM
  - Cambios en registro de licencias
  - Se superó el tiempo de espera para recibir respuesta del operador
 
 
  Al recibir mensaje de RabbitMQ iniciando Solicitud de Reactivación de procesos desde el editor topológico
 
  INICIA SUSPENSIÓN DE PROCESOS
 
  Registrar evento de solicitud de reactivación
  Descartar que no exista una suspensión vigente
 
  Registrar todos los procesos en curso como terminados
 
  A través de API, reportar módulo reactivado
  Registrar estado del módulo como activo
  Registrar proceso de suspensión en curso
 
  TERMINA SUSPENSIÓN DE PROCESOS