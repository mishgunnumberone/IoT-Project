import 'dart:async';

import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;


void main() => runApp (IoTApp());

const ip  = "http://172.16.0.69:5000";

class IoTApp extends StatefulWidget {
  @override
    State<StatefulWidget> createState() => IoTAppState();
}

class IoTAppState extends State<IoTApp> {

  IoTAppState() {
      Timer _timer = new Timer.periodic(new Duration(seconds: 5), (Timer timer) {loadData();});
  }
  

  void timerCallback(Timer timer) {
    print('hello');
    //loadData();
  }

  var _jsonAirInfo = {};
  var _waterLevel = {};

  var temperatureHumidityText = "";
  var waterLevelText = "";

  loadData() async {
    String airDataUrl = ip +"/air";
    try {
      http.Response airDataResponse = await http.get(airDataUrl);

      if (airDataResponse.statusCode == 200 && airDataResponse.body.isNotEmpty) {
      _jsonAirInfo = json.decode(airDataResponse.body);
      setState(() {
              temperatureHumidityText = "🌡Temperature: ${_jsonAirInfo['temperature']}℃   💧Humidity: ${_jsonAirInfo['humidity']}%";
            });
    } else {
      setState(() {
        temperatureHumidityText = "Lost connection or no data available";
      });
    }
    } catch (e) {
      setState(() {
        temperatureHumidityText = "Lost connection or no data available";
      });
    }

    String waterLeveldataUrl = ip + "/water_level";
    try {
      http.Response waterLevelDataResponse = await http.get(waterLeveldataUrl);
      if (waterLevelDataResponse.statusCode == 200 && waterLevelDataResponse.body.isNotEmpty) {
        _waterLevel = json.decode(waterLevelDataResponse.body);
        if (_waterLevel["waterLevel"] > 300) {
          setState(() {
            waterLevelText = "ATTENTION! Water leak detected. Water level: ${_waterLevel['waterLevel']}";
            });
        } else {
          setState(() {
            waterLevelText = "No any leaks detected. Water level: ${_waterLevel['waterLevel']}";
            });
        }
      } else {
        setState(() {
          waterLevelText = "Lost connection or no data available";
          });
      }
    } catch (e) {
      setState(() {
        waterLevelText = "Lost connection or no data available";
        });
    }
    //}
  }

  turnLedOn(context) async {
    String turnLedOnUrl = ip + "/turn_light_on";
    try {
      http.Response turnLedOnResponse = await http.get(turnLedOnUrl);

      if (turnLedOnResponse.statusCode == 200) {
        Scaffold.of(context).showSnackBar(new SnackBar(
                content: new Text("Led is ON"),
              ));
      } else {
        Scaffold.of(context).showSnackBar(new SnackBar(
                content: new Text("No connection"),
              ));
      }
    } catch (e) {
      Scaffold.of(context).showSnackBar(new SnackBar(
                content: new Text("No connection"),
              ));
    }
  }

  turnLedOff(context) async {
    String turnLedOffUrl = ip +"/turn_light_off";
    try {
      http.Response turnLedOffResponse = await http.get(turnLedOffUrl);

      if (turnLedOffResponse.statusCode == 200) {
        Scaffold.of(context).showSnackBar(new SnackBar(
                content: new Text("Led is OFF"),
              ));
      } else {
        Scaffold.of(context).showSnackBar(new SnackBar(
                content: new Text("No connection"),
              ));
      }
    } catch (e) {
      Scaffold.of(context).showSnackBar(new SnackBar(
                content: new Text("No connection"),
              ));
    }
  }

  void initState() {
      super.initState();

      loadData();
    }

  @override
    Widget build(BuildContext context) {
      return new MaterialApp(
        debugShowCheckedModeBanner: false,
        title: 'IoT',
        theme: ThemeData(
          primarySwatch: Colors.green
        ),
        home: new Scaffold(
          appBar: new AppBar(
            title: new Text('IoT'),
            actions: <Widget>[
              IconButton(
                icon: Icon(Icons.refresh),
                onPressed: loadData,
              )
            ],
          ),
          body: Builder(
        builder: (context) =>
          
           Center(
            child: new Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: <Widget>[
                new Text("Air Conditions", style: new TextStyle(fontWeight: FontWeight.bold, fontSize: 25.0)),
                new Text(temperatureHumidityText, style: new TextStyle(fontSize: 18.0)),
                const SizedBox(height: 20,),
                new Text("Water Level", style: new TextStyle(fontWeight: FontWeight.bold, fontSize: 25.0)),
                new Text(waterLevelText, style: new TextStyle(fontSize: 18.0)),
                const SizedBox(height: 20,),  
                new Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: <Widget>[
                    new RaisedButton(
                      child: new Text("Led On"),
                      onPressed: () {turnLedOn(context);},
                    ),
                    new SizedBox(width: 40,),
                    new RaisedButton(
                      child: new Text("Led Off"),
                      onPressed: () {turnLedOff(context);},
                    )
                  ],
                )
                
                
              ],
            ) 
          ),
          ),
        ),
      );
    }     
}

