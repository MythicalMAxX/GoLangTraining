# Implement an Order Processing System
## 🛒 Task Overview
Build an Order Management API where users can:
- Place orders (atomic transactions for order creation and inventory update)
- Fetch order details (efficient querying with joins)
- Cancel orders (soft delete with status tracking)
- Process orders asynchronously (using background workers)
- Store order history in MongoDB
## 🔹 Assignment Details
### 1️⃣ Database Schema Setup
- You will create three tables in PostgreSQL and integrate MongoDB for storing order history.
#### **Users Table**
- Use the **users** table


### Orders Table
| Column     | Type          | Description                                      |
|------------|--------------|--------------------------------------------------|
| id         | uuid         | Unique order ID                                  |
| user_id    | uuid         | Foreign key to users table                      |
| amount     | decimal      | Order amount                                    |
| status     | varchar(20)  | (pending, processing, completed, canceled)      |
| created_at | timestamp    | Order creation time                             |

### Inventory Table
| Column | Type   | Description          |
|--------|--------|----------------------|
| id     | uuid   | Unique product ID    |
| name   | varchar| Product name         |
| stock  | int    | Available stock      |

### MongoDB Collection: Order History
| Field    | Type      | Description                                      |
|----------|----------|--------------------------------------------------|
| _id      | ObjectId | Unique order history ID                          |
| order_id | uuid     | Reference to order ID in PostgreSQL             |
| status   | string   | Order status (pending, completed, etc.)         |
| log      | array    | List of status updates with timestamps          |


### 2️⃣ Make API Endpoints for these scenarios
#### Description
- Place a new order
- Get order details (optimized query)
- Cancel an order (soft delete)
- Fetch orders for a user
- Get available inventory
- Retrieve order history from MongoDB
3️⃣ Features to Implement
### ✅ Atomic Transactions:
- When a new order is placed:
    - Check inventory availability
    - Deduct stock
    - Insert order entry
    - Store order history in MongoDB
    - All steps should be inside a transaction
### ✅ Asynchronous Order Processing:
- Orders in "pending" status should be picked up by a background worker
- The worker should update the order status to "completed"
- MongoDB should log each status update
### ✅ Optimized Querying & Caching:
- Fetch orders with joined user details (JOIN query in PostgreSQL)
- Store order history logs in MongoDB for quick access
### ✅ MongoDB Integration:
- Every order status update should be logged in MongoDB
- Implement an API endpoint to fetch order history from MongoDB