-- +goose Up

CREATE TABLE orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    symbol VARCHAR(100) NOT NULL,
    side ENUM('buy', 'sell') NOT NULL,
    type ENUM('limit', 'market') NOT NULL,
    price DECIMAL(10, 4) NULL, 
    quantity DECIMAL(15, 6) NOT NULL,
    remaining_quantity DECIMAL(15, 6) NOT NULL,
    status ENUM('open', 'completed', 'partial', 'canceled') NOT NULL DEFAULT 'open',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_symbol_side_price_time (symbol, side, price, created_at),
    INDEX idx_symbol_status (symbol, status),
    INDEX idx_created_at (created_at)
);

-- +goose Down

DROP TABLE orders;
