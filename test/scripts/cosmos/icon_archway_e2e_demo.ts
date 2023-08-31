import { sendMessageFromDAppCosmos } from "./archway";
import { GetDataInBytes, GetDest, GetSrc, strToHex } from "./helper";
import { sendMessageFromDAppIcon } from "./icon";

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
//   .then(() => sendCallMessage(SRC, DST))
//   .then(() => sendCallMessage(DST, SRC))
//   .then(() => sendCallMessage(SRC, DST, "checkSuccessResponse", true))
//   .then(() => sendCallMessage(DST, SRC, "checkSuccessResponse", true))
//   .then(() => sendCallMessage(SRC, DST, "rollback", true))
//   .then(() => sendCallMessage(DST, SRC, "rollback", true))
  .catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });

async function sendCallMessage(
  src: string,
  dst: string,
  msgData?: string,
  needRollback?: boolean
){
  const testName = sendCallMessage.name + (needRollback ? "WithRollback" : "");
  console.log(`\n### ${testName}: ${src} => ${dst}`);
  if (!msgData) {
    msgData = `${testName}_${src}_${dst}`;
  }
  const rollbackData = needRollback ? `ThisIsRollbackMessage_${src}_${dst}` : undefined;
  let step = 1;

  console.log(`[${step++}] send message from DApp`);
  const sendMessageReceipt = await sendMessageFromDApp(src, msgData, rollbackData);
}
async function sendMessageFromDApp(src: string, msg: string, rollback?: string) {
    const isRollback = rollback ? true : false;
    if (src === "icon") {
        const hexMsg = strToHex(msg)
        const receipt = await sendMessageFromDAppIcon(hexMsg, rollback, isRollback);
        return receipt
    } else if (src === "archway") {
        const bytesData = GetDataInBytes(msg);
        const receipt = await sendMessageFromDAppCosmos(bytesData,rollback)
        console.log(receipt)
    }
}

