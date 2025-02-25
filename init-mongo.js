// init-mongo.js

// Start in the 'admin' database to query existing databases
db = db.getSiblingDB('admin');

// Check if 'ledger' database exists
var dbList = db.adminCommand({ listDatabases: 1 }).databases;
var ledgerExists = dbList.some(function (dbInfo) { return dbInfo.name === 'ledger'; });

if (!ledgerExists) {
    print("Database 'ledger' does not exist, creating it...");
}

// Switch to the 'ledger' database (creates it implicitly if it doesn't exist when we act on it)
db = db.getSiblingDB('ledger');

// Check if 'transactions' collection exists
var collections = db.getCollectionNames();
var transactionsExists = collections.indexOf('transactions') !== -1;

if (!transactionsExists) {
    print("Collection 'transactions' does not exist, creating it...");
    db.createCollection('transactions');
} else {
    print("Collection 'transactions' already exists, skipping creation...");
}
