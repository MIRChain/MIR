pragma solidity ^0.8.14;

pragma experimental ABIEncoderV2;

contract Simple {
    
    uint public value;
    
    function setValue(uint v) public {
        value = v;
    }
    
    function getValue() public view returns (uint) {
        return value;
    }

}