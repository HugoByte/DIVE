import { ethers } from "hardhat";
import { Deployments } from "./config";
const deployments = Deployments.getDefault();



async function getBlockNumber(currentChain: string){

    var depolyment_data = deployments.get(currentChain);

    const blockNum = await ethers.provider.getBlockNumber();

    depolyment_data.blockNum = blockNum;



    deployments.set(currentChain, depolyment_data);
    deployments.save();
}

async function main() {
    var data = JSON.parse(process.env.params?.toString());
  
    var currentChain = data.current_chain;
  
    await getBlockNumber(
      currentChain,
    );
  }
  
  main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });