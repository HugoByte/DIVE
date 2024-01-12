import IconService from 'icon-sdk-js';
import {ethers} from 'hardhat';
import {IconNetwork, getBtpAddress} from "./icon";
import {XCall, DAppProxy} from "./icon";
import {BaseContract, BigNumber, ContractReceipt} from "ethers";
import {Deployments} from "./setup/config";
import {TypedEvent, TypedEventFilter} from "../typechain-types/common";

const {IconConverter} = IconService;
const config = process.env.CONFIG_FILE || "dive.json";
const deployments = Deployments.getDefault(config);

function sleep(millis: number) {
  return new Promise(resolve => setTimeout(resolve, millis));
}

async function waitEvent<TEvent extends TypedEvent>(
    ctr : BaseContract,
    filter: TypedEventFilter<TEvent>,
    fromBlock?: number
): Promise<Array<TEvent>> {
  let height = await ctr.provider.getBlockNumber();
  let next = height + 1;
  if (fromBlock && fromBlock < height) {
    height = fromBlock;
  }
  while (true) {
    for (;height < next;height++){
      const events = await ctr.queryFilter(filter, height);
      if (events.length > 0) {
        return events as Array<TEvent>;
      }
    }
    await sleep(1000);
    next = await ctr.provider.getBlockNumber() + 1;
  }
}

function filterEvent<TEvent extends TypedEvent>(
    ctr : BaseContract,
    filter: TypedEventFilter<TEvent>,
    receipt: ContractReceipt) : Array<TEvent> {
  const inf = ctr.interface;
  const address = ctr.address;
  const topics = filter.topics || [];
  if (receipt.events && typeof topics[0] === "string") {
    const fragment = inf.getEvent(topics[0]);
    return receipt.events
        .filter((evt) => {
          if (evt.address == address) {
            return topics.every((v, i) => {
              if (!v) {
                return true
              } else if (typeof v === "string") {
                return v == evt.topics[i]
              } else {
                return v.includes(evt.topics[i])
              }
            })
          }
          return false
        })
        .map((evt) => {
           return { args : inf.decodeEventLog(fragment, evt.data, evt.topics) } as TEvent
        });
  }
  return [];
}

function hexToString(data: string) {
  const hexArray = ethers.utils.arrayify(data);
  let msg = '';
  for (let i = 0; i < hexArray.length; i++) {
    msg += String.fromCharCode(hexArray[i]);
  }
  return msg;
}

function isIconChain(chain: any) {
  return chain.network.includes('icon');
}

function isEVMChain(chain: any) {
  return chain.network.includes('hardhat') || chain.network.includes('eth');
}

async function sendMessageFromDApp(src: string, srcChain: any, dstChain: any, msg: string, srcContracts: any, dstContracts: any,
                                   rollback?: string) {
  const isRollback = rollback ? true : false;
  if (isIconChain(srcChain)) {
    const iconNetwork = IconNetwork.getNetwork(srcChain, src);
    const xcallSrc = new XCall(iconNetwork, srcContracts.xcall);
    const fee = await xcallSrc.getFee(dstChain.network, isRollback);
    console.log('fee=' + fee);

    const dappSrc = new DAppProxy(iconNetwork, srcContracts.dapp);
    const to = getBtpAddress(dstChain.network, dstContracts.dapp);
    const data = IconConverter.toHex(msg);
    const rbData = rollback ? IconConverter.toHex(rollback) : undefined;

    return await dappSrc.sendMessage(to, data, rbData, fee)
      .then((txHash) => dappSrc.getTxResult(txHash))
      .then((receipt) => {
        if (receipt.status != 1) {
          throw new Error(`DApp: failed to sendMessage: ${receipt.txHash}`);
        }
        return receipt;
      });
  } else if (isEVMChain(srcChain)) {
    const xcallSrc = await ethers.getContractAt('CallService', srcContracts.xcall);
    const fee = await xcallSrc.getFee(dstChain.network, isRollback);
    console.log('fee=' + fee);

    const dappSrc = await ethers.getContractAt('DAppProxySample', srcContracts.dapp);
    const to = getBtpAddress(dstChain.network, dstContracts.dapp);
    const data = IconConverter.toHex(msg);
    const rbData = rollback ? IconConverter.toHex(rollback) : "0x";

    return await dappSrc.sendMessage(to, data, rbData, {value: fee})
      .then((tx) => tx.wait(1))
      .then((receipt) => {
        if (receipt.status != 1) {
          throw new Error(`DApp: failed to sendMessage: ${receipt.transactionHash}`);
        }
        return receipt;
      })
  } else {
    throw new Error(`DApp: unknown source chain: ${srcChain}`);
  }
}

async function verifyCallMessageSent(src: string, srcChain: any, receipt: any, srcContracts: any) {
  let event;
  if (isIconChain(srcChain)) {
    const iconNetwork = IconNetwork.getNetwork(srcChain, src);
    const xcallSrc = new XCall(iconNetwork, srcContracts.xcall);
    const logs = xcallSrc.filterEvent(receipt.eventLogs,
        'CallMessageSent(Address,str,int,int)', xcallSrc.address);
    if (logs.length == 0) {
      throw new Error(`DApp: could not find event: "CallMessageSent"`);
    }
    console.log(logs);
    const indexed = logs[0].indexed || [];
    const data = logs[0].data || [];
    event = {
      _from: indexed[1],
      _to: indexed[2],
      _sn: BigNumber.from(indexed[3]),
      _nsn: BigNumber.from(data[0])
    };
  } else if (isEVMChain(srcChain)) {
    const xcallSrc = await ethers.getContractAt('CallService', srcContracts.xcall);
    const logs = filterEvent(xcallSrc, xcallSrc.filters.CallMessageSent(), receipt);
    if (logs.length == 0) {
      throw new Error(`DApp: could not find event: "CallMessageSent"`);
    }
    console.log(logs);
    event = logs[0].args;
  } else {
    throw new Error(`DApp: unknown source chain: ${srcChain}`);
  }
  console.log(`serialNum=${event._sn}`);
  return event._sn;
}

async function checkCallMessage(dst: string, srcChain: any, dstChain: any, sn: BigNumber, msg: string, srcContracts:any, dstContracts:any) {
  if (isEVMChain(dstChain)) {
    const xcallDst = await ethers.getContractAt('CallService', dstContracts.xcall);
    const filterCM = xcallDst.filters.CallMessage(
      getBtpAddress(srcChain.network, srcContracts.dapp),
      dstContracts.dapp,
      sn
    )
    const logs = await waitEvent(xcallDst, filterCM);
    if (logs.length == 0) {
      throw new Error(`DApp: could not find event: "CallMessage"`);
    }
    console.log(logs[0]);
    const calldata = hexToString(logs[0].args._data)
    if (msg !== calldata) {
      throw new Error(`DApp: the calldata is different from the sent message`);
    }
    return {
      _reqId: logs[0].args._reqId,
      _data: logs[0].args._data
    };
  } else if (isIconChain(dstChain)) {
    const iconNetwork = IconNetwork.getNetwork(dstChain, dst);
    const xcallDst = new XCall(iconNetwork, dstContracts.xcall);
    const {block, events} = await xcallDst.waitEvent("CallMessage(str,str,int,int,bytes)");
    if (events.length == 0) {
      throw new Error(`DApp: could not find event: "CallMessage"`);
    }
    console.log(events[0]);
    const indexed = events[0].indexed || [];
    const data = events[0].data || [];
    const event = {
      _from: indexed[1],
      _to: indexed[2],
      _sn: BigNumber.from(indexed[3]),
      _reqId: BigNumber.from(data[0]),
      _data: data[1]
    };
    if (!sn.eq(event._sn)) {
      throw new Error(`DApp: serial number mismatch (${sn} != ${event._sn})`);
    }
    const calldata = hexToString(event._data)
    if (msg !== calldata) {
      throw new Error(`DApp: the calldata is different from the sent message`);
    }
    return {
      _reqId: event._reqId,
      _data: event._data
    };
  } else {
    throw new Error(`DApp: unknown destination chain: ${dstChain}`);
  }
}

async function invokeExecuteCall(dst: string, dstChain: any, reqId: BigNumber, data: string, dstContracts:any) {
  if (isEVMChain(dstChain)) {
    const xcallDst = await ethers.getContractAt('CallService', dstContracts.xcall);
    return await xcallDst.executeCall(reqId, data, {gasLimit: 300000})
      .then((tx) => tx.wait(1))
      .then((receipt) => {
        if (receipt.status != 1) {
          throw new Error(`DApp: failed to executeCall: ${receipt.transactionHash}`);
        }
        return receipt;
      })
  } else if (isIconChain(dstChain)) {
    const iconNetwork = IconNetwork.getNetwork(dstChain, dst);
    const xcallDst = new XCall(iconNetwork, dstContracts.xcall);
    return await xcallDst.executeCall(reqId.toHexString(), data)
      .then((txHash) => xcallDst.getTxResult(txHash))
      .then((receipt) => {
        if (receipt.status != 1) {
          throw new Error(`DApp: failed to executeCall: ${receipt.txHash}`);
        }
        return receipt;
      });
  } else {
    throw new Error(`DApp: unknown destination chain: ${dstChain}`);
  }
}

async function verifyReceivedMessage(dst: string, dstChain: any, receipt: any, msg: string, dstContracts:any) {
  let event;
  if (isEVMChain(dstChain)) {
    const dappDst = await ethers.getContractAt('DAppProxySample', dstContracts.dapp);
    const logs = filterEvent(dappDst, dappDst.filters.MessageReceived(), receipt);
    if (logs.length == 0) {
      throw new Error(`DApp: could not find event: "MessageReceived"`);
    }
    console.log(logs);
    event = logs[0].args;
  } else if (isIconChain(dstChain)) {
    const iconNetwork = IconNetwork.getNetwork(dstChain,dst);
    const dappDst = new DAppProxy(iconNetwork, dstContracts.dapp);
    const logs = dappDst.filterEvent(receipt.eventLogs,'MessageReceived(str,bytes)', dappDst.address);
    if (logs.length == 0) {
      throw new Error(`DApp: could not find event: "MessageReceived"`);
    }
    console.log(logs);
    const data = logs[0].data || [];
    event = {_from: data[0], _data: data[1]}
  } else {
    throw new Error(`DApp: unknown destination chain: ${dstChain}`);
  }

  const receivedMsg = hexToString(event._data)
  console.log(`From: ${event._from}`);
  console.log(`Data: ${event._data}`);
  console.log(`Msg: ${receivedMsg}`);
  if (msg !== receivedMsg) {
    throw new Error(`DApp: received message is different from the sent message`);
  }
}

async function checkCallExecuted(dst: string, dstChain: any, receipt: any, reqId: BigNumber, expectRevert: boolean, dstContracts:any) {
  let event;
  if (isEVMChain(dstChain)) {
    const xcallDst = await ethers.getContractAt('CallService', dstContracts.xcall);
    const logs = filterEvent(xcallDst, xcallDst.filters.CallExecuted(), receipt);
    if (logs.length == 0) {
      throw new Error(`DApp: could not find event: "CallExecuted"`);
    }
    console.log(logs);
    event = logs[0].args;
  } else if (isIconChain(dstChain)) {
    const iconNetwork = IconNetwork.getNetwork(dstChain, dst);
    const xcallDst = new XCall(iconNetwork, dstContracts.xcall);
    const logs = xcallDst.filterEvent(receipt.eventLogs,'CallExecuted(int,int,str)', xcallDst.address);
    if (logs.length == 0) {
      throw new Error(`DApp: could not find event: "CallExecuted"`);
    }
    console.log(logs);
    const indexed = logs[0].indexed || [];
    const data = logs[0].data || [];
    event = {
      _reqId: BigNumber.from(indexed[1]),
      _code: BigNumber.from(data[0]),
      _msg: data[1]
    }
  } else {
    throw new Error(`DApp: unknown destination chain: ${dstChain}`);
  }
  if (!reqId.eq(event._reqId) ||
    (expectRevert && event._code.isZero()) || (!expectRevert && !event._code.isZero())) {
    throw new Error(`DApp: not the expected execution result`);
  }
}

async function checkResponseMessage(src: string, srcChain: any, sn: BigNumber, expectRevert: boolean, srcContracts:any) {
  let event, blockNum;
  if (isIconChain(srcChain)) {
    const iconNetwork = IconNetwork.getNetwork(srcChain, src);
    const xcallSrc = new XCall(iconNetwork, srcContracts.xcall);
    const {block, events} = await xcallSrc.waitEvent("ResponseMessage(int,int,str)");
    if (events.length == 0) {
      throw new Error(`DApp: could not find event: "ResponseMessage"`);
    }
    console.log(events);
    const indexed = events[0].indexed || [];
    const data = events[0].data || [];
    event = {
      _sn: BigNumber.from(indexed[1]),
      _code: BigNumber.from(data[0]),
      _msg: data[1]
    }
    blockNum = block.height;
  } else if (isEVMChain(srcChain)) {
    const xcallSrc = await ethers.getContractAt('CallService', srcContracts.xcall);
    const events = await waitEvent(xcallSrc, xcallSrc.filters.ResponseMessage());
    if (events.length == 0) {
      throw new Error(`DApp: could not find event: "ResponseMessage"`);
    }
    console.log(events)
    event = events[0].args;
    blockNum = (await events[0].getBlock()).number;
  } else {
    throw new Error(`DApp: unknown source chain: ${srcChain}`);
  }
  if (!sn.eq(event._sn)) {
    throw new Error(`DApp: received serial number (${event._sn}) is different from the sent one (${sn})`);
  }
  if ((expectRevert && event._code.isZero()) || (!expectRevert && !event._code.isZero())) {
    throw new Error(`DApp: not the expected response message`);
  }
  return blockNum;
}

async function checkRollbackMessage(src: string, srcChain: any, blockNum: number, srcContracts:any) {
  if (isEVMChain(srcChain)) {
    const xcallSrc = await ethers.getContractAt('CallService', srcContracts.xcall);
    const logs = await waitEvent(xcallSrc, xcallSrc.filters.RollbackMessage(), blockNum);
    if (logs.length == 0) {
      throw new Error(`DApp: could not find event: "RollbackMessage"`);
    }
    console.log(logs[0]);
    return logs[0].args._sn;
  } else if (isIconChain(srcChain)) {
    const iconNetwork = IconNetwork.getNetwork(srcChain, src);
    const xcallSrc = new XCall(iconNetwork, srcContracts.xcall);
    const {block, events} = await xcallSrc.waitEvent("RollbackMessage(int)", blockNum);
    if (events.length == 0) {
      throw new Error(`DApp: could not find event: "RollbackMessage"`);
    }
    console.log(events[0]);
    const indexed = events[0].indexed || [];
    return BigNumber.from(indexed[1]);
  } else {
    throw new Error(`DApp: unknown source chain: ${srcChain}`);
  }
}

async function invokeExecuteRollback(src: string, srcChain: any, sn: BigNumber, srcContracts:any) {
  if (isEVMChain(srcChain)) {
    const xcallSrc = await ethers.getContractAt('CallService', srcContracts.xcall);
    return await xcallSrc.executeRollback(sn, {gasLimit: 300000})
      .then((tx) => tx.wait(1))
      .then((receipt) => {
        if (receipt.status != 1) {
          throw new Error(`DApp: failed to executeRollback: ${receipt.transactionHash}`);
        }
        return receipt;
      });
  } else if (isIconChain(srcChain)) {
    const iconNetwork = IconNetwork.getNetwork(srcChain, src);
    const xcallSrc = new XCall(iconNetwork, srcContracts.xcall);
    return await xcallSrc.executeRollback(sn.toHexString())
      .then((txHash) => xcallSrc.getTxResult(txHash))
      .then((receipt) => {
        if (receipt.status != 1) {
          throw new Error(`DApp: failed to executeRollback: ${receipt.txHash}`);
        }
        return receipt;
      });
  } else {
    throw new Error(`DApp: unknown source chain: ${srcChain}`);
  }
}

async function verifyRollbackDataReceivedMessage(src: string, srcChain: any, receipt: any, srcContracts:any, rollback: string | undefined) {
  let event;
  if (isEVMChain(srcChain)) {
    const dappSrc = await ethers.getContractAt('DAppProxySample', srcContracts.dapp);
    const logs = filterEvent(dappSrc, dappSrc.filters.RollbackDataReceived(), receipt);
    if (logs.length == 0) {
      throw new Error(`DApp: could not find event: "RollbackDataReceived"`);
    }
    console.log(logs)
    event = logs[0].args;
  } else if (isIconChain(srcChain)) {
    const iconNetwork = IconNetwork.getNetwork(srcChain, src);
    const dappSrc = new DAppProxy(iconNetwork, srcContracts.dapp);
    const logs = dappSrc.filterEvent(receipt.eventLogs,"RollbackDataReceived(str,int,bytes)", dappSrc.address);
    if (logs.length == 0) {
      throw new Error(`DApp: could not find event: "RollbackDataReceived"`);
    }
    console.log(logs)
    const data = logs[0].data || [];
    event = {_from: data[0], _ssn: data[1], _rollback: data[2]}
  } else {
    throw new Error(`DApp: unknown source chain: ${srcChain}`);
  }

  const receivedRollback = hexToString(event._rollback)
  console.log(`From: ${event._from}`);
  console.log(`Ssn: ${event._ssn}`);
  console.log(`Data: ${event._rollback}`);
  console.log(`Rollback: ${receivedRollback}`);
  if (rollback !== receivedRollback) {
    throw new Error(`DApp: received rollback is different from the sent data`);
  }
}

async function checkRollbackExecuted(src: string, srcChain: any, receipt: any, sn: BigNumber, srcContracts:any) {
  let event;
  if (isIconChain(srcChain)) {
    const iconNetwork = IconNetwork.getNetwork(srcChain, src);
    const xcallSrc = new XCall(iconNetwork, srcContracts.xcall);
    const logs = xcallSrc.filterEvent(receipt.eventLogs, "RollbackExecuted(int,int,str)", xcallSrc.address);
    if (logs.length == 0) {
      throw new Error(`DApp: could not find event: "RollbackExecuted"`);
    }
    console.log(logs);
    const indexed = logs[0].indexed || [];
    const data = logs[0].data || [];
    event = {
      _sn: BigNumber.from(indexed[1]),
      _code: BigNumber.from(data[0]),
      _msg: data[1]
    }
  } else if (isEVMChain(srcChain)) {
    const xcallSrc = await ethers.getContractAt('CallService', srcContracts.xcall);
    const logs = filterEvent(xcallSrc, xcallSrc.filters.RollbackExecuted(), receipt);
    if (logs.length == 0) {
      throw new Error(`DApp: could not find event: "RollbackExecuted"`);
    }
    console.log(logs)
    event = logs[0].args;
  } else {
    throw new Error(`DApp: unknown source chain: ${srcChain}`);
  }
  if (!sn.eq(event._sn)) {
    throw new Error(`DApp: received serial number (${event._sn}) is different from the sent one (${sn})`);
  }
  if (!event._code.isZero()) {
    throw new Error(`DApp: not the expected execution result`);
  }
}

async function sendCallMessage(src: string, dst: string, msgData?: string, needRollback?: boolean) {
  const srcChain = deployments.get(src);
  const dstChain = deployments.get(dst);
  const srcContracts = deployments.getContracts(src);
  const dstContracts = deployments.getContracts(dst);

  const testName = sendCallMessage.name + (needRollback ? "WithRollback" : "");
  console.log(`\n### ${testName}: ${src} => ${dst}`);
  if (!msgData) {
    msgData = `${testName}_${src}_${dst}`;
  }
  const rollbackData = needRollback ? `ThisIsRollbackMessage_${src}_${dst}` : undefined;
  const expectRevert = (msgData === "revertMessage");
  let step = 1;

  console.log(`[${step++}] send message from DApp`);
  const sendMessageReceipt = await sendMessageFromDApp(src, srcChain, dstChain, msgData, srcContracts, dstContracts,rollbackData);
  const sn = await verifyCallMessageSent(src, srcChain, sendMessageReceipt, srcContracts);

  console.log(`[${step++}] check CallMessage event on ${dst} chain`);
  // #TODO: Update endpoint in hardhat config during runtime from services.json
  const callMsgEvent = await checkCallMessage(dst, srcChain, dstChain, sn, msgData, srcContracts, dstContracts);
  const reqId = callMsgEvent._reqId;
  const callData = callMsgEvent._data;

  console.log(`[${step++}] invoke executeCall with reqId=${reqId}`);
  const executeCallReceipt = await invokeExecuteCall(dst, dstChain, reqId, callData, dstContracts);

  if (!expectRevert) {
    console.log(`[${step++}] verify the received message`);
    await verifyReceivedMessage(dst, dstChain, executeCallReceipt, msgData, dstContracts);
  }
  console.log(`[${step++}] check CallExecuted event on ${dst} chain`);
  await checkCallExecuted(dst, dstChain, executeCallReceipt, reqId, expectRevert, dstContracts);

  if (needRollback) {
    console.log(`[${step++}] check ResponseMessage event on ${src} chain`);
    const responseHeight = await checkResponseMessage(src, srcChain, sn, expectRevert, srcContracts);

    if (expectRevert) {
      console.log(`[${step++}] check RollbackMessage event on ${src} chain`);
      const sn = await checkRollbackMessage(src, srcChain, responseHeight, srcContracts);

      console.log(`[${step++}] invoke executeRollback with sn=${sn}`);
      const executeRollbackReceipt = await invokeExecuteRollback(src, srcChain, sn, srcContracts);

      console.log(`[${step++}] verify rollback data received message`);
      await verifyRollbackDataReceivedMessage(src, srcChain, executeRollbackReceipt, srcContracts, rollbackData);

      console.log(`[${step++}] check RollbackExecuted event on ${src} chain`);
      await checkRollbackExecuted(src, srcChain, executeRollbackReceipt, sn, srcContracts);
    }
  }
}

async function show_banner() {

  const banner = `
       ___           __
  ___ |__ \\___  ____/ /__  ____ ___  ____
 / _ \\__/ / _ \\/ __  / _ \\/ __ \`__ \\/ __ \\
/  __/ __/  __/ /_/ /  __/ / / / / / /_/ /
\\___/____\\___/\\__,_/\\___/_/ /_/ /_/\\____/
`;
  console.log(banner);
}

const SRC = deployments.getSrc();
const DST = deployments.getDst();


show_banner()
  .then(() => sendCallMessage(SRC, DST))
  .then(() => sendCallMessage(DST, SRC))
  .then(() => sendCallMessage(SRC, DST, "checkSuccessResponse", true))
  .then(() => sendCallMessage(DST, SRC, "checkSuccessResponse", true))
  .then(() => sendCallMessage(SRC, DST, "revertMessage", true))
  .then(() => sendCallMessage(DST, SRC, "revertMessage", true))
  .catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
