import dotenv from "dotenv";
import {
  CreateSigningClient,
  GetChainInfo,
  GetCosmosContracts,
  GetDataInBytes,
  GetIconChainInfo,
  GetIconContracts,
  getHeight,
  getStake,
  getTestAccountWithStake,
} from "./helper";
import {
  Event,
  fromTendermintEvent,
  GasPrice,
  calculateFee,
} from "@cosmjs/stargate";
import { SigningCosmWasmClient } from "@cosmjs/cosmwasm-stargate";

// configure dotenv
dotenv.config();

// Define Global Variables
const defaultGasPrice = GasPrice.fromString("0stake");
let signingClient: SigningCosmWasmClient;
let accountAddress: string;

async function Setup(chainName: string): Promise<[SigningCosmWasmClient, string]> {
  // Chain Constants, modify as required in contracts.json
  const chain1 = {
    chainId: GetChainInfo(chainName, "chainId"),
    endpoint: GetChainInfo(chainName, "endpoint"),
    prefix: GetChainInfo(chainName, "prefix"),
  };

  // Create signing client and account address
  const mnemonic1 = process.env.MNEMONIC1 as string;
  [signingClient, accountAddress] = await CreateSigningClient(
    mnemonic1,
    chain1.prefix,
    chain1.endpoint
  );

  // Get Test Account with stake
  const testAccount = await getTestAccountWithStake(chainName);
  const testAddress = testAccount.substring(8, testAccount.length).trim();

  // To Get balance of given account address and transfer balance if 0
  const bal = await signingClient.getBalance(accountAddress, "stake");
  if (bal.amount == "0") {
    console.log(
      "No Balance in Signer account, Transferring balance to Signer account"
    );
    await getStake(testAddress!, accountAddress, chainName);
  }
  await new Promise((f) => setTimeout(f, 5000));
  return [signingClient, accountAddress];
}

export async function sendMessageFromDAppCosmos(
  data: number[],
  chainName: string,
  rollbackData?: string
) {
  [signingClient, accountAddress] = await Setup(chainName);
  return await sendMessageFromDapp(
    accountAddress,
    signingClient,
    data,
    chainName,
    rollbackData
  );
}

async function sendMessageFromDapp(
  accountAddress: string,
  signingClient: SigningCosmWasmClient,
  data: number[],
  chainName: string,
  rbData?: string
) {
  const dapp = await GetCosmosContracts("dapp", chainName);
  const iconDappAddress = await GetIconContracts("dapp");
  const DestNetwork = GetIconChainInfo("network");
  const execMsg = rbData
    ? {
        send_call_message: {
          to: DestNetwork + "/" + iconDappAddress,
          data: data,
          rollback: GetDataInBytes(rbData),
        },
      }
    : {
        send_call_message: {
          to: DestNetwork + "/" + iconDappAddress,
          data: data,
        },
      };
  console.log("params: ", execMsg);

  // To Execute Contract
  const defaultExecuteFee = calculateFee(1_500_000, defaultGasPrice);
  const exeResult = await signingClient.execute(
    accountAddress,
    dapp,
    execMsg,
    defaultExecuteFee
  );
  return exeResult;
}

export async function verifyCallMessageSentEventCosmos(exeResult: any) {
  const txResult = await signingClient.getTx(exeResult.transactionHash);
  const events = txResult?.events;
  for (const event of events!) {
    if (event.type === "wasm-CallMessageSent") {
      const decodedEvent = fromTendermintEvent(event);
      console.log(decodedEvent);
    }
  }
}

function sleep(millis: number) {
  return new Promise((resolve) => setTimeout(resolve, millis));
}

export async function verifyCallMessageEventCosmos(chainName: string) {
  [signingClient, accountAddress] = await Setup(chainName);
  return verifyCallMessageEvent(signingClient);
}

async function verifyCallMessageEvent(signingClient: SigningCosmWasmClient) {
  const event = await waitForEvent(signingClient, "wasm-CallMessage");
  return event?.[0];
}

async function waitForEvent(
  signingClient: SigningCosmWasmClient,
  eventName: string,
  Target_height?: number
): Promise<[Event, number] | undefined> {
  let height = await signingClient.getHeight();
  if (Target_height != undefined) {
    height = Target_height;
  }
  let decodedEvent: Event;
  let flag = false;
  while (!flag) {
    let tmp = height;
    const query = `tx.height=` + height;
    await sleep(5000);
    const txs = await signingClient.searchTx(query);
    if (txs.length > 0) {
      for (const tx of txs) {
        const events = tx.events;
        for (const event of events) {
          if (event.type === eventName) {
            decodedEvent = fromTendermintEvent(event);
            flag = true;
            return [decodedEvent, height];
          }
        }
      }
    }
    while (height < tmp + 1) {
      height = await signingClient.getHeight();
    }
  }
}

export async function executeCallCosmos(reqId: any, data: any, chainName: string) {
  const xcall = await GetCosmosContracts("xcall", chainName);
  const execMsg = {
    execute_call: {
      request_id: reqId.toString(),
      data: JSON.parse(data),
    },
  };

  // To Execute Contract
  const defaultExecuteFee = calculateFee(1_500_000, defaultGasPrice);
  const exeResult = await signingClient.execute(
    accountAddress,
    xcall,
    execMsg,
    defaultExecuteFee
  );
  return exeResult;
}

export async function verifyCallExecutedEventCosmos() {
  const response = await waitForEvent(signingClient, "wasm-CallExecuted");
  const event = response![0]
  const decodedEvent = fromTendermintEvent(event)
  console.log(decodedEvent);
  return response![1]
}

export async function verifyReceivedMessageCosmos(height: number) {
  let decodedEvent:Event
  const query = `tx.height=` + height;
    await sleep(5000);
    const txs = await signingClient.searchTx(query);
    if (txs.length > 0) {
      for (const tx of txs) {
        const events = tx.events;
        for (const event of events) {
          if (event.type === "wasm"){
            if (event.attributes[1].key === "data") {
              decodedEvent = fromTendermintEvent(event);
              return decodedEvent.attributes[1].value
            }
          }      
        }
      }
    }
}

export async function verifyResponseMessageEventCosmos(): Promise<[string, number]> {
  const values = await waitForEvent(signingClient, "wasm-ResponseMessage");
  console.log(values?.[0]);
  const seqNo = values?.[0].attributes.find((item) => item.key === "sn");
  return [seqNo!.value, values![1]];
}

export async function verifyRollbackMessageEventCosmos(height: number) {
  const event = await waitForEvent(
    signingClient,
    "wasm-RollbackMessage",
    height
  );
  console.log(event?.[0]);
}

export async function executeRollbackCosmos(seqNo: any, chainName: string) {
  const xcall = await GetCosmosContracts("xcall", chainName);
  const execMsg = {
    execute_rollback: {
      sequence_no: seqNo.toString(),
    },
  };

  // To Execute Contract
  const defaultExecuteFee = calculateFee(1_500_000, defaultGasPrice);
  const exeResult = await signingClient.execute(
    accountAddress,
    xcall,
    execMsg,
    defaultExecuteFee
  );
  console.log(exeResult);
  return exeResult;
}

export async function verifyRollbackExecutedEventCosmos() {
  const event = await waitForEvent(signingClient, "wasm-RollbackExecuted");
  console.log(event?.[0]);
}
