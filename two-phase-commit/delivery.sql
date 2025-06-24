CREATE Delivery
  (DeliveryID INT PRIMARY KEY,
   OrderID INT NOT NULL,
   DeliveryDate DATE NOT NULL,
   FOREIGN KEY (OrderID) REFERENCES Orders(OrderID))