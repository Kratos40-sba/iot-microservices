
#include "DHT.h"
#include <ESP8266WiFi.h>
#include <Ticker.h>
#include <AsyncMqttClient.h>

#define WIFI_SSID "-----------"
#define WIFI_PASSWORD "-----------"

// Raspberry Pi Mosquitto MQTT Broker
#define MQTT_HOST IPAddress(192, 168, 1, X)
// For a cloud MQTT broker, type the domain name
//#define MQTT_HOST "broker.hivemq.com"
#define MQTT_PORT 1883

// Temperature MQTT Topics
//#define MQTT_PUB_TEMP "esp/dht/temperature"
//#define MQTT_PUB_HUM "esp/dht/humidity"
#define MQTT_PUB "esp/sensor"
#define MQTT_SUB "esp/action"
// Digital pin connected to the DHT sensor
#define DHTPIN 0
// Digital pin connected to the soil sensor
#define SOILPIN A0
// Digital pin connected to the led
//#define LEDPIN 5
const int led = 5 ;
// Uncomment whatever DHT sensor type you're using
#define DHTTYPE DHT11   // DHT 11
//#define DHTTYPE DHT22   // DHT 22  (AM2302), AM2321
//#define DHTTYPE DHT21   // DHT 21 (AM2301)

// Initialize DHT sensor
DHT dht(DHTPIN, DHTTYPE);

// Variables to hold sensor readings
float temp;
float hum;
float percentage;
AsyncMqttClient mqttClient;
Ticker mqttReconnectTimer;

WiFiEventHandler wifiConnectHandler;
WiFiEventHandler wifiDisconnectHandler;
Ticker wifiReconnectTimer;

unsigned long previousMillis = 0;   // Stores last time temperature was published
const long interval = 10000;        // Interval at which to publish sensor readings

void connectToWifi() {
  Serial.println("Connecting to Wi-Fi...");
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
}

void onWifiConnect(const WiFiEventStationModeGotIP& event) {
  Serial.println("Connected to Wi-Fi.");
  connectToMqtt();
}

void onWifiDisconnect(const WiFiEventStationModeDisconnected& event) {
  Serial.println("Disconnected from Wi-Fi.");
  mqttReconnectTimer.detach(); // ensure we don't reconnect to MQTT while reconnecting to Wi-Fi
  wifiReconnectTimer.once(2, connectToWifi);
}

void connectToMqtt() {
  Serial.println("Connecting to MQTT...");
  mqttClient.connect();
}

void onMqttConnect(bool sessionPresent) {
  Serial.println("Connected to MQTT.");
  Serial.print("Session present: ");
  Serial.println(sessionPresent);
}

void onMqttDisconnect(AsyncMqttClientDisconnectReason reason) {
  Serial.println("Disconnected from MQTT.");

  if (WiFi.isConnected()) {
    mqttReconnectTimer.once(2, connectToMqtt);
  }
}

void onMqttSubscribe(uint16_t packetId, uint8_t qos) {
  Serial.println("Subscribe acknowledged.");
  Serial.print("  packetId: ");
  Serial.println(packetId);
  Serial.print("  qos: ");
  Serial.println(qos);
}


void onMqttPublish(uint16_t packetId) {
  Serial.print("Publish acknowledged.");
  Serial.print("  packetId: ");
  Serial.println(packetId);
}

void setup() {

  Serial.begin(9600);//Serial.begin(115200);
  Serial.println();
  pinMode(led , OUTPUT) ;
  delay(100);
  digitalWrite(led ,HIGH);

  dht.begin();
  wifiConnectHandler = WiFi.onStationModeGotIP(onWifiConnect);
  wifiDisconnectHandler = WiFi.onStationModeDisconnected(onWifiDisconnect);

  mqttClient.onConnect(onMqttConnect);
  mqttClient.onDisconnect(onMqttDisconnect);
  mqttClient.onSubscribe(onMqttSubscribe);
  mqttClient.onMessage(onMqttMessage);
  //mqttClient.onUnsubscribe(onMqttUnsubscribe);
  mqttClient.onPublish(onMqttPublish);

  mqttClient.setServer(MQTT_HOST, MQTT_PORT);
  // If your broker requires authentication (username and password), set them below
  //mqttClient.setCredentials("REPlACE_WITH_YOUR_USER", "REPLACE_WITH_YOUR_PASSWORD");

  connectToWifi();
}
// handle incoming message
void onMqttMessage(char* topic, char* payload, AsyncMqttClientMessageProperties properties, size_t len, size_t index, size_t total) {
  Serial.print("Action arrived in topic: ");
  Serial.println(topic);
  Serial.print("Action:");
   const  char* one = "1";
  if(strcmp(payload,one)==0){
    digitalWrite(led, HIGH);
    Serial.print("LED IS ON") ;
  }else{
    digitalWrite(led, LOW);
  }
}

void loop() {
  unsigned long currentMillis = millis();
  // Every X number of seconds (interval = 10 seconds)
  // it publishes a new MQTT message
  if (currentMillis - previousMillis >= interval) {
    // Save the last time a new reading was published
    previousMillis = currentMillis;
    // subscribe
      mqttClient.subscribe(MQTT_SUB,1);
    // New DHT sensor readings
    hum = dht.readHumidity();
    // Read temperature as Celsius (the default)
    temp = dht.readTemperature();
    // Read Soil Value
    percentage =  ( 100.00 - ( (analogRead(SOILPIN)/1023.00) * 100.00 ) );

    // Publish an MQTT message on topic esp/sensor
    char *msg = "";
    snprintf (msg, 75, "%.2f||%.2f||%.2f",temp,hum,percentage);
    // Publish an MQTT message on topic esp/dht/humidity
    uint16_t packetIdPub = mqttClient.publish(MQTT_PUB, 1, true, msg);
    Serial.printf("Publishing on topic %s at QoS 1, packetId %i: \n", MQTT_PUB, packetIdPub);
    Serial.println(msg);
    // Subscribe to an MQTT message on topic esp/action
  }
}