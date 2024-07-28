pragma solidity ^0.8.0;
contract GasEstimation {
    function PureRevert() public { revert(); }
    function Revert() public { revert("revert reason");}
    function OOG() public { for (uint i = 0; ; i++) {}}
    function Assert() public { assert(false);}
    function Valid() public {}
}