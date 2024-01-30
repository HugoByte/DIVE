import { ApiPromise, WsProvider, Keyring } from "@polkadot/api";
import { Builder } from "@paraspell/sdk";

async function testRelayToPara() {
  console.log("*".repeat(100));
  console.log("Demo on Sending tokens from Relaychain to Parachain. \n");
  console.log("*".repeat(100), "\n");

  console.log("Establishing connection to Relaychain....");
  let relayWsProvider = new WsProvider("ws://127.0.0.1:25885");
  const relayaAPI = await ApiPromise.create({ provider: relayWsProvider });

  const relayBlock = await relayaAPI.rpc.chain.getBlock();
  console.log(
    "Latest Relaychain Block Height: ",
    relayBlock.block.header.number.toHuman(),
    "\n"
  );

  console.log("Establishing connection to Parachain....");
  let paraWsProvider = new WsProvider("ws://127.0.0.1:42233");
  const paraAPI = await ApiPromise.create({
    provider: paraWsProvider,
    noInitWarn: true,
  });

  const paraBlock = await paraAPI.rpc.chain.getBlock();
  console.log(
    "Latest Para Chain Block Height: ",
    paraBlock.block.header.number.toHuman(),
    "\n"
  );
  const accountAddress = "gXCcrjjFX3RPyhHYgwZDmw8oe4JFpd5anko3nTY8VrmnJpe";
  console.log(
    "Destination Para Chain Account Address : ",
    accountAddress,
    "\n"
  );

  const initialBal = await paraAPI.query.tokens.accounts(accountAddress, {
    token: "KSM",
  });
  console.log(
    "Initial KSM Balance on Destination chain : ",
    initialBal.free.toHuman()
  );

  console.log(
    "Building an XCM call to transfer asset from Relaychain to Parachain...\n"
  );
  const call = Builder(relayaAPI)
    .to("Karura", 2000) // Destination Parachain and Para ID
    .amount(10000000000000) // Token amount
    .address("oQQwUS5xJwHYbx97jiU1YrnHN1L7PYaD4Uof8un6Hua5EqV") // AccountId32 or AccountKey20 address
    .build(); // Function called to build call

  console.log("Getting Alice address to sign and send the transaction.. \n");
  const keyring = new Keyring({ type: "sr25519" });
  const alice = keyring.addFromUri("//Alice");
  console.log("Alice Address : ", alice.address, "\n");

  const hash = await call.signAndSend(alice);
  console.log(
    "Transaction Successfully submitted. \nHash: ",
    JSON.stringify(hash)
  );
  relayaAPI.disconnect();

  // TODO: Listen for events and then check final Balance
  await new Promise((f) => setTimeout(f, 25000));

  const finalBal = await paraAPI.query.tokens.accounts(accountAddress, {
    token: "KSM",
  });
  console.log(
    "KSM Balance on Destination chain after transfer: ",
    finalBal.free.toHuman()
  );

  paraAPI.disconnect();
}

testRelayToPara();