# A financial service

## Description
A financial service api that exposes the ability to increase or decrease a user's total balance. A user has 1 available balances

1. Savings account
    
    * Represents the total balance of a users savings account

## Functionality
1. Increase a users chequing balance
2. Decrease a users chequing balance

This service will be written in go. It will require a postgres database. A table will store transactions. Balances will be calculated by summing the total of all transactions.

The functionality this service will expose is:
1. Changing balance
2. Check balance
