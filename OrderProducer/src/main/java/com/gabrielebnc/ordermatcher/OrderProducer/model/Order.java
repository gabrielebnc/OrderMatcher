package com.gabrielebnc.ordermatcher.OrderProducer.model;

import java.time.Instant;

public class Order {
    private final int orderId;
    private final int traderId;
    private final String symbol;
    private final char side;
    private final char orderType;
    private final int quantity;
    private final int price; //The price is multiplied by 100, so 845 means 8.45
    private final Instant timestamp;
    private final byte action;

    private Order(final int orderId,
                 final int traderId,
                 final String symbol,
                 final char side,
                 final char orderType,
                 final int quantity,
                 final int price,
                 final Instant timestamp,
                 final byte action) {
        this.orderId = orderId;
        this.traderId = traderId;
        this.symbol = symbol;
        this.side = side;
        this.orderType = orderType;
        this.quantity = quantity;
        this.price = price;
        this.timestamp = timestamp;
        this.action = action;
    }

    public Order(final int orderId,
                 final int traderId,
                 final String symbol,
                 final char side,
                 final char orderType,
                 final int quantity,
                 final int price,
                 final byte action) {
        this.orderId = orderId;
        this.traderId = traderId;
        this.symbol = symbol;
        this.side = side;
        this.orderType = orderType;
        this.quantity = quantity;
        this.price = price;
        this.timestamp = Instant.now();
        this.action = action;
    }

}
