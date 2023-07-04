pwd=$PWD
readonly JavaScoreContracts="https://github.com/icon-project/btp2-java.git"
readonly SoldityContracts="https://github.com/icon-project/btp2-solidity.git"

mkdir -p tmp
git clone $JavaScoreContracts -q tmp/java && git clone $SoldityContracts -q tmp/solidity


cd tmp/java

./gradlew bmc:optimizedJar bmv:bridge:optimizedJar  bmv:btpblock:optimizedJar dapp-sample:optimizedJar xcall:optimizedJar

mkdir -p $pwd/artifacts/jvm/contracts/

rsync bmv/bridge/build/libs/bmv-bridge-0.1.0-optimized.jar $pwd/artifacts/jvm/contracts/bmv-bridge.jar
rsync bmv/btpblock/build/libs/bmv-btpblock-0.1.0-optimized.jar $pwd/artifacts/jvm/contracts/bmv-btpblock.jar
rsync bmc/build/libs/bmc-0.1.0-optimized.jar $pwd/artifacts/jvm/contracts/bmc.jar
rsync xcall/build/libs/xcall-0.6.2-optimized.jar $pwd/artifacts/jvm/contracts/xcall.jar
rsync dapp-sample/build/libs/dapp-sample-0.1.0-optimized.jar $pwd/artifacts/jvm/contracts/dapp-sample.jar


cd ../solidity
mkdir -p $pwd/artifacts/evm/contracts
rsync -av --progress  bmc/contracts/ $pwd/artifacts/evm/contracts/bmc/ --exclude test
rsync -av --progress  bmv/contracts/ $pwd/artifacts/evm/contracts/bmv/ --exclude test
rsync -av --progress  bmv-bridge/contracts/ $pwd/artifacts/evm/contracts/bmv-bridge/ --exclude test
rsync -av --progress  xcall/contracts/ $pwd/artifacts/evm/contracts/xcall/ --exclude test

curl --create-dirs -o $pwd/artifacts/evm/contracts/dapp-sample/DAppProxySample.sol https://raw.githubusercontent.com/icon-project/btp2/main/e2edemo/solidity/contracts/dapp-sample/DAppProxySample.sol

rm -rf $pwd/tmp
