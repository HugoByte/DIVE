import { ApiPromise, WsProvider, Keyring } from "@polkadot/api";
import { Builder } from "@paraspell/sdk";

async function testRelayToPara() {
  console.log("*".repeat(100));
  console.log("Demo on Sending tokens from Parachain to Relaychain. \n");
  console.log("*".repeat(100), "\n");

  console.log("Establishing connection to Parachain....");
  let ParaWsProvider = new WsProvider("ws://127.0.0.1:32886");
  const API = await ApiPromise.create({
    provider: ParaWsProvider,
    noInitWarn: true,
  });

  const relayBlock = await API.rpc.chain.getBlock();
  console.log(
    "Latest Parachain Block Height: ",
    relayBlock.block.header.number.toHuman(),
    "\n"
  );
  
  // You can use ss58Format: ? as an argument when initalizing keyring to get exact alice address as mentioned for that particular parachain
  const keyring = new Keyring({ type: "sr25519" });
  const alice = keyring.addFromUri("//Alice");
  console.log("Alice Address : ", alice.address, "\n");

  const call = Builder(API)
    .from('Karura') 
    .amount(5000000000000) // Token amount
    .address("5CAqjCqo2CfzrbhuuqdGrBuEFrsKycxGd3KGaNn1qN89eFFG") // AccountId32 or AccountKey20 address
    .build()  

  const hash  = await call.signAndSend(alice);
  console.log(
    "Transaction Successfully submitted. \nHash: ",
    (hash.toHex())
  );

  API.disconnect();
}

testRelayToPara();
