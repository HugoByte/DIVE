if (process.argv.length < 3) {
  console.error('Usage: npx ts-node icon_archway_e2e_demo.ts <chainName, eg "archway"/"neutron">');
  process.exit(1);
}
const chainName = process.argv[2];

import {
  executeCallCosmos,
  executeRollbackCosmos,
  sendMessageFromDAppCosmos,
  verifyCallExecutedEventCosmos,
  verifyCallMessageEventCosmos,
  verifyCallMessageSentEventCosmos,
  verifyReceivedMessageCosmos,
  verifyResponseMessageEventCosmos,
  verifyRollbackExecutedEventCosmos,
  verifyRollbackMessageEventCosmos,
} from "./cosmos";
import { GetDataInBytes, GetDest, GetSrc, strToHex } from "./helper";
import { Deployments } from "../setup/config";
import {
  executeCallIcon,
  executeRollbackIcon,
  sendMessageFromDAppIcon,
  verifyCallExecutedEventIcon,
  verifyCallMessageEventIcon,
  verifyCallMessageSentEventIcon,
  verifyReceivedMessageIcon,
  verifyResponseMessageEventIcon,
  verifyRollbackExecutedEventIcon,
  verifyRollbackMessageEventIcon,
} from "./icon";


async function show_banner() {
  const banner = `
    ___           __
___ |__ \\___  ____/ /__  ____ ___  ____
/ _ \\__/ / _ \\/ __  / _ \\/ __ \`__ \\/ __ \\
/  __/ __/  __/ /_/ /  __/ / / / / / /_/ /
\\___/____\\___/\\__,_/\\___/_/ /_/ /_/\\____/
`;
  console.log(banner);
}

const SRC = GetSrc();
const DST = GetDest();

show_banner()
  .then(() => sendCallMessage(SRC, DST))
  .then(() => sendCallMessage(DST, SRC))
  .then(() => sendCallMessage(SRC, DST, "checkSuccessResponse", true))
  .then(() => sendCallMessage(DST, SRC, "checkSuccessResponse", true))
  .then(() => sendCallMessage(SRC, DST, "rollback", true))
  .then(() => sendCallMessage(DST, SRC, "rollback", true))
  .catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });

async function sendCallMessage(
  src: string,
  dst: string,
  msgData?: string,
  needRollback?: boolean
) {

  const testName = sendCallMessage.name + (needRollback ? "WithRollback" : "");
  console.log(`\n### ${testName}: ${src} => ${dst}`);
  if (!msgData) {
    msgData = `${testName}_${src}_${dst}`;
  }
  const expectRevert = msgData === "rollback";
  const rollbackData = needRollback
    ? `ThisIsRollbackMessage_${src}_${dst}`
    : undefined;
  let step = 1;

  console.log(`[${step++}] send message from DApp`);
  const sendMessageReceipt: any = await sendMessageFromDApp(
    src,
    msgData,
    rollbackData,
  );
  const sn = await verifyCallMessageSent(src, sendMessageReceipt!);

  console.log(`\n[${step++}] check CallMessage event on ${dst} chain`);
  const [reqId, callData]: any = await checkCallMessage(dst);

  console.log(`\n[${step++}] invoke executeCall with reqId=${reqId}`);
  const executeCallReceipt = await invokeExecuteCall(dst, reqId, callData);

  console.log(`\n[${step++}] check CallExecuted event on ${dst} chain`);
  const height = await checkCallExecuted(dst, executeCallReceipt, reqId);

  // Verify if correct meesage is received)
  if (!expectRevert){
  await verifyMessageReceived(dst, height!, msgData)
  }

  if (needRollback) {
    console.log(`\n[${step++}] check ResponseMessage event on ${src} chain`);
    const [responseHeight, seqNo]: any = await checkResponseMessage(
      src,
      expectRevert
    );

    if (expectRevert) {
      console.log(`\n[${step++}] check RollbackMessage event on ${src} chain`);
      const sn = await checkRollbackMessage(src, responseHeight);

      console.log(`\n[${step++}] invoke executeRollback with sn=${seqNo}`);
      const executeRollbackReceipt = await invokeExecuteRollback(src, seqNo);

      console.log(`\n[${step++}] check RollbackExecuted event on ${src} chain`);
      await checkRollbackExecuted(src);
    }
  }
}

async function sendMessageFromDApp(
  src: string,
  msg: string,
  rollback?: string,
) {
  const isRollback = rollback ? true : false;
  if (src.includes("icon")) {
    const hexMsg = strToHex(msg);
    return sendMessageFromDAppIcon(hexMsg, chainName, rollback, isRollback);
  } else if (src.includes("constantine-3") || src.includes("neutron-node-test-chain1")) {
    const bytesData = GetDataInBytes(msg);
    return await sendMessageFromDAppCosmos(bytesData, chainName, rollback);
  } else {
    throw new Error(`unknown source chain: ${src}`);
  }
}

async function verifyCallMessageSent(src: string, sendMessageReceipt: string) {
  console.log("**** Verify CallMessageSent Event ****");
  if (src.includes("icon")) {
    await verifyCallMessageSentEventIcon(sendMessageReceipt);
  } else if (src.includes("constantine-3") || src.includes("neutron-node-test-chain1")) {
    await verifyCallMessageSentEventCosmos(sendMessageReceipt);
  }
}

async function checkCallMessage(dst: string) {
  console.log("**** CallMessage Event ****");
  if (dst.includes("constantine-3")|| dst.includes("neutron-node-test-chain1")) {
    const eventLogs = await verifyCallMessageEventCosmos(chainName);
    console.log(eventLogs);
    const reqIdObject = eventLogs?.attributes.find(
      (item) => item.key === "reqId"
    );
    const dataObject = eventLogs?.attributes.find(
      (item) => item.key === "data"
    );
    return [reqIdObject!.value, dataObject!.value];
  } else if (dst.includes("icon")) {
    const eventLogs = await verifyCallMessageEventIcon();
    return [eventLogs!._reqId, eventLogs!._data];
  }
}

async function invokeExecuteCall(dst: string, reqId: any, callData: any) {
  console.log("**** Execute Call ****");
  if (dst.includes("constantine-3")|| dst.includes("neutron-node-test-chain1")) {
    console.log(await executeCallCosmos(reqId, callData, chainName));
  } else if (dst.includes("icon")) {
    console.log(await executeCallIcon(reqId, callData, chainName));
  }
}

async function checkCallExecuted(
  dst: string,
  executeCallReceipt: any,
  reqId: any
) {
  console.log("**** Verify CallExecuted Event ****");
  if (dst.includes("constantine-3")|| dst.includes("neutron-node-test-chain1")) {
    return await verifyCallExecutedEventCosmos();
  } else if (dst.includes("icon")) {
    return await verifyCallExecutedEventIcon();
  }
}

async function verifyMessageReceived(
  dst: string,
  height: number,
  msgData: string
) {
  let executedMsg:string | undefined;
  if (dst.includes("constantine-3")|| dst.includes("neutron-node-test-chain1")) {
    executedMsg = await verifyReceivedMessageCosmos(height!);
  } else if (dst.includes("icon")) {
    executedMsg = await verifyReceivedMessageIcon(height)
  }  
  if (executedMsg! === msgData) {
  } else {
    throw new Error(
      "Received Different Message. Message sent from source is : " + msgData
    );
  }
}


async function checkResponseMessage(
  src: string,
  expectRevert: boolean
): Promise<[number, any] | undefined> {
  console.log("**** Verify ResponseMessage Event ****");
  if (src.includes("icon")) {
    const [seqNo, height] = await verifyResponseMessageEventIcon();
    return [height, seqNo];
  } else if (src.includes("constantine-3") || src.includes("neutron-node-test-chain1")) {
    const [seqNo, height] = await verifyResponseMessageEventCosmos();
    return [height, seqNo];
  }
}

async function checkRollbackMessage(src: string, height: number) {
  console.log("**** Verify RollbackMessage Event ****");
  if (src.includes("icon")) {
    await verifyRollbackMessageEventIcon(height);
  } else if (src.includes("constantine-3") || src.includes("neutron-node-test-chain1")) {
    await verifyRollbackMessageEventCosmos(height);
  }
}

async function invokeExecuteRollback(src: string, seqNo: number) {
  console.log("**** Execute Rollback ****");
  if (src.includes("icon")) {
    await executeRollbackIcon(seqNo, chainName);
  } else if (src.includes("constantine-3") || src.includes("neutron-node-test-chain1")) {
    await executeRollbackCosmos(seqNo, chainName);
  }
}

async function checkRollbackExecuted(src: string) {
  console.log("**** Verify RollbackExecuted event ****");
  if (src.includes("icon")) {
    await verifyRollbackExecutedEventIcon();
  } else if (src.includes("constantine-3") || src.includes("neutron-node-test-chain1")) {
    await verifyRollbackExecutedEventCosmos();
  }
}
