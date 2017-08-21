contract Escrow {

  address public buyer;
  address public seller;
  address public arbiter;

  // Function with same name as contract is the constructor.
  // Will run once and only once, at the creation of the contract.
  function Escrow(address _seller, address _arbiter) {
    // msg object = contains info about txn calling into current contract.
    // in constructor, contains info on txn which created contract.
    buyer = msg.sender;
    seller = _seller;
    arbiter = _arbiter;
  }
}

