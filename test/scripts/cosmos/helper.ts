import { GasPrice } from "@cosmjs/stargate";
import { exec } from "child_process";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing"
import { SigningCosmWasmClient } from "@cosmjs/cosmwasm-stargate";
import { Deployments } from "../setup/config";

const defaultGasPrice = GasPrice.fromString("0stake");

export async function CreateSigningClient(
  mnemonic: string,
  prefix: string,
  endpoint: string
): Promise<[SigningCosmWasmClient, string]> {

  const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, { prefix: prefix });

  // Get account address
  const accounts = await wallet.getAccounts();
  const accountAddress = accounts[0].address;

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

export async function getStake(testaddress: string, destaddress: string, chainName: string) {
  const dockerID = await getContainerIdByPartialName(chainName);
  
  let command = "";
  
  if (chainName == "archway") {
    command = `docker exec ${dockerID} archwayd tx bank send ${testaddress} ${destaddress} 9000000stake --keyring-backend test \
    --chain-id constantine-3 -y`;
  } else if (chainName == "neutron"){
    command = `docker exec ${dockerID} neutrond tx bank send ${testaddress} ${destaddress} 9000000stake --keyring-backend test \
    --chain-id test-chain1 --home ./data/test-chain1 -y`;
  }

  exec(command, (error, stdout, stderr) => {
    if (error) {
      console.error(`Error executing docker exec: ${error.message}`);
      return;
    }

    console.log("Command executed successfully!");
    console.log("Output:", stdout);
  });
}

export async function getTestAccountWithStake(chainName: string): Promise<string> {
  const dockerID = await getContainerIdByPartialName(chainName);
  let command: string;
  
  if (chainName == "archway") {
    command = `docker exec ${dockerID} archwayd keys list --keyring-backend test |  grep 'address:'`;
  } else if (chainName == "neutron"){
    command = `docker exec ${dockerID} neutrond keys list --keyring-backend test --home ./data/test-chain1 | grep -A 4 "name: test-key" | awk '/address:/ {print $2}'`;
  }

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

export async function getContainerIdByPartialName(chainName: string): Promise<string> {
  let command: string;
  
  if (chainName == "archway") {
    command = 'docker ps -aqf "name=constantine-3"';
  } else if (chainName == "neutron"){
    command = 'docker ps -aqf "name=neutron-node-test-chain1"';
  }

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

const config = process.env.CONFIG_FILE || "dive.json";
const deployments = Deployments.getDefault(config);
const SRC = deployments.getSrc();
const DST = deployments.getDst();
const srcChain = deployments.get(SRC);
const dstChain = deployments.get(DST);
const srcContracts = deployments.getContracts(SRC);
const dstContracts = deployments.getContracts(DST);

export function GetCosmosContracts(contract: string) {
  if (srcChain.service_name.includes("constantine-3")){
    return srcContracts[contract]
  } else {
    return dstContracts[contract]
  }
}

export function GetIconContracts(contract: string) {
  if (srcChain.service_name.includes("icon")){
    return srcContracts[contract]
  } else {
    return dstContracts[contract]
  }
}

export function GetChainInfo(chainName:string, args: string){
  if (chainName.includes("icon")){
    return srcChain[args]
  } else {
    return dstChain[args]
  }
}

export function GetCosmosChainInfo(chainID:string, args: string){
  if (srcChain.chain_id == chainID){
    return srcChain[args]
  } else {
    return dstChain[args]
  }
}


export function GetIconChainInfo(args: string){
  if (srcChain.service_name.includes("icon")){
    return srcChain[args]
  } else {
    return dstChain[args]
  }
}

export function GetSrc(){
  return SRC
}

export function GetDest(){
  return DST
}

export function GetDataInBytes(msg: string) {
  const bytes: number[] = [];

  for (let i = 0; i < msg.length; i++) {
    const charCode = msg.charCodeAt(i);
    bytes.push(charCode);
  }
  return bytes;
}

export function strToHex(str: string) {
  let hex = "";
  for (let i = 0; i < str.length; i++) {
    hex += "" + str.charCodeAt(i).toString(16);
  }
  return hex;
}

export function hexToString(hex: string): string {
  let str = '';
  for (let i = 0; i < hex.length; i += 2) {
    const charCode = parseInt(hex.slice(i, i + 2), 16);
    str += String.fromCharCode(charCode);
  }
  return str;
}


