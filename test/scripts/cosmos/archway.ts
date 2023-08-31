import dotenv from "dotenv";
import {
  CreateSigningClient,
  GetArchwayChainInfo,
  GetCosmosContracts,
  GetDataInBytes,
  GetIconChainInfo,
  GetIconContracts,
  getHeight,
  getStake,
  getTestAccountWithStake,
} from "./helper";
import { fromTendermintEvent, GasPrice, calculateFee } from "@cosmjs/stargate";
import { SigningCosmWasmClient } from "@cosmjs/cosmwasm-stargate";

// configure dotenv
dotenv.config();

const defaultGasPrice = GasPrice.fromString("0stake");
let signingClient: SigningCosmWasmClient;
let accountAddress: string;

async function Setup(): Promise<[SigningCosmWasmClient, string]> {
  // Chain Constants, modify as required
  const chain1 = {
    chainId: GetArchwayChainInfo("chainId"),
    endpoint: GetArchwayChainInfo("endpoint"),
    prefix: GetArchwayChainInfo("prefix"),
  };

  // Create signing client and account address
  const mnemonic1 = process.env.MNEMONIC1 as string;
  [signingClient, accountAddress] = await CreateSigningClient(
    mnemonic1,
    chain1.prefix,
    chain1.endpoint
  );

  // Get Test Account with stake
  const testAccount = await getTestAccountWithStake();
  const testAddress = testAccount.substring(8, testAccount.length).trim();

  // To Get balance of given account address and transfer balance if 0
  const bal = await signingClient.getBalance(accountAddress, "stake");
  if (bal.amount == "0") {
    console.log(
      "No Balance in Signer account, Transferring balance to Signer account"
    );
    await getStake(testAddress!, accountAddress);
  }
  await new Promise((f) => setTimeout(f, 5000));
  return [signingClient, accountAddress];
}

async function main() {
  // Chain Constants, modify as required
  const chain1 = {
    chainId: GetArchwayChainInfo("chainId"),
    endpoint: GetArchwayChainInfo("endpoint"),
    prefix: GetArchwayChainInfo("prefix"),
  };

  // Create signing client and account address
  const mnemonic1 = process.env.MNEMONIC1 as string;
  console.log("");
  console.log("Chain 1: " + chain1.chainId);
  const [signingClient, accountAddress] = await CreateSigningClient(
    mnemonic1,
    chain1.prefix,
    chain1.endpoint
  );

  // To Check if the client is connected to local chain
  getHeight(signingClient, chain1.chainId);

  // Get Test Account with stake
  const testAccount = await getTestAccountWithStake();
  const testAddress = testAccount.substring(8, testAccount.length).trim();

  // To Get balance of given account address and transfer balance if 0
  const bal = await signingClient.getBalance(accountAddress, "stake");
  if (bal.amount == "0") {
    console.log(
      "No Balance in Signer account, Transferring balance to Signer account"
    );
    await getStake(testAddress!, accountAddress);
  }
  await new Promise((f) => setTimeout(f, 5000));
  const data = GetDataInBytes("Sending message from Cosmos to Icon0");
  const rbData = GetDataInBytes("RollBack Data");
  const receipt = await sendMessageFromDapp(
    accountAddress,
    signingClient,
    data,
    "RollBack Data"
  );
  verifyCallMessageSentEvent(signingClient, receipt);
  const [reqId, dataObject] = await verifyCallMessageEvent(signingClient);
  await executeCall(signingClient, reqId, dataObject, accountAddress);
  await verifyCallExecutedEvent(signingClient);
  const seqNo = await verifyResponseMessageEvent(signingClient);
  await verifyRollbackMessageEvent(signingClient);
  await executeRollback(signingClient, accountAddress, seqNo);
  await rollbackExecutedEvent(signingClient);
}

export async function sendMessageFromDAppCosmos(
  data: number[],
  rollbackData?: string
) {
  [signingClient, accountAddress] = await Setup();
  return await sendMessageFromDapp(
    accountAddress,
    signingClient,
    data,
    rollbackData
  );
}

export async function verifyCallMessageSentEventArchway(receipt: string) {
  await verifyCallMessageSentEvent(signingClient, receipt);
}

async function sendMessageFromDapp(
  accountAddress: string,
  signingClient: SigningCosmWasmClient,
  data: number[],
  rbData?: string
) {
  const dapp = await GetCosmosContracts("dapp");
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

async function verifyCallMessageSentEvent(
  signingClient: SigningCosmWasmClient,
  exeResult: any
) {
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

async function verifyCallMessageEvent(signingClient: SigningCosmWasmClient) {
  const event = await waitForEvent(signingClient, "wasm-CallMessage");
  console.log("*****");
  console.log(event?.attributes);
  const reqIdObject = event?.attributes.find((item) => item.key === "reqId");
  const dataObject = event?.attributes.find((item) => item.key === "data");
  return [reqIdObject?.value, dataObject?.value];
}

async function waitForEvent(
  signingClient: SigningCosmWasmClient,
  eventName: string
) {
  let height = await signingClient.getHeight();
  let flag = false;
  while (!flag) {
    console.log(height);
    let tmp = height;
    const query = `tx.height=` + height;
    await sleep(5000);
    const txs = await signingClient.searchTx(query);
    if (txs.length > 0) {
      for (const tx of txs) {
        const events = tx.events;
        console.log(events);
        for (const event of events) {
          if (event.type === eventName) {
            const decodedEvent = fromTendermintEvent(event);
            flag = true;
            return decodedEvent;
          }
        }
      }
    }
    while (height < tmp + 1) {
      height = await signingClient.getHeight();
    }
  }
}

async function executeCall(
  signingClient: SigningCosmWasmClient,
  reqId: any,
  data: any,
  accountAddress: string
) {
  const xcall = await GetCosmosContracts("xcall");
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
  console.log("executeCall transactioin Hash: " + exeResult.transactionHash);
  return exeResult;
}

async function verifyCallExecutedEvent(signingClient: SigningCosmWasmClient) {
  const event = await waitForEvent(signingClient, "wasm-CallExecuted");
  console.log(event);
}

async function verifyResponseMessageEvent(
  signingClient: SigningCosmWasmClient
) {
  console.log("************ ResponseMEssage Event*****************");
  const event = await waitForEvent(signingClient, "wasm-ResponseMessage");
  console.log(event);
  const seqNo = event?.attributes.find((item) => item.key === "sn");
  return seqNo?.value;
}

async function verifyRollbackMessageEvent(
  signingClient: SigningCosmWasmClient
) {
  console.log("************ RollbackMEssage Event*****************");
  const event = await waitForEvent(signingClient, "wasm-RollbackMessage");
  console.log(event);
}

async function executeRollback(
  signingClient: SigningCosmWasmClient,
  accountAddress: string,
  seqNo: any
) {
  console.log("************ Execute rollback message*****************");
  const xcall = await GetCosmosContracts("xcall");
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
  console.log(
    "executeRollback transactioin Hash: " + exeResult.transactionHash
  );
  return exeResult;
}

async function rollbackExecutedEvent(signingClient: SigningCosmWasmClient) {
  console.log("************ RollbackMEssage Event*****************");
  const event = await waitForEvent(signingClient, "wasm-RollbackExecuted");
  console.log(event);
}
