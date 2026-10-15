-- 0001_init.sql — snapshot skema dev (SQLite, dari AutoMigrate).
-- Untuk Postgres/Supabase: review tipe (AUTOINCREMENT->GENERATED, DATETIME->TIMESTAMPTZ) sebelum dipakai prod.

CREATE TABLE `admin_audits` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`admin_id` integer NOT NULL,`action` text NOT NULL,`target_type` text,`target_id` integer,`note` text);

CREATE TABLE `banners` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`title` text NOT NULL,`subtitle` text,`cta_label` text,`cta_href` text,`sort` integer,`active` numeric DEFAULT true);

CREATE TABLE `cart_items` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`user_id` integer NOT NULL,`product_id` integer NOT NULL);

CREATE TABLE `chat_messages` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`thread_id` integer NOT NULL,`sender_id` integer NOT NULL,`body` text NOT NULL);

CREATE TABLE `chat_threads` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`order_id` integer,`product_id` integer NOT NULL,`buyer_id` integer NOT NULL,`seller_id` integer NOT NULL);

CREATE TABLE `condition_scores` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`product_id` integer NOT NULL,`score` integer NOT NULL,`method` text,`note` text,`checked_at` datetime NOT NULL);

CREATE TABLE `confirm_codes` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`order_id` integer NOT NULL,`role` text NOT NULL,`code` text NOT NULL,`expires_at` datetime NOT NULL,`used_at` datetime,`attempts` integer DEFAULT 0);

CREATE TABLE `defects` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`product_id` integer NOT NULL,`photo_id` integer,`x` integer NOT NULL,`y` integer NOT NULL,`label` text NOT NULL,`sort` integer);

CREATE TABLE `donation_points` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`name` text NOT NULL,`lat` real,`lng` real,`address_label` text,`active` numeric DEFAULT true);

CREATE TABLE `orders` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`buyer_id` integer NOT NULL,`seller_id` integer NOT NULL,`product_id` integer NOT NULL,`total` integer NOT NULL,`pay_method` text NOT NULL,`mode` text NOT NULL,`status` text NOT NULL);

CREATE TABLE `point_ledgers` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`user_id` integer NOT NULL,`order_id` integer,`delta` integer NOT NULL,`reason` text NOT NULL,`balance_after` integer NOT NULL);

CREATE TABLE `product_photos` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`product_id` integer NOT NULL,`slot` text NOT NULL,`url` text NOT NULL,`sort` integer);

CREATE TABLE `product_sizes` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`product_id` integer NOT NULL,`dada` real,`pinggang` real,`pinggul` real,`tinggi` real,`torso` real);

CREATE TABLE `products` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`seller_id` integer NOT NULL,`name` text NOT NULL,`category` text,`price` integer NOT NULL,`description` text,`lat` real,`lng` real,`location_label` text,`status` text DEFAULT "draft");

CREATE TABLE `redemptions` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`user_id` integer NOT NULL,`reward` text NOT NULL,`cost` integer NOT NULL,`status` text NOT NULL);

CREATE TABLE `reports` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`reporter_id` integer NOT NULL,`target_user_id` integer,`product_id` integer,`reason` text NOT NULL,`detail` text,`status` text DEFAULT "open");

CREATE TABLE `reviews` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`order_id` integer NOT NULL,`reviewer_id` integer NOT NULL,`seller_id` integer NOT NULL,`stars` integer NOT NULL,`body` text);

CREATE TABLE `size_preferences` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`user_id` integer NOT NULL,`dada` real,`pinggang` real,`pinggul` real,`tinggi` real,`torso` real);

CREATE TABLE `users` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`name` text NOT NULL,`email` text NOT NULL,`password_hash` text,`avatar_url` text,`handle` text,`lat` real,`lng` real,`location_label` text,`level` text DEFAULT "Pemula Hijau",`email_verified` numeric,`is_admin` numeric DEFAULT false);
