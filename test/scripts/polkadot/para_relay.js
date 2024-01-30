import { ApiPromise, WsProvider, Keyring } from "@polkadot/api";
import { Builder } from "@paraspell/sdk";

async function testRelayToPara() {
  console.log("*".repeat(100));
  console.log("Demo on Sending tokens from Parachain to Relaychain. \n");
  console.log("*".repeat(100), "\n");

  console.log("Establishing connection to Parachain....");
  let ParaWsProvider = new WsProvider("ws://127.0.0.1:42233");
  const paraAPI = await ApiPromise.create({
    provider: ParaWsProvider,
    noInitWarn: true,
  });

  const paraBlock = await paraAPI.rpc.chain.getBlock();
  console.log(
    "Latest Parachain Block Height: ",
    paraBlock.block.header.number.toHuman(),
    "\n"
  );

  console.log("Establishing connection to relaychain....");
  let RelayWsProvider = new WsProvider("ws://127.0.0.1:25885");
  const relayAPI = await ApiPromise.create({
    provider: ParaWsProvider,
    noInitWarn: true,
  });

  const relayBlock = await relayAPI.rpc.chain.getBlock();
  console.log(
    "Latest Relaychain Block Height: ",
    relayBlock.block.header.number.toHuman(),
    "\n"
  );

  const accountAddress = "gXCcrjjFX3RPyhHYgwZDmw8oe4JFpd5anko3nTY8VrmnJpe";
  console.log(
    "Destination relay Chain Account Address : ",
    accountAddress,
    "\n"
  );

    const initialBal = await relayAPI.query.tokens.accounts(accountAddress, {
    token: "KSM",
  });
  console.log(
    "Initial KSM Balance on Destination chain : ",
    initialBal.free.toHuman()
  );
  
  // You can use ss58Format: ? as an argument when initalizing keyring to get exact alice address as mentioned for that particular parachain
  const keyring = new Keyring({ type: "sr25519" });
  const alice = keyring.addFromUri("//Alice");
  console.log("Alice Address : ", alice.address, "\n");

  const call = Builder(paraAPI)
    .from('Karura') 
    .amount(5000000000000) // Token amount
    .address("5CAqjCqo2CfzrbhuuqdGrBuEFrsKycxGd3KGaNn1qN89eFFG") // AccountId32 or AccountKey20 address
    .build()  

  const hash  = await call.signAndSend(alice);
  console.log(
    "Transaction Successfully submitted. \nHash: ",
    (hash.toHex())
  );

  paraAPI.disconnect();

// TODO: Listen for events and then check final Balance
await new Promise((f) => setTimeout(f, 25000));


  const finalBal = await relayAPI.query.tokens.accounts(accountAddress, {
    token: "KSM",
  });
  console.log(
    "KSM Balance on Destination chain after transfer: ",
    finalBal.free.toHuman()
  );

 
  relayAPI.disconnect();
}

testRelayToPara();