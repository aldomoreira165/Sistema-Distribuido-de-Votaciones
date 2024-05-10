# Proyecto 2: Sistema Distribuido de Votaciones

## Introducción
El proyecto se centra en la implementación de un sistema distribuido de votaciones para un concurso de bandas de música guatemalteca. La arquitectura se basa en microservicios desplegados en Kubernetes, con sistemas de mensajería para encolar servicios y dashboards en Grafana para visualizar datos en tiempo real. Se utilizan tecnologías como gRPC, Rust, Redis, MongoDB, entre otras, para lograr un sistema eficiente y escalable.

Dicha arquitectura está pensada para asegurar el envío más rápido y eficiente de los datos por medio de dos producers, por lo cual cada uno de ellos forma una ruta distinta y se desea determinar cual de estas es más rápida y efectiva.

## Objetivos
1. Implementar un sistema distribuido con microservicios en Kubernetes.
2. Utilizar sistemas de mensajería para encolar distintos servicios.
3. Utilizar Grafana como interfaz gráfica de dashboards.
4. Desplegar servicios en Cloud Run para la visualización de registros de MongoDB.

## Tecnologías Utilizadas
- **Kubernetes (K8S)**
Kubernetes es una plataforma de orquestación de contenedores que permite automatizar el despliegue, escalado y gestión de aplicaciones en contenedores. Proporciona una manera eficiente de administrar los recursos de manera distribuida, facilitando la implementación de microservicios y la gestión de contenedores a escala.

- **gRPC**
gRPC es un framework de comunicación de servicios remotos desarrollado por Google. Utiliza el protocolo HTTP/2 para la serialización de datos y ofrece soporte para múltiples lenguajes de programación. Permite la definición de servicios y mensajes mediante archivos .proto y proporciona una comunicación eficiente y bidireccional entre clientes y servidores.

- **Rust**
Rust es un lenguaje de programación de sistemas diseñado para la seguridad y el rendimiento. Se enfoca en prevenir errores de memoria y garantizar la seguridad de los programas mediante un sistema de tipos robusto y un control estricto sobre el manejo de memoria. Rust es especialmente adecuado para el desarrollo de sistemas de bajo nivel, como controladores de dispositivos, sistemas operativos y servidores web.

- **Redis**
Redis es un sistema de almacenamiento en memoria que se utiliza como base de datos en caché, almacenamiento de sesiones, cola de mensajes y más. Es conocido por su alta velocidad y capacidad de almacenamiento de datos en estructuras de datos como strings, hashes, listas, sets y sorted sets. En este proyecto, se utiliza Redis para almacenar contadores de votaciones en tiempo real.

- **MongoDB**
MongoDB es una base de datos NoSQL orientada a documentos que ofrece escalabilidad y flexibilidad para el almacenamiento de datos no estructurados. Permite el almacenamiento de datos en formato JSON-like, lo que facilita la manipulación y consulta de datos. Se utiliza en este proyecto para almacenar registros de logs generados por el sistema.

- **Grafana**
Grafana es una plataforma de análisis y visualización de datos que proporciona dashboards personalizables para monitorear y analizar métricas en tiempo real. Permite la conexión con una variedad de fuentes de datos, incluidas bases de datos, sistemas de monitoreo y servicios en la nube. En este proyecto, se utiliza Grafana para visualizar contadores de votaciones en tiempo real almacenados en Redis.

- **Cloud Run**
Cloud Run es un servicio de Google Cloud que permite ejecutar contenedores de forma gestionada y escalable. Proporciona una plataforma sin servidor para desplegar aplicaciones en contenedores y ofrece escalado automático según la demanda. En este proyecto, se utiliza Cloud Run para desplegar una API en Node.js y una webapp en Vue.js para visualizar registros de logs de MongoDB.

- **Vue JS**
Vue.js es un framework progresivo de JavaScript utilizado para construir interfaces de usuario interactivas y de una sola página. Ofrece una arquitectura basada en componentes y una sintaxis intuitiva que facilita el desarrollo de aplicaciones web complejas. En este proyecto, se utiliza Vue.js para desarrollar la interfaz de usuario de la webapp que visualiza registros de logs de MongoDB en Cloud Run.

- **Kafka**
Apache Kafka es una plataforma de streaming distribuido que permite la transmisión de datos de manera segura, duradera y tolerante a fallos. Actúa como un sistema de mensajería distribuida, donde los productores pueden enviar mensajes a los topics de Kafka, y los consumidores pueden suscribirse a estos topics para recibir los mensajes. Kafka es altamente escalable y se utiliza en aplicaciones de streaming en tiempo real para el procesamiento de eventos.

- **Go**
Go, también conocido como Golang, es un lenguaje de programación de código abierto desarrollado por Google. Se caracteriza por su sintaxis concisa, su rendimiento eficiente y su soporte para concurrencia y paralelismo. Go es ampliamente utilizado en el desarrollo de servicios de backend y sistemas distribuidos debido a su facilidad de uso y su potente estándar bibliotecas.

- **Python**
Python es un lenguaje de programación de alto nivel conocido por su sintaxis clara y legible. Es ampliamente utilizado en una variedad de aplicaciones, incluyendo desarrollo web, análisis de datos, inteligencia artificial, scripting, entre otros. Python cuenta con una amplia comunidad de desarrolladores y una extensa biblioteca estándar que facilita el desarrollo de aplicaciones de forma rápida y eficiente.

- **Locust**
Locust es una herramienta de código abierto utilizada para probar la escalabilidad y el rendimiento de aplicaciones web. Permite simular la carga de usuarios enviando solicitudes HTTP a una aplicación y registrando el tiempo de respuesta y otros datos relevantes. Locust se utiliza comúnmente para realizar pruebas de carga y estrés en sistemas distribuidos, como el que se implementa en este proyecto.

## Deployments y Services
- Namespace:
El archivo namespace.yaml define un Namespace en Kubernetes con el nombre kafka. Un Namespace proporciona un alcance para los nombres de los recursos dentro del clúster, lo que ayuda a organizar y gestionar los recursos de manera más eficiente. En este caso, el Namespace kafka se utiliza para agrupar todos los recursos relacionados con el proyecto dentro de un mismo ámbito, lo que facilita su gestión y mantenimiento.

    ```yaml
    apiVersion: v1
    kind: Namespace
    metadata:
      name: kafka
    ```

 - Ingress: 
El archivo ingress.yaml define las reglas de enrutamiento del Ingress en Kubernetes para el namespace kafka. Establece que las solicitudes con prefijo /grpc/ se redirigen al servicio producers-grpc-service en el puerto 3000, mientras que las solicitudes con prefijo /rust/ se redirigen al servicio producers-rust-service en el puerto 8000. Además, configura el redireccionamiento del tráfico HTTP a HTTPS, especificando los puertos en los que se escuchará el tráfico HTTP y HTTPS, y utilizando el controlador de Ingress de nginx.

    ```yaml
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: my-ingress
      namespace: kafka
      annotations:
        nginx.ingress.kubernetes.io/target-type: ip
        nginx.org/listen-ports: '[80,3000,3001,8000,8080]'
        nginx.org/listen-ports-ssl: '[443]'
        ingress.kubernetes.io/ssl-redirect: "true"
        nginx.ingress.kubernetes.io/service-upstream: "true"
    spec:
      ingressClassName: nginx
      rules:
      - http:
          paths:
          - path: /grpc/
            pathType: Prefix
            backend:
              service:
                name: producers-grpc-service
                port:
                  number: 3000
          - path: /rust/
            pathType: Prefix
            backend:
              service:
                name: producers-rust-service
                port:
                  number: 8000
    ```

- gRPC:
El archivo producers-grpc.yaml define un Deployment y un Service en Kubernetes para los servicios gRPC en el namespace kafka. El Deployment producers-grpc especifica un único pod con dos contenedores: uno para el cliente y otro para el servidor gRPC, utilizando las imágenes aldomoreirav/grpc-client:latest y aldomoreirav/grpc-server:latest, respectivamente. El Service producers-grpc-service expone los puertos 3000 y 3001 para el cliente y el servidor gRPC, respectivamente, con un tipo de servicio ClusterIP y selecciona los pods con la etiqueta role: producers-grpc.
    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: producers-grpc
      namespace: kafka
    spec:
      selector:
        matchLabels:
          role: producers-grpc
      replicas: 1
      template:
        metadata:
          labels:
            role: producers-grpc
        spec:
          containers:
          - name: client
            image: aldomoreirav/grpc-client:latest
            ports:
            - containerPort: 3000
          - name: server
            image: aldomoreirav/grpc-server:latest
            ports:
            - containerPort: 3001
    apiVersion: v1
    kind: Service
    metadata:
      name: producers-grpc-service
      namespace: kafka
    spec:
        type: ClusterIP
        ports:
        - name: client
          port: 3000
          targetPort: 3000
        - name: server
          port: 3001
          targetPort: 3001
        selector:
          role: producers-grpc
    ---      
    apiVersion: v1
    kind: Service
    metadata:
      name: producers-grpc-service
      namespace: kafka
    spec:
        type: ClusterIP
        ports:
        - name: client
          port: 3000
          targetPort: 3000
        - name: server
          port: 3001
          targetPort: 3001
        selector:
          role: producers-grpc
          
    ```

- Rust
El archivo producers-rust.yaml define un Deployment y un Service en Kubernetes para los servicios Rust en el namespace kafka. El Deployment producers-rust especifica un único pod con dos contenedores: uno para el cliente Rust y otro para el servidor Rust, utilizando las imágenes aldomoreirav/rust-client y aldomoreirav/rust-server, respectivamente. El Service producers-rust-service expone los puertos 8000 y 8080 para el cliente y el servidor Rust, respectivamente, con un tipo de servicio ClusterIP y selecciona los pods con la etiqueta role: producers-rust.

    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: producers-rust
      namespace: kafka
    spec:
      selector:
        matchLabels:
          role: producers-rust
      replicas: 1
      template:
        metadata:
          labels:
            role: producers-rust
    
        spec:
          containers:
          - name: rust-client
            image: aldomoreirav/rust-client
            ports:
            - containerPort: 8000
          - name: rust-server
            image: aldomoreirav/rust-server
            ports:
            - containerPort: 8080
    ---
    apiVersion: v1
    kind: Service
    metadata:
      name: producers-rust-service
      namespace: kafka
    spec:
      type: ClusterIP
      ports:
      - name: rust-client-port
        port: 8000
        targetPort: 8000
      - name: rust-server-port
        port: 8080
        targetPort: 8080
      selector:
        role: producers-rust
    ```

- Kafka
El archivo topic-kafka.yaml define un KafkaTopic en Kubernetes en el namespace kafka. Un KafkaTopic es un recurso de Strimzi que representa un tema en un clúster de Apache Kafka. En este caso, se define un tema con el nombre test y se asocia al clúster my-cluster utilizando las etiquetas strimzi.io/cluster: my-cluster. Se especifica que el tema tendrá 1 partición y 1 réplica, y se configuran parámetros adicionales como la retención de mensajes (retention.ms) y el tamaño del segmento (segment.bytes). Este recurso define la configuración del tema en Kafka y permite su gestión dentro del clúster Kubernetes utilizando Strimzi.
    ```yaml
    apiVersion: kafka.strimzi.io/v1beta2
    kind: KafkaTopic
    metadata:
      namespace: kafka
      name: test
      labels:
        strimzi.io/cluster: my-cluster
    spec:
      partitions: 1
      replicas: 1
      config:
        retention.ms: 7200000
        segment.bytes: 1073741824
    ```

- Consumer
El archivo consumer.yaml define un Deployment y un Service en Kubernetes para el consumidor en el namespace kafka. El Deployment consumer-deployment especifica un único pod con un contenedor llamado consumer, utilizando la imagen aldomoreirav/consumer:latest. El Service consumer-service expone el puerto 3003 para el consumo del consumidor, con un tipo de servicio ClusterIP y selecciona los pods con la etiqueta role: consumer-deployment. Este recurso permite que el consumidor esté disponible para recibir solicitudes dentro del clúster Kubernetes y se pueda acceder a él a través del Service consumer-service.

    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: consumer-deployment
      namespace: kafka
    spec:
      selector:
        matchLabels:
          role: consumer-deployment
      replicas: 1
      template:
        metadata:
          labels:
            role: consumer-deployment
        spec:
          containers:
          - name: consumer
            image: aldomoreirav/consumer:latest
    
    ---
    
    apiVersion: v1
    kind: Service
    metadata:
      name: consumer-service
      namespace: kafka
    spec:
      ports:
      - name: consumer
        port: 3003
        targetPort: 3003
      selector:
        role: consumer-deployment
    ```

- Redis
El archivo redis.yaml define un Deployment y un Service en Kubernetes para Redis en el namespace kafka. El Deployment redis especifica un único pod con un contenedor llamado redisdb, utilizando la imagen redis:latest y exponiendo el puerto 6379 para la comunicación con Redis. El Service redis expone el puerto 6379 utilizando un balanceador de carga (LoadBalancer) para permitir el acceso externo a Redis desde fuera del clúster Kubernetes. Ambos recursos tienen la etiqueta role: redis, lo que permite que el Service enrute el tráfico al pod correspondiente en el Deployment. Esto facilita el acceso y el uso de Redis dentro del clúster Kubernetes.

    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: redis
      namespace: kafka
    spec:
      replicas: 1
      selector:
        matchLabels:
          role: redis
      template:
        metadata:
          labels:
            role: redis
        spec:
          containers:
          - name: redisdb
            image: redis:latest
            ports:
            - containerPort: 6379
    ---
    apiVersion: v1
    kind: Service
    metadata:
      name: redis
      namespace: kafka
    spec:
      type: LoadBalancer
      ports:
      - port: 6379
        targetPort: 6379
      selector:
        role: redis
    ```

- Grafana
El archivo grafana.yaml define un Deployment y un Service en Kubernetes para Grafana en el namespace kafka. El Deployment grafana especifica un único pod con un contenedor llamado grafana, utilizando la imagen grafana/grafana:8.4.4 y exponiendo el puerto 3000 para acceder a Grafana. El Service grafana expone el puerto 3000 utilizando un balanceador de carga (LoadBalancer) para permitir el acceso externo a Grafana desde fuera del clúster Kubernetes. Ambos recursos tienen la etiqueta app: grafana, lo que permite que el Service enrute el tráfico al pod correspondiente en el Deployment. Esto facilita el acceso y la configuración de Grafana para la visualización de datos en tiempo real.

    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: grafana
      namespace: kafka
    spec:
      replicas: 1
      selector:
        matchLabels:
          app: grafana
      template:
        metadata:
          name: grafana
          labels:
            app: grafana
        spec:
          containers:
          - name: grafana
            image: grafana/grafana:8.4.4
            ports:
            - name: grafana
              containerPort: 3000
    
    ---
    apiVersion: v1
    kind: Service
    metadata:
      name: grafana
      namespace: kafka
    spec:
      type: LoadBalancer
      ports:
      - port: 3000
        targetPort: 3000
      selector:
        app: grafana
    ```    
    
- MongoDB
El archivo mongodb.yaml define un Deployment y un Service en Kubernetes para MongoDB en el namespace kafka. El Deployment mongodb especifica un único pod con un contenedor llamado mongodb, utilizando la imagen mongo:latest y exponiendo el puerto 27017 para la comunicación con MongoDB. El Service mongodb expone el puerto 27017 utilizando un balanceador de carga (LoadBalancer) para permitir el acceso externo a MongoDB desde fuera del clúster Kubernetes. Ambos recursos tienen la etiqueta role: mongodb, lo que permite que el Service enrute el tráfico al pod correspondiente en el Deployment. Esto facilita el acceso y el uso de MongoDB dentro del clúster Kubernetes.

    ```yaml
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: mongodb
      namespace: kafka
    spec:
      replicas: 1
      selector:
        matchLabels:
          role: mongodb
      template:
        metadata:
          labels:
            role: mongodb
        spec:
          containers:
          - name: mongodb
            image: mongo:latest
            ports:
            - containerPort: 27017
    ---
    apiVersion: v1
    kind: Service
    metadata:
      name: mongodb
      namespace: kafka
    spec:
      type: LoadBalancer
      ports:
      - port: 27017
        targetPort: 27017
      selector:
        role: mongodb
    ```

## Ejemplos
1. Envío de Votaciones: Locust genera tráfico simulado enviando datos de votaciones a los servicios gRPC y Wasm.
![locust](https://imgur.com/UcqpnL0.png)

2. Encolado de Datos: Los servicios gRPC y Wasm envían los datos a Kafka para su encolado.
Procesamiento de Datos: El consumidor Daemon procesa los datos encolados por Kafka y los envía a Redis y MongoDB.
![mongo](https://imgur.com/joISUHo.png)
![redis](https://imgur.com/BaPRpXa.png)

3. Visualización en Grafana: Grafana muestra en tiempo real los contadores de votaciones almacenados en Redis.
![grafana](https://imgur.com/lqheiyY.png)

4. Visualización en Cloud Run: La API y la webapp desplegadas en Cloud Run permiten visualizar los registros de logs de MongoDB.
![vue](https://imgur.com/3KvEwLa.png)

## Conclusiones

El proyecto de implementación del Sistema Distribuido de Votaciones ha sido exitoso en la integración de diversas tecnologías y prácticas para lograr un sistema eficiente y escalable. La utilización de Kubernetes como plataforma de orquestación de contenedores ha permitido gestionar de manera efectiva los recursos y desplegar los microservicios de forma distribuida. La incorporación de tecnologías como gRPC y Web Assembly ha facilitado la comunicación entre los distintos componentes del sistema, garantizando un intercambio eficiente de datos.

La inclusión de Kafka como sistema de mensajería ha mejorado la capacidad de procesamiento de datos, permitiendo la encolación y distribución de mensajes de manera robusta y tolerante a fallos. Además, la integración con Redis y MongoDB ha proporcionado soluciones de almacenamiento adecuadas para datos en tiempo real y logs, respectivamente. La visualización de métricas en tiempo real a través de Grafana ha brindado una interfaz intuitiva y poderosa para monitorear el flujo de votaciones y analizar el rendimiento del sistema.

El despliegue de servicios en Cloud Run ha ofrecido una solución sin servidor para la visualización de registros de logs, garantizando la disponibilidad y escalabilidad de la aplicación web. La adopción de tecnologías como Vue.js para el frontend y Go y Rust para el backend ha permitido desarrollar aplicaciones modernas y eficientes, asegurando una experiencia de usuario óptima.

En resumen, el proyecto ha demostrado la viabilidad y eficacia de la implementación de sistemas distribuidos utilizando tecnologías avanzadas y prácticas de desarrollo modernas. La combinación de herramientas y tecnologías seleccionadas ha resultado en un sistema robusto, escalable y altamente disponible, capaz de satisfacer los requerimientos del concurso de bandas de música guatemalteca de manera eficiente y efectiva.