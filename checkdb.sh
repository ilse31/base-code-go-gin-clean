#!/bin/bash

# Set default values
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-}
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-postgres}

echo "Testing database connection with:"
echo "User: $DB_USER"
echo "Host: $DB_HOST"
echo "Port: $DB_PORT"
echo "Database: $DB_NAME"

# Check if psql is installed
if ! command -v psql &> /dev/null; then
    echo "Error: psql command not found. Please install PostgreSQL client tools."
    exit 1
fi

# Test connection
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1 AS test;"

if [ $? -eq 0 ]; then
    echo "✅ Database connection successful!"
else
    echo "❌ Failed to connect to database"
    echo "Tried to connect to: postgres://$DB_USER:****@$DB_HOST:$DB_PORT/$DB_NAME"
    echo ""
    echo "Troubleshooting steps:"
    echo "1. Make sure PostgreSQL is running"
    echo "2. Check your connection parameters"
    echo "3. Verify the database exists and the user has access"
    echo "4. Check if your firewall allows connections to port $DB_PORT"
    echo ""
    echo "You can set environment variables like this:"
    echo "  export DB_USER=your_username"
    echo "  export DB_PASSWORD=your_password"
    echo "  export DB_HOST=your_host"
    echo "  export DB_PORT=5432"
    echo "  export DB_NAME=your_database"
    echo "  ./checkdb.sh"
fi
