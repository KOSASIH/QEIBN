pragma solidity ^0.8.0;

import "https://github.com/OpenZeppelin/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "https://github.com/OpenZeppelin/openzeppelin-solidity/contracts/token/ERC20/SafeERC20.sol";

contract QEIBNContract {
    using SafeMath for uint256;
    using SafeERC20 for ERC20;

    // Mapping of user addresses to their QEIBN accounts
    mapping (address => QEIBNAccount) public accounts;

    // Mapping of QEIBN account IDs to their corresponding user addresses
    mapping (uint256 => address) public accountOwners;

    // QEIBN account struct
    struct QEIBNAccount {
        uint256 id;
        uint256 balance;
        uint256[] transactionHistory;
    }

    // Event emitted when a new QEIBN account is created
    event NewAccount(address indexed user, uint256 accountId);

    // Event emitted when a QEIBN account is updated
    event UpdateAccount(address indexed user, uint256 accountId, uint256 newBalance);

    // Event emitted when a transaction is processed
    event TransactionProcessed(address indexed from, address indexed to, uint256 amount);

    // Function to create a new QEIBN account
    function createAccount() public {
        uint256 newAccountId = accounts[msg.sender].id++;
        accounts[msg.sender] = QEIBNAccount(newAccountId, 0, new uint256[](0));
        accountOwners[newAccountId] = msg.sender;
        emit NewAccount(msg.sender, newAccountId);
    }

    // Function to deposit funds into a QEIBN account
    function deposit(uint256 amount) public {
        require(accounts[msg.sender].id != 0, "Account does not exist");
        accounts[msg.sender].balance = accounts[msg.sender].balance.add(amount);
        emit UpdateAccount(msg.sender, accounts[msg.sender].id, accounts[msg.sender].balance);
    }

    // Function to withdraw funds from a QEIBN account
    function withdraw(uint256 amount) public {
        require(accounts[msg.sender].id != 0, "Account does not exist");
        require(accounts[msg.sender].balance >= amount, "Insufficient balance");
        accounts[msg.sender].balance = accounts[msg.sender].balance.sub(amount);
        emit UpdateAccount(msg.sender, accounts[msg.sender].id, accounts[msg.sender].balance);
    }

    // Function to process a transaction between two QEIBN accounts
    function processTransaction(address to, uint256 amount) public {
        require(accounts[msg.sender].id != 0, "Account does not exist");
        require(accounts[to].id != 0, "Recipient account does not exist");
        require(accounts[msg.sender].balance >= amount, "Insufficient balance");
        accounts[msg.sender].balance = accounts[msg.sender].balance.sub(amount);
        accounts[to].balance = accounts[to].balance.add(amount);
        emit TransactionProcessed(msg.sender, to, amount);
    }

    // Function to get the balance of a QEIBN account
    function getBalance() public view returns (uint256) {
        return accounts[msg.sender].balance;
    }

    // Function to get the transaction history of a QEIBN account
    function getTransactionHistory() public view returns (uint256[] memory) {
        return accounts[msg.sender].transactionHistory;
    }
}
