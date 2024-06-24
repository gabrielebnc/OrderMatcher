package com.gabrielebnc.ordermatcher.OrderMatcherServer;

public class OrderMatcherServer {

    /**
     *  The OrderMatcherServer
     *  the server will receive and process orders.
     *  Processing an order implies:
     *      - receiving, validating and keeping track of incoming and pending orders on multiple order books
     *      - logging incoming, semi-fulfilled and fulfilled orders
     *      - persist incoming, semi-fulfilled and fulfilled orders
     *      - communicate back to the client for confirm of validation of the order, semi-fulfillment of order and
     *          fulfillment of the order
     * <p>
     *  Logging:
     *      the logging will both happen on console level and on rolling files
     */
    public static void main(String[] args) {

    }
}
