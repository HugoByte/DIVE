import { ethers } from "hardhat";
import { Deployments } from "./config";
const deployments = Deployments.getDefault();

async function deploy_dapp_solidity(target: string, xCallAddress: string) {
  var deployments_data = deployments.get(target);

  const DAppSample = await ethers.getContractFactory("DAppProxySample");
  const dappSol = await DAppSample.deploy();
  await dappSol.deployed();
  await dappSol.initialize(xCallAddress);
  console.log(`${target} DApp: deployed to ${dappSol.address}`);

  deployments_data.contracts.dapp = dappSol.address;

  deployments.set(target, deployments_data);

  deployments.save();
}

async function main() {
  var data = JSON.parse(process.env.params?.toString());

  var deployments_data = deployments.get(data.name);


  await deploy_dapp_solidity(data.name, deployments_data.contracts.xcall);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
