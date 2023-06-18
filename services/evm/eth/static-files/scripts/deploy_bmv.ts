import { ethers } from "hardhat";
import { Deployments } from "./config";
const deployments = Deployments.getDefault();


async function deploy_bmv(
  currentChain: string,
  srcChainfirstBlockHeader: string,
  srcChainNetwok: string,
  srcChainNetwokTypeId: string
) {
  var depolyment_data = deployments.get(currentChain);

  var currentChainBmc = depolyment_data.contracts.bmc;


  const BMVBtp = await ethers.getContractFactory("BtpMessageVerifier");

  const bmvBtp = await BMVBtp.deploy(
    currentChainBmc,
    srcChainNetwok,
    srcChainNetwokTypeId,
    srcChainfirstBlockHeader,
    "0x0"
  );
  await bmvBtp.deployed();
  depolyment_data.contracts.bmv = bmvBtp.address;

  deployments.set(currentChain, depolyment_data);

  console.log(`${currentChain}: BMV: deployed to ${bmvBtp.address}`);

  deployments.save();
}

async function setup_link(
  currentChain: string,
  srcChainBmcAddr: string,
  srcChainNetwork: string
) {
  var depolyment_data = deployments.get(currentChain);

  var currentChainBmcm = depolyment_data.contracts.bmcm;

  var currentChainBmv = depolyment_data.contracts.bmv;

  const bmcm = await ethers.getContractAt("BMCManagement", currentChainBmcm);

  var dstBmcAddr = `btp://${srcChainNetwork}/${srcChainBmcAddr}`;

  console.log(`${currentChain}: addVerifier for ${srcChainNetwork}`);

  var bmvAddress = currentChainBmv;

  await bmcm.addVerifier(srcChainNetwork, bmvAddress).then((tx) => {
    return tx.wait(1);
  });
  console.log(`${currentChain}: addLink: ${dstBmcAddr}`);
  await bmcm.addLink(dstBmcAddr).then((tx) => {
    return tx.wait(1);
  });
  console.log(`${currentChain}: addRelay`);
  const signers = await ethers.getSigners();
  await bmcm.addRelay(dstBmcAddr, signers[0].getAddress()).then((tx) => {
    return tx.wait(1);
  });
}

async function main() {
  var data = JSON.parse(process.env.params?.toString());

  var currentChain = data.current_chain.name;
  var srcFirstBlockHeader = data.src_chain.firstBlockHeader;
  var srcChainNetwork = data.src_chain.network;
  var srcChainNetworkTypeId = data.src_chain.networkTypeId;
  var srcChainBmcAddr = data.src_chain.bmc;

  await deploy_bmv(
    currentChain,
    srcFirstBlockHeader,
    srcChainNetwork,
    srcChainNetworkTypeId
  );

  await setup_link(currentChain, srcChainBmcAddr, srcChainNetwork);
}


main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
