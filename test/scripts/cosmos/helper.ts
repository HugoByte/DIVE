import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { SigningCosmWasmClient } from "@cosmjs/cosmwasm-stargate";
import { GasPrice } from "@cosmjs/stargate";
import { exec } from "child_process";
import fs from "fs";

const defaultGasPrice = GasPrice.fromString("0stake");

export async function CreateSigningClient(
  mnemonic: string,
  prefix: string,
  endpoint: string
): Promise<[SigningCosmWasmClient, string]> {
  const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
    prefix: prefix,
  });

  // Print account address
  const accounts = await wallet.getAccounts();
  const accountAddress = accounts[0].address;
  console.log("Test Account address : " + accountAddress);

  // Create an signing client with created wallet
  const signingClient = await SigningCosmWasmClient.connectWithSigner(
    endpoint,
    wallet,
    {
      gasPrice: defaultGasPrice,
    }
  );
  return [signingClient, accountAddress];
}

export async function getHeight(
  client: SigningCosmWasmClient,
  chainId: string
) {
  // To Check if the client is connected to local chain
  const height = await client.getHeight();
  console.log("Current Block Height on " + chainId + " :" + height + "\n");
}

export async function getStake(testaddress: string, destaddress: string) {
  const dockerID = await getContainerIdByPartialName();
  console.log(dockerID);

  const commando = `docker exec ${dockerID} archwayd tx bank send ${testaddress} ${destaddress} 9000000stake --keyring-backend test \
    --chain-id constantine-3 -y`;

  exec(commando, (error, stdout, stderr) => {
    if (error) {
      console.error(`Error executing docker exec: ${error.message}`);
      return;
    }

    console.log("Command executed successfully!");
    console.log("Output:", stdout);
  });
}

export async function getTestAccountWithStake(): Promise<string> {
  const dockerID = await getContainerIdByPartialName();
  const command = `docker exec ${dockerID} archwayd keys list --keyring-backend test |  grep 'address:'`;
  return new Promise<string>((resolve, reject) => {
    exec(command, (error, stdout, stderr) => {
      if (error) {
        console.error(`Error executing docker exec: ${error.message}`);
        reject(error);
        return;
      }
      const output = stdout.trim();
      resolve(output);
    });
  });
}

export async function getContainerIdByPartialName(): Promise<string> {
  const command = 'docker ps -aqf "name=constantine-3"';

  return new Promise<string>((resolve, reject) => {
    exec(command, (error, stdout, stderr) => {
      if (error) {
        console.error(`Error executing docker exec: ${error.message}`);
        reject(error);
        return;
      }
      const dockerID = stdout.trim();
      resolve(dockerID);
    });
  });
}

export async function getBalance(
  client: SigningCosmWasmClient,
  address: string
) {
  const balance = await client.getBalance(address, "stake");
  return balance.amount;
}

export async function GetCosmosContracts(contract: string) {
  var dataArray = JSON.parse(fs.readFileSync("contracts.json", "utf-8"));
  return dataArray["node-service-constantine-3"]["contracts"][contract];
}

export async function GetIconContracts(contract: string) {
  var dataArray = JSON.parse(fs.readFileSync("contracts.json", "utf-8"));
  return dataArray["icon-node-0xacbc4e"]["contracts"][contract];
}

export function GetDataInBytes() {
  const msg = "Sending message from Cosmos to Icon0";
  const bytes: number[] = [];

  for (let i = 0; i < msg.length; i++) {
    const charCode = msg.charCodeAt(i);
    bytes.push(charCode);
  }
  return bytes;
}
