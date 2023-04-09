pragma solidity ^0.8.0;

contract Simple {
    
    bytes16 public value;
    
    function setValue(bytes16 v) public {
        value = v;
    }
    
    function getValue() public view returns (bytes16) {
        return value;
    }

}