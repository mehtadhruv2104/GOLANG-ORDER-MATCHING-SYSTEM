-- +goose Up
CREATE TABLE trades (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(100) NOT NULL,
    buy_order_id BIGINT NOT NULL,
    sell_order_id BIGINT NOT NULL,
    price DECIMAL(10, 4) NOT NULL,
    quantity DECIMAL(15, 6) NOT NULL,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (buy_order_id) REFERENCES orders(id),
    FOREIGN KEY (sell_order_id) REFERENCES orders(id),
    INDEX idx_symbol_time (symbol, executed_at),
    INDEX idx_buy_order (buy_order_id),
    INDEX idx_sell_order (sell_order_id)
);                                  

-- +goose Down

DROP TABLE trades;
