import IconService, {
  BigNumber,
  Block,
  ConfirmedTransaction,
  KeyStore,
  TransactionResult,
} from "icon-sdk-js";
import {
  GetChainInfo,
  GetCosmosContracts,
  GetIconChainInfo,
  GetIconContracts,
  hexToString,
  strToHex,
} from "./helper";

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

const ICON_RPC_URL = GetIconChainInfo("endpoint");
const NID = GetIconChainInfo("nid");
const ICON_XCALL = GetIconContracts("xcall");
const ICON_DAPP = GetIconContracts("dapp");

const callMessageSentSignature = "CallMessageSent(Address,str,int)";
const callMessageSignature = "CallMessage(str,str,int,int,bytes)";
const callExecutedSignature = "CallExecuted(int,int,str)";
const responseMessageSignature = "ResponseMessage(int,int)";
const rollbackMessageSignature = "RollbackMessage(int)";
const MessageReceivedSignature = "MessageReceived(str,bytes)";

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

async function sendMessage(
  _to: string,
  _data: string,
  chainName: string,
  _rollback?: string,
  isRollback?: boolean
) {
  try {
    const fee = await getFee(isRollback, chainName);
    const _params = _rollback
      ? { _to: _to, _data: _data, _rollback: IconConverter.toHex(_rollback) }
      : { _to: _to, _data: _data };
    console.log("params: \n", _params, "\n");
    const txObj = new CallTransactionBuilder()
      .from(ICON_WALLET.getAddress())
      .to(ICON_DAPP)
      .method("sendMessage")
      .params(_params)
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

async function getFee(useRollback = false, chainName: string) {
  try {
    const params = {
      _net: GetChainInfo(chainName, "chain_id"),
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

export async function waitEvent(
  sig: string,
  contract_address: string
): Promise<[EventLog[], number]> {
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
        return [events, height];
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

export async function verifyCallMessageEventIcon() {
  let [events, height] = await waitEvent(callMessageSignature, ICON_XCALL);
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

export async function executeCallIcon(reqId: number, data: string, chainName: string) {
  try {
    const fee = await getFee(false, chainName);

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
    const receipt = await ICON_SERVICE.sendTransaction(signedTx).execute();
    await sleep(5000);
    return await ICON_SERVICE.getTransactionResult(receipt).execute();
  } catch (e) {
    console.log(e);
    throw new Error("Error calling contract method");
  }
}

export async function verifyCallExecutedEventIcon() {
  let [events, height] = await waitEvent(callExecutedSignature, ICON_XCALL);
  console.log(events[0]);
  let event;
  if (events.length > 0) {
    const indexed = events[0].indexed || [];
    const data = events[0].data || [];
    event = {
      _reqId: IconConverter.toNumber(indexed[1]),
      _code: IconConverter.toNumber(data[0]),
      _msg: data[1],
    };
  }
  console.log(event);
  return height;
}

export async function verifyResponseMessageEventIcon(): Promise<
  [number, number]
> {
  let [events, height] = await waitEvent(responseMessageSignature, ICON_XCALL);
  const indexed = events[0].indexed || [];
  const data = events[0].data || [];
  const event = {
    _sn: IconConverter.toNumber(indexed[1]),
    _code: IconConverter.toNumber(data[0]),
  };
  console.log(events);
  const seqNo = event._sn;
  return [seqNo!, height!];
}

export async function verifyRollbackMessageEventIcon(height: number) {
  const block = await ICON_SERVICE.getBlockByHeight(
    IconConverter.toBigNumber(height)
  ).execute();
  let events = await filterEventFromBlock(
    block,
    rollbackMessageSignature,
    ICON_XCALL
  );
  console.log(events);
}

export async function verifyReceivedMessageIcon(height: number) {
  const block = await ICON_SERVICE.getBlockByHeight(
    IconConverter.toBigNumber(height)
  ).execute();
  let events = await filterEventFromBlock(
    block,
    MessageReceivedSignature,
    ICON_DAPP
  );
  const data = events[0].data || [];
  const msg = data[1]
  return hexToString(msg.slice(2))
}

export async function executeRollbackIcon(seqNo: number, chainName: string) {
  try {
    const fee = await getFee(false, chainName);
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
    const receipt = await ICON_SERVICE.sendTransaction(signedTx).execute();
    await sleep(5000);
    const transaction = await ICON_SERVICE.getTransaction(receipt).execute();
    console.log(transaction);
  } catch (e) {
    console.log(e);
    throw new Error("Error calling contract method");
  }
}

export async function verifyRollbackExecutedEventIcon() {
  let [events, height] = await waitEvent("RollbackExecuted(int)", ICON_XCALL);
  if (events.length > 0) {
    const indexed = events[0].indexed || [];
    const event = {
      _sn: IconConverter.toNumber(indexed[1]),
    };
    console.log(events);
    return {
      _sn: event._sn,
    };
  }
}

export async function sendMessageFromDAppIcon(
  data: string,
  chainName:string,
  rollbackData?: string,
  isRollback?: boolean
) {
  const _to = `${GetChainInfo(chainName, "chain_id")}/${GetCosmosContracts("dapp")}`;
  return await sendMessage(_to, data, chainName, rollbackData, isRollback);
}

export async function verifyCallMessageSentEventIcon(receipt: string) {
  await sleep(5000);
  const txResult = await ICON_SERVICE.getTransactionResult(receipt).execute();
  const filteredEvent = filterEvent(
    txResult.eventLogs,
    callMessageSentSignature,
    ICON_XCALL
  );
  console.log(filteredEvent);
}
