package com.gabrielebnc.ordermatcher.OrderProducer.factory;

import com.gabrielebnc.ordermatcher.OrderProducer.model.Order;

import java.util.Random;

public class RandomOrderFactory implements OrderFactory{

    private final Random RANDOM;
    private final int meanPrice;
    private final int stdPrice;
    private final String[] symbols;
    private final Integer[] traderIds;



    public RandomOrderFactory(final int meanPrice, final int stdPrice, final String[] symbols, final Integer[] traderIds, final long seed){
        this.meanPrice = meanPrice;
        this.stdPrice = stdPrice;
        this.symbols = symbols;
        this.traderIds = traderIds;
        this.RANDOM = new Random(seed);
    }

    public RandomOrderFactory(final int meanPrice, final int stdPrice, final String[] symbols, final Integer[] traderIds){
        this.meanPrice = meanPrice;
        this.stdPrice = stdPrice;
        this.symbols = symbols;
        this.traderIds = traderIds;
        this.RANDOM = new Random();
    }

    @Override
    public Order generateOrder() {

        return new Order(
                this.RANDOM.nextInt(),
                (int) this.randomChose(traderIds),
                this.randomChose(symbols),
                this.RANDOM.nextBoolean() ? 'B' : 'S',
                this.RANDOM.nextBoolean() ? 'M' : 'L',
                this.RANDOM.nextInt(1000),
                this.generateRandomPrice(),
                (byte) 0);
    }

    private int generateRandomPrice(){
        return (int) Math.round(this.RANDOM.nextGaussian(meanPrice, stdPrice));
    }

    private <T> T randomChose(final T[] objs){
        return objs[this.RANDOM.nextInt(objs.length)];
    }


}
