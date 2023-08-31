import IconService, {
  BigNumber,
  Block,
  ConfirmedTransaction,
  KeyStore,
  TransactionResult,
} from "icon-sdk-js";
import { GetArchwayChainInfo, GetCosmosContracts, GetIconChainInfo, GetIconContracts } from "./helper";

const {
  IconBuilder,
  IconConverter,
  SignedTransaction,
  HttpProvider,
  IconWallet,
} = IconService;

export class EventLog {
  scoreAddress: string | undefined;
  indexed: string[] | undefined;
  data: string[] | undefined;
}

const { CallTransactionBuilder, CallBuilder } = IconBuilder;

const ICON_RPC_URL = GetIconChainInfo("endpoint")
const NID = GetIconChainInfo("nid")
const ICON_XCALL = GetIconContracts("xcall");
const ICON_DAPP = GetIconContracts("dapp");
const NETWORK_LABEL_DESTINATION = GetArchwayChainInfo("network")
const DESTINATION_DAPP = GetCosmosContracts("dapp");

const callMessageSentSignature = "CallMessageSent(Address,str,int)"
const callMessageSignature = "CallMessage(str,str,int,int,bytes)"
const callExecutedSignature = "CallExecuted(int,int,str)"
const responseMessageSignature = "ResponseMessage(int,int)"
const rollbackMessageSignature = "RollbackMessage(int)"
const rollbackExecutedSignature = "RollbackExecuted(int)"

const HTTP_PROVIDER = new HttpProvider(ICON_RPC_URL);
const ICON_SERVICE = new IconService(HTTP_PROVIDER);
const ks = {
  address: "hxb6b5791be0b5ef67063b3c10b840fb81514db2fd",
  id: "87323a66-289a-4ce2-88e4-00278deb5b84",
  version: 3,
  coinType: "icx",
  crypto: {
    cipher: "aes-128-ctr",
    cipherparams: {
      iv: "069e46aaefae8f1c1f840d6b09144999",
    },
    ciphertext:
      "f35ff7cf4f5759cb0878088d0887574a896f7f0fc2a73898d88be1fe52977dbd",
    kdf: "scrypt",
    kdfparams: {
      dklen: 32,
      n: 65536,
      r: 8,
      p: 1,
      salt: "0fc9c3b24cdb8175",
    },
    mac: "1ef4ff51fdee8d4de9cf59e160da049eb0099eb691510994f5eca492f56c817a",
  },
};
const ICON_WALLET = IconWallet.loadKeystore(ks as KeyStore, "gochain", false);

async function sendMessage(_to: string, _data: string, _rollback:string) {
  try {
    const fee = await getFee(true);

    const params = {
      _to: _to,
      _data: _data,
      _rollback: _rollback
    };
    const txObj = new CallTransactionBuilder()
      .from(ICON_WALLET.getAddress())
      .to(ICON_DAPP)
      .method("sendMessage")
      .params(params)
      .stepLimit(IconConverter.toBigNumber(5000000000))
      .nid(IconConverter.toBigNumber(NID))
      .nonce(IconConverter.toBigNumber(1))
      .version(IconConverter.toBigNumber(3))
      .timestamp(new Date().getTime() * 1000)
      .value(fee)
      .build();

    const signedTx = new SignedTransaction(txObj, ICON_WALLET);
    return await ICON_SERVICE.sendTransaction(signedTx).execute();
  } catch (e) {
    console.log(e);
    throw new Error("Error calling contract method");
  }
}

async function getFee(useRollback = false) {
  try {
    const params = {
      _net: NETWORK_LABEL_DESTINATION,
      _rollback: useRollback ? "0x1" : "0x0",
    };

    const txObj = new CallBuilder()
      .to(ICON_XCALL)
      .method("getFee")
      .params(params)
      .build();

    return await ICON_SERVICE.call(txObj).execute();
  } catch (e) {
    console.log("error getting fee", e);
    throw new Error("Error getting fee");
  }
}

function strToHex(str: string) {
  let hex = "";
  for (let i = 0; i < str.length; i++) {
    hex += "" + str.charCodeAt(i).toString(16);
  }
  return hex;
}

function sleep(millis: number) {
  return new Promise((resolve) => setTimeout(resolve, millis));
}

function filterEvent(
  eventLogs: any,
  sig: string,
  address?: string
): Array<EventLog> {
  return (<Array<EventLog>>eventLogs).filter(
    (eventLog) =>
      eventLog.indexed &&
      eventLog.indexed[0] === sig &&
      (!address || address === eventLog.scoreAddress)
  );
}

function parseCallMessageSentEvent(event: any) {
  const indexed = event[0].indexed || [];
  const data = event[0].data || [];
  return {
    _from: indexed[1],
    _to: indexed[2],
    _sn: indexed[3],
  };
}

export async function waitEvent(sig: string, contract_address: string) {
  let latest = await ICON_SERVICE.getLastBlock().execute();
  let height = latest.height - 1;
  const heights = BigNumber.isBigNumber(height)
    ? height
    : new BigNumber(height as number);
  let block = await ICON_SERVICE.getBlockByHeight(heights).execute();
  while (true) {
    while (height < latest.height) {
      const events = await filterEventFromBlock(block, sig, contract_address);
      if (events.length > 0) {
        return events;
      }
      height++;
      if (height === latest.height) {
        block = latest;
      } else {
        const heightss = BigNumber.isBigNumber(height)
          ? height
          : new BigNumber(height as number);
        block = await ICON_SERVICE.getBlockByHeight(heightss).execute();
      }
    }
    await new Promise((resolve) => setTimeout(resolve, 100));
    latest = await ICON_SERVICE.getLastBlock().execute();
  }
}

async function filterEventFromBlock(
  block: Block,
  sig: string,
  address?: string | undefined
): Promise<EventLog[]> {
  return Promise.all(
    block
      .getTransactions()
      .map((tx: ConfirmedTransaction) =>
        ICON_SERVICE.getTransactionResult(tx.txHash).execute()
      )
  ).then((results) => {
    return results
      .map((result: TransactionResult) => {
        return filterEvent(result.eventLogs as Array<EventLog>, sig, address);
      })
      .flat();
  });
}

async function verifyCallMessageEventIcon() {
  let events = await waitEvent(
    callMessageSignature,
    ICON_XCALL
  );
  if (events.length > 0) {
    const indexed = events[0].indexed || [];
    const data = events[0].data || [];
    const event = {
      _from: indexed[1],
      _to: indexed[2],
      _sn: IconConverter.toNumber(indexed[3]),
      _reqId: IconConverter.toNumber(data[0]),
      _data: data[1],
    };
    console.log(event);
    return {
      _reqId: event._reqId,
      _data: event._data,
    };
  }
}

async function executeCall(reqId: number, data: string) {
  try {
    const fee = await getFee();

    const params = {
      _reqId: `${reqId.toString()}`,
      _data: data,
    };
    const txObj = new CallTransactionBuilder()
      .from(ICON_WALLET.getAddress())
      .to(ICON_XCALL)
      .method("executeCall")
      .params(params)
      .stepLimit(IconConverter.toBigNumber(5000000000))
      .nid(IconConverter.toBigNumber(NID))
      .nonce(IconConverter.toBigNumber(1))
      .version(IconConverter.toBigNumber(3))
      .timestamp(new Date().getTime() * 1000)
      .value(fee)
      .build();

    const signedTx = new SignedTransaction(txObj, ICON_WALLET);
    return await ICON_SERVICE.sendTransaction(signedTx).execute();
  } catch (e) {
    console.log(e);
    throw new Error("Error calling contract method");
  }
}

async function verifyCallExecutedEventIcon(eventLogs:TransactionResult["eventLogs"]) {
  const filtereCallExecute = filterEvent(
    eventLogs,
    callExecutedSignature,
    ICON_XCALL
  );  
  console.log(filtereCallExecute);
}

async function verifyResponseMessageEventIcon() {
  let events = await waitEvent(
    responseMessageSignature,
    ICON_XCALL
  );
  if (events.length > 0) {
    const indexed = events[0].indexed || [];
    const data = events[0].data || [];
    const event = {
      _sn: IconConverter.toNumber(indexed[1]),
      _code: IconConverter.toNumber(data[0]),
    }
    console.log(events);
    return {
      _sn: event._sn,
    };
  }
}

async function verifyRollbackMessageEventIcon(){
  let events = await waitEvent(
    rollbackMessageSignature,
    ICON_XCALL
  );
  if (events.length > 0) {
    const indexed = events[0].indexed || [];
    const event = {
      _sn: IconConverter.toNumber(indexed[1]),
    }
    console.log(events);
  }

}

async function executeRollback(seqNo: number) {
  try {
    const fee = await getFee();
    const params = {
      _sn: `${seqNo.toString()}`,
    };
    const txObj = new CallTransactionBuilder()
      .from(ICON_WALLET.getAddress())
      .to(ICON_XCALL)
      .method("executeRollback")
      .params(params)
      .stepLimit(IconConverter.toBigNumber(5000000000))
      .nid(IconConverter.toBigNumber(NID))
      .nonce(IconConverter.toBigNumber(1))
      .version(IconConverter.toBigNumber(3))
      .timestamp(new Date().getTime() * 1000)
      .value(fee)
      .build();

    const signedTx = new SignedTransaction(txObj, ICON_WALLET);
    return await ICON_SERVICE.sendTransaction(signedTx).execute();
  } catch (e) {
    console.log(e);
    throw new Error("Error calling contract method");
  }
}

async function verifyRollbackExecutedEventIcon(){
  let events = await waitEvent(
    "RollbackExecuted(int)",
    ICON_XCALL
  );
  if (events.length > 0) {
    const indexed = events[0].indexed || [];
    const event = {
      _sn: IconConverter.toNumber(indexed[1]),
    }
    console.log(events);
    return {
      _sn: event._sn,
    };
  }
}

async function main() {
  const _to = `${NETWORK_LABEL_DESTINATION}/${DESTINATION_DAPP}`;
  const _data = strToHex("rollback");
  const _rollback = strToHex("This is the rollback meesage to be executed")

  const receipt = await sendMessage(_to, _data, _rollback);
  console.log(receipt);
  await sleep(5000);
  const txResult = await ICON_SERVICE.getTransactionResult(receipt).execute();
  const filteredEvent = filterEvent(
    txResult.eventLogs,
    callMessageSentSignature,
    ICON_XCALL
  );
  // parsing the CallMessageSent event logs
  const parsedEvent = parseCallMessageSentEvent(filteredEvent);
  console.log("parsedEvent", parsedEvent);

  // // Verify CallMessage event
  const callMsgEvent = await verifyCallMessageEventIcon();
  const request_id = callMsgEvent!._reqId;
  const Data = callMsgEvent!._data;

  // // Execute Call
  const execReceipt = await executeCall(request_id, Data);
  await sleep(5000);
  const execResult = await ICON_SERVICE.getTransactionResult(
    execReceipt
  ).execute();
  console.log(execResult);

  // // verify Call Executed Event
  await verifyCallExecutedEventIcon(execResult.eventLogs);

  // verify Response Message Event
  const seqNo = await verifyResponseMessageEventIcon()
  console.log("seqNo: ", seqNo?._sn);

  // verify Rollback Message Event
  await verifyRollbackMessageEventIcon()

  // Execute Rollback
  const execRollback = await executeRollback(seqNo?._sn!)
  console.log(execRollback);

  //verify rollbackExecuted event 
  await verifyRollbackExecutedEventIcon()
}

main();
