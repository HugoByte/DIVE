import { ethers } from "hardhat";
import { Deployments } from "./config";
const deployments = Deployments.getDefault();

async function deploy_xcall_solidity(
  target: string,
  bmcAddress: string,
  bmcmAddress: string
) {
  var deployments_data = deployments.get(target);
  const CallSvc = await ethers.getContractFactory("CallService");
  const xcallSol = await CallSvc.deploy({gasPrice: 50000000000});
  await xcallSol.deployed();
  await xcallSol.initialize(bmcAddress);

  console.log(`${target}: xCall: deployed to ${xcallSol.address}`);

  deployments_data.contracts.xcall = xcallSol.address;

  console.log(`${target}: register xCall to BMC`);
  const bmcm = await ethers.getContractAt("BMCManagement", bmcmAddress);
  await bmcm.addService("xcall", deployments_data.contracts.xcall);

  deployments.set(target, deployments_data);

  deployments.save();
}

async function main() {
  var data = JSON.parse(process.env.params?.toString());

  var deployments_data = deployments.get(data.name);

  await deploy_xcall_solidity(
    data.name,
    deployments_data.contracts.bmc,
    deployments_data.contracts.bmcm
  );
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
