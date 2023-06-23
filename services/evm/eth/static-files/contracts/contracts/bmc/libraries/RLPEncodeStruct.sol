// SPDX-License-Identifier: Apache-2.0
pragma solidity >=0.8.0;
pragma abicoder v2;

import "@iconfoundation/btp2-solidity-library/contracts/utils/RLPEncode.sol";
import "./Types.sol";

library RLPEncodeStruct {
    using RLPEncode for bytes;
    using RLPEncode for string;
    using RLPEncode for uint256;
    using RLPEncode for int256;
    using RLPEncode for address;
    using RLPEncode for bytes[];

    uint8 internal constant LIST_SHORT_START = 0xc0;
    uint8 internal constant LIST_LONG_START = 0xf7;

    function encodeFeeInfo(Types.FeeInfo memory _fi)
    internal
    pure
    returns (bytes memory)
    {
        bytes[] memory _items = new bytes[](_fi.values.length);
        for (uint256 i = 0; i < _fi.values.length; i++) {
            _items[i] = _fi.values[i].encodeUint();
        }
        bytes memory _rlp = abi.encodePacked(
            _fi.network.encodeString(),
            _items.encodeList()
        );
        return _rlp.encodeList();
    }

    function encodeBMCMessage(Types.BMCMessage memory _bs)
        internal
        pure
        returns (bytes memory)
    {
        bytes memory _rlp =
            abi.encodePacked(
                _bs.msgType.encodeString(),
                _bs.payload.encodeBytes());
        return _rlp.encodeList();
    }

    function encodeBTPMessage(Types.BTPMessage memory _bm)
        internal
        pure
        returns (bytes memory)
    {
        bytes memory _rlp =
            abi.encodePacked(
                _bm.src.encodeString(),
                _bm.dst.encodeString(),
                _bm.svc.encodeString(),
                _bm.sn.encodeInt(),
                _bm.message.encodeBytes(),
                _bm.nsn.encodeInt(),
                encodeFeeInfo(_bm.feeInfo)
            );
        return _rlp.encodeList();
    }

    function encodeResponseMessage(Types.ResponseMessage memory _res)
        internal
        pure
        returns (bytes memory)
    {
        bytes memory _rlp =
            abi.encodePacked(
                _res.code.encodeUint(),
                _res.message.encodeString()
            );
        return _rlp.encodeList();
    }

    function encodeInitMessage(string[] memory _links)
        internal
        pure
        returns (bytes memory)
    {
        bytes[] memory _items = new bytes[](_links.length);
        for (uint256 i = 0; i < _links.length; i++) {
            _items[i] = _links[i].encodeString();
        }
        return _items.encodeList().encodeList();
    }

    function encodePropagateMessage(string memory _link)
        internal
        pure
        returns (bytes memory)
    {
        return _link.encodeString().encodeList();
    }

    function encodeClaimMessage(Types.ClaimMessage memory _cm)
    internal
    pure
    returns (bytes memory)
    {
        bytes memory _rlp =
        abi.encodePacked(
            _cm.amount.encodeUint(),
            _cm.receiver.encodeString());
        return _rlp.encodeList();
    }
}
