import { GasPrice, SigningStargateClient } from "@cosmjs/stargate";
import { exec } from "child_process";
import fs from "fs";
import { Secp256k1HdWallet } from "@cosmjs/launchpad";

const defaultGasPrice = GasPrice.fromString("0stake");

export async function CreateSigningClient(
  mnemonic: string,
  prefix: string,
  endpoint: string
): Promise<[SigningStargateClient, string]> {

  const wallet = await Secp256k1HdWallet.fromMnemonic(mnemonic, { prefix: prefix });

  // Get account address
  const accounts = await wallet.getAccounts();
  const accountAddress = accounts[0].address;

  // Create an signing client with created wallet
  const signingClient = await SigningStargateClient.connectWithSigner(
    endpoint,
    wallet,
    {
      gasPrice: defaultGasPrice,
    }
  );
  return [signingClient, accountAddress];
}

export async function getHeight(
  client: SigningStargateClient,
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
  client: SigningStargateClient,
  address: string
) {
  const balance = await client.getBalance(address, "stake");
  return balance.amount;
}

export function GetCosmosContracts(contract: string, chainName: string) {
  var dataArray = JSON.parse(fs.readFileSync("contracts.json", "utf-8"));
  return dataArray[chainName]["contracts"][contract];
}

export function GetIconContracts(contract: string) {
  var dataArray = JSON.parse(fs.readFileSync("contracts.json", "utf-8"));
  return dataArray["icon"]["contracts"][contract];
}

export function GetChainInfo(chainName: string,args: string){
  var dataArray = JSON.parse(fs.readFileSync("contracts.json", "utf-8"))
  return dataArray[chainName][args];
}

export function GetNeutronChainInfo(args: string){
  var dataArray = JSON.parse(fs.readFileSync("contracts.json", "utf-8"))
  return dataArray["neutron"][args];
}


export function GetIconChainInfo(args: string){
  var dataArray = JSON.parse(fs.readFileSync("contracts.json", "utf-8"))
  return dataArray["icon"][args];
}

export function GetSrc(){
  var dataArray = JSON.parse(fs.readFileSync("contracts.json", "utf-8"))
  return dataArray["path"]["src"];
}

export function GetDest(){
  var dataArray = JSON.parse(fs.readFileSync("contracts.json", "utf-8"))
  return dataArray["path"]["dest"];
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