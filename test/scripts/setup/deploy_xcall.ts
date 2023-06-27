import {ethers} from 'hardhat';
import {Contract, IconNetwork, Jar, BMC} from "../icon";
import {Deployments, chainType} from "./config";

const {JAVASCORE_PATH} = process.env
const deployments = Deployments.getDefault();

async function deploy_xcall_java(target: string, chain: any) {
  const iconNetwork = IconNetwork.getNetwork(target);
  const content = Jar.readFromFile(JAVASCORE_PATH, "xcall", "0.6.2");
  const xcall = new Contract(iconNetwork)
  const deployTxHash = await xcall.deploy({
    content: content,
    params: {
      _bmc: chain.contracts.bmc,
    }
  })
  const result = await xcall.getTxResult(deployTxHash);
  if (result.status != 1) {
    throw new Error(`xCall deployment failed: ${result.txHash}`);
  }
  chain.contracts.xcall = xcall.address;
  console.log(`${target}: xCall: deployed to ${xcall.address}`);

  console.log(`${target}: register xCall to BMC`);
  const bmc = new BMC(iconNetwork, chain.contracts.bmc)
  await bmc.addService('xcall', chain.contracts.xcall)
    .then((txHash) => bmc.getTxResult(txHash))
    .then((result) => {
      if (result.status != 1) {
        throw new Error(`${target}: failed to register xCall to BMC: ${result.txHash}`);
      }
    })
}

async function deploy_xcall_solidity(target: string, chain: any) {
  const CallSvc = await ethers.getContractFactory("CallService")
  const xcallSol = await CallSvc.deploy()
  await xcallSol.deployed()
  await xcallSol.initialize(chain.contracts.bmc)
  chain.contracts.xcall = xcallSol.address
  console.log(`${target}: xCall: deployed to ${xcallSol.address}`);

  console.log(`${target}: register xCall to BMC`);
  const bmcm = await ethers.getContractAt('BMCManagement', chain.contracts.bmcm)
  await bmcm.addService('xcall', chain.contracts.xcall);
}

async function main() {
  const src = deployments.getSrc();
  const dst = deployments.getDst();
  const srcChain = deployments.get(src);
  const dstChain = deployments.get(dst);

  // deploy to src network first
  await deploy_xcall_java(src, srcChain);

  // deploy to dst network
  switch (chainType(dstChain)) {
    case 'icon':
      await deploy_xcall_java(dst, dstChain);
      break;
    case 'hardhat':
      await deploy_xcall_solidity(dst, dstChain);
      break;
    default:
      throw new Error(`Unknown chain type: ${chainType(dstChain)}`);
  }

  // update deployments
  deployments.set(src, srcChain);
  deployments.set(dst, dstChain);
  deployments.save();
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
