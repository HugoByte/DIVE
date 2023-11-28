import IconService from "icon-sdk-js";
import {IconNetwork, Chain, Gov} from "../icon";

const {CI_WORKFLOW} = process.env
const {IconAmount} = IconService;
const iconNetwork = IconNetwork.getDefault();
const chain = new Chain(iconNetwork);
const gov = new Gov(iconNetwork);
const prepAddress = iconNetwork.wallet.getAddress()

async function ensure_decentralization() {
  const mainPReps = await chain.getMainPReps();
  console.log(mainPReps)
  const prep = await chain.getPRep(prepAddress)
    .catch((error) => {
      console.log('Need to register PRep and get power first')
    });
  if (mainPReps.preps.length == 0 && prep == undefined) {
    const totalSupply = await iconNetwork.getTotalSupply()
    const minDelegated = totalSupply.div(500)
    const bondAmount = IconAmount.of(100_000, IconAmount.Unit.ICX).toLoop()

    console.log(`ICON: registerPRep`)
    const name = `node_${prepAddress}`
    await chain.registerPRep(name)
      .then((txHash) => chain.getTxResult(txHash))
      .then((result) => {
        if (result.status != 1) {
          throw new Error(`ICON: failed to registerPrep: ${result.txHash}`);
        }
      })

    console.log(`ICON: setStake`)
    await chain.setStake(minDelegated.plus(bondAmount))
      .then((txHash) => chain.getTxResult(txHash))
      .then((result) => {
        if (result.status != 1) {
          throw new Error(`ICON: failed to setStake: ${result.txHash}`);
        }
      })

    console.log(`ICON: setDelegation`)
    await chain.setDelegation(prepAddress, minDelegated)
      .then((txHash) => chain.getTxResult(txHash))
      .then((result) => {
        if (result.status != 1) {
          throw new Error(`ICON: failed to setDelegation: ${result.txHash}`);
        }
      })

    console.log(`ICON: setBonderList`)
    await chain.setBonderList(prepAddress)
      .then((txHash) => chain.getTxResult(txHash))
      .then((result) => {
        if (result.status != 1) {
          throw new Error(`ICON: failed to setBonderList: ${result.txHash}`);
        }
      })

    console.log(`ICON: setBond`)
    await chain.setBond(prepAddress, bondAmount)
      .then((txHash) => chain.getTxResult(txHash))
      .then((result) => {
        if (result.status != 1) {
          throw new Error(`ICON: failed to setBond: ${result.txHash}`);
        }
      })
  }

  if (mainPReps.preps.length == 0) {
    throw new Error(`ICON: need to wait until the next term for decentralization`);
  }
}

async function setup_node() {
  // ensure BTP revision
  const BTP_REVISION = 21
  const rev = parseInt(await chain.getRevision(), 16);
  console.log(`ICON: revision: ${rev}`)
  if (rev < BTP_REVISION) {
    console.log(`ICON: Set revision to ${BTP_REVISION}`)
    await gov.setRevision(BTP_REVISION)
      .then((txHash) => gov.getTxResult(txHash))
      .then((result) => {
        if (result.status != 1) {
          throw new Error(`ICON: failed to set revision: ${result.txHash}`);
        }
      })
  }

  // ensure public key registration
  const pubkey = await chain.getPRepNodePublicKey(prepAddress)
    .catch((error) => {
      console.log(`error: ${error}`)
    })
  console.log(`ICON: pubkey: ${pubkey}`)
  if (pubkey == undefined) {
    console.log('ICON: register PRep node publicKey')
    // register node publicKey in compressed form
    const pkey = iconNetwork.wallet.getPublicKey(true);
    await chain.registerPRepNodePublicKey(prepAddress, pkey)
      .then((txHash) => chain.getTxResult(txHash))
      .then((result) => {
        if (result.status != 1) {
          throw new Error(`ICON: failed to registerPRepNodePublicKey: ${result.txHash}`);
        }
      })
  }
}

function sleep(millis: number) {
  return new Promise(resolve => setTimeout(resolve, millis));
}

async function main() {
  let success = false;
  for (let i = 0; i < 21; i++) {
    success = await ensure_decentralization()
      .then(() => {
        return true;
      })
      .catch((error) => {
        if (CI_WORKFLOW == "true") {
          console.log(error);
          return false;
        }
        throw error;
      });
    if (success) {
      await setup_node()
        .then(() => {
          console.log('ICON: node setup completed')
        });
      break;
    }
    console.log(`... wait 10 seconds (${i})`)
    await sleep(10000);
  }
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
