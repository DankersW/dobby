# Dobby

Dobby serves the house. Dobby is the bridge between all the magical wireless devices in the house.
Curretly it can talk to Thread wireless nodes via the [WSN Gateway](https://github.com/DankersW/wsn/tree/main/app/thread_gateway) using an UART CLI interface of a NRF5240 micro-controller that is plugged into a free USB slot of the device.

Dobby decodes WSN protobuf messages and plublishes the parsed data to other services that are interested.
