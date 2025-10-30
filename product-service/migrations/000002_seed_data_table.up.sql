INSERT INTO products (
    user_id, 
    name, 
    description, 
    price, 
    image_url, 
    stock, 
    category,
    created_at,
    updated_at
) VALUES
(
    1,
    'Classic White T-Shirt',
    'Soft cotton crew-neck t-shirt. Perfect for everyday wear. Machine washable.',
    19.99,
    'https://cdn.example.com/images/tshirt-white.jpg',
    150,
    '{"Men", "Women"}',
    '2025-10-25 10:30:00+00',
    '2025-10-25 10:30:00+00'
),
(
    2,
    'Slim Fit Denim Jeans',
    'Stretch denim with slim fit. Mid-rise, 5-pocket design. Available in dark wash.',
    59.99,
    'https://cdn.example.com/images/jeans-slim.jpg',
    80,
    '{"Men"}',
    '2025-10-24 14:20:00+00',
    '2025-10-26 09:15:00+00'
),
(
    3,
    'Floral Summer Dress',
    'Lightweight chiffon dress with floral print. Knee-length, perfect for warm days.',
    44.50,
    'https://cdn.example.com/images/dress-floral.jpg',
    45,
    '{"Women"}',
    '2025-10-23 11:00:00+00',
    '2025-10-23 11:00:00+00'
),
(
    18,
    'Kids Dinosaur Hoodie',
    'Cozy fleece hoodie with T-Rex print and kangaroo pocket. Fun and warm!',
    32.99,
    'https://cdn.example.com/images/hoodie-dino.jpg',
    120,
    '{"Boy"}',
    '2025-10-22 16:45:00+00',
    '2025-10-27 08:30:00+00'
),
(
    13,
    'Unicorn Sparkle Leggings',
    'Stretchy leggings with glitter unicorn print. Elastic waistband for comfort.',
    24.99,
    'https://cdn.example.com/images/leggings-unicorn.jpg',
    90,
    '{"Girl"}',
    '2025-10-21 09:10:00+00',
    '2025-10-21 09:10:00+00'
),
(
    16,
    'Leather Crossbody Bag',
    'Genuine leather mini bag with adjustable strap. Ideal for daily essentials.',
    89.00,
    'https://cdn.example.com/images/bag-crossbody.jpg',
    30,
    '{"Women"}',
    '2025-10-20 13:25:00+00',
    '2025-10-28 12:00:00+00'
),
(
    17,
    'Striped Polo Shirt',
    'Breathable cotton polo with contrast stripes. Button placket and ribbed collar.',
    27.50,
    'https://cdn.example.com/images/polo-striped.jpg',
    65,
    '{"Men", "Boy"}',
    '2025-10-19 17:55:00+00',
    '2025-10-19 17:55:00+00'
),
(
    18,
    'Butterfly Graphic Tee',
    '100% cotton t-shirt featuring colorful butterfly artwork. Soft and durable.',
    22.00,
    'https://cdn.example.com/images/tee-butterfly.jpg',
    200,
    '{"Girl", "Women"}',
    '2025-10-18 08:40:00+00',
    '2025-10-29 10:20:00+00'
),
(
    19,
    'Cargo Jogger Pants',
    'Relaxed fit joggers with multiple cargo pockets. Drawstring waist and cuffs.',
    49.99,
    'https://cdn.example.com/images/joggers-cargo.jpg',
    55,
    '{"Men"}',
    '2025-10-17 12:15:00+00',
    '2025-10-17 12:15:00+00'
),
(
    20,
    'Rainbow Tutu Skirt',
    'Fluffy layered tulle skirt in rainbow colors. Elastic waist for easy wear.',
    28.75,
    'https://cdn.example.com/images/skirt-tutu.jpg',
    110,
    '{"Girl"}',
    '2025-10-16 15:30:00+00',
    '2025-10-30 07:45:00+00'
);


--  curl -i localhost:5000/v1/products
-- HTTP/1.1 200 OK
-- Content-Type: application/json
-- Vary: Authorization
-- Date: Thu, 30 Oct 2025 12:12:47 GMT
-- Transfer-Encoding: chunked

-- {
--         "metadata": {
--                 "current_page": 1,
--                 "page_size": 20,
--                 "first_page": 1,
--                 "last_page": 1,
--                 "total_records": 10
--         },
--         "products": [
--                 {
--                         "id": 11,
--                         "user_id": 0,
--                         "name": "Men's Classic White Shirt",
--                         "description": "A crisp white dress shirt made from premium cotton, perfect for formal occasions.",
--                         "price": "$ 49.99",
--                         "image_url": "https://example.com/images/mens-white-shirt.jpg",
--                         "stock": 50,
--                         "category": [
--                                 "men"
--                         ],
--                         "created_at": "2025-10-17T17:29:27+01:00",
--                         "updated_at": "2025-10-17T17:29:27+01:00",
--                         "version": 1
--                 },
--                 {
--                         "id": 12,
--                         "user_id": 0,
--                         "name": "Women's Floral Maxi Dress",
--                         "description": "A flowy maxi dress with vibrant floral patterns, ideal for summer outings.",
--                         "price": "$ 79.99",
--                         "image_url": "https://example.com/images/floral-maxi-dress.jpg",
--                         "stock": 30,
--                         "category": [
--                                 "women"
--                         ],
--                         "created_at": "2025-10-17T17:29:27+01:00",
--                         "updated_at": "2025-10-17T17:29:27+01:00",
--                         "version": 1
--                 },
--                 {
--                         "id": 13,
--                         "user_id": 0,
--                         "name": "Unisex Black Hoodie",
--                         "description": "Comfortable black hoodie with a minimalist design, suitable for all genders.",
--                         "price": "$ 39.99",
--                         "image_url": "https://example.com/images/black-hoodie.jpg",
--                         "stock": 100,
--                         "category": [
--                                 "unisex"
--                         ],
--                         "created_at": "2025-10-17T17:29:27+01:00",
--                         "updated_at": "2025-10-17T17:29:27+01:00",
--                         "version": 1
--                 },
--                 {
--                         "id": 14,
--                         "user_id": 0,
--                         "name": "Men's Leather Jacket",
--                         "description": "Stylish black leather jacket with a modern fit, perfect for cooler weather.",
--                         "price": "$ 129.99",
--                         "image_url": "https://example.com/images/mens-leather-jacket.jpg",
--                         "stock": 20,
--                         "category": [
--                                 "men"
--                         ],
--                         "created_at": "2025-10-17T17:29:27+01:00",
--                         "updated_at": "2025-10-17T17:29:27+01:00",
--                         "version": 1
--                 },
--                 {
--                         "id": 15,
--                         "user_id": 0,
--                         "name": "Women's High-Waist Jeans",
--                         "description": "Trendy high-waist blue jeans with a slim fit, made from stretch denim.",
--                         "price": "$ 59.99",
--                         "image_url": "https://example.com/images/high-waist-jeans.jpg",
--                         "stock": 40,
--                         "category": [
--                                 "women"
--                         ],
--                         "created_at": "2025-10-17T17:29:27+01:00",
--                         "updated_at": "2025-10-17T17:29:27+01:00",
--                         "version": 1
--                 },
--                 {
--                         "id": 16,
--                         "user_id": 0,
--                         "name": "Unisex Sneakers",
--                         "description": "White canvas sneakers with a durable sole, great for casual wear.",  
--                         "price": "$ 69.99",
--                         "image_url": "https://example.com/images/unisex-sneakers.jpg",
--                         "stock": 60,
--                         "category": [
--                                 "unisex"
--                         ],
--                         "created_at": "2025-10-17T17:29:27+01:00",
--                         "updated_at": "2025-10-17T17:29:27+01:00",
--                         "version": 1
--                 },
--                 {
--                         "id": 17,
--                         "user_id": 0,
--                         "name": "Men's Slim Fit Chinos",
--                         "description": "Versatile navy chinos with a slim fit, suitable for both casual and semi-formal settings.",
--                         "price": "$ 54.99",
--                         "image_url": "https://example.com/images/mens-chinos.jpg",
--                         "stock": 45,
--                         "category": [
--                                 "men"
--                         ],
--                         "created_at": "2025-10-17T17:29:27+01:00",
--                         "updated_at": "2025-10-17T17:29:27+01:00",
--                         "version": 1
--                 },
--                 {
--                         "id": 18,
--                         "user_id": 0,
--                         "name": "Women's Silk Scarf",
--                         "description": "Elegant silk scarf with a geometric pattern, perfect as an accessory.",
--                         "price": "$ 29.99",
--                         "image_url": "https://example.com/images/silk-scarf.jpg",
--                         "stock": 25,
--                         "category": [
--                                 "women"
--                         ],
--                         "created_at": "2025-10-17T17:29:27+01:00",
--                         "updated_at": "2025-10-17T17:29:27+01:00",
--                         "version": 1
--                 },
--                 {
--                         "id": 19,
--                         "user_id": 0,
--                         "name": "Unisex Baseball Cap",
--                         "description": "Adjustable black baseball cap with a simple logo, ideal for everyday 
-- use.",
--                         "price": "$ 19.99",
--                         "image_url": "https://example.com/images/baseball-cap.jpg",
--                         "stock": 80,
--                         "category": [
--                                 "unisex"
--                         ],
--                         "created_at": "2025-10-17T17:29:27+01:00",
--                         "updated_at": "2025-10-17T17:29:27+01:00",
--                 {
--                         "id": 20,
--                         "user_id": 0,
--                         "name": "Men's Wool Coat",
--                         "description": "Warm wool coat in charcoal gray, designed for winter elegance.",
--                         "price": "$ 149.99",                                           inter elegance.",     
--                         "image_url": "https://example.com/images/mens-wool-coat.jpg",  
--                         "stock": 15,
--                         "category": [
--                                 "men"
--                         ],
--                         "created_at": "2025-10-17T17:29:27+01:00",
--                         "updated_at": "2025-10-17T17:29:27+01:00",
--                         "version": 1
--                 }
--         ]
-- }