# Inventory Management Server

## Description

- This is a simple inventory management server that allows users to create, read, update, and delete items from an inventory.

- The server is build using GO, gorilla/mux.

## Database Backend

- Mongo DB

## API Endpoints

- **GET** `/inventory` - Returns all items in the inventory
- **GET** `/inventory/{id}` - Returns a single item in the inventory
- **POST** `/inventory` - Creates a new item in the inventory
  - Body:
    - `name` - Name of the item
    - `description` - Description of the item
    - `quantity` - Quantity of the item
    - `price` - Price of the item
- **PUT** `/inventory/{id}` - Updates an item in the inventory
  - Body:
    - `name` - Name of the item
    - `description` - Description of the item
    - `quantity` - Quantity of the item
    - `price` - Price of the item
- **DELETE** `/inventory/{id}` - Deletes an item in the inventory
