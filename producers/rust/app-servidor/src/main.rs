use rocket::response::status::BadRequest;
use rocket::serde::json::{json, Value as JsonValue};
use rocket::serde::json::Json;
use rocket::config::SecretKey;
use rocket_cors::{AllowedOrigins, CorsOptions};
use serde::Serialize; // Importa Serialize desde serde
use rdkafka::config::ClientConfig;
use rdkafka::producer::{FutureProducer, FutureRecord};
use std::sync::Arc;
use std::time::Duration;

#[derive(rocket::serde::Deserialize, Serialize)]
struct Data {
    name: String,
    album: String,
    year: String,
    rank: String,
}

async fn send_to_kafka(data: &Data) -> Result<(), Box<dyn std::error::Error>> {
    let producer: FutureProducer = ClientConfig::new()
        .set("bootstrap.servers", "my-cluster-kafka-bootstrap.kafka.svc:9092")
        .create()?;
    let payload = serde_json::to_string(data)?;
    let record = FutureRecord::to("test")
        .key(&data.album)
        .payload(&payload);

    match producer.send(record, Duration::from_secs(1)).await {
        Ok(_) => Ok(()),
        Err((e, _)) => Err(Box::new(e)),
    }
}

#[rocket::post("/data", data = "<data>")]
async fn receive_data(data: Json<Data>) -> Result<String, BadRequest<String>> {
    let received_data = data.into_inner();
    let response = JsonValue::from(json!({
        "message": format!("Received data: Name: {}, Album: {}, Year: {}, Rank: {}", received_data.name, received_data.album, received_data.year, received_data.rank)
    }));
    // print la data recibida
    println!("Received data: Name: {}, Album: {}, Year: {}, Rank: {}", received_data.name, received_data.album, received_data.year, received_data.rank);

    // Enviar datos a Kafka
    if let Err(e) = send_to_kafka(&received_data).await {
        eprintln!("Failed to send to Kafka: {:?}", e);
        return Err(BadRequest("Failed to send data to Kafka".to_string()))
    }

    Ok(response.to_string())
}

#[rocket::main]
async fn main() {
    let secret_key = SecretKey::generate(); // Genera una nueva clave secreta

    // Configuración de opciones CORS
    let cors = CorsOptions::default()
        .allowed_origins(AllowedOrigins::all())
        .to_cors()
        .expect("failed to create CORS fairing");

    let producer = Arc::new(
        ClientConfig::new()
            .set("bootstrap.servers", "my-cluster-kafka-bootstrap.votacion.svc.cluster.local:9092")
            .set("message.timeout.ms", "5000")
            .create::<FutureProducer<_>>()
            .expect("Failed to create Kafka producer")
    );

    let config = rocket::Config {
        address: "0.0.0.0".parse().unwrap(),
        port: 8080,
        secret_key: secret_key.unwrap(), // Desempaqueta la clave secreta generada
        ..rocket::Config::default()
    };

    // Montar la aplicación Rocket con el middleware CORS y el estado del productor Kafka
    rocket::custom(config)
        .manage(producer)
        .attach(cors)
        .mount("/", rocket::routes![receive_data])
        .launch()
        .await
        .unwrap();
}
