# Jarmuz-Message

Messaging JMOD for Jablko

This project is a GroupMe messaging bot for Jablko that other modules can send post requests to, which in turn send messages to a GroupMe chat. This is useful for any mods that involve long-running status updates or alerts. A quick example: A JMOD that provides a weather dashboard (that could be linked to temperature/humidity sensors or OpenWeatherMap) that sends a weather update at a given time of day.


## Consuming in another JMOD

Simply send post request to:

`http://127.0.0.1:<Jablko Core Port>/jmod/sendMessage?github.com/ccoverstreet/Jarmuz-Message` 

with the following body:

```json
{"message": "Some message"}
```
