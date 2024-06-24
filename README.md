# Order Matcher
Order Marching Engine simulator

--------------

### Introduction

The scope of the project is to simulate a client-server simulation of an Order Matching Engine,
with low-latency communication achieved through [Chronicle Queue](https://github.com/OpenHFT/Chronicle-Queue)
<br>
<br>
#### OrderClient
The client(s) will emulate an incoming order traffic to process.\
Each client will pseudo-randomly generate orders on the base some metrics and configurations.
<br>
<br>

#### OrderMatcherServer

The server will receive and process orders.\
Processing an order implies:
- receiving, validating and keeping track of incoming and pending orders on multiple order books
- logging incoming, semi-fulfilled and fulfilled orders
- persist incoming, semi-fulfilled and fulfilled orders
- communicate back to the client for confirm of validation of the order, semi-fulfillment of order and
         fulfillment of the order
  <br>
  <br>
#### Logging 
Logging will both happen on console level and on rolling files

-----------------

## Orders

At the moment, the types and variables for each order are restricted for simplicity
### <a name="ordertypes"><a/>Order Types
Each order type can be both Buy or Sell
- **Market Order:** Buy or Sell at the current market price (lowest ask / highest bid)
- **Limit Order:** Buy os Sell at a specific price
### Order fields
Each Order will share the following informations, keeping them minimal to minimize message size
- **Order ID:**
- **TraderID:**
- **Symbol:** the asset object of the order
- **Side:** buy or sell
- **Order Type:** see [order types](#ordertypes)
- **Quantity:**
- **Price:**
- **Timestamp:** when the order is placed (client timestamp)
- **Action [Optional]**: specific request about the order (cancel order, modify order..) 
