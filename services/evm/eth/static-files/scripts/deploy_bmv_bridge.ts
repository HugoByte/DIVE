import { ethers } from "hardhat";
import { Deployments } from "./config";
const deployments = Deployments.getDefault();

async function deploy_bmv_bridge(
  currentChain: string,
  srcNetLastblockHeight: string,
  srcChainNetwok: string
) {

  var depolyment_data = deployments.get(currentChain);
 
  var currentChainNetwork = depolyment_data.network;

  var currentChainBmcAddress = depolyment_data.contracts.bmc;

 
  const BMVBridge = await ethers.getContractFactory("BMV");
  const bmvb = await BMVBridge.deploy(
    currentChainBmcAddress,
    srcChainNetwok,
    srcNetLastblockHeight
  );
  await bmvb.deployed();

  depolyment_data.contracts.bmvb = bmvb.address;
  

  deployments.set(currentChain, depolyment_data);
  deployments.save();

  console.log(`${currentChain}: BMV-Bridge: deployed to ${bmvb.address}`);
}

async function setup_link(
  currentChain: string,
  srcBmcAddress: string,
  srcChainNetwork: string
) {
  var depolyment_data = deployments.get(currentChain);

  var currentChainBmcmAddress = depolyment_data.contracts.bmcm;

  const bmcm = await ethers.getContractAt(
    "BMCManagement",
    currentChainBmcmAddress
  );

  var srcChainBmcAddress = `btp://${srcChainNetwork}/${srcBmcAddress}`;

  console.log(`${currentChain}: addVerifier for ${srcChainNetwork}`);

  var currentBmvAddress = depolyment_data.contracts.bmvb;

  await bmcm.addVerifier(srcChainNetwork, currentBmvAddress).then((tx) => {
    return tx.wait(1);
  });
  console.log(`${currentChain}: addLink: ${srcChainBmcAddress}`);
  await bmcm.addLink(srcChainBmcAddress).then((tx) => {
    return tx.wait(1);
  });
  console.log(`${currentChain}: addRelay`);
  const signers = await ethers.getSigners();
  await bmcm
    .addRelay(srcChainBmcAddress, signers[0].getAddress())
    .then((tx) => {
      return tx.wait(1);
    });
}

async function main() {
  var data = JSON.parse(process.env.params?.toString());

  var currentChain = data.current_chain.name;
  var srcNetLastblockHeight = data.src.lastBlockHeight;
  var srcChainBmcAddress = data.src.bmc;
  var srcChainNetwork = data.src.network;

  await deploy_bmv_bridge(currentChain, srcNetLastblockHeight, srcChainNetwork);

  await setup_link(currentChain, srcChainBmcAddress, srcChainNetwork);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
