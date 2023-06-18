import { ethers } from "hardhat";
import { Deployments } from "./config";
const deployments = new Deployments(new Map());

async function deploy_solidity(target: string, chain_network: string) {
  const BMCManagement = await ethers.getContractFactory("BMCManagement");
  const bmcm = await BMCManagement.deploy();
  await bmcm.deployed();
  await bmcm.initialize();

  const BMCService = await ethers.getContractFactory("BMCService");
  const bmcs = await BMCService.deploy();
  await bmcs.deployed();
  await bmcs.initialize(bmcm.address);

  const BMCPeriphery = await ethers.getContractFactory("BMCPeriphery");
  const bmcp = await BMCPeriphery.deploy();
  await bmcp.deployed();
  await bmcp.initialize(chain_network, bmcm.address, bmcs.address);

  await bmcm.setBMCPeriphery(bmcp.address).then((tx) => {
    return tx.wait(1);
  });

  await bmcm.setBMCService(bmcs.address).then((tx) => {
    return tx.wait(1);
  });

  await bmcs.setBMCPeriphery(bmcp.address).then((tx) => {
    return tx.wait(1);
  });

  var contract_address = {
    address: {
      bmcm: bmcm.address,
      bmcs: bmcs.address,
      bmc: bmcp.address,
    },
  };

  console.log(JSON.stringify(contract_address));

  deployments.set(target, {
    network: chain_network,
    contracts: {
      bmcm: bmcm.address,
      bmcs: bmcs.address,
      bmc: bmcp.address,
    },
  });
}

async function main() {
  var data = JSON.parse(process.env.params?.toString());

  await deploy_solidity(data.link, data.chainNetwork);

  deployments.save();
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
