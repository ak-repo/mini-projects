CREATE TABLE IF NOT EXISTS users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(150) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (name, email) VALUES
('Anil Kumar', 'anil.kumar@example.com'),
('Rahul Sharma', 'rahul.sharma@example.com'),
('Sneha Iyer', 'sneha.iyer@example.com');
