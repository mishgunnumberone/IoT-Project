#include "dht.h"
#define dht_apin A1

dht DHT;

const int waterSensorPin = A2;
const int ledPin = 2;
const int pResistorPin = A0;
const int buzzerPin = 3;




int ledInputPinValue = 1;

bool isLedTurnedOn = true;

void setup() {
  // put your setup code here, to run once:
  
  pinMode(waterSensorPin, INPUT);
  pinMode(ledPin, OUTPUT);
  pinMode(pResistorPin, INPUT);
  pinMode(buzzerPin, OUTPUT);

    
  Serial.begin(9600);
  delay(500);
  
  
}

int printSerial_pResistorValue() {
  int pResistorValue = analogRead(pResistorPin);
  Serial.print("\t");
  Serial.print(pResistorValue);
}


void alarmLeak() {
  int waterSensorValue = analogRead(waterSensorPin);
  if (waterSensorValue >= 300) {
    for (int i = 0; i < 10; i++) {
      digitalWrite(buzzerPin, HIGH);
      delay(5);
      digitalWrite(buzzerPin, LOW);
    }
  } 
}

void loop() {
  DHT.read11(dht_apin);
  

  int h = DHT.humidity;

  // Read temperature as Celsius (the default)

  int t = DHT.temperature;

  int waterSensorValue = analogRead(waterSensorPin);

  if (Serial.available() > 0) {
    ledInputPinValue = Serial.parseInt();
  }

  if (ledInputPinValue == 1) {
    isLedTurnedOn = true;
  } else if (ledInputPinValue == 0) {
    isLedTurnedOn = false;
    digitalWrite(ledPin, LOW);
  }

  if (isLedTurnedOn) {
    int pResistorValue = analogRead(pResistorPin);
    if (pResistorValue < 380) {
      digitalWrite(ledPin, LOW);
    } else {
      digitalWrite(ledPin, HIGH);
    }
  }

  alarmLeak();


  //digitalWrite(2, HIGH);
  Serial.print(h);

  Serial.print("\t"); // for splitting

  

  Serial.print(t);

  Serial.print("\t"); // for splitting


  Serial.print(waterSensorValue);

  Serial.print("\n"); // for new line

  delay(950);

  
}
