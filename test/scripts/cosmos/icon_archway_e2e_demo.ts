import dotenv from "dotenv";
import {
  CreateSigningClient,
  GetCosmosContracts,
  GetDataInBytes,
  GetIconContracts,
  getHeight,
  getStake,
  getTestAccountWithStake,
} from "./helper";
import { fromTendermintEvent, GasPrice, calculateFee } from "@cosmjs/stargate";
import { SigningCosmWasmClient } from "@cosmjs/cosmwasm-stargate";
import IconService from "icon-sdk-js";

const { IconConverter } = IconService;

// configure dotenv
dotenv.config();

const defaultGasPrice = GasPrice.fromString("0stake");

async function main() {
  // Chain Constants, modify as required
  const chain1 = {
    chainId: "archway-node-0",
    endpoint: "http://localhost:4564",
    prefix: "archway",
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
  const data = GetDataInBytes();
  const receipt = await sendMessageFromDapp(accountAddress, signingClient, data);
  // verifyCallMessageEvent(signingClient, receipt);
}

async function sendMessageFromDapp(
  accountAddress: string,
  signingClient: SigningCosmWasmClient,
  data: number[]
) {
  const dapp = await GetCosmosContracts("dapp");
  const iconDappAddress = await GetIconContracts("dapp");
  const execMsg = {
    send_call_message: {
      to: "0x3.icon/" + iconDappAddress,
      data: data,
    },
  };

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



main();
