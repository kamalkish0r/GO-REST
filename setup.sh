#!/bin/bash

echo "This script will create an .env file and set up DATABASE_URL."

# Prompt user for the database URL
read -p "Enter your database URL (e.g., postgresql://username:password@localhost:5432/dbname?sslmode=disable): " db_url

# Create .env file
echo "DATABASE_URL=\"$db_url\"" > .env

echo "Environment variables set up successfully in .env file."
